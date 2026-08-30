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
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/gang-svc/internal/api/rest"
	"github.com/snisid/platform/services/gang-svc/internal/domain"
	"github.com/snisid/platform/services/gang-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/gang-svc/internal/service"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dbHost := getEnv("GANG_DB_HOST", "postgres://snisid:snisid@localhost:5432/snisid_gang")
	dbPool, err := pgxpool.New(ctx, dbHost)
	if err != nil {
		logger.Fatal("cannot connect to database", zap.Error(err))
	}
	defer dbPool.Close()

	orgRepo := postgres.NewOrganizationRepo(dbPool)
	incRepo := postgres.NewIncidentRepo(dbPool)
	alliRepo := postgres.NewAllianceRepo(dbPool)
	var eventPub domain.EventPublisher = &noopEventPublisher{}

	orgSvc := service.NewOrganizationService(orgRepo, incRepo, alliRepo, eventPub)
	incidentSvc := service.NewIncidentService(incRepo, orgRepo, eventPub)
	allianceSvc := service.NewAllianceService(alliRepo, orgRepo, eventPub)

	orgHandler := rest.NewOrganizationHandler(orgSvc, incidentSvc, allianceSvc)

	r := gin.Default()
	v1 := r.Group("/api/v1/gangs")
	{
		v1.POST("", orgHandler.CreateOrganization)
		v1.GET("", orgHandler.ListOrganizations)
		v1.GET("/:id", orgHandler.GetOrganization)
		v1.POST("/:id/incidents", orgHandler.CreateIncident)
		v1.GET("/:id/incidents", orgHandler.GetIncidentsByGang)
		v1.GET("/by-dept/:code", orgHandler.GetByDeptCode)
		v1.GET("/alliances/map", orgHandler.GetAllianceMap)
		v1.GET("/sanctioned", orgHandler.GetSanctioned)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", getEnv("GANG_SERVICE_PORT", "8095")),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen error", zap.Error(err))
		}
	}()

	log.Println("GANG-HT service started on port", getEnv("GANG_SERVICE_PORT", "8095"))

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
