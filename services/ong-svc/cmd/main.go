package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/ong-svc/internal/handler"
	"github.com/snisid/platform/services/ong-svc/internal/kafka"
	"github.com/snisid/platform/services/ong-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/ong-svc/internal/service"
)

func runMigrations(pool *pgxpool.Pool) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS ong_organizations (
		org_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		national_ong_id VARCHAR(25) UNIQUE NOT NULL,
		org_name VARCHAR(200) NOT NULL,
		org_name_local VARCHAR(200),
		acronym VARCHAR(20),
		org_type VARCHAR(30) NOT NULL,
		registration_status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
		mjsp_registration_number VARCHAR(50),
		headquarter_country CHAR(3) NOT NULL,
		headquarter_city VARCHAR(100),
		haiti_office_dept CHAR(2),
		haiti_office_address TEXT,
		operating_depts TEXT[] DEFAULT '{}',
		sectors TEXT[] DEFAULT '{}',
		annual_budget_usd DECIMAL(15,2),
		haiti_staff_count INTEGER DEFAULT 0,
		expat_staff_count INTEGER DEFAULT 0,
		director_name VARCHAR(200),
		director_snisid_id UUID,
		director_nationality CHAR(3),
		contact_email VARCHAR(200),
		contact_phone VARCHAR(30),
		risk_flag VARCHAR(40) NOT NULL DEFAULT 'NONE',
		risk_notes TEXT,
		is_access_restricted BOOLEAN DEFAULT FALSE,
		created_by UUID NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE TABLE IF NOT EXISTS ong_staff_registry (
		staff_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		org_id UUID NOT NULL REFERENCES ong_organizations(org_id),
		full_name VARCHAR(200) NOT NULL,
		nationality CHAR(3) NOT NULL,
		role VARCHAR(100),
		is_expatriate BOOLEAN DEFAULT FALSE,
		passport_number VARCHAR(50),
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE TABLE IF NOT EXISTS ong_field_access_requests (
		request_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		org_id UUID NOT NULL REFERENCES ong_organizations(org_id),
		access_type VARCHAR(30),
		requested_zones TEXT[] DEFAULT '{}',
		access_date DATE NOT NULL,
		purpose TEXT NOT NULL,
		staff_count SMALLINT DEFAULT 1,
		status VARCHAR(20) DEFAULT 'PENDING',
		pnh_escort_required BOOLEAN DEFAULT FALSE,
		approved_by UUID,
		approval_notes TEXT,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`)
	return err
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/snisid_ong?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		logger.Fatal("failed to ping database", zap.Error(err))
	}

	if err := runMigrations(pool); err != nil {
		logger.Fatal("auto-migration failed", zap.Error(err))
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "kafka:9092"
	}
	producer := kafka.NewProducer(kafkaBrokers, logger)
	defer producer.Close()

	repo := postgres.NewONGRepo(pool)
	svc := service.NewONGService(repo, logger)
	h := handler.NewONGHandler(svc, logger)

	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "ong-svc", "uptime": time.Since(startTime).String()})
	})

	api := r.Group("/api/v1/ong")
	{
		api.POST("/organizations", h.RegisterOrg)
		api.GET("/organizations", h.ListOrgs)
		api.GET("/organizations/:id", h.GetOrg)
		api.POST("/organizations/:id/screen", h.ScreenOrg)
		api.POST("/staff", h.RegisterStaff)
		api.POST("/access-requests", h.RequestAccess)
		api.PATCH("/access-requests/:id", h.ApproveAccess)
		api.GET("/flagged", h.ListFlagged)
		api.GET("/unregistered", h.ListUnregistered)
	}

	port := os.Getenv("ONG_SERVICE_PORT")
	if port == "" {
		port = ":8133"
	}

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		logger.Info("starting ong-svc", zap.String("addr", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}
	logger.Info("server exited")
}

var startTime = time.Now()
