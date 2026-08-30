package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/snisid/platform/services/gang/internal/api/rest"
	"github.com/snisid/platform/services/gang/internal/repository"
	"github.com/snisid/platform/services/gang/internal/service"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	gangRepo := repository.NewInMemoryGangRepo()
	memberRepo := repository.NewInMemoryMemberRepo()
	incidentRepo := repository.NewInMemoryIncidentRepo()
	territoryRepo := repository.NewInMemoryTerritoryRepo()

	gangSvc := service.NewGangService(gangRepo)
	memberSvc := service.NewMemberService(memberRepo, gangRepo)
	incidentSvc := service.NewIncidentService(incidentRepo, gangRepo)
	territorySvc := service.NewTerritoryService(territoryRepo, gangRepo)

	router := rest.NewRouter(gangSvc, memberSvc, incidentSvc, territorySvc)

	srv := &http.Server{
		Addr:    ":8095",
		Handler: router,
	}

	go func() {
		logger.Info("GANG-HT démarré", zap.String("port", "8095"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Erreur serveur HTTP", zap.Error(err))
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}
