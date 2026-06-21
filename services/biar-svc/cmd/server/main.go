package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/biar-svc/internal/handler"
	"github.com/snisid/platform/services/biar-svc/internal/kafka"
	"github.com/snisid/platform/services/biar-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/biar-svc/internal/service"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func runMigrations(pool *pgxpool.Pool, logger *zap.Logger) {
	ctx := context.Background()
	migrations := []string{
		`CREATE TYPE IF NOT EXISTS biar_recovery_context AS ENUM (
			'POLICE_OPERATION','CHECKPOINT','PORT_SEIZURE','AIRPORT_SEIZURE',
			'COMMUNITY_SURRENDER','CRIME_SCENE','RAID','BORDER_SEIZURE','OTHER'
		)`,
		`CREATE TYPE IF NOT EXISTS biar_weapon_disposition AS ENUM (
			'HELD_AS_EVIDENCE','DESTROYED','RETURNED_TO_OWNER',
			'TRANSFERRED_TO_POLICE','SENT_TO_INTERPOL','PENDING'
		)`,
		`CREATE TABLE IF NOT EXISTS biar_illicit_weapons (
			weapon_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			national_biar_id VARCHAR(25) UNIQUE NOT NULL,
			serial_number VARCHAR(100), serial_obliterated BOOLEAN DEFAULT FALSE,
			make VARCHAR(100), model VARCHAR(100), caliber VARCHAR(30),
			weapon_type VARCHAR(50) NOT NULL, manufacture_country CHAR(3),
			estimated_manufacture_year SMALLINT,
			recovery_date TIMESTAMPTZ NOT NULL,
			recovery_context biar_recovery_context NOT NULL,
			recovery_location VARCHAR(300), recovery_dept_code CHAR(2),
			recovery_commune VARCHAR(100),
			recovery_lat DECIMAL(10,7), recovery_lng DECIMAL(10,7),
			seizing_unit VARCHAR(50) NOT NULL, seizing_officer UUID,
			case_reference VARCHAR(100), from_person_id UUID,
			gang_id UUID, crime_category VARCHAR(50),
			associated_cases TEXT[] DEFAULT '{}',
			origin_country CHAR(3), transit_countries CHAR(3)[] DEFAULT '{}',
			trafficking_route TEXT, import_method TEXT,
			iarms_ref VARCHAR(50), atf_etrace_ref VARCHAR(50),
			reported_to_interpol BOOLEAN DEFAULT FALSE,
			interpol_reported_at TIMESTAMPTZ,
			disposition biar_weapon_disposition DEFAULT 'HELD_AS_EVIDENCE',
			disposal_date TIMESTAMPTZ, disposal_auth UUID,
			quantity_ammunition INTEGER DEFAULT 0,
			ammunition_type VARCHAR(50),
			photos_refs TEXT[] DEFAULT '{}', notes TEXT,
			created_by UUID NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS biar_batch_seizures (
			batch_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			batch_reference VARCHAR(50) UNIQUE NOT NULL,
			operation_name TEXT, seizure_date TIMESTAMPTZ NOT NULL,
			location_desc VARCHAR(300), dept_code CHAR(2),
			total_weapons INTEGER NOT NULL,
			weapon_ids UUID[] DEFAULT '{}',
			seizing_unit VARCHAR(50) NOT NULL, lead_officer UUID,
			partnering_agencies TEXT[] DEFAULT '{}', notes TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS biar_iarms_sync_log (
			sync_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			weapon_id UUID REFERENCES biar_illicit_weapons(weapon_id),
			direction VARCHAR(10) NOT NULL, iarms_ref VARCHAR(50),
			sync_status VARCHAR(20) DEFAULT 'PENDING',
			synced_at TIMESTAMPTZ, error_message TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
	}
	for _, m := range migrations {
		if _, err := pool.Exec(ctx, m); err != nil {
			logger.Warn("migration warning", zap.Error(err))
		}
	}
	logger.Info("migrations completed")
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	dbHost := getEnv("BIAR_DB_HOST", "localhost")
	dbPort := getEnv("BIAR_DB_PORT", "26257")
	dbName := getEnv("BIAR_DB_NAME", "snisid_biar")
	dbUser := getEnv("BIAR_DB_USER", "root")
	dbURL := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", dbUser, dbHost, dbPort, dbName)
	if u := os.Getenv("BIAR_DATABASE_URL"); u != "" {
		dbURL = u
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		logger.Fatal("failed to ping database", zap.Error(err))
	}

	runMigrations(pool, logger)

	repo := postgres.NewWeaponRepo(pool)
	svc := service.NewBIARService(repo, logger)

	r := handler.SetupRouter(svc, logger)

	kafkaBrokers := strings.Split(getEnv("BIAR_KAFKA_BROKERS", "kafka:9092"), ",")
	kafkaTopic := getEnv("BIAR_KAFKA_TOPIC", "biar.events")
	kafkaProducer := kafka.NewProducer(kafkaBrokers, kafkaTopic, logger)
	defer kafkaProducer.Close()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	port := getEnv("BIAR_SERVICE_PORT", ":8103")
	if port != "" && port[0] != ':' {
		port = ":" + port
	}

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		logger.Info("starting biar-svc", zap.String("addr", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}
	logger.Info("server exited")
}
