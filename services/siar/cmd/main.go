package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/snisid/platform/services/siar/internal/api/rest"
	"github.com/snisid/platform/services/siar/internal/service"
)

func main() {
	firearmSvc := service.NewFirearmService()
	licenseSvc := service.NewLicenseService()
	transferSvc := service.NewTransferService()
	dealerSvc := service.NewDealerService()

	router := rest.NewRouter(firearmSvc, licenseSvc, transferSvc, dealerSvc)

	port := os.Getenv("SIAR_SERVICE_PORT")
	if port == "" {
		port = "8102"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("SIAR-HT listening on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("SIAR-HT listen error: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down SIAR-HT server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("SIAR-HT forced shutdown: %s", err)
	}
	log.Println("SIAR-HT server exited")
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}
