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

	"github.com/snisid/platform/services/corr-svc/internal/handler"
	"github.com/snisid/platform/services/corr-svc/internal/kafka"
	"github.com/snisid/platform/services/corr-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/corr-svc/internal/service"
)

func runMigrations(pool *pgxpool.Pool) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS corr_integrity_cases (
		case_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		national_corr_id VARCHAR(25) UNIQUE NOT NULL,
		officer_snisid_id UUID NOT NULL,
		officer_badge VARCHAR(30),
		officer_unit VARCHAR(50),
		officer_rank VARCHAR(50),
		allegation_type VARCHAR(30) NOT NULL,
		severity VARCHAR(20) NOT NULL,
		status VARCHAR(30) NOT NULL DEFAULT 'REPORTED',
		allegation_summary TEXT NOT NULL,
		incident_date_from TIMESTAMPTZ,
		incident_date_to TIMESTAMPTZ,
		gang_id UUID,
		financial_gain_usd DECIMAL(15,2),
		reported_by_type VARCHAR(30),
		reporting_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		is_whistleblower BOOLEAN DEFAULT FALSE,
		created_by UUID NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE TABLE IF NOT EXISTS corr_whistleblower_reports (
		report_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		report_token VARCHAR(64) UNIQUE NOT NULL,
		allegation_type VARCHAR(30) NOT NULL,
		severity_estimate VARCHAR(20),
		description TEXT NOT NULL,
		submission_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		processed BOOLEAN DEFAULT FALSE,
		integrity_case_id UUID,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE TABLE IF NOT EXISTS corr_behavioral_alerts (
		alert_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		officer_snisid_id UUID NOT NULL,
		alert_type VARCHAR(50) NOT NULL,
		description TEXT NOT NULL,
		module_source VARCHAR(30),
		risk_score SMALLINT,
		auto_generated BOOLEAN DEFAULT TRUE,
		reviewed BOOLEAN DEFAULT FALSE,
		is_false_positive BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE TABLE IF NOT EXISTS corr_asset_declarations (
		declaration_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		officer_snisid_id UUID NOT NULL,
		declaration_year SMALLINT NOT NULL,
		real_estate_usd DECIMAL(15,2) DEFAULT 0,
		vehicles_usd DECIMAL(15,2) DEFAULT 0,
		bank_accounts_usd DECIMAL(15,2) DEFAULT 0,
		other_assets_usd DECIMAL(15,2) DEFAULT 0,
		is_flagged BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`)
	return err
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/snisid_corr?sslmode=disable"
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

	repo := postgres.NewIntegrityRepo(pool)
	svc := service.NewIntegrityService(repo, logger)
	h := handler.NewIntegrityHandler(svc, logger)

	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "corr-svc", "uptime": time.Since(startTime).String()})
	})

	api := r.Group("/api/v1/corr")
	{
		api.POST("/cases", h.OpenCase)
		api.GET("/cases/:id", h.GetCase)
		api.GET("/cases/active", h.ListActiveCases)
		api.POST("/whistleblower", h.SubmitWhistleblower)
		api.GET("/whistleblower/:token", h.TrackWhistleblower)
		api.GET("/alerts/behavioral", h.ListBehavioralAlerts)
		api.POST("/declarations", h.SubmitDeclaration)
		api.GET("/declarations/flagged", h.ListFlaggedDeclarations)
	}

	port := os.Getenv("CORR_SERVICE_PORT")
	if port == "" {
		port = ":8130"
	}

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		logger.Info("starting corr-svc", zap.String("addr", port))
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
