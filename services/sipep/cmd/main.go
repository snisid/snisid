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

	"github.com/snisid/platform/services/sipep/internal/api/rest"
	"github.com/snisid/platform/services/sipep/internal/service"
)

func main() {
	inmateSvc := service.NewInmateService()
	facilitySvc := service.NewFacilityService()
	movementSvc := service.NewMovementService()

	router := gin.Default()

	v1 := router.Group("/api/v1/sipep")
	{
		inmateHandler := rest.NewInmateHandler(inmateSvc)
		inmateHandler.RegisterRoutes(v1)

		facilityHandler := rest.NewFacilityHandler(facilitySvc, inmateSvc)
		facilityHandler.RegisterRoutes(v1)

		movementHandler := rest.NewMovementHandler(movementSvc, inmateSvc)
		movementHandler.RegisterRoutes(v1)
	}

	port := os.Getenv("SIPEP_SERVICE_PORT")
	if port == "" {
		port = "8092"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("SIPEP-HT listening on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down SIPEP-HT server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %s", err)
	}
	log.Println("SIPEP-HT server exited")
}
