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

	"github.com/snisid/platform/services/mar-svc/internal/handler"
	"github.com/snisid/platform/services/mar-svc/internal/kafka"
	"github.com/snisid/platform/services/mar-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/mar-svc/internal/service"
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
		`CREATE TYPE IF NOT EXISTS mar_vessel_type AS ENUM (
			'CARGO_SHIP','TANKER','FISHING_BOAT','GO_FAST',
			'SAILBOAT','YACHT','FERRY','PATROL_BOAT',
			'WOODEN_BOAT','CANOE','UNKNOWN'
		)`,
		`CREATE TYPE IF NOT EXISTS mar_vessel_status AS ENUM (
			'REGISTERED','STOLEN','SUSPECTED','DETAINED',
			'SUNK','DESTROYED','MISSING','INTERPOL_ALERT'
		)`,
		`CREATE TYPE IF NOT EXISTS mar_incident_type AS ENUM (
			'DRUG_SEIZURE','ARMS_SEIZURE','MIGRANT_INTERDICTION',
			'SMUGGLING','SUSPICIOUS_ACTIVITY','DISTRESS',
			'PIRACY','ILLEGAL_FISHING','HUMAN_TRAFFICKING'
		)`,
		`CREATE TABLE IF NOT EXISTS mar_vessels (
			vessel_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			national_mar_id VARCHAR(25) UNIQUE NOT NULL,
			vessel_name VARCHAR(150), imo_number VARCHAR(20),
			mmsi VARCHAR(15), call_sign VARCHAR(15),
			vessel_type mar_vessel_type NOT NULL,
			flag_country CHAR(3), hull_color VARCHAR(50),
			length_m DECIMAL(8,2), tonnage_gt INTEGER,
			engine_count SMALLINT, horsepower INTEGER,
			owner_name VARCHAR(200), owner_snisid_id UUID,
			registration_number VARCHAR(50),
			registration_port VARCHAR(100),
			status mar_vessel_status NOT NULL DEFAULT 'REGISTERED',
			gang_id UUID, interpol_svd_ref VARCHAR(50),
			notes TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS mar_ais_sightings (
			sighting_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			vessel_id UUID REFERENCES mar_vessels(vessel_id),
			mmsi VARCHAR(15), vessel_name VARCHAR(150),
			sighting_timestamp TIMESTAMPTZ NOT NULL,
			lat DECIMAL(10,7) NOT NULL, lng DECIMAL(10,7) NOT NULL,
			speed_knots DECIMAL(5,2), heading_degrees SMALLINT,
			destination VARCHAR(100),
			source_type VARCHAR(30), zone_code VARCHAR(20),
			alert_triggered BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS mar_incidents (
			incident_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			vessel_id UUID REFERENCES mar_vessels(vessel_id),
			incident_type mar_incident_type NOT NULL,
			incident_date TIMESTAMPTZ NOT NULL,
			lat DECIMAL(10,7), lng DECIMAL(10,7),
			zone_desc VARCHAR(100), responding_unit VARCHAR(50),
			outcome TEXT, persons_involved INTEGER DEFAULT 0,
			snisid_person_ids UUID[] DEFAULT '{}',
			drug_types TEXT[] DEFAULT '{}',
			drug_weight_kg DECIMAL(12,3),
			weapons_found BOOLEAN DEFAULT FALSE,
			weapons_count INTEGER DEFAULT 0,
			migrants_count INTEGER DEFAULT 0,
			biar_refs UUID[] DEFAULT '{}',
			case_reference VARCHAR(100), photo_refs TEXT[] DEFAULT '{}',
			created_by UUID NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS mar_watch_vessels (
			watch_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			vessel_id UUID REFERENCES mar_vessels(vessel_id),
			mmsi VARCHAR(15), vessel_name VARCHAR(150),
			watch_reason TEXT NOT NULL,
			alert_level VARCHAR(20) DEFAULT 'CAUTION',
			requesting_unit VARCHAR(50),
			is_active BOOLEAN DEFAULT TRUE,
			expiry_date TIMESTAMPTZ, created_by UUID NOT NULL,
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

	dbHost := getEnv("MAR_DB_HOST", "localhost")
	dbPort := getEnv("MAR_DB_PORT", "26257")
	dbName := getEnv("MAR_DB_NAME", "snisid_mar")
	dbUser := getEnv("MAR_DB_USER", "root")
	dbURL := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", dbUser, dbHost, dbPort, dbName)
	if u := os.Getenv("MAR_DATABASE_URL"); u != "" {
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

	repo := postgres.NewMaritimeRepo(pool)
	svc := service.NewMaritimeService(repo, logger)

	r := handler.SetupRouter(svc, logger)

	kafkaBrokers := strings.Split(getEnv("MAR_KAFKA_BROKERS", "kafka:9092"), ",")
	kafkaTopic := getEnv("MAR_KAFKA_TOPIC", "mar.events")
	kafkaProducer := kafka.NewProducer(kafkaBrokers, kafkaTopic, logger)
	defer kafkaProducer.Close()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	port := getEnv("MAR_SERVICE_PORT", ":8107")
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
		logger.Info("starting mar-svc", zap.String("addr", port))
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
