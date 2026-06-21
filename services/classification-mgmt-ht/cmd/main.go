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

	"github.com/snisid/classification-mgmt-ht/internal/handler"
	"github.com/snisid/classification-mgmt-ht/internal/kafka"
	"github.com/snisid/classification-mgmt-ht/internal/repository"
	"github.com/snisid/classification-mgmt-ht/internal/service"
)

func main() {
	dbHost := getEnv("CLASS_DB_HOST", "localhost")
	dbPort := getEnv("CLASS_DB_PORT", "26257")
	dbName := getEnv("CLASS_DB_NAME", "snisid_classification")
	dbUser := getEnv("CLASS_DB_USER", "root")
	dbSSLMode := getEnv("CLASS_DB_SSLMODE", "disable")
	kafkaBrokers := getEnv("CLASS_KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("CLASS_KAFKA_TOPIC", "snisid.classification.events")
	port := getEnv("CLASS_SERVICE_PORT", "8313")

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
	svc := service.NewClassificationService(repo, producer)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := handler.NewHandler(svc)
	api := r.Group("/api/v1/classification")
	h.RegisterRoutes(api)

	srv := &http.Server{Addr: ":" + port, Handler: r}
	go func() {
		log.Printf("classification-mgmt-ht service started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run classification-mgmt-ht: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down classification-mgmt-ht...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TYPE IF NOT EXISTS class_sensitivity AS ENUM ('PUBLIC','INTERNAL','CONFIDENTIAL','SECRET','TOP_SECRET')`,
		`CREATE TYPE IF NOT EXISTS class_audit_action AS ENUM ('CLASSIFY','DOWNGRADE','UPGRADE','DECLASSIFY','DESTROY')`,
		`CREATE TABLE IF NOT EXISTS classification_rules (
			id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			data_type           VARCHAR(255) NOT NULL,
			sensitivity_level   class_sensitivity NOT NULL,
			handling_caveats    TEXT[] DEFAULT '{}',
			dissemination_limit VARCHAR(500),
			encryption_required BOOLEAN NOT NULL DEFAULT FALSE,
			access_control_mfa  BOOLEAN NOT NULL DEFAULT FALSE,
			audit_logging       BOOLEAN NOT NULL DEFAULT FALSE,
			retention_days      INT NOT NULL DEFAULT 0,
			destruction_required BOOLEAN NOT NULL DEFAULT FALSE,
			created_by          UUID NOT NULL,
			active              BOOLEAN NOT NULL DEFAULT TRUE,
			created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_class_rules_data_type ON classification_rules(data_type)`,
		`CREATE TABLE IF NOT EXISTS classification_tags (
			id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			resource_uri         VARCHAR(500) NOT NULL,
			classification_top_level class_sensitivity NOT NULL,
			classification_atomic VARCHAR(50),
			handling_caveats     TEXT[] DEFAULT '{}',
			owner_agency         VARCHAR(255) NOT NULL,
			tagged_by            UUID NOT NULL,
			tagged_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			expires_at           TIMESTAMPTZ
		)`,
		`CREATE INDEX IF NOT EXISTS idx_class_tags_uri ON classification_tags(resource_uri)`,
		`CREATE TABLE IF NOT EXISTS classification_audit (
			id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			resource_uri          VARCHAR(500) NOT NULL,
			action                class_audit_action NOT NULL,
			from_level            VARCHAR(50),
			to_level              VARCHAR(50),
			rationale             TEXT,
			authorized_by         UUID NOT NULL,
			classification_authority VARCHAR(200) NOT NULL,
			timestamp             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			ip_address            VARCHAR(45) NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_class_audit_timestamp ON classification_audit(timestamp DESC)`,
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
