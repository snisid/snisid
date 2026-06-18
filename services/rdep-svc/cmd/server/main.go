package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/rdep-svc/internal/api/rest"
	"github.com/snisid/platform/services/rdep-svc/internal/domain"
	"github.com/snisid/platform/services/rdep-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/rdep-svc/internal/service"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dbHost := getEnv("RDEP_DB_HOST", "postgres://snisid:snisid@localhost:5432/snisid_rdep")
	dbPool, err := pgxpool.New(ctx, dbHost)
	if err != nil {
		logger.Fatal("cannot connect to database", zap.Error(err))
	}
	defer dbPool.Close()

	deporteeRepo := postgres.NewDeporteeRepo(dbPool)
	foreignRepo := postgres.NewForeignRecordRepo(dbPool)
	monitoringEventRepo := postgres.NewMonitoringEventRepo(dbPool)

	var eventPub domain.EventPublisher = &noopEventPublisher{}
	var fbiClient domain.FBIRecordClient = &noopFBIClient{}
	var interpClient domain.InterpolClient = &noopInterpolClient{}
	var afisClient domain.AFISClient = &noopAFISClient{}

	intakeSvc := service.NewIntakeService(deporteeRepo, foreignRepo, eventPub)
	screeningSvc := service.NewScreeningService(deporteeRepo, foreignRepo, fbiClient, interpClient, afisClient, eventPub)
	monitoringSvc := service.NewMonitoringService(deporteeRepo, monitoringEventRepo, eventPub)

	intakeHandler := rest.NewIntakeHandler(intakeSvc, screeningSvc)
	monitoringHandler := rest.NewMonitoringHandler(monitoringSvc)

	r := gin.Default()
	v1 := r.Group("/api/v1/rdep")
	{
		v1.POST("/intake", intakeHandler.ProcessIntake)
		v1.POST("/:id/screen", intakeHandler.ScreenDeportee)
		v1.GET("/:id", intakeHandler.GetDeportee)
		v1.GET("/high-risk", intakeHandler.GetHighRisk)
		v1.GET("/gang-affiliated", intakeHandler.GetGangAffiliated)
		v1.GET("/stats/by-country", intakeHandler.GetStatsByCountry)
		v1.POST("/:id/monitoring/events", monitoringHandler.RecordEvent)
		v1.GET("/:id/monitoring/events", monitoringHandler.GetEvents)
		v1.PUT("/:id/address", monitoringHandler.UpdateAddress)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", getEnv("RDEP_SERVICE_PORT", "8094")),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen error", zap.Error(err))
		}
	}()

	log.Println("RDEP-HT service started on port", getEnv("RDEP_SERVICE_PORT", "8094"))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	srv.Shutdown(shutdownCtx)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

type noopEventPublisher struct{}

func (n *noopEventPublisher) Publish(topic string, event interface{}) error {
	log.Printf("Event published to %s: %v", topic, event)
	return nil
}

type noopFBIClient struct{}

func (n *noopFBIClient) GetRecord(ctx context.Context, fbiNumber string) (*domain.ForeignRecord, error) {
	return nil, nil
}

type noopInterpolClient struct{}

func (n *noopInterpolClient) CheckNotices(ctx context.Context, personID uuid.UUID) ([]string, error) {
	return nil, nil
}

type noopAFISClient struct{}

func (n *noopAFISClient) CheckPrint(ctx context.Context, fingerprintData string) (*domain.AFISHit, error) {
	return nil, nil
}
