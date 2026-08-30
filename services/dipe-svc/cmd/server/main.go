package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/dipe-svc/internal/api/rest"
	"github.com/snisid/platform/services/dipe-svc/internal/handler"
	"github.com/snisid/platform/services/dipe-svc/internal/kafka"
	"github.com/snisid/platform/services/dipe-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/dipe-svc/internal/service"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	port := os.Getenv("DIPE_SERVICE_PORT")
	if port == "" {
		port = ":8118"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://root@localhost:26257/snisid_dipe?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		logger.Fatal("failed to ping database", zap.Error(err))
	}

	brokersStr := os.Getenv("KAFKA_BROKERS")
	if brokersStr == "" {
		brokersStr = "kafka:9092"
	}
	brokers := strings.Split(brokersStr, ",")

	kafkaProducer, err := kafka.NewProducer(brokers, logger)
	if err != nil {
		logger.Warn("kafka producer not available, continuing without it", zap.Error(err))
	}
	if kafkaProducer != nil {
		defer kafkaProducer.Close()
	}

	repo := postgres.NewMissingRepo(pool)
	svc := service.NewMissingPersonService(repo, logger)
	r := rest.NewRouter(svc, logger)
	healthHandler := handler.NewHealthHandler(pool, logger)

	r.GET("/healthz", healthHandler.Healthz)
	r.GET("/metrics", handler.MetricsHandler())

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		logger.Info("starting dipe-svc", zap.String("addr", port))
		fmt.Printf("dipe-svc listening on %s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("shutting down dipe-svc...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}
	logger.Info("server exited")
}
