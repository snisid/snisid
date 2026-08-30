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

	"github.com/snisid/age-verification-svc/internal/handler"
	"github.com/snisid/age-verification-svc/internal/kafka"
	"github.com/snisid/age-verification-svc/internal/repository"
	"github.com/snisid/age-verification-svc/internal/service"
)

func main() {
	dbHost := getEnv("AGE_DB_HOST", "localhost")
	dbPort := getEnv("AGE_DB_PORT", "26257")
	dbName := getEnv("AGE_DB_NAME", "snisid_age_verification")
	dbUser := getEnv("AGE_DB_USER", "root")
	dbSSLMode := getEnv("AGE_DB_SSLMODE", "disable")
	kafkaBrokers := getEnv("AGE_KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("AGE_KAFKA_TOPIC", "snisid.age-verification.events")
	port := getEnv("AGE_SERVICE_PORT", "8095")

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
	svc := service.NewAgeVerificationService(repo, producer)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := handler.NewHandler(svc)
	api := r.Group("/api/v1")
	h.RegisterRoutes(api)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("age-verification-svc service started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run age-verification-svc: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down age-verification-svc...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS age_attestations (
			attestation_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			identity_id    UUID NOT NULL,
			date_of_birth  DATE NOT NULL,
			issued_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			expires_at     TIMESTAMPTZ NOT NULL,
			is_revoked     BOOLEAN DEFAULT FALSE,
			revoked_at     TIMESTAMPTZ
		)`,
		`CREATE TABLE IF NOT EXISTS age_claims (
			claim_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			attestation_id UUID NOT NULL REFERENCES age_attestations(attestation_id),
			verifier_id    VARCHAR(200) NOT NULL,
			bracket        VARCHAR(20) NOT NULL,
			is_satisfied   BOOLEAN NOT NULL,
			claimed_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS verifier_requests (
			request_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			verifier_id  VARCHAR(200) NOT NULL,
			bracket      VARCHAR(20) NOT NULL,
			requested_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			is_approved  BOOLEAN DEFAULT FALSE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_age_attestations_identity ON age_attestations(identity_id)`,
		`CREATE INDEX IF NOT EXISTS idx_age_claims_attestation ON age_claims(attestation_id)`,
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
