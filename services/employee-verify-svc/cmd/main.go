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

	"github.com/snisid/employee-verify-svc/internal/handler"
	"github.com/snisid/employee-verify-svc/internal/kafka"
	"github.com/snisid/employee-verify-svc/internal/repository"
	"github.com/snisid/employee-verify-svc/internal/service"
)

func main() {
	dbHost := getEnv("EV_DB_HOST", "localhost")
	dbPort := getEnv("EV_DB_PORT", "26257")
	dbName := getEnv("EV_DB_NAME", "snisid_employee_verify")
	dbUser := getEnv("EV_DB_USER", "root")
	dbSSLMode := getEnv("EV_DB_SSLMODE", "disable")
	kafkaBrokers := getEnv("EV_KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("EV_KAFKA_TOPIC", "snisid.employee-verify.events")
	port := getEnv("EV_SERVICE_PORT", "8096")

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
	svc := service.NewEmployeeVerifyService(repo, producer)

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
		log.Printf("employee-verify-svc service started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run employee-verify-svc: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down employee-verify-svc...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS employers (
			employer_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			company_name   VARCHAR(200) NOT NULL,
			ein            VARCHAR(20) UNIQUE NOT NULL,
			address        TEXT NOT NULL,
			contact_email  VARCHAR(200) NOT NULL,
			contact_phone  VARCHAR(30) NOT NULL,
			registered_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			is_active      BOOLEAN DEFAULT TRUE
		)`,
		`CREATE TABLE IF NOT EXISTS verification_cases (
			tcn             VARCHAR(50) PRIMARY KEY,
			employer_id     UUID NOT NULL REFERENCES employers(employer_id),
			employee_name   VARCHAR(200) NOT NULL,
			document_number VARCHAR(100) NOT NULL,
			document_type   VARCHAR(50) NOT NULL,
			status          VARCHAR(30) NOT NULL DEFAULT 'PENDING',
			created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS verification_results (
			result_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tcn          VARCHAR(50) NOT NULL REFERENCES verification_cases(tcn),
			ssa_match    BOOLEAN NOT NULL,
			dhs_match    BOOLEAN NOT NULL,
			is_eligible  BOOLEAN NOT NULL,
			reason       TEXT,
			completed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			status       VARCHAR(30) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS case_history (
			history_id  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tcn         VARCHAR(50) NOT NULL REFERENCES verification_cases(tcn),
			action      VARCHAR(100) NOT NULL,
			actioned_by VARCHAR(200) NOT NULL,
			actioned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			details     TEXT
		)`,
		`CREATE INDEX IF NOT EXISTS idx_verification_cases_employer ON verification_cases(employer_id)`,
		`CREATE INDEX IF NOT EXISTS idx_verification_cases_status ON verification_cases(status)`,
		`CREATE INDEX IF NOT EXISTS idx_verification_results_tcn ON verification_results(tcn)`,
		`CREATE INDEX IF NOT EXISTS idx_case_history_tcn ON case_history(tcn)`,
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
