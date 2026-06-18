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
	"github.com/snisid/platform/services/ai-scheduler"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	port := getEnv("PORT", "8082")
	sched := &scheduler.AIScheduler{}

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	r.POST("/v1/optimize", func(c *gin.Context) {
		var req struct {
			Pods  []scheduler.Pod  `json:"pods"`
			Nodes []scheduler.Node `json:"nodes"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		placement := sched.OptimizePlacement(req.Pods, req.Nodes)
		c.JSON(http.StatusOK, gin.H{"placement": placement})
	})

	srv := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		logger.Info(ctx, "ai-scheduler starting", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, "ai-scheduler http server failed", err)
		}
	}()

	<-ctx.Done()
	logger.Info(ctx, "ai-scheduler shutting down")
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
