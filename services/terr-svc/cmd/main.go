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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/terr-svc/internal/handler"
	"github.com/snisid/platform/services/terr-svc/internal/kafka"
	"github.com/snisid/platform/services/terr-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/terr-svc/internal/service"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	port := getEnv("TERR_SERVICE_PORT", "8098")
	dbURL := getEnv("TERR_DB_URL", "postgresql://root@localhost:26257/snisid_terr?sslmode=disable")
	kafkaBrokers := getEnv("TERR_KAFKA_BROKERS", "localhost:9092")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Fatal("failed to ping database", zap.Error(err))
	}
	log.Println("connected to CockroachDB")

	runMigrations(ctx, pool)

	producer := kafka.NewProducer([]string{kafkaBrokers})
	defer producer.Close()
	_ = producer

	repo := postgres.NewTerritoryRepo(pool)
	svc := service.NewTerritoryService(repo, logger)

	r := handler.NewRouter(svc, logger)

	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		logger.Info("TERR-HT service starting", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}
	logger.Info("server exited")
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func runMigrations(ctx context.Context, pool *pgxpool.Pool) {
	data, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		log.Printf("warning: could not read migration file: %v", err)
		return
	}
	if _, err := pool.Exec(ctx, string(data)); err != nil {
		log.Printf("warning: migration error: %v", err)
	}
}
