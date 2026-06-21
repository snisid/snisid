package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/snisid/platform/services/chef-svc/internal/handler"
	"github.com/snisid/platform/services/chef-svc/internal/kafka"
	"github.com/snisid/platform/services/chef-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/chef-svc/internal/service"
)

func main() {
	port := getEnv("CHEF_SERVICE_PORT", "8097")
	dbURL := getEnv("CHEF_DB_URL", "postgresql://root@localhost:26257/snisid_chef?sslmode=disable")
	kafkaBrokers := getEnv("CHEF_KAFKA_BROKERS", "localhost:9092")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("connected to CockroachDB")

	runMigrations(db)

	producer := kafka.NewProducer([]string{kafkaBrokers})
	defer producer.Close()

	memberRepo := postgres.NewMemberRepo(db)
	intelRepo := postgres.NewIntelNoteRepo(db)
	sightRepo := postgres.NewSightingRepo(db)

	publisher := &kafkaEventPublisher{producer: producer}

	svc := service.NewMemberService(memberRepo, intelRepo, sightRepo, publisher)

	h := handler.NewHandler(svc)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	v1 := r.Group("/api/v1/chef")
	h.RegisterRoutes(v1)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Printf("CHEF-HT service starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server exited")
}

type kafkaEventPublisher struct {
	producer *kafka.Producer
}

func (p *kafkaEventPublisher) PublishEvent(eventType string, payload interface{}) error {
	return p.producer.Publish(context.Background(), eventType, "", payload)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func runMigrations(db *sql.DB) {
	data, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		log.Printf("warning: could not read migration file: %v", err)
		return
	}
	if _, err := db.Exec(string(data)); err != nil {
		log.Printf("warning: migration error: %v", err)
	}
}
