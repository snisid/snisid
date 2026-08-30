package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/snisid/radiation-safety-svc/internal/handler"
	"github.com/snisid/radiation-safety-svc/internal/kafka"
	"github.com/snisid/radiation-safety-svc/internal/repository"
	"github.com/snisid/radiation-safety-svc/internal/service"
)

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	port := getEnv("RADIATION_PORT", "8308")
	dbDSN := getEnv("RADIATION_DB_DSN", "postgres://localhost:5432/snisid_radiation?sslmode=disable")
	kafkaBrokers := strings.Split(getEnv("RADIATION_KAFKA_BROKERS", "localhost:9092"), ",")

	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.PingContext(context.Background()); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	prod := kafka.NewProducer(kafkaBrokers, "snisid.radiation.events")
	repo := repository.NewRadiationRepository(db)
	svc := service.NewRadiationService(repo, prod)
	h := handler.NewRadiationHandler(svc)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h.RegisterRoutes(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	_ = db.Close()
}
