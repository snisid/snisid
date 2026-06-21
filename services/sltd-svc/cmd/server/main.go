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

	"github.com/snisid/platform/services/sltd-svc/internal/handler"
	"github.com/snisid/platform/services/sltd-svc/internal/kafka"
	"github.com/snisid/platform/services/sltd-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/sltd-svc/internal/service"
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
		`CREATE TYPE IF NOT EXISTS sltd_doc_type AS ENUM (
			'PASSPORT','NATIONAL_ID','TRAVEL_DOCUMENT',
			'VISA','RESIDENCE_PERMIT','REFUGEE_DOCUMENT','LAISSEZ_PASSER'
		)`,
		`CREATE TYPE IF NOT EXISTS sltd_doc_status AS ENUM (
			'LOST','STOLEN','REVOKED','EXPIRED','FOUND','RECOVERED','CANCELLED'
		)`,
		`CREATE TABLE IF NOT EXISTS sltd_documents (
			doc_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			national_sltd_id VARCHAR(25) UNIQUE NOT NULL,
			doc_type sltd_doc_type NOT NULL,
			document_number VARCHAR(100) NOT NULL,
			issuing_country CHAR(3) NOT NULL DEFAULT 'HTI',
			holder_name VARCHAR(200), holder_snisid_id UUID,
			holder_dob DATE, holder_nationality CHAR(3) DEFAULT 'HTI',
			issue_date DATE, expiry_date DATE,
			status sltd_doc_status NOT NULL,
			reported_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			reported_by UUID NOT NULL, reporting_dept_code CHAR(2),
			theft_context TEXT, found_date TIMESTAMPTZ,
			found_location VARCHAR(300), interpol_sltd_ref VARCHAR(50),
			reported_to_interpol BOOLEAN DEFAULT FALSE,
			interpol_reported_at TIMESTAMPTZ,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS sltd_check_log (
			check_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			document_number VARCHAR(100) NOT NULL,
			doc_type sltd_doc_type, checked_by UUID NOT NULL,
			check_location VARCHAR(100), post_id UUID,
			result VARCHAR(20) NOT NULL,
			source VARCHAR(20) NOT NULL,
			sltd_doc_id UUID REFERENCES sltd_documents(doc_id),
			checked_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
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

	dbHost := getEnv("SLTD_DB_HOST", "localhost")
	dbPort := getEnv("SLTD_DB_PORT", "26257")
	dbName := getEnv("SLTD_DB_NAME", "snisid_sltd")
	dbUser := getEnv("SLTD_DB_USER", "root")
	dbURL := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", dbUser, dbHost, dbPort, dbName)
	if u := os.Getenv("SLTD_DATABASE_URL"); u != "" {
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

	repo := postgres.NewDocumentRepo(pool)
	svc := service.NewSLTDService(repo, logger)

	r := handler.SetupRouter(svc, logger)

	kafkaBrokers := strings.Split(getEnv("SLTD_KAFKA_BROKERS", "kafka:9092"), ",")
	kafkaTopic := getEnv("SLTD_KAFKA_TOPIC", "sltd.events")
	kafkaProducer := kafka.NewProducer(kafkaBrokers, kafkaTopic, logger)
	defer kafkaProducer.Close()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	port := getEnv("SLTD_SERVICE_PORT", ":8108")
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
		logger.Info("starting sltd-svc", zap.String("addr", port))
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
