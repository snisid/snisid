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

	"github.com/snisid/governance-svc/internal/handler"
	"github.com/snisid/governance-svc/internal/kafka"
	"github.com/snisid/governance-svc/internal/repository"
	"github.com/snisid/governance-svc/internal/service"
)

func main() {
	dbHost := getEnv("GOV_DB_HOST", "localhost")
	dbPort := getEnv("GOV_DB_PORT", "26257")
	dbName := getEnv("GOV_DB_NAME", "snisid_governance")
	dbUser := getEnv("GOV_DB_USER", "root")
	dbSSLMode := getEnv("GOV_DB_SSLMODE", "disable")
	kafkaBrokers := getEnv("GOV_KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("GOV_KAFKA_TOPIC", "snisid.governance.events")
	port := getEnv("GOV_SERVICE_PORT", "8105")

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
	svc := service.NewGovernanceService(repo, producer)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := handler.NewHandler(svc)
	api := r.Group("/api/v1/governance")
	h.RegisterRoutes(api)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("governance-svc started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run governance-svc: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down governance-svc...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TYPE IF NOT EXISTS license_type AS ENUM (
			'OSI_APPROVED', 'PROPRIETARY', 'CREATIVE_COMMONS', 'PUBLIC_DOMAIN', 'OTHER'
		)`,
		`CREATE TYPE IF NOT EXISTS compliance_status AS ENUM (
			'COMPLIANT', 'NON_COMPLIANT', 'PENDING_REVIEW', 'EXEMPTED'
		)`,
		`CREATE TABLE IF NOT EXISTS software_licenses (
			license_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name            VARCHAR(200) NOT NULL,
			spdx_id         VARCHAR(100) NOT NULL,
			license_type    license_type NOT NULL,
			version         VARCHAR(50) NOT NULL,
			publisher       VARCHAR(200) NOT NULL,
			is_osi_approved BOOLEAN DEFAULT FALSE,
			text            TEXT,
			registered_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS governance_policies (
			policy_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name        VARCHAR(200) NOT NULL,
			description TEXT,
			is_active   BOOLEAN DEFAULT TRUE,
			created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS policy_rules (
			rule_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			policy_id  UUID NOT NULL REFERENCES governance_policies(policy_id),
			rule_type  VARCHAR(50) NOT NULL,
			condition  TEXT NOT NULL,
			action     VARCHAR(100) NOT NULL,
			priority   INTEGER DEFAULT 0,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS license_audits (
			audit_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			license_id  UUID NOT NULL REFERENCES software_licenses(license_id),
			policy_id   UUID NOT NULL REFERENCES governance_policies(policy_id),
			status      compliance_status NOT NULL DEFAULT 'PENDING_REVIEW',
			findings    TEXT,
			audited_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			reviewed_by VARCHAR(150)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_licenses_type ON software_licenses(license_type)`,
		`CREATE INDEX IF NOT EXISTS idx_licenses_spdx ON software_licenses(spdx_id)`,
		`CREATE INDEX IF NOT EXISTS idx_audits_license ON license_audits(license_id)`,
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
