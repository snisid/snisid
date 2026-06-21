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

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/blkl-svc/internal/handler"
	"github.com/snisid/platform/services/blkl-svc/internal/kafka"
	"github.com/snisid/platform/services/blkl-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/blkl-svc/internal/service"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func runMigrations(pool *pgxpool.Pool, logger *zap.Logger) {
	ctx := context.Background()
	migrations := []string{
		`CREATE TYPE IF NOT EXISTS blkl_restriction_type AS ENUM (
			'ENTRY_BAN','EXIT_BAN','BOTH_BAN','CONDITIONAL_BAN'
		)`,
		`CREATE TYPE IF NOT EXISTS blkl_source AS ENUM (
			'JUDICIAL_ORDER','WANTED_WARRANT','UN_SANCTIONS',
			'OFAC_SANCTIONS','MINISTERIAL_ORDER','EXPULSION',
			'OPR_TRAVEL_RESTRICTION','INTERPOL_NOTICE'
		)`,
		`CREATE TABLE IF NOT EXISTS blkl_blacklist (
			entry_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			national_blkl_id VARCHAR(25) UNIQUE NOT NULL,
			snisid_person_id UUID NOT NULL,
			restriction_type blkl_restriction_type NOT NULL,
			source blkl_source NOT NULL, source_record_id UUID,
			reason TEXT NOT NULL, court_order_ref VARCHAR(100),
			ordered_by VARCHAR(150),
			effective_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			expiry_date TIMESTAMPTZ, is_permanent BOOLEAN DEFAULT FALSE,
			is_active BOOLEAN DEFAULT TRUE,
			alert_level VARCHAR(20) DEFAULT 'WANTED',
			armed_dangerous BOOLEAN DEFAULT FALSE,
			created_by UUID NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS blkl_alerts_log (
			alert_log_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			entry_id UUID NOT NULL REFERENCES blkl_blacklist(entry_id),
			triggered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			post_code VARCHAR(10), direction VARCHAR(10),
			action_taken TEXT, officer_id UUID, outcome TEXT
		)`,
	}
	for _, m := range migrations {
		if _, err := pool.Exec(ctx, m); err != nil {
			logger.Warn("migration warning", zap.Error(err))
		}
	}
	logger.Info("migrations completed")
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	dbHost := getEnv("BLKL_DB_HOST", "localhost")
	dbPort := getEnv("BLKL_DB_PORT", "26257")
	dbName := getEnv("BLKL_DB_NAME", "snisid_blkl")
	dbUser := getEnv("BLKL_DB_USER", "root")
	dbURL := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", dbUser, dbHost, dbPort, dbName)
	if u := os.Getenv("BLKL_DATABASE_URL"); u != "" {
		dbURL = u
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		logger.Fatal("failed to ping database", zap.Error(err))
	}

	runMigrations(pool, logger)

	repo := postgres.NewBlacklistRepo(pool)
	svc := service.NewBLKLService(repo, logger)

	r := handler.SetupRouter(svc, logger)

	kafkaBrokers := strings.Split(getEnv("BLKL_KAFKA_BROKERS", "kafka:9092"), ",")
	kafkaTopic := getEnv("BLKL_KAFKA_TOPIC", "blkl.events")
	kafkaProducer := kafka.NewProducer(kafkaBrokers, kafkaTopic, logger)
	defer kafkaProducer.Close()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	port := getEnv("BLKL_SERVICE_PORT", ":8110")
	if port != "" && port[0] != ':' {
		port = ":" + port
	}

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		logger.Info("starting blkl-svc", zap.String("addr", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}
	logger.Info("server exited")
}
