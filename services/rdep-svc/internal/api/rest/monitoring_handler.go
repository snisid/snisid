package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/rdep-svc/internal/service"
)

type MonitoringHandler struct {
	monitoringSvc *service.MonitoringService
}

func NewMonitoringHandler(monitoringSvc *service.MonitoringService) *MonitoringHandler {
	return &MonitoringHandler{monitoringSvc: monitoringSvc}
}

func (h *MonitoringHandler) RecordEvent(c *gin.Context) {
	var req service.MonitoringEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := h.monitoringSvc.RecordEvent(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

func (h *MonitoringHandler) GetEvents(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID invalide"})
		return
	}

	events, err := h.monitoringSvc.GetEvents(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

func (h *MonitoringHandler) UpdateAddress(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID invalide"})
		return
	}

	var req struct {
		Address  string `json:"address" binding:"required"`
		Commune  string `json:"commune" binding:"required"`
		DeptCode string `json:"dept_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.monitoringSvc.UpdateAddress(c.Request.Context(), id, req.Address, req.Commune, req.DeptCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Adresse mise à jour"})
}
