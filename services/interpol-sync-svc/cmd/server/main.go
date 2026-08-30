package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/snisid/interpol-sync-svc/internal/handler"
	"github.com/snisid/interpol-sync-svc/internal/scheduler"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	dbHost := getEnv("SIVC_DB_HOST", "localhost")
	dbPort := getEnv("SIVC_DB_PORT", "5432")
	dbName := getEnv("SIVC_DB_NAME", "snisid_sivc")
	dbUser := getEnv("SIVC_DB_USER", "sivc_svc")
	dbPassword := getEnv("SIVC_DB_PASSWORD", "")

	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		dbHost, dbPort, dbName, dbUser, dbPassword)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	interpolURL := getEnv("INTERPOL_GATEWAY_URL", "https://i247-gateway.pnh.gov.ht/api")
	interpolKey := getEnv("INTERPOL_API_KEY", "")

	smvHandler := handler.NewSMVHandler(db, interpolURL, interpolKey, logger)
	sadHandler := handler.NewSADHandler(db, interpolURL, interpolKey, logger)

	syncScheduler := scheduler.NewSyncScheduler(smvHandler, sadHandler, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go syncScheduler.Start(ctx)

	logger.Info("INTERPOL sync service started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down INTERPOL sync service...")
	cancel()
	time.Sleep(2 * time.Second)
	fmt.Println("INTERPOL sync service stopped")
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
