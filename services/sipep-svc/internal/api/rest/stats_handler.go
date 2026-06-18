package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/sipep-svc/internal/service"
)

type StatsHandler struct {
	overcrowdingSvc *service.OvercrowdingService
}

func NewStatsHandler(overcrowdingSvc *service.OvercrowdingService) *StatsHandler {
	return &StatsHandler{overcrowdingSvc: overcrowdingSvc}
}

func (h *StatsHandler) GetFacilityOccupancy(c *gin.Context) {
	thresholdStr := c.DefaultQuery("threshold", "1.5")
	threshold, err := strconv.ParseFloat(thresholdStr, 64)
	if err != nil {
		threshold = 1.5
	}

	occupancies, err := h.overcrowdingSvc.GetFacilityOccupancy(c.Request.Context(), threshold)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"facilities": occupancies})
}

func (h *StatsHandler) GetOvercrowdingAlerts(c *gin.Context) {
	thresholdStr := c.DefaultQuery("threshold", "1.5")
	threshold, err := strconv.ParseFloat(thresholdStr, 64)
	if err != nil {
		threshold = 1.5
	}

	alerts, err := h.overcrowdingSvc.GetOvercrowdingAlerts(c.Request.Context(), threshold)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"alerts": alerts})
}

func (h *StatsHandler) GetPreventiveDetentionStats(c *gin.Context) {
	stats, err := h.overcrowdingSvc.GetPreventiveDetentionStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
