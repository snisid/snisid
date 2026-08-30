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

	"github.com/snisid/critical-infra-protection-ht/internal/handler"
	"github.com/snisid/critical-infra-protection-ht/internal/kafka"
	"github.com/snisid/critical-infra-protection-ht/internal/repository"
	"github.com/snisid/critical-infra-protection-ht/internal/service"
)

func main() {
	dbHost := getEnv("INFRAPROT_DB_HOST", "localhost")
	dbPort := getEnv("INFRAPROT_DB_PORT", "26257")
	dbName := getEnv("INFRAPROT_DB_NAME", "snisid_infraprot")
	dbUser := getEnv("INFRAPROT_DB_USER", "root")
	dbSSLMode := getEnv("INFRAPROT_DB_SSLMODE", "disable")
	kafkaBrokers := getEnv("INFRAPROT_KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("INFRAPROT_KAFKA_TOPIC", "snisid.infraprot.events")
	port := getEnv("INFRAPROT_SERVICE_PORT", "8311")

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
	svc := service.NewInfraProtService(repo, producer)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := handler.NewHandler(svc)
	api := r.Group("/api/v1/infraprot")
	h.RegisterRoutes(api)

	srv := &http.Server{Addr: ":" + port, Handler: r}
	go func() {
		log.Printf("critical-infra-protection-ht service started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run critical-infra-protection-ht: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down critical-infra-protection-ht...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TYPE IF NOT EXISTS infraprot_sector AS ENUM ('ENERGY','TELECOM','WATER','BANKING','TRANSPORT','HEALTH','GOVERNMENT','FOOD')`,
		`CREATE TYPE IF NOT EXISTS infraprot_criticality AS ENUM ('CRITICAL','HIGH','MEDIUM','LOW')`,
		`CREATE TYPE IF NOT EXISTS infraprot_incident_type AS ENUM ('CYBER_ATTACK','PHYSICAL_BREACH','NATURAL_DISASTER','SABOTAGE','OUTAGE')`,
		`CREATE TYPE IF NOT EXISTS infraprot_incident_status AS ENUM ('REPORTED','RESPONDING','CONTAINED','RESOLVED')`,
		`CREATE TABLE IF NOT EXISTS infraprot_assets (
			id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			asset_name          VARCHAR(255) NOT NULL,
			sector              infraprot_sector NOT NULL,
			owner_entity        VARCHAR(255) NOT NULL,
			location_lat        DECIMAL(10,7) NOT NULL,
			location_lng        DECIMAL(10,7) NOT NULL,
			region              VARCHAR(255) NOT NULL,
			dept_code           CHAR(2) NOT NULL,
			criticality         infraprot_criticality NOT NULL,
			cyber_maturity_score   DECIMAL(3,2) DEFAULT 0,
			physical_security_score DECIMAL(3,2) DEFAULT 0,
			last_cisa_assessment_at TIMESTAMPTZ,
			contact_name        VARCHAR(255) NOT NULL,
			contact_phone       VARCHAR(50) NOT NULL,
			has_backup_generator BOOLEAN NOT NULL DEFAULT FALSE,
			has_cyber_insurance  BOOLEAN NOT NULL DEFAULT FALSE,
			created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_infraprot_assets_sector ON infraprot_assets(sector)`,
		`CREATE TABLE IF NOT EXISTS infraprot_incidents (
			id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			asset_id        UUID NOT NULL REFERENCES infraprot_assets(id),
			incident_type   infraprot_incident_type NOT NULL,
			severity        VARCHAR(50) NOT NULL,
			description     TEXT NOT NULL,
			impact_assessment TEXT,
			downtime_hours  DECIMAL(8,2),
			estimated_loss_usd DECIMAL(15,2),
			responded_by    UUID,
			status          infraprot_incident_status NOT NULL DEFAULT 'REPORTED',
			created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_infraprot_incidents_status ON infraprot_incidents(status)`,
		`CREATE INDEX IF NOT EXISTS idx_infraprot_incidents_asset ON infraprot_incidents(asset_id)`,
		`CREATE TABLE IF NOT EXISTS infraprot_sector_assessments (
			id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			sector           infraprot_sector NOT NULL,
			assessment_date  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			overall_risk_score INT NOT NULL CHECK (overall_risk_score >= 1 AND overall_risk_score <= 10),
			top_threats      TEXT[] DEFAULT '{}',
			vulnerabilities  TEXT[] DEFAULT '{}',
			recommendations  TEXT[] DEFAULT '{}',
			assessor_agency  VARCHAR(255) NOT NULL,
			next_assessment_due TIMESTAMPTZ,
			created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_infraprot_assessments_sector ON infraprot_sector_assessments(sector)`,
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
