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
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/fir-svc/internal/api/rest"
	"github.com/snisid/platform/services/fir-svc/internal/domain"
	"github.com/snisid/platform/services/fir-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/fir-svc/internal/service"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dbHost := getEnv("FIR_DB_HOST", "postgres://snisid:snisid@localhost:5432/snisid_fir")
	dbPool, err := pgxpool.New(ctx, dbHost)
	if err != nil {
		logger.Fatal("cannot connect to database", zap.Error(err))
	}
	defer dbPool.Close()

	recordRepo := postgres.NewCriminalRecordRepo(dbPool)
	arrestRepo := postgres.NewArrestRepo(dbPool)
	convictionRepo := postgres.NewConvictionRepo(dbPool)

	var eventPub domain.EventPublisher = &noopEventPublisher{}
	var snisidClient domain.SNISIDClient = &noopSNISIDClient{}

	recordSvc := service.NewRecordService(recordRepo, arrestRepo, convictionRepo, eventPub, snisidClient)
	certSvc := service.NewCertificateService(recordRepo)
	expungementSvc := service.NewExpungementService(recordRepo, eventPub)
	_ = expungementSvc

	recordHandler := rest.NewRecordHandler(recordSvc)
	certHandler := rest.NewCertificateHandler(certSvc)
	searchHandler := rest.NewSearchHandler(recordSvc)

	r := gin.Default()
	v1 := r.Group("/api/v1/fir")
	{
		v1.POST("/records", recordHandler.CreateRecord)
		v1.GET("/records/:person_id", recordHandler.GetRecord)
		v1.POST("/records/:id/arrests", recordHandler.AddArrest)
		v1.POST("/records/:id/convictions", recordHandler.AddConviction)
		v1.GET("/records/:id/arrests", recordHandler.GetArrests)
		v1.GET("/records/:id/convictions", recordHandler.GetConvictions)
		v1.POST("/certificates/issue", certHandler.IssueCertificate)
		v1.GET("/certificates/verify/:num", certHandler.VerifyCertificate)
		v1.GET("/search", searchHandler.Search)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", getEnv("FIR_SERVICE_PORT", "8093")),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen error", zap.Error(err))
		}
	}()

	log.Println("FIR-HT service started on port", getEnv("FIR_SERVICE_PORT", "8093"))

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

type noopEventPublisher struct{}

func (n *noopEventPublisher) Publish(topic string, event interface{}) error {
	log.Printf("Event published to %s: %v", topic, event)
	return nil
}

type noopSNISIDClient struct{}

func (n *noopSNISIDClient) GetPerson(personID uuid.UUID) (*domain.PersonInfo, error) {
	return &domain.PersonInfo{
		PersonID:    personID,
		FullName:    "Test User",
		Nationality: "HTI",
	}, nil
}
