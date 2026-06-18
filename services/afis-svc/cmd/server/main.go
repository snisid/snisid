package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	milvusSDK "github.com/milvus-io/milvus-sdk-go/v2/client"

	"github.com/snisid/platform/services/afis-svc/internal/api/rest"
	"github.com/snisid/platform/services/afis-svc/internal/matcher"
	"github.com/snisid/platform/services/afis-svc/internal/nfiq2"
	"github.com/snisid/platform/services/afis-svc/internal/repository/milvus"
	minioRepo "github.com/snisid/platform/services/afis-svc/internal/repository/minio"
	"github.com/snisid/platform/services/afis-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/afis-svc/internal/service"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dbPool, err := pgxpool.New(ctx, getEnv("AFIS_DB_HOST", "postgres://snisid:snisid@localhost:5432/snisid_afis"))
	if err != nil {
		logger.Fatal("cannot connect to database", zap.Error(err))
	}
	defer dbPool.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: getEnv("AFIS_REDIS_ADDR", "redis-master:6379"),
	})

	minioClient, err := minio.New(getEnv("AFIS_MINIO_ENDPOINT", "minio:9000"), &minio.Options{
		Creds:  credentials.NewStaticV4(getEnv("MINIO_ACCESS_KEY", "snisid"), getEnv("MINIO_SECRET_KEY", "snisid123"), ""),
		Secure: false,
	})
	if err != nil {
		logger.Fatal("cannot connect to minio", zap.Error(err))
	}

	milvusClient, err := milvusSDK.NewClient(ctx, milvusSDK.Config{
		Address: getEnv("AFIS_MILVUS_ADDR", "milvus:19530"),
	})
	if err != nil {
		logger.Fatal("cannot connect to milvus", zap.Error(err))
	}

	fpRepo := postgres.NewFingerprintRepo(dbPool)
	subjectRepo := postgres.NewSubjectRepo(dbPool)
	latentRepo := postgres.NewLatentRepo(dbPool)
	vectorRepo := milvus.NewVectorRepo(milvusClient, getEnv("AFIS_MILVUS_COLLECTION", "afis_fingerprints"), 512)
	imageRepo := minioRepo.NewImageRepo(minioClient, getEnv("AFIS_MINIO_BUCKET", "afis-biometric"))

	qualitySvc := service.NewQualityService(60)
	nfiq2Scorer := nfiq2.NewScorer("/models/nfiq2")
	_ = nfiq2Scorer

	minutiaeMatcher := matcher.NewMinutiaeMatcher(0.85)

	enrollSvc := service.NewEnrollmentService(fpRepo, subjectRepo, imageRepo, vectorRepo, qualitySvc, nil)
	searchSvc := service.NewSearchService(vectorRepo, fpRepo, subjectRepo, minutiaeMatcher, nil)
	latentSvc := service.NewLatentService(latentRepo, imageRepo, searchSvc)

	enrollHandler := rest.NewEnrollHandler(enrollSvc)
	searchHandler := rest.NewSearchHandler(searchSvc)
	latentHandler := rest.NewLatentHandler(latentSvc)
	qualityHandler := rest.NewQualityHandler(qualitySvc)

	r := gin.Default()
	v1 := r.Group("/api/v1/afis")
	{
		v1.POST("/enroll", enrollHandler.Enroll)
		v1.POST("/search/tenprint", searchHandler.SearchTenprint)
		v1.POST("/search/latent", latentHandler.SearchLatent)
		v1.GET("/subjects/:id", enrollHandler.GetSubject)
		v1.POST("/latents", latentHandler.Submit)
		v1.PATCH("/latents/:id/match", latentHandler.ConfirmMatch)
		v1.GET("/quality/check", qualityHandler.CheckQuality)
		v1.GET("/stats", qualityHandler.Stats)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", getEnv("AFIS_SERVICE_PORT", "8091")),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen error", zap.Error(err))
		}
	}()

	_ = rdb

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	srv.Shutdown(shutdownCtx)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
