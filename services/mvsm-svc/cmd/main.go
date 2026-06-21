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

	"github.com/snisid/platform/services/mvsm-svc/internal/handler"
	"github.com/snisid/platform/services/mvsm-svc/internal/kafka"
	"github.com/snisid/platform/services/mvsm-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/mvsm-svc/internal/service"
)

func runMigrations(pool *pgxpool.Pool) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS mvsm_events (
		event_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		national_mvsm_id VARCHAR(25) UNIQUE NOT NULL,
		event_type VARCHAR(30) NOT NULL,
		event_name VARCHAR(200),
		risk_level VARCHAR(20) NOT NULL DEFAULT 'LOW',
		status VARCHAR(20) DEFAULT 'PLANNED',
		organizer_name VARCHAR(200),
		gang_id UUID,
		scheduled_date TIMESTAMPTZ NOT NULL,
		actual_start TIMESTAMPTZ,
		actual_end TIMESTAMPTZ,
		location_desc VARCHAR(300),
		dept_code CHAR(2),
		commune VARCHAR(100),
		lat DECIMAL(10,7),
		lng DECIMAL(10,7),
		estimated_crowd INTEGER,
		peak_crowd INTEGER,
		incidents_during INTEGER DEFAULT 0,
		casualties INTEGER DEFAULT 0,
		arrests_made INTEGER DEFAULT 0,
		weapons_found INTEGER DEFAULT 0,
		created_by UUID NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE TABLE IF NOT EXISTS mvsm_real_time_updates (
		update_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		event_id UUID NOT NULL REFERENCES mvsm_events(event_id),
		update_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		current_crowd_est INTEGER,
		situation TEXT NOT NULL,
		risk_change VARCHAR(20),
		action_taken TEXT,
		reported_by UUID NOT NULL,
		lat DECIMAL(10,7),
		lng DECIMAL(10,7)
	);`)
	return err
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/snisid_mvsm?sslmode=disable"
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

	repo := postgres.NewEventRepo(pool)
	svc := service.NewEventService(repo, logger)
	h := handler.NewEventHandler(svc, logger)

	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "mvsm-svc", "uptime": time.Since(startTime).String()})
	})

	api := r.Group("/api/v1/mvsm")
	{
		api.POST("/events", h.CreateEvent)
		api.GET("/events/upcoming", h.ListUpcoming)
		api.POST("/events/:id/updates", h.AddUpdate)
		api.GET("/events/active", h.ListActive)
		api.PATCH("/events/:id/risk", h.UpdateRiskLevel)
	}

	port := os.Getenv("MVSM_SERVICE_PORT")
	if port == "" {
		port = ":8127"
	}

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		logger.Info("starting mvsm-svc", zap.String("addr", port))
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
