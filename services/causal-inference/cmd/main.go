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
	"github.com/snisid/platform/services/causal-inference"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	port := getEnv("PORT", "8086")
	engine := &causalinference.CausalEngine{}

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	r.POST("/v1/estimate-effect", func(c *gin.Context) {
		var req struct {
			Feature string  `json:"feature"`
			Delta   float64 `json:"delta"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		effect := engine.EstimateEffect(req.Feature, req.Delta)
		c.JSON(http.StatusOK, gin.H{"effect": effect})
	})

	r.POST("/v1/recommend-intervention", func(c *gin.Context) {
		var req struct {
			SubjectID string `json:"subject_id"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		action := engine.RecommendIntervention(req.SubjectID)
		c.JSON(http.StatusOK, gin.H{"action": action})
	})

	srv := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		logger.Info(ctx, "causal-inference starting", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, "causal-inference http server failed", err)
		}
	}()

	<-ctx.Done()
	logger.Info(ctx, "causal-inference shutting down")
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
