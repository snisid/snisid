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

	"github.com/snisid/service-desk-svc/internal/handler"
	"github.com/snisid/service-desk-svc/internal/kafka"
	"github.com/snisid/service-desk-svc/internal/repository"
	"github.com/snisid/service-desk-svc/internal/service"
)

func main() {
	dbHost := getEnv("SVC_DB_HOST", "localhost")
	dbPort := getEnv("SVC_DB_PORT", "26257")
	dbName := getEnv("SVC_DB_NAME", "snisid_service_desk")
	dbUser := getEnv("SVC_DB_USER", "root")
	dbSSLMode := getEnv("SVC_DB_SSLMODE", "disable")
	kafkaBrokers := getEnv("SVC_KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("SVC_KAFKA_TOPIC", "snisid.service-desk.events")
	port := getEnv("SVC_SERVICE_PORT", "8102")

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
	svc := service.NewServiceDeskService(repo, producer)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := handler.NewHandler(svc)
	api := r.Group("/api/v1/service-desk")
	h.RegisterRoutes(api)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("service-desk-svc started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run service-desk-svc: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down service-desk-svc...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TYPE IF NOT EXISTS case_status AS ENUM (
			'OPEN', 'IN_PROGRESS', 'RESOLVED', 'CLOSED'
		)`,
		`CREATE TYPE IF NOT EXISTS recovery_method AS ENUM (
			'EMAIL', 'SMS', 'DOCUMENT', 'BIOMETRIC', 'IN_PERSON'
		)`,
		`CREATE TABLE IF NOT EXISTS support_cases (
			case_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			citizen_id UUID NOT NULL,
			subject    VARCHAR(300) NOT NULL,
			description TEXT,
			status     case_status NOT NULL DEFAULT 'OPEN',
			assigned_to VARCHAR(150),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			resolved_at TIMESTAMPTZ
		)`,
		`CREATE TABLE IF NOT EXISTS verification_challenges (
			challenge_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			case_id      UUID NOT NULL REFERENCES support_cases(case_id),
			method       recovery_method NOT NULL,
			challenge    VARCHAR(300) NOT NULL,
			expires_at   TIMESTAMPTZ NOT NULL,
			is_resolved  BOOLEAN DEFAULT FALSE,
			created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS identity_recovery_requests (
			request_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			case_id          UUID NOT NULL REFERENCES support_cases(case_id),
			citizen_id       UUID NOT NULL,
			preferred_method recovery_method NOT NULL,
			verified_methods recovery_method[],
			is_verified      BOOLEAN DEFAULT FALSE,
			created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			resolved_at      TIMESTAMPTZ
		)`,
		`CREATE TABLE IF NOT EXISTS case_notes (
			note_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			case_id    UUID NOT NULL REFERENCES support_cases(case_id),
			author     VARCHAR(150) NOT NULL,
			content    TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS case_resolutions (
			resolution_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			case_id       UUID NOT NULL REFERENCES support_cases(case_id),
			action        VARCHAR(100) NOT NULL,
			details       TEXT,
			resolved_by   VARCHAR(150),
			created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_support_cases_citizen ON support_cases(citizen_id)`,
		`CREATE INDEX IF NOT EXISTS idx_support_cases_status ON support_cases(status)`,
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
