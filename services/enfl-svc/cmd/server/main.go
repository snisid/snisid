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

	"github.com/snisid/platform/services/enfl-svc/internal/api/rest"
	"github.com/snisid/platform/services/enfl-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/enfl-svc/internal/service"
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

	repo := postgres.NewChildRepo(pool)
	svc := service.NewChildService(repo, logger)
	handler := rest.NewChildHandler(svc, logger)

	r := gin.Default()

	api := r.Group("/api/v1/enfl")
	{
		api.POST("/children", handler.RegisterChild)
		api.GET("/children/:id", handler.GetChild)
		api.GET("/missing", handler.ListMissing)
		api.GET("/restaveks", handler.ListRestaveks)
		api.POST("/children/:id/locate", handler.LocateChild)
		api.GET("/gang-recruited", handler.ListGangRecruited)
	}

	port := os.Getenv("ENFL_SERVICE_PORT")
	if port == "" {
		port = ":8119"
	}

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		logger.Info("starting enfl-svc", zap.String("addr", port))
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
