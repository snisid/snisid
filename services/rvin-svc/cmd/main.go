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

	"github.com/snisid/platform/services/rvin-svc/internal/handler"
	"github.com/snisid/platform/services/rvin-svc/internal/kafka"
	"github.com/snisid/platform/services/rvin-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/rvin-svc/internal/service"
)

func runMigrations(pool *pgxpool.Pool) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS rvin_unidentified_remains (
		remains_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		national_rvin_id VARCHAR(25) UNIQUE NOT NULL,
		discovery_date TIMESTAMPTZ NOT NULL,
		discovery_location VARCHAR(300) NOT NULL,
		dept_code CHAR(2) NOT NULL,
		commune VARCHAR(100),
		lat DECIMAL(10,7),
		lng DECIMAL(10,7),
		discovery_source VARCHAR(30) NOT NULL,
		status VARCHAR(30) NOT NULL DEFAULT 'UNIDENTIFIED',
		estimated_sex VARCHAR(10),
		estimated_age_min SMALLINT,
		estimated_age_max SMALLINT,
		estimated_height_cm SMALLINT,
		skin_tone VARCHAR(30),
		distinguishing_marks TEXT,
		decomposition_level SMALLINT,
		dna_sample_taken BOOLEAN DEFAULT FALSE,
		dna_sample_ref VARCHAR(100),
		morgue_location VARCHAR(200),
		case_reference VARCHAR(100),
		examiner_id UUID NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE TABLE IF NOT EXISTS rvin_dna_comparisons (
		comparison_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		remains_id UUID NOT NULL REFERENCES rvin_unidentified_remains(remains_id),
		reference_dna_ref VARCHAR(100),
		comparison_date TIMESTAMPTZ NOT NULL,
		match_probability DECIMAL(10,8),
		is_match BOOLEAN DEFAULT FALSE,
		lab_reference VARCHAR(100),
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`)
	return err
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/snisid_rvin?sslmode=disable"
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

	repo := postgres.NewRemainsRepo(pool)
	svc := service.NewRemainsService(repo, logger)
	h := handler.NewRemainsHandler(svc, logger)

	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "rvin-svc", "uptime": time.Since(startTime).String()})
	})

	api := r.Group("/api/v1/rvin")
	{
		api.POST("/remains", h.RegisterRemains)
		api.GET("/remains/:id", h.GetRemains)
		api.POST("/remains/:id/dna", h.SubmitDNA)
		api.GET("/unidentified", h.ListUnidentified)
		api.GET("/stats/by-source", h.GetStatsBySource)
	}

	port := os.Getenv("RVIN_SERVICE_PORT")
	if port == "" {
		port = ":8120"
	}

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		logger.Info("starting rvin-svc", zap.String("addr", port))
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
