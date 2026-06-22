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

	"github.com/snisid/hsm-svc/internal/handler"
	"github.com/snisid/hsm-svc/internal/kafka"
	"github.com/snisid/hsm-svc/internal/repository"
	"github.com/snisid/hsm-svc/internal/service"
)

func main() {
	dbHost := getEnv("HSM_DB_HOST", "localhost")
	dbPort := getEnv("HSM_DB_PORT", "26257")
	dbName := getEnv("HSM_DB_NAME", "snisid_hsm")
	dbUser := getEnv("HSM_DB_USER", "root")
	dbSSLMode := getEnv("HSM_DB_SSLMODE", "disable")
	kafkaBrokers := getEnv("HSM_KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("HSM_KAFKA_TOPIC", "snisid.hsm.events")
	port := getEnv("HSM_SERVICE_PORT", "8090")

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
	svc := service.NewHSMService(repo, producer)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := handler.NewHandler(svc)
	api := r.Group("/api/v1/hsm")
	h.RegisterRoutes(api)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("hsm-svc started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run hsm-svc: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down hsm-svc...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TYPE IF NOT EXISTS key_algorithm AS ENUM (
			'RSA', 'EC', 'PQC', 'AES', 'HMAC'
		)`,
		`CREATE TYPE IF NOT EXISTS key_state AS ENUM (
			'ACTIVE', 'DEACTIVATED', 'COMPROMISED', 'DESTROYED', 'PENDING_ROTATE'
		)`,
		`CREATE TABLE IF NOT EXISTS hsm_keys (
			key_id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			key_label           VARCHAR(200) NOT NULL,
			algorithm           key_algorithm NOT NULL,
			key_size            INTEGER NOT NULL,
			state               key_state NOT NULL DEFAULT 'ACTIVE',
			usages              TEXT[] NOT NULL DEFAULT '{}',
			slot_id             INTEGER NOT NULL DEFAULT 0,
			is_extractable      BOOLEAN DEFAULT FALSE,
			public_key_pem      TEXT,
			key_hash            VARCHAR(64) NOT NULL,
			rotated_at          TIMESTAMPTZ,
			expires_at          TIMESTAMPTZ,
			created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			created_by          VARCHAR(150) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS hsm_slots (
			slot_id             INTEGER PRIMARY KEY,
			label               VARCHAR(200) NOT NULL,
			manufacturer        VARCHAR(200),
			model               VARCHAR(200),
			serial_number       VARCHAR(200),
			firmware_version    VARCHAR(50),
			is_logged_in        BOOLEAN DEFAULT FALSE,
			token_present       BOOLEAN DEFAULT FALSE,
			hardware_model      VARCHAR(200)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_hsm_keys_algorithm ON hsm_keys(algorithm)`,
		`CREATE INDEX IF NOT EXISTS idx_hsm_keys_state ON hsm_keys(state)`,
		`CREATE INDEX IF NOT EXISTS idx_hsm_keys_label ON hsm_keys(key_label)`,
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
