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

	"github.com/snisid/audit-svc/internal/handler"
	"github.com/snisid/audit-svc/internal/kafka"
	"github.com/snisid/audit-svc/internal/repository"
	"github.com/snisid/audit-svc/internal/service"
)

func main() {
	dbHost := getEnv("AUDIT_DB_HOST", "localhost")
	dbPort := getEnv("AUDIT_DB_PORT", "26257")
	dbName := getEnv("AUDIT_DB_NAME", "snisid_audit")
	dbUser := getEnv("AUDIT_DB_USER", "root")
	dbSSLMode := getEnv("AUDIT_DB_SSLMODE", "disable")
	kafkaBrokers := getEnv("AUDIT_KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("AUDIT_KAFKA_TOPIC", "snisid.audit.events")
	port := getEnv("AUDIT_SERVICE_PORT", "8103")

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
	svc := service.NewAuditService(repo, producer)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := handler.NewHandler(svc)
	api := r.Group("/api/v1/audit")
	h.RegisterRoutes(api)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("audit-svc started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run audit-svc: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down audit-svc...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TYPE IF NOT EXISTS event_source AS ENUM (
			'SERVICE_DESK', 'CIVIL_HT', 'BIO_ADN', 'GOVERNANCE', 'SLA', 'IDENTITY', 'SYSTEM'
		)`,
		`CREATE TYPE IF NOT EXISTS audit_event_type AS ENUM (
			'CREATED', 'UPDATED', 'DELETED', 'ACCESSED', 'VERIFIED', 'BREACH', 'ESCALATE'
		)`,
		`CREATE TYPE IF NOT EXISTS audit_category AS ENUM (
			'IDENTITY', 'SECURITY', 'COMPLIANCE', 'OPERATIONAL', 'GOVERNANCE'
		)`,
		`CREATE TABLE IF NOT EXISTS audit_events (
			event_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			source      event_source NOT NULL,
			event_type  audit_event_type NOT NULL,
			category    audit_category NOT NULL,
			actor_id    UUID,
			resource_id VARCHAR(300) NOT NULL,
			action      VARCHAR(100) NOT NULL,
			payload     JSONB,
			hash        VARCHAR(64) NOT NULL,
			prev_hash   VARCHAR(64) NOT NULL DEFAULT '',
			timestamp   TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS immutable_audit_chain (
			entry_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			event_id   UUID NOT NULL REFERENCES audit_events(event_id),
			hash       VARCHAR(64) NOT NULL,
			prev_hash  VARCHAR(64) NOT NULL DEFAULT '',
			data       JSONB,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_events_source ON audit_events(source)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_events_category ON audit_events(category)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_events_timestamp ON audit_events(timestamp)`,
		`CREATE INDEX IF NOT EXISTS idx_immutable_chain_hash ON immutable_audit_chain(hash)`,
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
