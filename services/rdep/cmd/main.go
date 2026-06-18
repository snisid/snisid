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
	"github.com/snisid/platform/services/rdep/internal/api/rest"
	"github.com/snisid/platform/services/rdep/internal/service"
)

func main() {
	deporteeSvc := service.NewDeporteeService()
	extraditionSvc := service.NewExtraditionService()
	flightSvc := service.NewFlightService()

	deporteeHandler := rest.NewDeporteeHandler(deporteeSvc)
	extraditionHandler := rest.NewExtraditionHandler(extraditionSvc)
	flightHandler := rest.NewFlightHandler(flightSvc)

	router := gin.Default()

	v1 := router.Group("/api/v1/rdep")
	{
		v1.POST("/intake", deporteeHandler.Intake)
		v1.GET("/:id", deporteeHandler.GetDeportee)
		v1.POST("/:id/screen", deporteeHandler.ScreenDeportee)
		v1.POST("/:id/monitoring/events", deporteeHandler.AddMonitoringEvent)
		v1.GET("/high-risk", deporteeHandler.ListHighRisk)
		v1.GET("/gang-affiliated", deporteeHandler.ListGangAffiliated)
		v1.GET("/stats/by-country", deporteeHandler.StatsByCountry)

		v1.POST("/extraditions", extraditionHandler.Create)
		v1.GET("/extraditions/:id", extraditionHandler.Get)
		v1.GET("/extraditions", extraditionHandler.List)
		v1.PUT("/extraditions/:id/status", extraditionHandler.UpdateStatus)

		v1.POST("/flights", flightHandler.Create)
		v1.GET("/flights/:id", flightHandler.Get)
		v1.GET("/flights", flightHandler.List)
		v1.GET("/flights/by-number/:flight_number", flightHandler.GetByNumber)
	}

	srv := &http.Server{
		Addr:    ":8094",
		Handler: router,
	}

	go func() {
		log.Printf("RDEP-HT listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down RDEP-HT server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %s", err)
	}
	log.Println("RDEP-HT server exited")
}
