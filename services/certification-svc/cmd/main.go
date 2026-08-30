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

	"github.com/snisid/certification-svc/internal/handler"
	"github.com/snisid/certification-svc/internal/kafka"
	"github.com/snisid/certification-svc/internal/repository"
	"github.com/snisid/certification-svc/internal/service"
)

func main() {
	dbHost := getEnv("CERT_DB_HOST", "localhost")
	dbPort := getEnv("CERT_DB_PORT", "26257")
	dbName := getEnv("CERT_DB_NAME", "snisid_certification")
	dbUser := getEnv("CERT_DB_USER", "root")
	dbSSLMode := getEnv("CERT_DB_SSLMODE", "disable")
	kafkaBrokers := getEnv("CERT_KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("CERT_KAFKA_TOPIC", "snisid.certification.events")
	port := getEnv("CERT_SERVICE_PORT", "8093")

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
	svc := service.NewCertificationService(repo, producer)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := handler.NewHandler(svc)
	api := r.Group("/api/v1/certification")
	h.RegisterRoutes(api)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("certification-svc started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run certification-svc: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down certification-svc...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TYPE IF NOT EXISTS ial_level AS ENUM (
			'IAL_NONE', 'IAL1', 'IAL2', 'IAL3'
		)`,
		`CREATE TYPE IF NOT EXISTS aal_level AS ENUM (
			'AAL_NONE', 'AAL1', 'AAL2', 'AAL3'
		)`,
		`CREATE TYPE IF NOT EXISTS fal_level AS ENUM (
			'FAL_NONE', 'FAL1', 'FAL2', 'FAL3'
		)`,
		`CREATE TABLE IF NOT EXISTS certification_profiles (
			profile_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			identity_id   UUID NOT NULL UNIQUE,
			ial           ial_level NOT NULL DEFAULT 'IAL_NONE',
			aal           aal_level NOT NULL DEFAULT 'AAL_NONE',
			fal           fal_level NOT NULL DEFAULT 'FAL_NONE',
			is_active     BOOLEAN DEFAULT TRUE,
			valid_from    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			valid_until   TIMESTAMPTZ,
			last_assessed TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			assessor_id   VARCHAR(150) NOT NULL,
			assessor_org  VARCHAR(200),
			created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS certification_claims (
			claim_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			identity_id      UUID NOT NULL,
			framework_name   VARCHAR(100) NOT NULL,
			claim_type       VARCHAR(100) NOT NULL,
			claim_value      TEXT NOT NULL,
			issuer           VARCHAR(200) NOT NULL,
			issued_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			expires_at       TIMESTAMPTZ,
			is_verified      BOOLEAN DEFAULT FALSE,
			verification_ref VARCHAR(200),
			created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS certification_audit (
			audit_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			identity_id   UUID NOT NULL,
			action        VARCHAR(100) NOT NULL,
			field         VARCHAR(100),
			old_value     TEXT,
			new_value     TEXT,
			performed_by  VARCHAR(150) NOT NULL,
			performed_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			notes         TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS certification_compliance (
			check_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			identity_id   UUID NOT NULL,
			check_type    VARCHAR(100) NOT NULL,
			requirement   TEXT NOT NULL,
			is_compliant  BOOLEAN NOT NULL,
			details       TEXT,
			checked_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			checked_by    VARCHAR(150) NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_cert_profiles_identity ON certification_profiles(identity_id)`,
		`CREATE INDEX IF NOT EXISTS idx_cert_claims_identity ON certification_claims(identity_id)`,
		`CREATE INDEX IF NOT EXISTS idx_cert_audit_identity ON certification_audit(identity_id)`,
		`CREATE INDEX IF NOT EXISTS idx_cert_compliance_identity ON certification_compliance(identity_id)`,
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
