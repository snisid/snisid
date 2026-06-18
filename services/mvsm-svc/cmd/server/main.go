package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/mvsm-svc/internal/api/rest"
	"github.com/snisid/platform/services/mvsm-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/mvsm-svc/internal/service"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/snisid?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		logger.Fatal("failed to ping database", zap.Error(err))
	}

	repo := postgres.NewEventRepo(pool)
	svc := service.NewEventService(repo, logger)
	handler := rest.NewEventHandler(svc, logger)

	r := gin.Default()

	api := r.Group("/api/v1/mvsm")
	{
		api.POST("/events", handler.CreateEvent)
		api.GET("/events/upcoming", handler.ListUpcoming)
		api.POST("/events/:id/updates", handler.AddUpdate)
		api.GET("/events/active", handler.ListActive)
		api.PATCH("/events/:id/risk", handler.UpdateRiskLevel)
	}

	port := os.Getenv("MVSM_SERVICE_PORT")
	if port == "" {
		port = ":8127"
	}

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		logger.Info("starting mvsm-svc", zap.String("addr", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}
	logger.Info("server exited")
}
