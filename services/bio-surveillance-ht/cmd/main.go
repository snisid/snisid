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

	"github.com/snisid/bio-surveillance-ht/internal/handler"
	"github.com/snisid/bio-surveillance-ht/internal/kafka"
	"github.com/snisid/bio-surveillance-ht/internal/repository"
	"github.com/snisid/bio-surveillance-ht/internal/service"
)

func main() {
	dbHost := getEnv("BIOSURV_DB_HOST", "localhost")
	dbPort := getEnv("BIOSURV_DB_PORT", "26257")
	dbName := getEnv("BIOSURV_DB_NAME", "snisid_biosurv")
	dbUser := getEnv("BIOSURV_DB_USER", "root")
	dbSSLMode := getEnv("BIOSURV_DB_SSLMODE", "disable")
	kafkaBrokers := getEnv("BIOSURV_KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("BIOSURV_KAFKA_TOPIC", "snisid.biosurv.events")
	port := getEnv("BIOSURV_SERVICE_PORT", "8305")

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
	svc := service.NewBioSurveillanceService(repo, producer)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := handler.NewHandler(svc)
	api := r.Group("/api/v1/biosurv")
	h.RegisterRoutes(api)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("bio-surveillance-ht service started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run bio-surveillance-ht: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down bio-surveillance-ht...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TYPE IF NOT EXISTS biosurv_pathogen_type AS ENUM ('VIRUS', 'BACTERIA', 'PARASITE', 'FUNGUS', 'UNKNOWN')`,
		`CREATE TYPE IF NOT EXISTS biosurv_alert_level AS ENUM ('GREEN', 'YELLOW', 'ORANGE', 'RED')`,
		`CREATE TYPE IF NOT EXISTS biosurv_transmission_mode AS ENUM ('AIRBORNE', 'FOODBORNE', 'WATERBORNE', 'VECTOR', 'CONTACT', 'UNKNOWN')`,
		`CREATE TYPE IF NOT EXISTS biosurv_facility_type AS ENUM ('HOSPITAL', 'CLINIC', 'LAB', 'PHARMACY', 'EMERGENCY_POST')`,
		`CREATE TYPE IF NOT EXISTS biosurv_stock_status AS ENUM ('ADEQUATE', 'LOW', 'CRITICAL', 'OUT_OF_STOCK')`,
		`CREATE TABLE IF NOT EXISTS biosurv_disease_alerts (
			id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			disease_name      VARCHAR(200) NOT NULL,
			pathogen_type     biosurv_pathogen_type NOT NULL,
			icd10_code        VARCHAR(10) NOT NULL,
			alert_level       biosurv_alert_level NOT NULL DEFAULT 'GREEN',
			first_case_detected_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			symptoms_hallmark TEXT,
			transmission_mode biosurv_transmission_mode NOT NULL DEFAULT 'UNKNOWN',
			incubation_days   INT NOT NULL DEFAULT 0,
			fatality_rate     DECIMAL(5,2) NOT NULL DEFAULT 0.0,
			cases_confirmed   INT NOT NULL DEFAULT 0,
			cases_suspected   INT NOT NULL DEFAULT 0,
			cases_deaths      INT NOT NULL DEFAULT 0,
			affected_regions  TEXT[],
			source_lab        VARCHAR(200),
			who_alert_ref     VARCHAR(50),
			containment_measures TEXT,
			created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS biosurv_vaccination_campaigns (
			id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			campaign_name     VARCHAR(200) NOT NULL,
			target_disease    VARCHAR(200) NOT NULL,
			vaccine_type      VARCHAR(100) NOT NULL,
			target_population INT NOT NULL DEFAULT 0,
			doses_administered INT NOT NULL DEFAULT 0,
			coverage_pct      DECIMAL(5,2) NOT NULL DEFAULT 0.0,
			regions_active    TEXT[],
			start_date        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			end_date          TIMESTAMPTZ,
			coordinator_agency VARCHAR(200) NOT NULL DEFAULT '',
			created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS biosurv_facilities (
			id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			facility_name   VARCHAR(200) NOT NULL,
			facility_type   biosurv_facility_type NOT NULL,
			region          VARCHAR(100) NOT NULL,
			commune         VARCHAR(100) NOT NULL,
			dept_code       CHAR(2) NOT NULL,
			capacity_beds   INT NOT NULL DEFAULT 0,
			beds_available  INT NOT NULL DEFAULT 0,
			stock_status    biosurv_stock_status NOT NULL DEFAULT 'ADEQUATE',
			has_ventilators BOOLEAN NOT NULL DEFAULT FALSE,
			has_ambulance   BOOLEAN NOT NULL DEFAULT FALSE,
			last_report_at  TIMESTAMPTZ,
			created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_biosurv_alerts_level ON biosurv_disease_alerts(alert_level)`,
		`CREATE INDEX IF NOT EXISTS idx_biosurv_alerts_region ON biosurv_disease_alerts USING GIN(affected_regions)`,
		`CREATE INDEX IF NOT EXISTS idx_biosurv_campaigns_region ON biosurv_vaccination_campaigns USING GIN(regions_active)`,
		`CREATE INDEX IF NOT EXISTS idx_biosurv_facilities_region ON biosurv_facilities(region)`,
		`CREATE INDEX IF NOT EXISTS idx_biosurv_facilities_stock ON biosurv_facilities(stock_status)`,
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
