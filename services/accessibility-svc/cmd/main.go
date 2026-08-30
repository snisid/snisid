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

	"github.com/snisid/accessibility-svc/internal/handler"
	"github.com/snisid/accessibility-svc/internal/kafka"
	"github.com/snisid/accessibility-svc/internal/repository"
	"github.com/snisid/accessibility-svc/internal/service"
)

func main() {
	dbHost := getEnv("ACC_DB_HOST", "localhost")
	dbPort := getEnv("ACC_DB_PORT", "26257")
	dbName := getEnv("ACC_DB_NAME", "snisid_accessibility")
	dbUser := getEnv("ACC_DB_USER", "root")
	dbSSLMode := getEnv("ACC_DB_SSLMODE", "disable")
	kafkaBrokers := getEnv("ACC_KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("ACC_KAFKA_TOPIC", "snisid.accessibility.events")
	port := getEnv("ACC_SERVICE_PORT", "8101")

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
	svc := service.NewAccessibilityService(repo, producer)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := handler.NewHandler(svc)
	api := r.Group("/api/v1/accessibility")
	h.RegisterRoutes(api)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("accessibility-svc started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run accessibility-svc: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down accessibility-svc...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TYPE IF NOT EXISTS wcag_level AS ENUM ('A', 'AA', 'AAA')`,
		`CREATE TABLE IF NOT EXISTS acc_audit_runs (
			audit_run_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			target_url       VARCHAR(500) NOT NULL,
			wcag_level       wcag_level NOT NULL,
			status           VARCHAR(30) NOT NULL DEFAULT 'PENDING',
			total_violations INTEGER DEFAULT 0,
			passed           INTEGER DEFAULT 0,
			failed           INTEGER DEFAULT 0,
			started_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			completed_at     TIMESTAMPTZ,
			created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS acc_violations (
			violation_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			audit_run_id   UUID NOT NULL REFERENCES acc_audit_runs(audit_run_id),
			wcag_level     wcag_level NOT NULL,
			guideline      VARCHAR(100) NOT NULL,
			description    TEXT NOT NULL,
			element        JSONB NOT NULL,
			severity       VARCHAR(20) NOT NULL,
			remediated     BOOLEAN DEFAULT FALSE,
			remediated_at  TIMESTAMPTZ,
			created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS acc_remediation_tracks (
			track_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			violation_id   UUID NOT NULL REFERENCES acc_violations(violation_id),
			assigned_to    VARCHAR(200) NOT NULL,
			notes          TEXT,
			remediated_at  TIMESTAMPTZ,
			verified_by    VARCHAR(200),
			status         VARCHAR(30) NOT NULL DEFAULT 'OPEN'
		)`,
		`CREATE TABLE IF NOT EXISTS acc_audit_schedules (
			schedule_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			target_url     VARCHAR(500) NOT NULL,
			wcag_level     wcag_level NOT NULL,
			cron_expr      VARCHAR(100) NOT NULL,
			enabled        BOOLEAN DEFAULT TRUE,
			last_run_at    TIMESTAMPTZ,
			next_run_at    TIMESTAMPTZ,
			created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_acc_violations_audit ON acc_violations(audit_run_id)`,
		`CREATE INDEX IF NOT EXISTS idx_acc_violations_remediated ON acc_violations(remediated)`,
		`CREATE INDEX IF NOT EXISTS idx_acc_schedules_enabled ON acc_audit_schedules(enabled)`,
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
