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

	"github.com/snisid/executive-protection-ht/internal/handler"
	"github.com/snisid/executive-protection-ht/internal/kafka"
	"github.com/snisid/executive-protection-ht/internal/repository"
	"github.com/snisid/executive-protection-ht/internal/service"
)

func main() {
	dbHost := getEnv("EXECPROT_DB_HOST", "localhost")
	dbPort := getEnv("EXECPROT_DB_PORT", "26257")
	dbName := getEnv("EXECPROT_DB_NAME", "snisid_execprot")
	dbUser := getEnv("EXECPROT_DB_USER", "root")
	dbSSLMode := getEnv("EXECPROT_DB_SSLMODE", "disable")
	kafkaBrokers := getEnv("EXECPROT_KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("EXECPROT_KAFKA_TOPIC", "snisid.execprot.events")
	port := getEnv("EXECPROT_SERVICE_PORT", "8306")

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
	svc := service.NewExecutiveProtectionService(repo, producer)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := handler.NewHandler(svc)
	api := r.Group("/api/v1/execprot")
	h.RegisterRoutes(api)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("executive-protection-ht service started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run executive-protection-ht: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down executive-protection-ht...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TYPE IF NOT EXISTS execprot_protection_level AS ENUM ('PRESIDENT', 'PRIME_MINISTER', 'CABINET_MINISTER', 'JUDGE', 'DIPLOMAT', 'WITNESS')`,
		`CREATE TYPE IF NOT EXISTS execprot_risk_assessment AS ENUM ('LOW', 'MEDIUM', 'HIGH', 'CRITICAL')`,
		`CREATE TYPE IF NOT EXISTS execprot_transport_mode AS ENUM ('MOTORCADE', 'HELICOPTER', 'COMMERCIAL_FLIGHT')`,
		`CREATE TYPE IF NOT EXISTS execprot_movement_status AS ENUM ('DRAFT', 'APPROVED', 'ACTIVE', 'COMPLETED', 'CANCELLED')`,
		`CREATE TYPE IF NOT EXISTS execprot_threat_type AS ENUM ('DIRECT_THREAT', 'SOCIAL_MEDIA', 'KNOWN_GROUP', 'STALKER')`,
		`CREATE TYPE IF NOT EXISTS execprot_threat_status AS ENUM ('PENDING', 'ACTIVE', 'MITIGATED', 'FALSE_ALARM')`,
		`CREATE TABLE IF NOT EXISTS execprot_protectees (
			id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			full_name          VARCHAR(200) NOT NULL,
			official_title     VARCHAR(200) NOT NULL,
			protection_level   execprot_protection_level NOT NULL,
			risk_assessment    execprot_risk_assessment NOT NULL DEFAULT 'LOW',
			active_threats     INT NOT NULL DEFAULT 0,
			primary_agent_id   UUID NOT NULL,
			secondary_agents   UUID[],
			secure_vehicle_plate VARCHAR(20) NOT NULL,
			residence_location TEXT,
			workplace_location TEXT,
			daily_schedule_refs TEXT[],
			created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS execprot_movement_plans (
			id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			protectee_id      UUID NOT NULL REFERENCES execprot_protectees(id),
			event_name        VARCHAR(200) NOT NULL,
			date              TIMESTAMPTZ NOT NULL,
			departure_location VARCHAR(200) NOT NULL,
			arrival_location  VARCHAR(200) NOT NULL,
			transport_mode    execprot_transport_mode NOT NULL,
			route_plan        TEXT,
			advance_done      BOOLEAN NOT NULL DEFAULT FALSE,
			cleared_by        UUID,
			status            execprot_movement_status NOT NULL DEFAULT 'DRAFT',
			created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS execprot_threat_assessments (
			id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			protectee_id UUID NOT NULL REFERENCES execprot_protectees(id),
			threat_type  execprot_threat_type NOT NULL,
			threat_level execprot_risk_assessment NOT NULL,
			threat_detail TEXT,
			source_info  TEXT,
			assessed_by  UUID NOT NULL,
			mitigation   TEXT,
			status       execprot_threat_status NOT NULL DEFAULT 'PENDING',
			created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_execprot_protectees_risk ON execprot_protectees(risk_assessment)`,
		`CREATE INDEX IF NOT EXISTS idx_execprot_movements_date ON execprot_movement_plans(date)`,
		`CREATE INDEX IF NOT EXISTS idx_execprot_movements_status ON execprot_movement_plans(status)`,
		`CREATE INDEX IF NOT EXISTS idx_execprot_threats_protectee ON execprot_threat_assessments(protectee_id, status)`,
		`CREATE INDEX IF NOT EXISTS idx_execprot_threats_status ON execprot_threat_assessments(status)`,
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
