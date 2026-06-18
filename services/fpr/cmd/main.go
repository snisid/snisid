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
	"github.com/snisid/platform/services/fpr/internal/handler"
	"github.com/snisid/platform/services/fpr/internal/service"
)

func main() {
	matcher := service.NewFPRMatcher(40.0)
	tmpl := service.NewTemplateManager()
	h := handler.NewHTTPHandler(matcher, tmpl)

	router := gin.Default()
	h.RegisterRoutes(router)

	srv := &http.Server{
		Addr:    ":8122",
		Handler: router,
	}

	go func() {
		log.Printf("fpr listening on %s", srv.Addr)
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
