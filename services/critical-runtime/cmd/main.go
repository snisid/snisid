package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/internal/platform/logger"
	"github.com/snisid/platform/services/critical-runtime/healer"
	"github.com/snisid/platform/services/critical-runtime/monitor"
	"github.com/snisid/platform/services/critical-runtime/snapshot"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	port := getEnv("PORT", "8088")

	store := &snapshot.SnapshotStore{}
	h := healer.NewHealer("critical-runtime", store)
	checker := &monitor.RuntimeChecker{ID: "critical-runtime"}

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	r.POST("/v1/check", func(c *gin.Context) {
		var state monitor.SystemState
		if err := c.BindJSON(&state); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ok, msg := checker.CheckInvariant(state)
		if !ok {
			checker.OnViolation(msg)
		}
		c.JSON(http.StatusOK, gin.H{"valid": ok, "message": msg})
	})

	r.POST("/v1/heal", func(c *gin.Context) {
		var v healer.Violation
		if err := c.BindJSON(&v); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		plan := h.Heal(v)
		c.JSON(http.StatusOK, plan)
	})

	r.GET("/v1/healings", func(c *gin.Context) {
		c.JSON(http.StatusOK, h.GetActiveHealings())
	})

	r.POST("/v1/snapshot", func(c *gin.Context) {
		var state snapshot.ValidState
		if err := c.BindJSON(&state); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		store.Save(state)
		c.JSON(http.StatusOK, gin.H{"status": "saved"})
	})

	r.GET("/v1/snapshot/latest", func(c *gin.Context) {
		c.JSON(http.StatusOK, store.GetLastValid())
	})

	r.POST("/v1/resume", func(c *gin.Context) {
		h.Resume()
		c.JSON(http.StatusOK, gin.H{"status": "resumed"})
	})

	srv := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		logger.Info(ctx, "critical-runtime starting", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, "critical-runtime http server failed", err)
		}
	}()

	<-ctx.Done()
	logger.Info(ctx, "critical-runtime shutting down")
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

func init() {
	_ = json.RawMessage{}
}
