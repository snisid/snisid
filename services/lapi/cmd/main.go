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
	"github.com/snisid/platform/services/lapi/internal/handler"
	"github.com/snisid/platform/services/lapi/internal/service"
)

func main() {
	cache := service.NewOfflineCache()
	syncSvc := service.NewLAPISyncService(cache)
	h := handler.NewHTTPHandler(syncSvc, cache)

	router := gin.Default()
	h.RegisterRoutes(router)

	srv := &http.Server{
		Addr:    ":8121",
		Handler: router,
	}

	go func() {
		log.Printf("lapi listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %s", err)
	}
	log.Println("server exited")
}
