package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/sigeo-svc/internal/handler"
	"github.com/snisid/platform/services/sigeo-svc/internal/kafka"
	"github.com/snisid/platform/services/sigeo-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/sigeo-svc/internal/service"
)

func runMigrations(pool *pgxpool.Pool) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `
	CREATE EXTENSION IF NOT EXISTS postgis;
	CREATE TABLE IF NOT EXISTS sigeo_incidents_unified (
		event_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		source_module VARCHAR(20) NOT NULL,
		source_record_id UUID NOT NULL,
		event_type VARCHAR(50) NOT NULL,
		event_date TIMESTAMPTZ NOT NULL,
		lat DECIMAL(10,7) NOT NULL,
		lng DECIMAL(10,7) NOT NULL,
		dept_code CHAR(2),
		commune VARCHAR(100),
		severity SMALLINT,
		gang_id UUID,
		description TEXT,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE TABLE IF NOT EXISTS sigeo_checkpoints (
		cp_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		cp_type VARCHAR(30) NOT NULL,
		dept_code CHAR(2),
		road_number VARCHAR(10),
		description VARCHAR(300),
		controlling_gang_id UUID,
		is_active BOOLEAN DEFAULT TRUE,
		source_module VARCHAR(20),
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`)
	return err
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/snisid_sigeo?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		logger.Fatal("failed to ping database", zap.Error(err))
	}

	if err := runMigrations(pool); err != nil {
		logger.Fatal("auto-migration failed", zap.Error(err))
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "kafka:9092"
	}
	producer := kafka.NewProducer(kafkaBrokers, logger)
	defer producer.Close()

	repo := postgres.NewGeoRepo(pool)
	svc := service.NewGeoIntelService(repo, logger)
	h := handler.NewGeoHandler(svc, logger)

	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "sigeo-svc", "uptime": time.Since(startTime).String()})
	})

	api := r.Group("/api/v1/sigeo")
	{
		api.GET("/incidents/unified", h.ListIncidents)
		api.POST("/incidents/ingest", h.IngestIncident)
		api.GET("/checkpoints/active", h.ListCheckpoints)
		api.GET("/zone-report", h.GetZoneReport)
	}

	port := os.Getenv("SIGEO_SERVICE_PORT")
	if port == "" {
		port = ":8125"
	}

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		logger.Info("starting sigeo-svc", zap.String("addr", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}
	logger.Info("server exited")
}

var startTime = time.Now()
