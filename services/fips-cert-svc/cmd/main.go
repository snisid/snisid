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

	"github.com/snisid/fips-cert-svc/internal/handler"
	"github.com/snisid/fips-cert-svc/internal/kafka"
	"github.com/snisid/fips-cert-svc/internal/repository"
	"github.com/snisid/fips-cert-svc/internal/service"
)

func main() {
	dbHost := getEnv("FIPS_DB_HOST", "localhost")
	dbPort := getEnv("FIPS_DB_PORT", "26257")
	dbName := getEnv("FIPS_DB_NAME", "snisid_fips")
	dbUser := getEnv("FIPS_DB_USER", "root")
	dbSSLMode := getEnv("FIPS_DB_SSLMODE", "disable")
	kafkaBrokers := getEnv("FIPS_KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("FIPS_KAFKA_TOPIC", "snisid.fips.events")
	port := getEnv("FIPS_SERVICE_PORT", "8098")

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
	svc := service.NewFIPSService(repo, producer)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := handler.NewHandler(svc)
	api := r.Group("/api/v1/fips")
	h.RegisterRoutes(api)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("fips-cert-svc started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run fips-cert-svc: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down fips-cert-svc...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TYPE IF NOT EXISTS fips_level AS ENUM (
			'LEVEL_1', 'LEVEL_2', 'LEVEL_3', 'LEVEL_4'
		)`,
		`CREATE TYPE IF NOT EXISTS validation_status AS ENUM (
			'PENDING', 'IN_REVIEW', 'VALIDATED', 'REJECTED', 'EXPIRED'
		)`,
		`CREATE TABLE IF NOT EXISTS fips_crypto_modules (
			module_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name              VARCHAR(200) NOT NULL,
			version           VARCHAR(50) NOT NULL,
			vendor            VARCHAR(200) NOT NULL,
			fips_level        fips_level NOT NULL,
			algorithms        TEXT[],
			cert_number       VARCHAR(100),
			validation_date   TIMESTAMPTZ,
			expiry_date       TIMESTAMPTZ,
			status            validation_status NOT NULL DEFAULT 'PENDING',
			created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS fips_cve_results (
			scan_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			module_id   UUID NOT NULL REFERENCES fips_crypto_modules(module_id),
			cve_id      VARCHAR(50) NOT NULL,
			severity    VARCHAR(20) NOT NULL,
			discovered  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			patched     BOOLEAN,
			notes       TEXT
		)`,
		`CREATE INDEX IF NOT EXISTS idx_fips_modules_vendor ON fips_crypto_modules(vendor)`,
		`CREATE INDEX IF NOT EXISTS idx_fips_modules_status ON fips_crypto_modules(status)`,
		`CREATE INDEX IF NOT EXISTS idx_fips_cves_module ON fips_cve_results(module_id)`,
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
