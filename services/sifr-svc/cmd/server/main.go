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

	"github.com/snisid/platform/services/sifr-svc/internal/handler"
	"github.com/snisid/platform/services/sifr-svc/internal/kafka"
	"github.com/snisid/platform/services/sifr-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/sifr-svc/internal/service"
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
		`CREATE TYPE IF NOT EXISTS sifr_crossing_direction AS ENUM ('ENTRY', 'EXIT')`,
		`CREATE TYPE IF NOT EXISTS sifr_doc_type AS ENUM (
			'PASSPORT', 'NATIONAL_ID', 'LAISSEZ_PASSER',
			'BIRTH_CERTIFICATE', 'TRAVEL_DOCUMENT', 'NONE'
		)`,
		`CREATE TYPE IF NOT EXISTS sifr_alert_type AS ENUM (
			'WANTED_PERSON', 'STOLEN_DOCUMENT', 'BLACKLIST',
			'ACTIVE_WARRANT', 'SANCTIONS', 'CUSTOMS_ALERT'
		)`,
		`CREATE TABLE IF NOT EXISTS sifr_border_posts (
			post_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			post_code VARCHAR(10) UNIQUE NOT NULL,
			name VARCHAR(150) NOT NULL, dept_code CHAR(2) NOT NULL,
			border_country CHAR(3) NOT NULL DEFAULT 'DOM',
			post_lat DECIMAL(10,7), post_lng DECIMAL(10,7),
			is_official BOOLEAN DEFAULT TRUE,
			is_active BOOLEAN DEFAULT TRUE,
			lanes_count SMALLINT DEFAULT 2,
			has_biometric_scanner BOOLEAN DEFAULT FALSE,
			has_vehicle_scanner BOOLEAN DEFAULT FALSE,
			operating_hours VARCHAR(50), commanding_officer UUID,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS sifr_crossings (
			crossing_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			post_id UUID NOT NULL REFERENCES sifr_border_posts(post_id),
			direction sifr_crossing_direction NOT NULL,
			crossing_datetime TIMESTAMPTZ NOT NULL,
			snisid_person_id UUID,
			document_type sifr_doc_type NOT NULL DEFAULT 'PASSPORT',
			document_number VARCHAR(100), document_country CHAR(3),
			document_expiry DATE, traveler_name VARCHAR(200) NOT NULL,
			traveler_dob DATE, traveler_nationality CHAR(3),
			vehicle_plate VARCHAR(20), lane_number SMALLINT,
			processing_officer UUID NOT NULL,
			alert_triggered BOOLEAN DEFAULT FALSE,
			alert_type sifr_alert_type, alert_action_taken TEXT,
			processing_time_sec INTEGER,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS sifr_alerts_log (
			alert_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			crossing_id UUID REFERENCES sifr_crossings(crossing_id),
			post_id UUID NOT NULL, alert_type sifr_alert_type NOT NULL,
			snisid_person_id UUID, document_number VARCHAR(100),
			vehicle_plate VARCHAR(20), alert_source VARCHAR(50),
			source_record_id UUID, notified_units TEXT[] DEFAULT '{}',
			action_taken TEXT, resolved BOOLEAN DEFAULT FALSE,
			resolved_by UUID, resolved_at TIMESTAMPTZ,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS sifr_clandestine_crossings (
			report_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			location_desc VARCHAR(300), dept_code CHAR(2),
			lat DECIMAL(10,7), lng DECIMAL(10,7),
			reported_date TIMESTAMPTZ NOT NULL,
			crossing_type VARCHAR(50), estimated_persons INTEGER,
			gang_related BOOLEAN DEFAULT FALSE, gang_id UUID,
			trafficking_type TEXT, reported_by UUID NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
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

	dbHost := getEnv("SIFR_DB_HOST", "localhost")
	dbPort := getEnv("SIFR_DB_PORT", "26257")
	dbName := getEnv("SIFR_DB_NAME", "snisid_sifr")
	dbUser := getEnv("SIFR_DB_USER", "root")
	dbURL := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", dbUser, dbHost, dbPort, dbName)
	if u := os.Getenv("SIFR_DATABASE_URL"); u != "" {
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

	repo := postgres.NewBorderRepo(pool)
	svc := service.NewBorderService(repo, logger)

	r := handler.SetupRouter(svc, logger)

	kafkaBrokers := strings.Split(getEnv("SIFR_KAFKA_BROKERS", "kafka:9092"), ",")
	kafkaTopic := getEnv("SIFR_KAFKA_TOPIC", "sifr.events")
	kafkaProducer := kafka.NewProducer(kafkaBrokers, kafkaTopic, logger)
	defer kafkaProducer.Close()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	port := getEnv("SIFR_SERVICE_PORT", ":8106")
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
		logger.Info("starting sifr-svc", zap.String("addr", port))
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
