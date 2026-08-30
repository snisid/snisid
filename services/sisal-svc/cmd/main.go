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

	"github.com/snisid/platform/services/sisal-svc/internal/handler"
	"github.com/snisid/platform/services/sisal-svc/internal/kafka"
	"github.com/snisid/platform/services/sisal-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/sisal-svc/internal/service"
)

func runMigrations(pool *pgxpool.Pool) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS sisal_alerts (
		alert_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		national_sisal_id VARCHAR(25) UNIQUE NOT NULL,
		hazard_type VARCHAR(30) NOT NULL,
		severity VARCHAR(20) NOT NULL,
		title VARCHAR(200) NOT NULL,
		message_fr TEXT NOT NULL,
		message_ht TEXT NOT NULL,
		affected_depts TEXT[] DEFAULT '{}',
		affected_pop_est INTEGER,
		issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		valid_until TIMESTAMPTZ,
		source_agency VARCHAR(100) NOT NULL,
		source_event_id UUID,
		is_cancelled BOOLEAN DEFAULT FALSE,
		cancelled_at TIMESTAMPTZ,
		cancel_reason TEXT,
		created_by UUID NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE TABLE IF NOT EXISTS sisal_subscriptions (
		sub_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		snisid_person_id UUID,
		phone_number VARCHAR(30),
		email VARCHAR(200),
		dept_code CHAR(2),
		commune VARCHAR(100),
		min_severity VARCHAR(20) DEFAULT 'WARNING',
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`)
	return err
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/snisid_sisal?sslmode=disable"
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

	repo := postgres.NewAlertRepo(pool)
	svc := service.NewAlertService(repo, logger)
	h := handler.NewAlertHandler(svc, logger)

	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "sisal-svc", "uptime": time.Since(startTime).String()})
	})

	api := r.Group("/api/v1/sisal")
	{
		api.POST("/alerts", h.IssueAlert)
		api.GET("/alerts/active", h.ListActiveAlerts)
		api.GET("/alerts/history", h.ListHistory)
		api.POST("/alerts/:id/cancel", h.CancelAlert)
		api.POST("/subscribe", h.Subscribe)
	}

	port := os.Getenv("SISAL_SERVICE_PORT")
	if port == "" {
		port = ":8128"
	}

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		logger.Info("starting sisal-svc", zap.String("addr", port))
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
