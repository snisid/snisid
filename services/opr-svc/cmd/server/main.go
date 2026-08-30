package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/opr-svc/internal/api/rest"
	"github.com/snisid/platform/services/opr-svc/internal/domain"
	"github.com/snisid/platform/services/opr-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/opr-svc/internal/service"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dbHost := getEnv("OPR_DB_HOST", "postgres://snisid:snisid@localhost:5432/snisid_opr")
	dbPool, err := pgxpool.New(ctx, dbHost)
	if err != nil {
		logger.Fatal("cannot connect to database", zap.Error(err))
	}
	defer dbPool.Close()

	orderRepo := postgres.NewProtectionOrderRepo(dbPool)
	violRepo := postgres.NewViolationRepo(dbPool)
	var eventPub domain.EventPublisher = &noopEventPublisher{}

	oprSvc := service.NewOPRService(orderRepo, violRepo, eventPub)
	oprHandler := rest.NewOPRHandler(oprSvc)

	r := gin.Default()
	v1 := r.Group("/api/v1/opr")
	{
		v1.POST("/orders", oprHandler.CreateOrder)
		v1.GET("/check/:person_id", oprHandler.CheckSubject)
		v1.POST("/violations", oprHandler.RecordViolation)
		v1.GET("/expiring-soon", oprHandler.GetExpiringSoon)
		v1.GET("/orders/by-gang/:id", oprHandler.GetByGangID)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", getEnv("OPR_SERVICE_PORT", "8096")),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen error", zap.Error(err))
		}
	}()

	log.Println("OPR-HT service started on port", getEnv("OPR_SERVICE_PORT", "8096"))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	srv.Shutdown(shutdownCtx)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

type noopEventPublisher struct{}

func (n *noopEventPublisher) Publish(topic string, event interface{}) error {
	log.Printf("Event published to %s: %v", topic, event)
	return nil
}
