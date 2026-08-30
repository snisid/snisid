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
	"github.com/snisid/platform/services/governance-mesh/execution"
	"github.com/snisid/platform/services/governance-mesh/federation"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	port := getEnv("PORT", "8094")
	hub := federation.NewFederationHub("governance-hub", []string{"peer-ht", "peer-do", "peer-fr"})

	r := gin.Default()
	r.Use(execution.IdentityBinding())
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	r.POST("/v1/policies/distribute", func(c *gin.Context) {
		var pkg federation.PolicyPackage
		if err := c.BindJSON(&pkg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		results := hub.DistributePolicy(pkg)
		c.JSON(http.StatusOK, results)
	})

	r.POST("/v1/policies/acknowledge", func(c *gin.Context) {
		var req struct {
			PackageID string `json:"package_id"`
			Peer      string `json:"peer"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		hub.AcknowledgePolicy(req.PackageID, req.Peer)
		c.JSON(http.StatusOK, gin.H{"status": "acknowledged"})
	})

	r.POST("/v1/policies/revoke", func(c *gin.Context) {
		var req struct {
			PackageID string `json:"package_id"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := hub.RevokePolicy(req.PackageID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "revoked"})
	})

	r.GET("/v1/policies/status/:id", func(c *gin.Context) {
		id := c.Param("id")
		status := hub.GetDistributionStatus(id)
		c.JSON(http.StatusOK, gin.H{"distribution": status})
	})

	r.POST("/v1/peers", func(c *gin.Context) {
		var req struct {
			Peer string `json:"peer"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		hub.AddPeer(req.Peer)
		c.JSON(http.StatusOK, gin.H{"status": "added"})
	})

	r.DELETE("/v1/peers/:peer", func(c *gin.Context) {
		peer := c.Param("peer")
		hub.RemovePeer(peer)
		c.JSON(http.StatusOK, gin.H{"status": "removed"})
	})

	srv := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		logger.Info(ctx, "governance-mesh starting", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, "governance-mesh http server failed", err)
		}
	}()

	<-ctx.Done()
	logger.Info(ctx, "governance-mesh shutting down")
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
