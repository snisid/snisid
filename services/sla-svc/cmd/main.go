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

	"github.com/snisid/sla-svc/internal/handler"
	"github.com/snisid/sla-svc/internal/kafka"
	"github.com/snisid/sla-svc/internal/repository"
	"github.com/snisid/sla-svc/internal/service"
)

func main() {
	dbHost := getEnv("SLA_DB_HOST", "localhost")
	dbPort := getEnv("SLA_DB_PORT", "26257")
	dbName := getEnv("SLA_DB_NAME", "snisid_sla")
	dbUser := getEnv("SLA_DB_USER", "root")
	dbSSLMode := getEnv("SLA_DB_SSLMODE", "disable")
	kafkaBrokers := getEnv("SLA_KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("SLA_KAFKA_TOPIC", "snisid.sla.events")
	port := getEnv("SLA_SERVICE_PORT", "8104")

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
	svc := service.NewSLAService(repo, producer)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := handler.NewHandler(svc)
	api := r.Group("/api/v1/sla")
	h.RegisterRoutes(api)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("sla-svc started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run sla-svc: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down sla-svc...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS slas (
			sla_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name        VARCHAR(200) NOT NULL,
			description TEXT,
			owner       VARCHAR(150) NOT NULL,
			is_active   BOOLEAN DEFAULT TRUE,
			created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS slos (
			slo_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			sla_id          UUID NOT NULL REFERENCES slas(sla_id),
			name            VARCHAR(200) NOT NULL,
			target_value    DOUBLE PRECISION NOT NULL,
			threshold       DOUBLE PRECISION NOT NULL,
			time_window_days INTEGER NOT NULL DEFAULT 30,
			created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS sli_data (
			sli_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			slo_id      UUID NOT NULL REFERENCES slos(slo_id),
			sla_id      UUID NOT NULL REFERENCES slas(sla_id),
			name        VARCHAR(200) NOT NULL,
			value       DOUBLE PRECISION NOT NULL,
			recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS breach_records (
			breach_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			sla_id       UUID NOT NULL REFERENCES slas(sla_id),
			slo_id       UUID NOT NULL REFERENCES slos(slo_id),
			sli_value    DOUBLE PRECISION NOT NULL,
			threshold    DOUBLE PRECISION NOT NULL,
			detected_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			resolved_at  TIMESTAMPTZ,
			is_active    BOOLEAN DEFAULT TRUE
		)`,
		`CREATE TABLE IF NOT EXISTS uptime_windows (
			window_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			sla_id      UUID NOT NULL REFERENCES slas(sla_id),
			start_time  TIMESTAMPTZ NOT NULL,
			end_time    TIMESTAMPTZ,
			is_up       BOOLEAN DEFAULT TRUE,
			duration_ms BIGINT DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS escalation_policies (
			policy_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			sla_id         UUID NOT NULL REFERENCES slas(sla_id),
			escalate_after INTEGER NOT NULL DEFAULT 300,
			notify_channel VARCHAR(50) NOT NULL DEFAULT 'EMAIL',
			notify_target  VARCHAR(200) NOT NULL,
			is_active      BOOLEAN DEFAULT TRUE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_slos_sla ON slos(sla_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sli_data_slo ON sli_data(slo_id, recorded_at)`,
		`CREATE INDEX IF NOT EXISTS idx_breaches_sla ON breach_records(sla_id)`,
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
