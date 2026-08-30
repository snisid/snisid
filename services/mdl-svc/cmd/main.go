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

	"github.com/snisid/mdl-svc/internal/handler"
	"github.com/snisid/mdl-svc/internal/kafka"
	"github.com/snisid/mdl-svc/internal/repository"
	"github.com/snisid/mdl-svc/internal/service"
)

func main() {
	dbHost := getEnv("MDL_DB_HOST", "localhost")
	dbPort := getEnv("MDL_DB_PORT", "26257")
	dbName := getEnv("MDL_DB_NAME", "snisid_mdl")
	dbUser := getEnv("MDL_DB_USER", "root")
	dbSSLMode := getEnv("MDL_DB_SSLMODE", "disable")
	kafkaBrokers := getEnv("MDL_KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("MDL_KAFKA_TOPIC", "snisid.mdl.events")
	port := getEnv("MDL_SERVICE_PORT", "8094")

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
	svc := service.NewMDLService(repo, producer)

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
		log.Printf("mdl-svc service started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run mdl-svc: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down mdl-svc...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS mdl_issuances (
			issuance_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			identity_id   UUID NOT NULL,
			device_id     VARCHAR(200) NOT NULL,
			issued_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			expires_at    TIMESTAMPTZ NOT NULL,
			is_revoked    BOOLEAN DEFAULT FALSE,
			revoked_at    TIMESTAMPTZ
		)`,
		`CREATE TABLE IF NOT EXISTS mdl_presentations (
			presentation_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			issuance_id          UUID NOT NULL REFERENCES mdl_issuances(issuance_id),
			reader_id            VARCHAR(200) NOT NULL,
			presented_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			is_verified          BOOLEAN DEFAULT FALSE,
			verification_result  TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS mdl_data_elements (
			element_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			issuance_id   UUID NOT NULL REFERENCES mdl_issuances(issuance_id),
			element_name  VARCHAR(100) NOT NULL,
			element_value TEXT NOT NULL,
			is_mandatory  BOOLEAN DEFAULT FALSE
		)`,
		`CREATE TABLE IF NOT EXISTS mdl_device_engagements (
			engagement_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			issuance_id     UUID NOT NULL REFERENCES mdl_issuances(issuance_id),
			qr_payload      TEXT NOT NULL,
			engagement_code VARCHAR(50) NOT NULL,
			created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			expires_at      TIMESTAMPTZ NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS mdl_qr_barcodes (
			barcode_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			engagement_id UUID NOT NULL REFERENCES mdl_device_engagements(engagement_id),
			encoded_data  TEXT NOT NULL,
			format        VARCHAR(20) NOT NULL,
			generated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS mdl_trust_registry (
			entry_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			reader_id     VARCHAR(200) UNIQUE NOT NULL,
			reader_name   VARCHAR(200) NOT NULL,
			public_key    TEXT NOT NULL,
			is_trusted    BOOLEAN DEFAULT TRUE,
			registered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			expires_at    TIMESTAMPTZ
		)`,
		`CREATE INDEX IF NOT EXISTS idx_mdl_issuances_identity ON mdl_issuances(identity_id)`,
		`CREATE INDEX IF NOT EXISTS idx_mdl_presentations_issuance ON mdl_presentations(issuance_id)`,
		`CREATE INDEX IF NOT EXISTS idx_mdl_trust_registry_reader ON mdl_trust_registry(reader_id)`,
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
