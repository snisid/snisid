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

	"github.com/snisid/platform/services/sipci-svc/internal/handler"
	"github.com/snisid/platform/services/sipci-svc/internal/kafka"
	"github.com/snisid/platform/services/sipci-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/sipci-svc/internal/service"
)

func runMigrations(pool *pgxpool.Pool) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS sipci_assets (
		asset_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		national_sipci_id VARCHAR(25) UNIQUE NOT NULL,
		asset_name VARCHAR(200) NOT NULL,
		asset_category VARCHAR(30) NOT NULL,
		owner_entity VARCHAR(200),
		operating_org VARCHAR(200),
		dept_code CHAR(2) NOT NULL,
		commune VARCHAR(100),
		lat DECIMAL(10,7) NOT NULL,
		lng DECIMAL(10,7) NOT NULL,
		criticality_score SMALLINT,
		population_served INTEGER,
		single_point_failure BOOLEAN DEFAULT FALSE,
		current_threat_level VARCHAR(20) NOT NULL DEFAULT 'NORMAL',
		is_in_gang_zone BOOLEAN DEFAULT FALSE,
		controlling_gang_id UUID,
		under_extortion BOOLEAN DEFAULT FALSE,
		incident_count_12m INTEGER DEFAULT 0,
		protection_unit VARCHAR(50),
		site_manager_phone VARCHAR(30),
		created_by UUID NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE TABLE IF NOT EXISTS sipci_incidents (
		incident_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		asset_id UUID NOT NULL REFERENCES sipci_assets(asset_id),
		incident_type VARCHAR(50) NOT NULL,
		incident_date TIMESTAMPTZ NOT NULL,
		perpetrator_type VARCHAR(30),
		gang_id UUID,
		description TEXT NOT NULL,
		impact_severity SMALLINT,
		population_affected INTEGER,
		economic_loss_usd DECIMAL(15,2),
		case_reference VARCHAR(100),
		created_by UUID NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`)
	return err
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/snisid_sipci?sslmode=disable"
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

	repo := postgres.NewAssetRepo(pool)
	svc := service.NewInfraService(repo, logger)
	h := handler.NewInfraHandler(svc, logger)

	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "sipci-svc", "uptime": time.Since(startTime).String()})
	})

	api := r.Group("/api/v1/sipci")
	{
		api.GET("/assets", h.ListAssets)
		api.GET("/assets/:id", h.GetAsset)
		api.GET("/assets/critical", h.ListCritical)
		api.GET("/assets/under-threat", h.ListUnderThreat)
		api.POST("/assets", h.RegisterAsset)
		api.POST("/incidents", h.ReportIncident)
		api.GET("/incidents/recent", h.ListRecentIncidents)
		api.POST("/assets/:id/assess", h.AssessRisk)
	}

	port := os.Getenv("SIPCI_SERVICE_PORT")
	if port == "" {
		port = ":8131"
	}

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		logger.Info("starting sipci-svc", zap.String("addr", port))
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
