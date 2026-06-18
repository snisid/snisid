package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/snisid/platform/services/corr/internal/api/rest"
	"github.com/snisid/platform/services/corr/internal/repository"
	"github.com/snisid/platform/services/corr/internal/service"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	caseRepo := repository.NewInMemoryCaseRepo()
	officerRepo := repository.NewInMemoryOfficerRepo()
	evidenceRepo := repository.NewInMemoryEvidenceRepo()
	wbRepo := repository.NewInMemoryWhistleblowerRepo()
	alertRepo := repository.NewInMemoryAlertRepo()
	declRepo := repository.NewInMemoryDeclarationRepo()

	caseSvc := service.NewCaseService(caseRepo, officerRepo)
	investigationSvc := service.NewInvestigationService(caseRepo, evidenceRepo)
	evidenceSvc := service.NewEvidenceService(evidenceRepo)
	wbSvc := service.NewWhistleblowerService(wbRepo, caseRepo)
	alertSvc := service.NewAlertService(alertRepo, declRepo, caseRepo)

	router := rest.NewRouter(caseSvc, investigationSvc, evidenceSvc, wbSvc, alertSvc)

	srv := &http.Server{
		Addr:    ":8130",
		Handler: router,
	}

	go func() {
		logger.Info("CORR-HT démarré", zap.String("port", "8130"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Erreur serveur HTTP", zap.Error(err))
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}
