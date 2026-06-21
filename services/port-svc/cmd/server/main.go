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

	"github.com/snisid/platform/services/port-svc/internal/handler"
	"github.com/snisid/platform/services/port-svc/internal/kafka"
	"github.com/snisid/platform/services/port-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/port-svc/internal/service"
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
		`CREATE TYPE IF NOT EXISTS port_risk_level AS ENUM ('LOW','MEDIUM','HIGH','CRITICAL')`,
		`CREATE TYPE IF NOT EXISTS port_container_status AS ENUM (
			'PENDING_INSPECTION','CLEARED','HELD_FOR_INSPECTION',
			'SEIZED','RELEASED_AFTER_INSPECTION'
		)`,
		`CREATE TABLE IF NOT EXISTS port_vessels_arrivals (
			arrival_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			port_code VARCHAR(10) NOT NULL,
			vessel_imo VARCHAR(20), vessel_name VARCHAR(150) NOT NULL,
			flag_country CHAR(3), shipping_company VARCHAR(200),
			arrival_date TIMESTAMPTZ NOT NULL,
			origin_port VARCHAR(100), origin_country CHAR(3),
			container_count INTEGER DEFAULT 0, manifest_ref VARCHAR(100),
			mar_vessel_id UUID, risk_score SMALLINT DEFAULT 0,
			risk_level port_risk_level DEFAULT 'LOW',
			cbp_targeting_ref VARCHAR(50),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS port_containers (
			container_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			arrival_id UUID NOT NULL REFERENCES port_vessels_arrivals(arrival_id),
			container_number VARCHAR(20) NOT NULL,
			container_type VARCHAR(10),
			declared_content TEXT NOT NULL,
			declared_weight_kg DECIMAL(12,3),
			declared_value_usd DECIMAL(15,2),
			shipper_name VARCHAR(200), shipper_country CHAR(3),
			consignee_name VARCHAR(200), consignee_snisid_id UUID,
			status port_container_status NOT NULL DEFAULT 'PENDING_INSPECTION',
			risk_score SMALLINT DEFAULT 0, risk_level port_risk_level DEFAULT 'LOW',
			risk_flags TEXT[] DEFAULT '{}',
			selected_for_scan BOOLEAN DEFAULT FALSE,
			scan_date TIMESTAMPTZ, scan_result TEXT,
			seized BOOLEAN DEFAULT FALSE, seizure_description TEXT,
			case_reference VARCHAR(100),
			cbp_targeting_match BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS port_risk_factors (
			factor_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			container_id UUID NOT NULL REFERENCES port_containers(container_id),
			factor_type VARCHAR(50) NOT NULL, description TEXT NOT NULL,
			weight_score SMALLINT NOT NULL, source VARCHAR(50),
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

	dbHost := getEnv("PORT_DB_HOST", "localhost")
	dbPort := getEnv("PORT_DB_PORT", "26257")
	dbName := getEnv("PORT_DB_NAME", "snisid_port")
	dbUser := getEnv("PORT_DB_USER", "root")
	dbURL := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", dbUser, dbHost, dbPort, dbName)
	if u := os.Getenv("PORT_DATABASE_URL"); u != "" {
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

	repo := postgres.NewContainerRepo(pool)
	svc := service.NewPortService(repo, logger)

	r := handler.SetupRouter(svc, logger)

	kafkaBrokers := strings.Split(getEnv("PORT_KAFKA_BROKERS", "kafka:9092"), ",")
	kafkaTopic := getEnv("PORT_KAFKA_TOPIC", "port.events")
	kafkaProducer := kafka.NewProducer(kafkaBrokers, kafkaTopic, logger)
	defer kafkaProducer.Close()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	port := getEnv("PORT_SERVICE_PORT", ":8111")
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
		logger.Info("starting port-svc", zap.String("addr", port))
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
