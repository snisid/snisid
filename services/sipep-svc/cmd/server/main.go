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

	"github.com/snisid/platform/services/sipep-svc/internal/api/rest"
	"github.com/snisid/platform/services/sipep-svc/internal/domain"
	"github.com/snisid/platform/services/sipep-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/sipep-svc/internal/service"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dbHost := getEnv("SIPEP_DB_HOST", "postgres://snisid:snisid@localhost:5432/snisid_sipep")
	dbPool, err := pgxpool.New(ctx, dbHost)
	if err != nil {
		logger.Fatal("cannot connect to database", zap.Error(err))
	}
	defer dbPool.Close()

	inmateRepo := postgres.NewInmateRepo(dbPool)
	detentionRepo := postgres.NewDetentionRepo(dbPool)
	transferRepo := postgres.NewTransferRepo(dbPool)

	var eventPub domain.EventPublisher = &noopEventPublisher{}
	var snisidClient domain.SNISIDClient = &noopSNISIDClient{}

	intakeSvc := service.NewIntakeService(inmateRepo, detentionRepo, eventPub, snisidClient)
	releaseSvc := service.NewReleaseService(inmateRepo, detentionRepo, eventPub)
	transferSvc := service.NewTransferService(inmateRepo, detentionRepo, transferRepo, eventPub)
	overcrowdingSvc := service.NewOvercrowdingService(dbPool)

	intakeHandler := rest.NewIntakeHandler(intakeSvc)
	inmateHandler := rest.NewInmateHandler(intakeSvc, releaseSvc, transferSvc)
	statsHandler := rest.NewStatsHandler(overcrowdingSvc)

	r := gin.Default()
	v1 := r.Group("/api/v1/sipep")
	{
		v1.POST("/intake", intakeHandler.ProcessIntake)
		v1.GET("/inmates/:id", inmateHandler.GetInmate)
		v1.GET("/inmates/search", inmateHandler.SearchInmates)
		v1.POST("/release", inmateHandler.ProcessRelease)
		v1.POST("/transfers", inmateHandler.ProcessTransfer)
		v1.GET("/inmates/:id/transfers", inmateHandler.GetTransfers)
		v1.GET("/facilities/occupancy", statsHandler.GetFacilityOccupancy)
		v1.GET("/alerts/overcrowding", statsHandler.GetOvercrowdingAlerts)
		v1.GET("/stats/preventive-detention", statsHandler.GetPreventiveDetentionStats)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", getEnv("SIPEP_SERVICE_PORT", "8092")),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen error", zap.Error(err))
		}
	}()

	log.Println("SIPEP-HT service started on port", getEnv("SIPEP_SERVICE_PORT", "8092"))

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

type noopSNISIDClient struct{}

func (n *noopSNISIDClient) GetPerson(personID uuid.UUID) (*domain.PersonInfo, error) {
	return &domain.PersonInfo{
		PersonID:    personID,
		FullName:    "Test User",
		Nationality: "HTI",
	}, nil
}
