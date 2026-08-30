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

	"github.com/snisid/platform/services/aero-svc/internal/handler"
	"github.com/snisid/platform/services/aero-svc/internal/kafka"
	"github.com/snisid/platform/services/aero-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/aero-svc/internal/service"
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
		`CREATE TYPE IF NOT EXISTS aero_aircraft_type AS ENUM (
			'COMMERCIAL_JET','TURBOPROP','PISTON_SINGLE','PISTON_TWIN',
			'HELICOPTER','ULTRALIGHT','DRONE_LARGE','UNKNOWN'
		)`,
		`CREATE TYPE IF NOT EXISTS aero_strip_status AS ENUM (
			'ACTIVE','INACTIVE','DESTROYED','LEGALIZED','UNDER_SURVEILLANCE'
		)`,
		`CREATE TABLE IF NOT EXISTS aero_aircraft_registry (
			aircraft_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			registration_mark VARCHAR(20), icao_hex_code VARCHAR(10),
			aircraft_type aero_aircraft_type NOT NULL,
			make VARCHAR(100), model VARCHAR(100),
			manufacture_year SMALLINT, flag_country CHAR(3),
			owner_name VARCHAR(200), owner_snisid_id UUID,
			operator_name VARCHAR(200),
			is_registered BOOLEAN DEFAULT FALSE,
			is_suspected BOOLEAN DEFAULT FALSE,
			is_stolen BOOLEAN DEFAULT FALSE, gang_id UUID,
			drug_trafficking BOOLEAN DEFAULT FALSE,
			interpol_ref VARCHAR(50), faa_registry_ref VARCHAR(50),
			notes TEXT, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS aero_clandestine_strips (
			strip_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			strip_name VARCHAR(150), dept_code CHAR(2) NOT NULL,
			commune VARCHAR(100),
			lat DECIMAL(10,7) NOT NULL, lng DECIMAL(10,7) NOT NULL,
			length_m INTEGER, surface_type VARCHAR(30),
			status aero_strip_status NOT NULL DEFAULT 'ACTIVE',
			capable_aircraft TEXT[] DEFAULT '{}', gang_id UUID,
			first_detected DATE, last_activity_date DATE,
			source_intel TEXT, satellite_image_ref VARCHAR(500),
			created_by UUID NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS aero_suspicious_flights (
			flight_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			aircraft_id UUID REFERENCES aero_aircraft_registry(aircraft_id),
			registration_mark VARCHAR(20),
			flight_date TIMESTAMPTZ NOT NULL,
			origin_airport VARCHAR(10), destination_airport VARCHAR(10),
			origin_country CHAR(3), destination_country CHAR(3) DEFAULT 'HTI',
			landing_strip_id UUID REFERENCES aero_clandestine_strips(strip_id),
			landing_location VARCHAR(300), flight_type VARCHAR(30),
			cargo_suspected TEXT, source_radar VARCHAR(50),
			source_informant BOOLEAN DEFAULT FALSE,
			case_reference VARCHAR(100), created_by UUID NOT NULL,
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

	dbHost := getEnv("AERO_DB_HOST", "localhost")
	dbPort := getEnv("AERO_DB_PORT", "26257")
	dbName := getEnv("AERO_DB_NAME", "snisid_aero")
	dbUser := getEnv("AERO_DB_USER", "root")
	dbURL := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", dbUser, dbHost, dbPort, dbName)
	if u := os.Getenv("AERO_DATABASE_URL"); u != "" {
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

	repo := postgres.NewAircraftRepo(pool)
	svc := service.NewAeroService(repo, logger)

	r := handler.SetupRouter(svc, logger)

	kafkaBrokers := strings.Split(getEnv("AERO_KAFKA_BROKERS", "kafka:9092"), ",")
	kafkaTopic := getEnv("AERO_KAFKA_TOPIC", "aero.events")
	kafkaProducer := kafka.NewProducer(kafkaBrokers, kafkaTopic, logger)
	defer kafkaProducer.Close()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	port := getEnv("AERO_SERVICE_PORT", ":8109")
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
		logger.Info("starting aero-svc", zap.String("addr", port))
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
