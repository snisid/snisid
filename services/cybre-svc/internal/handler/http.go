package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/cybre-svc/internal/domain"
	"github.com/snisid/platform/services/cybre-svc/internal/service"
)

type CybreHandler struct {
	svc *service.CybreService
	log *zap.Logger
}

func NewCybreHandler(svc *service.CybreService, log *zap.Logger) *CybreHandler {
	return &CybreHandler{svc: svc, log: log}
}

func (h *CybreHandler) DeclareIncident(c *gin.Context) {
	var req domain.DeclareIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	incident, err := h.svc.DeclareIncident(&req)
	if err != nil {
		h.log.Error("declare incident failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, incident)
}

func (h *CybreHandler) GetIncident(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	incident, err := h.svc.GetIncident(id)
	if err != nil {
		h.log.Error("get incident failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, incident)
}

func (h *CybreHandler) ListRecentIntrusions(c *gin.Context) {
	intrusions, err := h.svc.ListRecentIntrusions()
	if err != nil {
		h.log.Error("list recent intrusions failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, intrusions)
}

func (h *CybreHandler) AddThreatIntel(c *gin.Context) {
	var req domain.AddThreatIntelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ti, err := h.svc.AddThreatIntel(&req)
	if err != nil {
		h.log.Error("add threat intel failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, ti)
}

func (h *CybreHandler) CheckIndicator(c *gin.Context) {
	indicatorType := c.Query("type")
	value := c.Query("value")
	if indicatorType == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type and value required"})
		return
	}
	ti, err := h.svc.CheckIndicator(indicatorType, value)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"found": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"found": true, "indicator": ti})
}

func (h *CybreHandler) GetStatsByType(c *gin.Context) {
	stats, err := h.svc.GetStatsByType()
	if err != nil {
		h.log.Error("get stats by type failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, stats)
}
