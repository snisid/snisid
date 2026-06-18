package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"github.com/snisid/vehicle-criminal-svc/configs"
	"github.com/snisid/vehicle-criminal-svc/internal/api/rest"
	"github.com/snisid/vehicle-criminal-svc/internal/repository/interpol"
	"github.com/snisid/vehicle-criminal-svc/internal/repository/kafka"
	"github.com/snisid/vehicle-criminal-svc/internal/repository/postgres"
	"github.com/snisid/vehicle-criminal-svc/internal/repository/redis"
	"github.com/snisid/vehicle-criminal-svc/internal/service"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg := configs.Load()

	db, err := sqlx.Connect("postgres", cfg.DatabaseDSN())
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})
	defer rdb.Close()

	kafkaWriter := &kafka.Writer{
		Addr:         kafka.TCP(cfg.KafkaBrokers),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
	}
	defer kafkaWriter.Close()

	alertRepo := postgres.NewCriminalAlertRepo(db)
	plateRepo := postgres.NewStolenPlateRepo(db)
	intelRepo := postgres.NewIntelReportRepo(db)
	interpolSyncRepo := postgres.NewInterpolSyncRepo(db)

	hotlistCache := redis.NewHotlistCache(rdb)
	eventPublisher := kafka.NewEventPublisher(kafkaWriter)
	interpolClient := interpol.NewSMVClient(cfg.InterpolGatewayURL, cfg.InterpolAPIKey, cfg.InterpolNCBCode)

	alertSvc := service.NewCriminalAlertService(alertRepo, hotlistCache, eventPublisher, interpolClient, logger)
	plateSvc := service.NewStolenPlateService(plateRepo, hotlistCache, eventPublisher, logger)
	hotlistSvc := service.NewHotlistService(alertRepo, hotlistCache, logger)
	interpolSvc := service.NewInterpolSyncService(interpolSyncRepo, alertRepo, interpolClient, logger)
	intelSvc := service.NewVehicleIntelService(intelRepo, logger)

	if err := hotlistSvc.RefreshHotlist(context.Background()); err != nil {
		logger.Error("Failed to load hotlist", zap.Error(err))
	}

	router := rest.NewRouter(alertSvc, plateSvc, intelSvc, interpolSvc)

	server := &http.Server{
		Addr:         ":" + cfg.ServicePort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("SIVC-HT server starting", zap.String("port", cfg.ServicePort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Server shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	fmt.Println("SIVC-HT server stopped gracefully")
}
