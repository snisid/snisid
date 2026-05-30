package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/backend/internal/domain/identity/entity"
	"github.com/snisid/platform/backend/internal/domain/identity/handler"
	"github.com/snisid/platform/backend/internal/domain/identity/repository"
	"github.com/snisid/platform/backend/internal/domain/identity/usecase"
	"github.com/snisid/platform/backend/internal/platform/events"
	"github.com/snisid/platform/backend/internal/platform/logger"
	"github.com/snisid/platform/backend/internal/platform/middleware"
	"github.com/snisid/platform/backend/internal/platform/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	
	// Swagger docs
	// _ "github.com/snisid/platform/backend/api/openapi/docs"
	// swaggerFiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
)

// @title SNISID Identity API
// @version 1.0
// @description The Core Identity Management Service for SNISID
// @host localhost:8081
// @BasePath /v1
func main() {
	jwtSecret := getEnv("JWT_SECRET", "dev-secret")
	broker := getEnv("KAFKA_BROKER", "localhost:9092")
	dbURL := getEnv("DATABASE_URL", "host=localhost user=snisid password=snisid dbname=snisid port=5432 sslmode=disable")
	port := getEnv("PORT", "8081")

	// 1. Initialize DB
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect database", err)
	}
	
	// Ensure auto-migration is off in prod, relying on raw SQL migrations
	if getEnv("ENV", "dev") == "dev" {
		db.AutoMigrate(&entity.Identity{}, &entity.IdentityHistory{}, &entity.BiometricReference{}, &entity.DocumentAssociation{})
	}

	// 2. Initialize Repositories
	repo := repository.NewPostgresRepository(db)

	// 3. Initialize Kafka Producer
	producer := events.NewProducer([]string{broker}, "identity.events")
	defer producer.Close()

	// 4. Initialize UseCases
	svc := usecase.NewIdentityService(repo, producer)

	// 5. Initialize Handlers
	httpHandler := handler.NewHttpHandler(svc)

	// Setup Router
	r := gin.Default()
	r.Use(middleware.RateLimit(30, 60))

	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	// Swagger Endpoint (uncomment when swag init is run)
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Protected API Group
	api := r.Group("/v1", middleware.Auth(jwtSecret))
	httpHandler.RegisterRoutes(api)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("failed to run identity service", err)
		}
	}()

	logger.Info("identity-service started on port "+port, nil)

	// Graceful Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("shutting down identity-service...", nil)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("gateway forced to shutdown", err)
	}
}

func getEnv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
