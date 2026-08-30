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
	"github.com/snisid/platform/services/graph-service/internal/service"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	port := getEnv("PORT", "8096")
	analyzer := graph.NewThreatAnalyzer()

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	r.POST("/v1/relationship", func(c *gin.Context) {
		var req struct {
			User   graph.UserNode   `json:"user"`
			Action graph.ActionEdge `json:"action"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		analyzer.MapRelationship(req.User, req.Action)
		c.JSON(http.StatusOK, gin.H{"status": "mapped"})
	})

	r.GET("/v1/threat/:uid", func(c *gin.Context) {
		uid := c.Param("uid")
		score := analyzer.DetectInsiderThreat(uid)
		c.JSON(http.StatusOK, gin.H{"uid": uid, "threat_score": score})
	})

	r.POST("/v1/shortest-path", func(c *gin.Context) {
		var query graph.GraphQuery
		if err := c.BindJSON(&query); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result := analyzer.FindShortestPath(query)
		if result == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "path not found"})
			return
		}
		c.JSON(http.StatusOK, result)
	})

	r.POST("/v1/clusters", func(c *gin.Context) {
		var req struct {
			MinSize int `json:"min_size"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		clusters := analyzer.DetectClusters(req.MinSize)
		c.JSON(http.StatusOK, clusters)
	})

	r.GET("/v1/anomalies/:uid", func(c *gin.Context) {
		uid := c.Param("uid")
		anomalies := analyzer.DetectAnomalousAccess(uid)
		c.JSON(http.StatusOK, anomalies)
	})

	srv := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		logger.Info(ctx, "graph-service starting", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, "graph-service http server failed", err)
		}
	}()

	<-ctx.Done()
	logger.Info(ctx, "graph-service shutting down")
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
