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
	"github.com/snisid/platform/services/governance-engine/digital-twin"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	port := getEnv("PORT", "8091")

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	r.POST("/v1/step", func(c *gin.Context) {
		var req struct {
			State  digitaltwin.WorldState `json:"state"`
			Action string                 `json:"action"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result := digitaltwin.Step(req.State, req.Action)
		c.JSON(http.StatusOK, result)
	})

	r.POST("/v1/forecast", func(c *gin.Context) {
		var req struct {
			Initial digitaltwin.WorldState `json:"initial"`
			Actions []string               `json:"actions"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		states := digitaltwin.RunImpactForecast(req.Initial, req.Actions)
		c.JSON(http.StatusOK, states)
	})

	srv := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		logger.Info(ctx, "governance-engine starting", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, "governance-engine http server failed", err)
		}
	}()

	<-ctx.Done()
	logger.Info(ctx, "governance-engine shutting down")
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
