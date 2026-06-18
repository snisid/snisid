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
	"github.com/snisid/platform/services/legal-oversight"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	port := getEnv("PORT", "8099")

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	r.POST("/v1/validate", func(c *gin.Context) {
		var decision legaloversight.Decision
		if err := c.BindJSON(&decision); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		approved, reason := legaloversight.ValidateDecision(decision)
		c.JSON(http.StatusOK, gin.H{"approved": approved, "reason": reason})
	})

	srv := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		logger.Info(ctx, "legal-oversight starting", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, "legal-oversight http server failed", err)
		}
	}()

	<-ctx.Done()
	logger.Info(ctx, "legal-oversight shutting down")
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
