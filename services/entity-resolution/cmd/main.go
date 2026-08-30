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
	"github.com/snisid/platform/services/entity-resolution/internal/handlers"
	"github.com/snisid/platform/services/entity-resolution/internal/matching"
	"github.com/snisid/platform/services/entity-resolution/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	dsn := getEnv("DATABASE_URL", "host=localhost port=5432 user=snisid password=snisid dbname=snisid sslmode=disable")
	port := getEnv("PORT", "8095")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&models.Identity{}, &models.ResolvedIdentity{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	lsh := matching.NewLSHIndex(100, 20, 5)

	var identities []models.Identity
	db.Find(&identities)
	records := make([]matching.IdentityStoreRecord, len(identities))
	for i := range identities {
		records[i] = identities[i]
	}
	indexable := matching.IdentitiesToIndexable(records)
	lsh.Build(indexable)
	log.Printf("LSH index built with %d identities", len(indexable))

	matchEngine := matching.NewCompositeEngine(db, lsh)
	handler := handlers.NewResolutionHandler(matchEngine, db)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/v1/resolution")
	{
		v1.POST("/match", handler.MatchHandler)
		v1.POST("/reconcile", handler.ReconcileHandler)
		v1.GET("/candidates/:id", handler.CandidatesHandler)
		v1.POST("/merge", handler.MergeHandler)
		v1.POST("/split", handler.SplitHandler)
		v1.GET("/stats", handler.StatsHandler)
	}

	srv := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		log.Printf("Entity Resolution Engine starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down entity-resolution...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
