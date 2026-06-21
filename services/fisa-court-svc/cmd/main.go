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

	"github.com/snisid/fisa-court-svc/internal/handler"
	"github.com/snisid/fisa-court-svc/internal/kafka"
	"github.com/snisid/fisa-court-svc/internal/repository"
	"github.com/snisid/fisa-court-svc/internal/service"
)

func main() {
	dbHost := getEnv("FISA_DB_HOST", "localhost")
	dbPort := getEnv("FISA_DB_PORT", "26257")
	dbName := getEnv("FISA_DB_NAME", "snisid_fisacourt")
	dbUser := getEnv("FISA_DB_USER", "root")
	dbSSLMode := getEnv("FISA_DB_SSLMODE", "disable")
	kafkaBrokers := getEnv("FISA_KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("FISA_KAFKA_TOPIC", "snisid.fisa.events")
	port := getEnv("FISA_SERVICE_PORT", "8312")

	dbURL := fmt.Sprintf("postgresql://%s@%s:%s/%s?sslmode=%s", dbUser, dbHost, dbPort, dbName, dbSSLMode)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to CockroachDB: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping CockroachDB: %v", err)
	}
	db.SetMaxOpenConns(25)

	if err := runMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	producer := kafka.NewProducer([]string{kafkaBrokers}, kafkaTopic)
	defer producer.Close()

	repo := repository.NewPostgresRepo(db)
	svc := service.NewFISAService(repo, producer)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := handler.NewHandler(svc)
	api := r.Group("/api/v1/fisa")
	h.RegisterRoutes(api)

	srv := &http.Server{Addr: ":" + port, Handler: r}
	go func() {
		log.Printf("fisa-court-svc service started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run fisa-court-svc: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down fisa-court-svc...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TYPE IF NOT EXISTS fisa_warrant_type AS ENUM ('TITLE_III_PHONE','TITLE_III_INTERNET','FISA_ELECTRONIC','FISA_PHYSICAL','PEN_REGISTER','TRAP_TRACE')`,
		`CREATE TYPE IF NOT EXISTS fisa_warrant_status AS ENUM ('DRAFT','PENDING','APPROVED','ACTIVE','EXPIRED','REVOKED')`,
		`CREATE TABLE IF NOT EXISTS fisa_warrants (
			id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			warrant_id          VARCHAR(50) UNIQUE NOT NULL,
			warrant_type        fisa_warrant_type NOT NULL,
			target_identity     VARCHAR(255) NOT NULL,
			target_details      TEXT,
			issuing_court       VARCHAR(255) NOT NULL,
			judge_name          VARCHAR(255) NOT NULL,
			applicant_agency    VARCHAR(255) NOT NULL,
			applicant_officer   UUID NOT NULL,
			probable_cause_summary TEXT,
			duration_days       INT NOT NULL,
			authorized_start    TIMESTAMPTZ,
			authorized_end      TIMESTAMPTZ,
			renewals            INT NOT NULL DEFAULT 0,
			status              fisa_warrant_status NOT NULL DEFAULT 'DRAFT',
			review_required_at  TIMESTAMPTZ,
			emergency_authorized BOOLEAN NOT NULL DEFAULT FALSE,
			emergency_approved_by UUID,
			created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_fisa_warrants_status ON fisa_warrants(status)`,
		`CREATE INDEX IF NOT EXISTS idx_fisa_warrants_warrant_id ON fisa_warrants(warrant_id)`,
		`CREATE TABLE IF NOT EXISTS fisa_reports (
			id                       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			warrant_id               UUID NOT NULL REFERENCES fisa_warrants(id),
			reporting_period_start   TIMESTAMPTZ NOT NULL,
			reporting_period_end     TIMESTAMPTZ NOT NULL,
			communications_intercepted INT NOT NULL DEFAULT 0,
			minimization_applied    BOOLEAN NOT NULL DEFAULT FALSE,
			incidental_collection   INT NOT NULL DEFAULT 0,
			us_person_identities    INT NOT NULL DEFAULT 0,
			results_summary         TEXT,
			submitted_by            UUID NOT NULL,
			submitted_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_fisa_reports_warrant ON fisa_reports(warrant_id)`,
		`CREATE TABLE IF NOT EXISTS fisa_dockets (
			id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			docket_number       VARCHAR(50) UNIQUE NOT NULL,
			court_term          VARCHAR(20) UNIQUE NOT NULL,
			judge_presiding     VARCHAR(255) NOT NULL,
			applications_filed   INT NOT NULL DEFAULT 0,
			applications_approved INT NOT NULL DEFAULT 0,
			applications_modified INT NOT NULL DEFAULT 0,
			applications_denied  INT NOT NULL DEFAULT 0,
			total_targets       INT NOT NULL DEFAULT 0,
			foreign_targets     INT NOT NULL DEFAULT 0,
			us_person_targets   INT NOT NULL DEFAULT 0,
			sealed_until        TIMESTAMPTZ,
			created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_fisa_dockets_term ON fisa_dockets(court_term)`,
	}
	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			return fmt.Errorf("migration: %s: %w", m[:60], err)
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
