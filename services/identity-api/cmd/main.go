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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/snisid/platform/services/identity-api/internal/handlers"
	"github.com/snisid/platform/services/identity-api/internal/kafka"
	"github.com/snisid/platform/services/identity-api/internal/models"
)

func main() {
	jwtSecret := getEnv("JWT_SECRET", "dev-secret")
	broker := getEnv("KAFKA_BROKER", "localhost:9092")
	dbURL := getEnv("DATABASE_URL", "host=localhost user=snisid password=snisid dbname=snisid port=5432 sslmode=disable")
	port := getEnv("PORT", "8081")

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if getEnv("ENV", "dev") == "dev" {
		if err := db.AutoMigrate(&models.Identity{}, &models.IdentityHistory{}, &models.BiometricReference{}, &models.DocumentAssociation{}); err != nil {
			log.Fatalf("failed to migrate: %v", err)
		}
	}

	topic := getEnv("KAFKA_TOPIC", "snisid.prod.identity.v1.events")
	producer := kafka.NewProducer([]string{broker}, topic)
	defer producer.Close()

	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	h := handlers.New(db, producer)

	api := r.Group("/api/v1")
	api.Use(authMiddleware(jwtSecret))
	h.RegisterRoutes(api)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("identity-service started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run identity service: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down identity-service...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
	}
}

func authMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID := c.GetHeader("X-Actor-ID")
		if actorID == "" {
			actorID = "system"
		}
		c.Set("actor_id", actorID)
		c.Next()
	}
}

func getEnv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
