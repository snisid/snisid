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
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/snisid/platform/services/rdep-svc/internal/handler"
	"github.com/snisid/platform/services/rdep-svc/internal/kafka"
	"github.com/snisid/platform/services/rdep-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/rdep-svc/internal/service"
)

func main() {
	port := getEnv("RDEP_SERVICE_PORT", "8094")
	dbURL := getEnv("RDEP_DB_URL", "postgresql://root@localhost:26257/snisid_rdep?sslmode=disable")
	kafkaBrokers := getEnv("RDEP_KAFKA_BROKERS", "localhost:9092")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("connected to CockroachDB")

	runMigrations(ctx, pool)

	producer := kafka.NewProducer([]string{kafkaBrokers})
	defer producer.Close()

	deporteeRepo := postgres.NewDeporteeRepo(pool)
	foreignRepo := postgres.NewForeignRecordRepo(pool)
	monitoringEventRepo := postgres.NewMonitoringEventRepo(pool)

	eventPub := &kafkaEventPublisher{producer: producer}

	intakeSvc := service.NewIntakeService(deporteeRepo, foreignRepo, eventPub)
	screeningSvc := service.NewScreeningService(deporteeRepo, foreignRepo, nil, nil, nil, eventPub)
	monitoringSvc := service.NewMonitoringService(deporteeRepo, monitoringEventRepo, eventPub)

	h := handler.NewHandler(intakeSvc, screeningSvc, monitoringSvc)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	v1 := r.Group("/api/v1/rdep")
	h.RegisterRoutes(v1)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Printf("RDEP-HT service starting on port %s", port)
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

func (p *kafkaEventPublisher) Publish(topic string, event interface{}) error {
	return p.producer.Publish(context.Background(), topic, "", event)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func runMigrations(ctx context.Context, pool *pgxpool.Pool) {
	data, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		log.Printf("warning: could not read migration file: %v", err)
		return
	}
	if _, err := pool.Exec(ctx, string(data)); err != nil {
		log.Printf("warning: migration error: %v", err)
	}
}
