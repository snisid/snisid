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

	"github.com/snisid/card-svc/internal/handler"
	"github.com/snisid/card-svc/internal/kafka"
	"github.com/snisid/card-svc/internal/repository"
	"github.com/snisid/card-svc/internal/service"
)

func main() {
	dbHost := getEnv("CARD_DB_HOST", "localhost")
	dbPort := getEnv("CARD_DB_PORT", "26257")
	dbName := getEnv("CARD_DB_NAME", "snisid_card")
	dbUser := getEnv("CARD_DB_USER", "root")
	dbSSLMode := getEnv("CARD_DB_SSLMODE", "disable")
	kafkaBrokers := getEnv("CARD_KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("CARD_KAFKA_TOPIC", "snisid.card.events")
	port := getEnv("CARD_SERVICE_PORT", "8092")

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
	svc := service.NewCardService(repo, producer)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := handler.NewHandler(svc)
	api := r.Group("/api/v1/card")
	h.RegisterRoutes(api)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("card-svc started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run card-svc: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down card-svc...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TYPE IF NOT EXISTS card_type AS ENUM (
			'NATIONAL_ID', 'RESIDENCE_PERMIT', 'DRIVERS_LICENSE', 'PASSPORT', 'EMPLOYEE_BADGE'
		)`,
		`CREATE TYPE IF NOT EXISTS card_status AS ENUM (
			'ORDERED', 'PERSONALIZED', 'ISSUED', 'ACTIVE', 'BLOCKED', 'EXPIRED', 'DESTROYED', 'LOST', 'STOLEN'
		)`,
		`CREATE TABLE IF NOT EXISTS card_profiles (
			profile_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			card_type     card_type NOT NULL,
			name          VARCHAR(150) NOT NULL,
			description   TEXT,
			form_factor   VARCHAR(50) NOT NULL,
			material      VARCHAR(50) NOT NULL,
			has_chip      BOOLEAN DEFAULT FALSE,
			has_mrz       BOOLEAN DEFAULT FALSE,
			valid_days    INTEGER NOT NULL DEFAULT 3650,
			created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS card_personalization (
			order_id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			profile_id      UUID NOT NULL REFERENCES card_profiles(profile_id),
			card_serial     VARCHAR(50) UNIQUE NOT NULL,
			citizen_id      UUID NOT NULL,
			full_name       VARCHAR(200) NOT NULL,
			date_of_birth   VARCHAR(10) NOT NULL,
			nationality     VARCHAR(100) NOT NULL,
			photo_data      TEXT,
			signature_data  TEXT,
			status          card_status NOT NULL DEFAULT 'ORDERED',
			ordered_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			personalized_at TIMESTAMPTZ,
			issued_at       TIMESTAMPTZ,
			activated_at    TIMESTAMPTZ,
			blocked_at      TIMESTAMPTZ,
			block_reason    TEXT,
			expires_at      TIMESTAMPTZ,
			created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS card_stock (
			stock_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			profile_id    UUID NOT NULL REFERENCES card_profiles(profile_id),
			serial_from   VARCHAR(50) NOT NULL,
			serial_to     VARCHAR(50) NOT NULL,
			quantity      INTEGER NOT NULL,
			available_qty INTEGER NOT NULL,
			location      VARCHAR(200),
			received_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS card_shipments (
			shipment_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			profile_id    UUID NOT NULL REFERENCES card_profiles(profile_id),
			serial_from   VARCHAR(50) NOT NULL,
			serial_to     VARCHAR(50) NOT NULL,
			quantity      INTEGER NOT NULL,
			tracking_ref  VARCHAR(200),
			vendor        VARCHAR(200) NOT NULL,
			received_by   VARCHAR(150),
			received_at   TIMESTAMPTZ,
			notes         TEXT,
			created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_card_personalization_serial ON card_personalization(card_serial)`,
		`CREATE INDEX IF NOT EXISTS idx_card_personalization_status ON card_personalization(status)`,
		`CREATE INDEX IF NOT EXISTS idx_card_personalization_citizen ON card_personalization(citizen_id)`,
		`CREATE INDEX IF NOT EXISTS idx_card_stock_profile ON card_stock(profile_id)`,
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
