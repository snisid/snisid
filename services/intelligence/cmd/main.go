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
	"github.com/snisid/platform/services/intelligence/featurestore"
	"github.com/snisid/platform/services/intelligence/ml"
	"github.com/snisid/platform/services/intelligence/oversight"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	port := getEnv("PORT", "8098")
	registry := ml.NewModelRegistry()
	store := featurestore.NewStore(getEnv("REDIS_ADDR", "localhost:6379"))
	ai := &oversight.OversightAI{PlatformID: "intelligence"}

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	r.POST("/v1/models/register", func(c *gin.Context) {
		var req struct {
			Name      string `json:"name"`
			Version   string `json:"version"`
			Algorithm string `json:"algorithm"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		registry.Register(req.Name, req.Version, req.Algorithm)
		c.JSON(http.StatusOK, gin.H{"status": "registered"})
	})

	r.GET("/v1/models/:name", func(c *gin.Context) {
		name := c.Param("name")
		meta, err := registry.Get(name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, meta)
	})

	r.POST("/v1/features", func(c *gin.Context) {
		var req struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := store.SaveFeature(ctx, req.Key, req.Value); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "saved"})
	})

	r.GET("/v1/features/:key", func(c *gin.Context) {
		key := c.Param("key")
		val, err := store.GetFeature(ctx, key)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"key": key, "value": val})
	})

	r.POST("/v1/validate", func(c *gin.Context) {
		var decision oversight.Decision
		if err := c.BindJSON(&decision); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ok, reason := ai.ValidateDecision(decision)
		ai.AuditLog(decision, map[bool]string{true: "APPROVED", false: "REJECTED"}[ok], reason)
		c.JSON(http.StatusOK, gin.H{"approved": ok, "reason": reason})
	})

	srv := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		logger.Info(ctx, "intelligence starting", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, "intelligence http server failed", err)
		}
	}()

	<-ctx.Done()
	logger.Info(ctx, "intelligence shutting down")
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
