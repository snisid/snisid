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

	"github.com/snisid/platform/services/siar-svc/internal/handler"
	"github.com/snisid/platform/services/siar-svc/internal/kafka"
	"github.com/snisid/platform/services/siar-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/siar-svc/internal/service"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func runMigrations(pool *pgxpool.Pool, logger *zap.Logger) error {
	ctx := context.Background()
	migrations := []string{
		`CREATE TYPE IF NOT EXISTS siar_weapon_type AS ENUM (
			'HANDGUN','RIFLE','SHOTGUN','SUBMACHINE_GUN','ASSAULT_RIFLE',
			'MACHINE_GUN','SNIPER','RPG','GRENADE','HOMEMADE','OTHER'
		)`,
		`CREATE TYPE IF NOT EXISTS siar_status AS ENUM (
			'REGISTERED','REPORTED_STOLEN','SEIZED','DESTROYED',
			'REPORTED_LOST','TRANSFERRED','DEACTIVATED'
		)`,
		`CREATE TYPE IF NOT EXISTS siar_registration_type AS ENUM (
			'CIVILIAN','POLICE','MILITARY','SECURITY_COMPANY',
			'EMBASSY','ILLEGAL_FOUND','HISTORICAL'
		)`,
		`CREATE TABLE IF NOT EXISTS siar_firearms (
			firearm_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			national_siar_id VARCHAR(25) UNIQUE NOT NULL,
			serial_number VARCHAR(100), make VARCHAR(100) NOT NULL,
			model VARCHAR(100) NOT NULL, caliber VARCHAR(30) NOT NULL,
			weapon_type siar_weapon_type NOT NULL, manufacture_year SMALLINT,
			manufacture_country CHAR(3), status siar_status NOT NULL DEFAULT 'REGISTERED',
			reg_type siar_registration_type NOT NULL,
			owner_snisid_id UUID, owner_entity_name VARCHAR(200),
			license_number VARCHAR(50), license_expiry DATE,
			import_date DATE, import_country CHAR(3),
			import_permit_ref VARCHAR(100), importer_name VARCHAR(200),
			customs_entry_ref VARCHAR(100),
			current_dept_code CHAR(2), storage_location TEXT,
			fir_record_id UUID, gang_id UUID,
			case_references TEXT[] DEFAULT '{}',
			iarms_ref VARCHAR(50), atf_etrace_ref VARCHAR(50),
			notes TEXT, created_by UUID NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS siar_licenses (
			license_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			license_number VARCHAR(50) UNIQUE NOT NULL,
			holder_snisid_id UUID NOT NULL,
			license_type VARCHAR(50) NOT NULL,
			firearms_authorized INTEGER DEFAULT 1,
			issue_date DATE NOT NULL, expiry_date DATE NOT NULL,
			issuing_authority VARCHAR(100) NOT NULL,
			is_active BOOLEAN DEFAULT TRUE,
			revocation_reason TEXT, revoked_at TIMESTAMPTZ,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS siar_transfers (
			transfer_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			firearm_id UUID NOT NULL REFERENCES siar_firearms(firearm_id),
			from_owner_id UUID, to_owner_id UUID,
			transfer_type VARCHAR(50),
			transfer_date DATE NOT NULL, permit_ref VARCHAR(100),
			authorized_by UUID, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS siar_seizures (
			seizure_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			firearm_id UUID REFERENCES siar_firearms(firearm_id),
			seizure_date TIMESTAMPTZ NOT NULL,
			seizing_unit VARCHAR(50) NOT NULL,
			seizing_officer UUID, location_desc VARCHAR(300),
			dept_code CHAR(2), context TEXT,
			from_person_id UUID, gang_id UUID,
			case_reference VARCHAR(100),
			disposed_of BOOLEAN DEFAULT FALSE,
			disposal_method VARCHAR(50),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
	}
	for _, m := range migrations {
		if _, err := pool.Exec(ctx, m); err != nil {
			logger.Warn("migration warning", zap.Error(err), zap.String("sql", m[:50]))
		}
	}
	logger.Info("migrations completed")
	return nil
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	dbHost := getEnv("SIAR_DB_HOST", "localhost")
	dbPort := getEnv("SIAR_DB_PORT", "26257")
	dbName := getEnv("SIAR_DB_NAME", "snisid_siar")
	dbUser := getEnv("SIAR_DB_USER", "root")
	dbURL := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", dbUser, dbHost, dbPort, dbName)
	if u := os.Getenv("SIAR_DATABASE_URL"); u != "" {
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

	if err := runMigrations(pool, logger); err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}

	repo := postgres.NewFirearmRepo(pool)
	svc := service.NewSIARService(repo, logger)

	r := handler.SetupRouter(svc, logger)

	kafkaBrokers := strings.Split(getEnv("SIAR_KAFKA_BROKERS", "kafka:9092"), ",")
	kafkaTopic := getEnv("SIAR_KAFKA_TOPIC", "siar.events")
	kafkaProducer := kafka.NewProducer(kafkaBrokers, kafkaTopic, logger)
	defer kafkaProducer.Close()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	port := getEnv("SIAR_SERVICE_PORT", ":8102")
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
		logger.Info("starting siar-svc", zap.String("addr", port))
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
