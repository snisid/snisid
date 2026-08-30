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

	"github.com/snisid/bug-bounty-svc/internal/handler"
	"github.com/snisid/bug-bounty-svc/internal/kafka"
	"github.com/snisid/bug-bounty-svc/internal/repository"
	"github.com/snisid/bug-bounty-svc/internal/service"
)

func main() {
	dbHost := getEnv("BB_DB_HOST", "localhost")
	dbPort := getEnv("BB_DB_PORT", "26257")
	dbName := getEnv("BB_DB_NAME", "snisid_bug_bounty")
	dbUser := getEnv("BB_DB_USER", "root")
	dbSSLMode := getEnv("BB_DB_SSLMODE", "disable")
	kafkaBrokers := getEnv("BB_KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("BB_KAFKA_TOPIC", "snisid.bug-bounty.events")
	port := getEnv("BB_SERVICE_PORT", "8100")

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
	svc := service.NewBugBountyService(repo, producer)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := handler.NewHandler(svc)
	api := r.Group("/api/v1/bug-bounty")
	h.RegisterRoutes(api)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("bug-bounty-svc started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run bug-bounty-svc: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down bug-bounty-svc...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TYPE IF NOT EXISTS severity_level AS ENUM (
			'CRITICAL', 'HIGH', 'MEDIUM', 'LOW', 'INFO'
		)`,
		`CREATE TABLE IF NOT EXISTS bb_programs (
			scope_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			program_id   UUID NOT NULL,
			target       VARCHAR(300) NOT NULL,
			scope_type   VARCHAR(50) NOT NULL,
			in_scope     BOOLEAN DEFAULT TRUE,
			reward_min   DECIMAL(12,2),
			reward_max   DECIMAL(12,2)
		)`,
		`CREATE TABLE IF NOT EXISTS bb_reports (
			report_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			program_id   UUID NOT NULL,
			submitter    VARCHAR(200) NOT NULL,
			title        VARCHAR(300) NOT NULL,
			description  TEXT NOT NULL,
			severity     severity_level NOT NULL,
			scope_id     UUID,
			status       VARCHAR(30) NOT NULL DEFAULT 'SUBMITTED',
			submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS bb_triage_results (
			triage_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			report_id    UUID NOT NULL REFERENCES bb_reports(report_id),
			triager      VARCHAR(200) NOT NULL,
			severity     severity_level NOT NULL,
			reproducible BOOLEAN DEFAULT FALSE,
			duplicate_of UUID,
			notes        TEXT,
			triaged_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS bb_rewards (
			reward_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			report_id    UUID NOT NULL REFERENCES bb_reports(report_id),
			amount       DECIMAL(12,2) NOT NULL,
			currency     VARCHAR(10) NOT NULL DEFAULT 'USD',
			paid_to      VARCHAR(200) NOT NULL,
			approved_by  VARCHAR(200) NOT NULL,
			paid_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS bb_pentest_engagements (
			engagement_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			program_id    UUID NOT NULL,
			title         VARCHAR(300) NOT NULL,
			scope         TEXT NOT NULL,
			start_date    DATE NOT NULL,
			end_date      DATE,
			team_lead     VARCHAR(200) NOT NULL,
			status        VARCHAR(30) NOT NULL DEFAULT 'SCHEDULED',
			created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS bb_retest_schedules (
			schedule_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			report_id      UUID NOT NULL REFERENCES bb_reports(report_id),
			scheduled_for  TIMESTAMPTZ NOT NULL,
			completed_at   TIMESTAMPTZ,
			assigned_to    VARCHAR(200) NOT NULL,
			status         VARCHAR(30) NOT NULL DEFAULT 'PENDING'
		)`,
		`CREATE INDEX IF NOT EXISTS idx_bb_reports_program ON bb_reports(program_id)`,
		`CREATE INDEX IF NOT EXISTS idx_bb_reports_status ON bb_reports(status)`,
		`CREATE INDEX IF NOT EXISTS idx_bb_pentests_program ON bb_pentest_engagements(program_id)`,
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
