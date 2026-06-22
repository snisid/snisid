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

	"github.com/snisid/pki-bridge-svc/internal/handler"
	"github.com/snisid/pki-bridge-svc/internal/kafka"
	"github.com/snisid/pki-bridge-svc/internal/repository"
	"github.com/snisid/pki-bridge-svc/internal/service"
)

func main() {
	dbHost := getEnv("PKI_BRIDGE_DB_HOST", "localhost")
	dbPort := getEnv("PKI_BRIDGE_DB_PORT", "26257")
	dbName := getEnv("PKI_BRIDGE_DB_NAME", "snisid_pki_bridge")
	dbUser := getEnv("PKI_BRIDGE_DB_USER", "root")
	dbSSLMode := getEnv("PKI_BRIDGE_DB_SSLMODE", "disable")
	kafkaBrokers := getEnv("PKI_BRIDGE_KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("PKI_BRIDGE_KAFKA_TOPIC", "snisid.pki-bridge.events")
	port := getEnv("PKI_BRIDGE_SERVICE_PORT", "8099")

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
	svc := service.NewPKIBridgeService(repo, producer)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := handler.NewHandler(svc)
	api := r.Group("/api/v1/pki-bridge")
	h.RegisterRoutes(api)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("pki-bridge-svc started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run pki-bridge-svc: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down pki-bridge-svc...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS pki_foreign_cas (
			ca_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name           VARCHAR(200) NOT NULL,
			country        VARCHAR(100) NOT NULL,
			public_key_pem TEXT NOT NULL,
			cert_policy    TEXT,
			registered_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			status         VARCHAR(20) NOT NULL DEFAULT 'ACTIVE'
		)`,
		`CREATE TABLE IF NOT EXISTS pki_cross_certs (
			cross_cert_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			subject         VARCHAR(300) NOT NULL,
			issuer_ca_id    UUID NOT NULL REFERENCES pki_foreign_cas(ca_id),
			serial_number   VARCHAR(100) NOT NULL,
			not_before      TIMESTAMPTZ NOT NULL,
			not_after       TIMESTAMPTZ NOT NULL,
			certificate_pem TEXT NOT NULL,
			created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS pki_trust_anchors (
			anchor_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			subject        VARCHAR(300) NOT NULL,
			public_key_pem TEXT NOT NULL,
			added_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			expires_at     TIMESTAMPTZ
		)`,
		`CREATE TABLE IF NOT EXISTS pki_path_validations (
			validation_id  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			path_id        UUID NOT NULL,
			result         BOOLEAN NOT NULL,
			errors         TEXT[],
			validated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS pki_bridge_agreements (
			agreement_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name           VARCHAR(200) NOT NULL,
			partner_ca     VARCHAR(200) NOT NULL,
			policy_id      UUID NOT NULL,
			signed_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			expires_at     TIMESTAMPTZ,
			status         VARCHAR(20) NOT NULL DEFAULT 'ACTIVE'
		)`,
		`CREATE INDEX IF NOT EXISTS idx_pki_cross_certs_subject ON pki_cross_certs(subject)`,
		`CREATE INDEX IF NOT EXISTS idx_pki_agreements_status ON pki_bridge_agreements(status)`,
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
