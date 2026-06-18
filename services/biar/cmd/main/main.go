package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/snisid/platform/services/biar/internal/api/rest"
	"github.com/snisid/platform/services/biar/internal/repository"
	"github.com/snisid/platform/services/biar/internal/service"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	weaponRepo := repository.NewInMemoryWeaponRepo()
	batchRepo := repository.NewInMemoryBatchRepo()
	syncRepo := repository.NewInMemorySyncRepo()

	iarmsClient := service.NewIARMSClient(
		os.Getenv("BIAR_IARMS_GATEWAY"),
		os.Getenv("BIAR_ATF_API_KEY"),
		"HTI",
	)

	weaponSvc := service.NewWeaponService(weaponRepo)
	batchSvc := service.NewBatchService(batchRepo, weaponRepo)
	statsSvc := service.NewStatsService(weaponRepo)
	syncSvc := service.NewSyncService(iarmsClient, weaponRepo, syncRepo)

	router := rest.NewRouter(weaponSvc, batchSvc, statsSvc, syncSvc)

	port := os.Getenv("BIAR_SERVICE_PORT")
	if port == "" {
		port = "8103"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		logger.Info("BIAR-HT démarré", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Erreur serveur HTTP", zap.Error(err))
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}
