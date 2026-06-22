package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/snisid/enrollment-svc/internal/handler"
	"github.com/snisid/enrollment-svc/internal/kafka"
	"github.com/snisid/enrollment-svc/internal/repository"
	"github.com/snisid/enrollment-svc/internal/service"
)

func main() {
	dbHost := getEnv("ENROLLMENT_DB_HOST", "localhost")
	dbPort := getEnv("ENROLLMENT_DB_PORT", "26257")
	dbName := getEnv("ENROLLMENT_DB_NAME", "snisid_enrollment")
	dbUser := getEnv("ENROLLMENT_DB_USER", "root")
	dbSSLMode := getEnv("ENROLLMENT_DB_SSLMODE", "disable")
	kafkaBrokers := getEnv("ENROLLMENT_KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("ENROLLMENT_KAFKA_TOPIC", "snisid.enrollment.events")
	port := getEnv("ENROLLMENT_SERVICE_PORT", "8091")

	dbURL := fmt.Sprintf("postgresql://%s@%s:%s/%s?sslmode=%s", dbUser, dbHost, dbPort, dbName, dbSSLMode)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to CockroachDB: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping CockroachDB: %v", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := runMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	producer := kafka.NewProducer([]string{kafkaBrokers}, kafkaTopic)
	defer producer.Close()

	repo := repository.NewPostgresRepo(db)
	svc := service.NewEnrollmentService(repo, producer)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := handler.NewHandler(svc)
	api := r.Group("/api/v1/enrollment")
	h.RegisterRoutes(api)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("enrollment-svc started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run enrollment-svc: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down enrollment-svc...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TYPE IF NOT EXISTS identity_proofing_level AS ENUM (
			'IAL_NONE', 'IAL1', 'IAL2', 'IAL3'
		)`,
		`CREATE TYPE IF NOT EXISTS enrollment_status AS ENUM (
			'DRAFT', 'PENDING_DOCUMENTS', 'DOCUMENTS_RECEIVED', 'PENDING_BIOMETRICS',
			'BIOMETRICS_CAPTURED', 'PENDING_REVIEW', 'APPROVED', 'REJECTED', 'EXPIRED'
		)`,
		`CREATE TYPE IF NOT EXISTS document_type AS ENUM (
			'PASSPORT', 'NATIONAL_ID', 'DRIVERS_LICENSE', 'BIRTH_CERTIFICATE', 'RESIDENCE_PERMIT'
		)`,
		`CREATE TABLE IF NOT EXISTS enrollment_requests (
			request_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			citizen_id         UUID,
			full_name          VARCHAR(200) NOT NULL,
			date_of_birth      VARCHAR(10) NOT NULL,
			nationality        VARCHAR(100) NOT NULL,
			email              VARCHAR(200),
			phone              VARCHAR(50),
			proofing_level     identity_proofing_level NOT NULL DEFAULT 'IAL2',
			status             enrollment_status NOT NULL DEFAULT 'DRAFT',
			submitted_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			assigned_officer   VARCHAR(150),
			remarks            TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS enrollment_documents (
			doc_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			request_id        UUID NOT NULL REFERENCES enrollment_requests(request_id),
			doc_type          document_type NOT NULL,
			doc_number        VARCHAR(100),
			issuing_auth      VARCHAR(200),
			issue_date        VARCHAR(10),
			expiry_date       VARCHAR(10),
			front_image       TEXT,
			back_image        TEXT,
			is_verified       BOOLEAN DEFAULT FALSE,
			verified_at       TIMESTAMPTZ,
			uploaded_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS enrollment_biometrics (
			sample_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			request_id        UUID NOT NULL REFERENCES enrollment_requests(request_id),
			sample_type       VARCHAR(50) NOT NULL,
			format            VARCHAR(20) NOT NULL,
			data              TEXT,
			quality           DOUBLE PRECISION DEFAULT 0,
			captured_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			device_id         VARCHAR(100),
			operator_id       VARCHAR(100)
		)`,
		`CREATE TABLE IF NOT EXISTS enrollment_reviews (
			review_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			request_id        UUID NOT NULL REFERENCES enrollment_requests(request_id),
			officer_id        VARCHAR(100) NOT NULL,
			officer_name      VARCHAR(200) NOT NULL,
			decision          VARCHAR(10) NOT NULL,
			reason            TEXT,
			verified_level    identity_proofing_level NOT NULL,
			reviewed_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_enrollment_status ON enrollment_requests(status)`,
		`CREATE INDEX IF NOT EXISTS idx_enrollment_citizen ON enrollment_requests(citizen_id)`,
		`CREATE INDEX IF NOT EXISTS idx_enrollment_docs_request ON enrollment_documents(request_id)`,
		`CREATE INDEX IF NOT EXISTS idx_enrollment_bio_request ON enrollment_biometrics(request_id)`,
	}

	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			return fmt.Errorf("migration failed: %s: %w", m[:60], err)
		}
	}
	return nil
}

func getEnv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
