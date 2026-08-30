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
	"github.com/snisid/platform/services/offline-sync-engine/internal/engine"
	"github.com/snisid/platform/services/offline-sync-engine/internal/handlers"
	"github.com/snisid/platform/services/offline-sync-engine/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	dbPath := getEnv("SQLITE_DB_PATH", "data/offline_sync.db")
	port := getEnv("PORT", "8090")

	if err := os.MkdirAll("data", 0750); err != nil {
		log.Fatalf("failed to create data directory: %v", err)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&models.OfflineEvent{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	syncEngine := engine.NewSyncEngine(db)
	if err := syncEngine.ResetStuckEvents(); err != nil {
		log.Printf("warning: failed to reset stuck events: %v", err)
	}

	handler := handlers.NewSyncEngineHandler(syncEngine)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/v1/offline")
	{
		v1.POST("/enqueue", handler.EnqueueHandler)
		v1.POST("/sync", handler.SyncHandler)
		v1.GET("/queue", handler.ListQueueHandler)
		v1.DELETE("/queue/:id", handler.DeleteHandler)
		v1.GET("/status", handler.StatusHandler)
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("Offline Sync Engine starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("offline sync engine failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down offline-sync-engine...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
