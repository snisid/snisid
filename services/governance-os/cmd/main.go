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
	"github.com/snisid/platform/services/governance-os/optimizer"
	"github.com/snisid/platform/services/governance-os/security-mesh"
	"github.com/snisid/platform/services/governance-os/self-healing"
	"github.com/snisid/platform/services/governance-os/simulator"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	port := getEnv("PORT", "8095")
	meshHealer := security.NewMeshHealer("gos-cluster")
	selfHealer := selfhealing.NewMeshHealer("gos-self-healer", 3)
	simEngine := simulator.NewSimulationEngine("production")
	opt := optimizer.NewPolicyOptimizer("gos-optimizer")

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	r.POST("/v1/mesh/monitor", func(c *gin.Context) {
		var nodes []security.MeshNode
		if err := c.BindJSON(&nodes); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		meshHealer.MonitorMesh(nodes)
		c.JSON(http.StatusOK, gin.H{"status": "monitored"})
	})

	r.GET("/v1/mesh/operations", func(c *gin.Context) {
		c.JSON(http.StatusOK, meshHealer.GetOperationLog())
	})

	r.GET("/v1/mesh/isolated", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"nodes": meshHealer.GetIsolatedNodes()})
	})

	r.POST("/v1/self-heal/monitor", func(c *gin.Context) {
		var nodes []selfhealing.NodeState
		if err := c.BindJSON(&nodes); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		selfHealer.MonitorTrust(nodes)
		c.JSON(http.StatusOK, gin.H{"status": "monitored"})
	})

	r.GET("/v1/self-heal/isolation-log", func(c *gin.Context) {
		c.JSON(http.StatusOK, selfHealer.GetIsolationLog())
	})

	r.POST("/v1/simulate", func(c *gin.Context) {
		var scenario simulator.PolicyScenario
		if err := c.BindJSON(&scenario); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result := simEngine.RunScenario(scenario)
		c.JSON(http.StatusOK, result)
	})

	r.POST("/v1/optimize", func(c *gin.Context) {
		var history map[string]interface{}
		if err := c.BindJSON(&history); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updates := opt.SuggestOptimizations(history)
		c.JSON(http.StatusOK, updates)
	})

	r.GET("/v1/optimize/stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, opt.GetStats())
	})

	srv := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		logger.Info(ctx, "governance-os starting", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, "governance-os http server failed", err)
		}
	}()

	<-ctx.Done()
	logger.Info(ctx, "governance-os shutting down")
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
