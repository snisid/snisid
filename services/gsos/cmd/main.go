package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/internal/platform/logger"
	"github.com/snisid/platform/services/gsos/ai"
	"github.com/snisid/platform/services/gsos/router"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	port := getEnv("PORT", "8097")
	fedRouter := router.NewFederationRouter([]string{"ht", "do", "fr", "gp", "mq", "gf"})
	correlationLayer := ai.NewCorrelationLayer("v2.1")

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	r.POST("/v1/route", func(c *gin.Context) {
		var event router.GSPEvent
		if err := c.BindJSON(&event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fedRouter.RouteEvent(event)
		c.JSON(http.StatusOK, gin.H{"status": "routed"})
	})

	r.POST("/v1/correlate", func(c *gin.Context) {
		var events []ai.SecurityEvent
		if err := c.BindJSON(&events); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		correlations := correlationLayer.AnalyzeGlobalThreats(events)
		c.JSON(http.StatusOK, correlations)
	})

	r.POST("/v1/ingest", func(c *gin.Context) {
		var event ai.SecurityEvent
		if err := c.BindJSON(&event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		correlationLayer.IngestEvent(event)
		c.JSON(http.StatusOK, gin.H{"status": "ingested"})
	})

	r.GET("/v1/correlation/stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, correlationLayer.GetStats())
	})

	r.GET("/v1/alerts", func(c *gin.Context) {
		alerts := fedRouter.GetAlertHistory(0)
		c.JSON(http.StatusOK, alerts)
	})

	srv := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		logger.Info(ctx, "gsos starting", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, "gsos http server failed", err)
		}
	}()

	<-ctx.Done()
	logger.Info(ctx, "gsos shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
