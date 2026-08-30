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
	"github.com/snisid/platform/services/chef/internal/api/rest"
	"github.com/snisid/platform/services/chef/internal/service"
)

func main() {
	repo := service.NewInMemoryRepository()
	memberSvc := service.NewMemberService(repo)
	h := rest.NewHTTPHandler(memberSvc)

	router := gin.Default()
	h.RegisterRoutes(router)

	srv := &http.Server{
		Addr:    ":8097",
		Handler: router,
	}

	go func() {
		log.Printf("CHEF-HT listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down CHEF-HT server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %s", err)
	}
	log.Println("CHEF-HT server exited")
}
