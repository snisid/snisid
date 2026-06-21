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

	"github.com/snisid/counterintel-ht/internal/handler"
	"github.com/snisid/counterintel-ht/internal/kafka"
	"github.com/snisid/counterintel-ht/internal/repository"
	"github.com/snisid/counterintel-ht/internal/service"
)

func main() {
	dbHost := getEnv("COUNTERINTEL_DB_HOST", "localhost")
	dbPort := getEnv("COUNTERINTEL_DB_PORT", "26257")
	dbName := getEnv("COUNTERINTEL_DB_NAME", "snisid_counterintel")
	dbUser := getEnv("COUNTERINTEL_DB_USER", "root")
	dbSSLMode := getEnv("COUNTERINTEL_DB_SSLMODE", "disable")
	kafkaBrokers := getEnv("COUNTERINTEL_KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("COUNTERINTEL_KAFKA_TOPIC", "snisid.counterintel.events")
	port := getEnv("COUNTERINTEL_SERVICE_PORT", "8310")

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
	svc := service.NewCounterintelService(repo, producer)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := handler.NewHandler(svc)
	api := r.Group("/api/v1/counterintel")
	h.RegisterRoutes(api)

	srv := &http.Server{Addr: ":" + port, Handler: r}
	go func() {
		log.Printf("counterintel-ht service started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run counterintel-ht: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down counterintel-ht...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TYPE IF NOT EXISTS counterintel_inv_type AS ENUM ('STANDARD','ENHANCED','TOP_SECRET','REINVESTIGATION')`,
		`CREATE TYPE IF NOT EXISTS counterintel_inv_status AS ENUM ('PENDING','IN_PROGRESS','FAVORABLE','UNFAVORABLE')`,
		`CREATE TYPE IF NOT EXISTS counterintel_clearance AS ENUM ('UNCLASSIFIED','CONFIDENTIAL','SECRET','TOP_SECRET')`,
		`CREATE TYPE IF NOT EXISTS counterintel_alert_type AS ENUM ('UNAUTHORIZED_ACCESS','DATA_EXFIL','PRIVILEGE_ESCALATION','BEHAVIORAL','COLLUSION')`,
		`CREATE TYPE IF NOT EXISTS counterintel_severity AS ENUM ('LOW','MEDIUM','HIGH','CRITICAL')`,
		`CREATE TYPE IF NOT EXISTS counterintel_threat_status AS ENUM ('OPEN','INVESTIGATING','CONFIRMED','FALSE_POSITIVE','MITIGATED')`,
		`CREATE TYPE IF NOT EXISTS counterintel_rel_type AS ENUM ('DIPLOMATIC','BUSINESS','ACADEMIC','FAMILY','PERSONAL')`,
		`CREATE TABLE IF NOT EXISTS counterintel_background_investigations (
			id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			subject_identity_ref VARCHAR(255) NOT NULL,
			investigation_type   counterintel_inv_type NOT NULL,
			status               counterintel_inv_status NOT NULL DEFAULT 'PENDING',
			criminal_record_check BOOLEAN NOT NULL DEFAULT FALSE,
			financial_check       BOOLEAN NOT NULL DEFAULT FALSE,
			foreign_contacts_check BOOLEAN NOT NULL DEFAULT FALSE,
			social_media_check    BOOLEAN NOT NULL DEFAULT FALSE,
			drug_test             BOOLEAN NOT NULL DEFAULT FALSE,
			psych_eval            BOOLEAN NOT NULL DEFAULT FALSE,
			adjudicator           UUID,
			adjudication_notes    TEXT,
			completed_at          TIMESTAMPTZ,
			clearance_level_grantED counterintel_clearance,
			expires_at            TIMESTAMPTZ,
			created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_counterintel_inv_status ON counterintel_background_investigations(status)`,
		`CREATE TABLE IF NOT EXISTS counterintel_insider_threats (
			id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			subject_id      VARCHAR(255) NOT NULL,
			alert_type      counterintel_alert_type NOT NULL,
			severity        counterintel_severity NOT NULL,
			description     TEXT NOT NULL,
			evidence_refs   TEXT[] DEFAULT '{}',
			detected_by     VARCHAR(100) NOT NULL,
			status          counterintel_threat_status NOT NULL DEFAULT 'OPEN',
			investigation_ref UUID,
			created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_counterintel_threat_status ON counterintel_insider_threats(status)`,
		`CREATE TABLE IF NOT EXISTS counterintel_foreign_contacts (
			id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			subject_id      VARCHAR(255) NOT NULL,
			contact_name    VARCHAR(255) NOT NULL,
			foreign_government VARCHAR(255) NOT NULL,
			relationship_type counterintel_rel_type NOT NULL,
			last_contact_at TIMESTAMPTZ,
			frequency       VARCHAR(100),
			approved_by     UUID,
			notes           TEXT,
			created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_counterintel_contacts_subject ON counterintel_foreign_contacts(subject_id)`,
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
