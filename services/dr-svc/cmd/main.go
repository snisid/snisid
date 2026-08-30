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

	"github.com/snisid/dr-svc/internal/handler"
	"github.com/snisid/dr-svc/internal/kafka"
	"github.com/snisid/dr-svc/internal/repository"
	"github.com/snisid/dr-svc/internal/service"
)

func main() {
	dbHost := getEnv("DR_DB_HOST", "localhost")
	dbPort := getEnv("DR_DB_PORT", "26257")
	dbName := getEnv("DR_DB_NAME", "snisid_dr")
	dbUser := getEnv("DR_DB_USER", "root")
	dbSSLMode := getEnv("DR_DB_SSLMODE", "disable")
	kafkaBrokers := getEnv("DR_KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("DR_KAFKA_TOPIC", "snisid.dr.events")
	port := getEnv("DR_SERVICE_PORT", "8097")

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
	svc := service.NewDRService(repo, producer)

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
		log.Printf("dr-svc service started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run dr-svc: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down dr-svc...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS dr_regions (
			region_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name         VARCHAR(100) UNIQUE NOT NULL,
			endpoint     VARCHAR(300) NOT NULL,
			is_active    BOOLEAN DEFAULT TRUE,
			health       VARCHAR(20) NOT NULL DEFAULT 'HEALTHY',
			last_checked TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS replication_status (
			replication_id  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			source_region   VARCHAR(100) NOT NULL,
			target_region   VARCHAR(100) NOT NULL,
			lag_seconds     INTEGER NOT NULL DEFAULT 0,
			is_healthy      BOOLEAN DEFAULT TRUE,
			last_checked_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS failover_plans (
			plan_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name          VARCHAR(200) NOT NULL,
			source_region VARCHAR(100) NOT NULL,
			target_region VARCHAR(100) NOT NULL,
			is_automated  BOOLEAN DEFAULT FALSE,
			created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			is_executed   BOOLEAN DEFAULT FALSE
		)`,
		`CREATE TABLE IF NOT EXISTS failover_executions (
			execution_id  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			plan_id       UUID NOT NULL REFERENCES failover_plans(plan_id),
			started_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			completed_at  TIMESTAMPTZ,
			is_successful BOOLEAN DEFAULT FALSE,
			error_message TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS backup_manifests (
			manifest_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			region         VARCHAR(100) NOT NULL,
			backup_path    TEXT NOT NULL,
			backup_size_mb BIGINT DEFAULT 0,
			started_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			completed_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			is_valid       BOOLEAN DEFAULT TRUE
		)`,
		`CREATE TABLE IF NOT EXISTS recovery_points (
			point_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			manifest_id   UUID NOT NULL REFERENCES backup_manifests(manifest_id),
			recovery_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			is_restored   BOOLEAN DEFAULT FALSE,
			restored_at   TIMESTAMPTZ
		)`,
		`CREATE TABLE IF NOT EXISTS dr_test_results (
			test_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			plan_id       UUID NOT NULL REFERENCES failover_plans(plan_id),
			test_name     VARCHAR(200) NOT NULL,
			started_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			completed_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			is_successful BOOLEAN DEFAULT FALSE,
			details       TEXT
		)`,
		`CREATE INDEX IF NOT EXISTS idx_dr_regions_name ON dr_regions(name)`,
		`CREATE INDEX IF NOT EXISTS idx_failover_executions_plan ON failover_executions(plan_id)`,
		`CREATE INDEX IF NOT EXISTS idx_backup_manifests_region ON backup_manifests(region)`,
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
