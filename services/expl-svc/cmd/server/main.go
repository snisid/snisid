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

	"github.com/snisid/platform/services/expl-svc/internal/handler"
	"github.com/snisid/platform/services/expl-svc/internal/kafka"
	"github.com/snisid/platform/services/expl-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/expl-svc/internal/service"
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
		`CREATE TYPE IF NOT EXISTS expl_type AS ENUM (
			'IED','GRENADE','RPG','MORTAR','LANDMINE','DYNAMITE',
			'BLASTING_CAP','AMMUNITION_BULK','MILITARY_ORDNANCE','UNKNOWN'
		)`,
		`CREATE TYPE IF NOT EXISTS expl_status AS ENUM (
			'RECOVERED','DESTROYED','DETONATED','STORED_EVIDENCE','TRANSFERRED'
		)`,
		`CREATE TABLE IF NOT EXISTS expl_incidents (
			incident_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			national_expl_id VARCHAR(25) UNIQUE NOT NULL,
			incident_type VARCHAR(30) NOT NULL,
			explosive_type expl_type NOT NULL,
			status expl_status NOT NULL DEFAULT 'RECOVERED',
			quantity INTEGER DEFAULT 1, weight_kg DECIMAL(10,3),
			manufacturer VARCHAR(100), lot_number VARCHAR(50),
			manufacture_country CHAR(3), estimated_date DATE,
			incident_date TIMESTAMPTZ NOT NULL,
			location_desc VARCHAR(300), dept_code CHAR(2),
			commune VARCHAR(100), lat DECIMAL(10,7), lng DECIMAL(10,7),
			responding_unit VARCHAR(50), eod_officer UUID,
			casualties SMALLINT DEFAULT 0, gang_id UUID,
			from_person_id UUID, case_reference VARCHAR(100),
			dna_sample_taken BOOLEAN DEFAULT FALSE,
			bio_sample_ref VARCHAR(100),
			photo_refs TEXT[] DEFAULT '{}',
			interpol_exploint_ref VARCHAR(50),
			notes TEXT, created_by UUID NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS expl_legal_stocks (
			stock_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			holder_entity VARCHAR(200) NOT NULL,
			holder_type VARCHAR(30),
			explosive_type expl_type NOT NULL,
			quantity_kg DECIMAL(12,3) NOT NULL,
			storage_location TEXT NOT NULL, dept_code CHAR(2),
			license_ref VARCHAR(50),
			last_audit_date DATE, next_audit_date DATE,
			is_secured BOOLEAN DEFAULT TRUE,
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

	dbHost := getEnv("EXPL_DB_HOST", "localhost")
	dbPort := getEnv("EXPL_DB_PORT", "26257")
	dbName := getEnv("EXPL_DB_NAME", "snisid_expl")
	dbUser := getEnv("EXPL_DB_USER", "root")
	dbURL := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", dbUser, dbHost, dbPort, dbName)
	if u := os.Getenv("EXPL_DATABASE_URL"); u != "" {
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

	repo := postgres.NewIncidentRepository(pool)
	svc := service.NewExplService(repo, logger)

	r := handler.SetupRouter(svc, logger)

	kafkaBrokers := strings.Split(getEnv("EXPL_KAFKA_BROKERS", "kafka:9092"), ",")
	kafkaTopic := getEnv("EXPL_KAFKA_TOPIC", "expl.events")
	kafkaProducer := kafka.NewProducer(kafkaBrokers, kafkaTopic, logger)
	defer kafkaProducer.Close()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	port := getEnv("EXPL_SERVICE_PORT", ":8104")
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
		logger.Info("starting expl-svc", zap.String("addr", port))
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
