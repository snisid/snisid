package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type HealthHandler struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewHealthHandler(pool *pgxpool.Pool, logger *zap.Logger) *HealthHandler {
	return &HealthHandler{pool: pool, logger: logger}
}

func (h *HealthHandler) Healthz(c *gin.Context) {
	if err := h.pool.Ping(context.Background()); err != nil {
		h.logger.Error("health check failed", zap.Error(err))
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func MetricsHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
