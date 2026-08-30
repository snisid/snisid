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

	"github.com/snisid/platform/services/cybre-svc/internal/handler"
	"github.com/snisid/platform/services/cybre-svc/internal/kafka"
	"github.com/snisid/platform/services/cybre-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/cybre-svc/internal/service"
)

func runMigrations(pool *pgxpool.Pool) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS cybre_incidents (
		incident_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		national_cybre_id VARCHAR(25) UNIQUE NOT NULL,
		crime_type VARCHAR(50) NOT NULL,
		severity VARCHAR(20) NOT NULL DEFAULT 'MEDIUM',
		status VARCHAR(20) DEFAULT 'OPEN',
		victim_count INTEGER DEFAULT 1,
		total_financial_loss_usd DECIMAL(15,2),
		incident_date TIMESTAMPTZ NOT NULL,
		reported_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		attack_vector TEXT,
		targeted_platform VARCHAR(100),
		suspect_phone TEXT[] DEFAULT '{}',
		suspect_email TEXT[] DEFAULT '{}',
		case_reference VARCHAR(100),
		created_by UUID NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE TABLE IF NOT EXISTS cybre_intrusion_attempts (
		attempt_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		incident_id UUID,
		target_system VARCHAR(100) NOT NULL,
		attack_timestamp TIMESTAMPTZ NOT NULL,
		attack_type VARCHAR(50),
		source_ip_hash VARCHAR(64),
		source_country CHAR(3),
		was_successful BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE TABLE IF NOT EXISTS cybre_threat_intelligence (
		threat_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		indicator_type VARCHAR(30) NOT NULL,
		indicator_value VARCHAR(500) NOT NULL,
		threat_category VARCHAR(50),
		confidence_score SMALLINT,
		source VARCHAR(100),
		is_active BOOLEAN DEFAULT TRUE,
		first_seen TIMESTAMPTZ,
		last_seen TIMESTAMPTZ,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`)
	return err
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/snisid_cybre?sslmode=disable"
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

	repo := postgres.NewCybreRepo(pool)
	svc := service.NewCybreService(repo, logger)
	h := handler.NewCybreHandler(svc, logger)

	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "cybre-svc", "uptime": time.Since(startTime).String()})
	})

	api := r.Group("/api/v1/cybre")
	{
		api.POST("/incidents", h.DeclareIncident)
		api.GET("/incidents/:id", h.GetIncident)
		api.GET("/intrusions/recent", h.ListRecentIntrusions)
		api.POST("/threat-intel", h.AddThreatIntel)
		api.GET("/threat-intel/check", h.CheckIndicator)
		api.GET("/stats/by-type", h.GetStatsByType)
	}

	port := os.Getenv("CYBRE_SERVICE_PORT")
	if port == "" {
		port = ":8132"
	}

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		logger.Info("starting cybre-svc", zap.String("addr", port))
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
