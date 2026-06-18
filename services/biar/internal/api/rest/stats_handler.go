package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/biar/internal/service"
)

type StatsHandler struct {
	svc *service.StatsService
}

func NewStatsHandler(svc *service.StatsService) *StatsHandler {
	return &StatsHandler{svc: svc}
}

func (h *StatsHandler) ByGang(c *gin.Context) {
	stats, err := h.svc.ByGang(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur requête"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *StatsHandler) ByOrigin(c *gin.Context) {
	stats, err := h.svc.ByOrigin(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur requête"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *StatsHandler) Routes(c *gin.Context) {
	routes, err := h.svc.Routes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur requête"})
		return
	}
	c.JSON(http.StatusOK, routes)
}
