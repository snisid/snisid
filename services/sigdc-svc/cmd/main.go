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

	"github.com/snisid/platform/services/sigdc-svc/internal/handler"
	"github.com/snisid/platform/services/sigdc-svc/internal/kafka"
	"github.com/snisid/platform/services/sigdc-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/sigdc-svc/internal/service"
)

func runMigrations(pool *pgxpool.Pool) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS sigdc_disasters (
		disaster_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		national_sigdc_id VARCHAR(25) UNIQUE NOT NULL,
		disaster_type VARCHAR(30) NOT NULL,
		disaster_name VARCHAR(200),
		alert_level VARCHAR(20) NOT NULL,
		status VARCHAR(20) DEFAULT 'ACTIVE',
		onset_date TIMESTAMPTZ NOT NULL,
		affected_depts TEXT[] DEFAULT '{}',
		epicenter_lat DECIMAL(10,7),
		epicenter_lng DECIMAL(10,7),
		magnitude DECIMAL(4,2),
		estimated_affected INTEGER,
		confirmed_dead INTEGER DEFAULT 0,
		confirmed_injured INTEGER DEFAULT 0,
		confirmed_missing INTEGER DEFAULT 0,
		response_agencies TEXT[] DEFAULT '{}',
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE TABLE IF NOT EXISTS sigdc_early_warnings (
		warning_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		disaster_type VARCHAR(30) NOT NULL,
		alert_level VARCHAR(20) NOT NULL,
		source_agency VARCHAR(100),
		message_text TEXT NOT NULL,
		affected_depts TEXT[] DEFAULT '{}',
		issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		expires_at TIMESTAMPTZ,
		channels_sent TEXT[] DEFAULT '{}'
	);
	CREATE TABLE IF NOT EXISTS sigdc_victim_registrations (
		registration_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		disaster_id UUID NOT NULL REFERENCES sigdc_disasters(disaster_id),
		full_name VARCHAR(200),
		status VARCHAR(30) NOT NULL,
		location_found VARCHAR(300),
		dept_code CHAR(2),
		registration_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		registered_by UUID NOT NULL
	);
	CREATE TABLE IF NOT EXISTS sigdc_resources (
		resource_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		disaster_id UUID NOT NULL REFERENCES sigdc_disasters(disaster_id),
		resource_type VARCHAR(50) NOT NULL,
		provider_org VARCHAR(150),
		quantity INTEGER,
		dept_code CHAR(2),
		status VARCHAR(20) DEFAULT 'AVAILABLE',
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`)
	return err
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/snisid_sigdc?sslmode=disable"
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

	repo := postgres.NewDisasterRepo(pool)
	svc := service.NewDisasterService(repo, logger)
	h := handler.NewDisasterHandler(svc, logger)

	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "sigdc-svc", "uptime": time.Since(startTime).String()})
	})

	api := r.Group("/api/v1/sigdc")
	{
		api.POST("/disasters", h.DeclareDisaster)
		api.GET("/disasters/active", h.ListActiveDisasters)
		api.POST("/warnings", h.IssueWarning)
		api.POST("/victims", h.RegisterVictim)
		api.GET("/resources/available", h.ListResources)
	}

	port := os.Getenv("SIGDC_SERVICE_PORT")
	if port == "" {
		port = ":8126"
	}

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		logger.Info("starting sigdc-svc", zap.String("addr", port))
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
