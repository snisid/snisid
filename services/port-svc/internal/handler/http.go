package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/port-svc/internal/domain"
	"github.com/snisid/platform/services/port-svc/internal/service"
)

type PortHandler struct {
	svc *service.PortService
	log *zap.Logger
}

func NewPortHandler(svc *service.PortService, log *zap.Logger) *PortHandler {
	return &PortHandler{svc: svc, log: log}
}

func (h *PortHandler) RecordArrival(c *gin.Context) {
	var req domain.RecordArrivalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	arrival, err := h.svc.RecordArrival(&req)
	if err != nil {
		h.log.Error("record arrival failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, arrival)
}

func (h *PortHandler) GetArrival(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid arrival id"})
		return
	}
	arrival, err := h.svc.GetArrival(id)
	if err != nil {
		h.log.Error("get arrival failed", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "arrival not found"})
		return
	}
	c.JSON(http.StatusOK, arrival)
}

func (h *PortHandler) GetHighRiskContainers(c *gin.Context) {
	containers, err := h.svc.GetHighRiskContainers()
	if err != nil {
		h.log.Error("get high risk containers failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, containers)
}

func (h *PortHandler) ScanContainer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid container id"})
		return
	}
	var req domain.ScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	container, err := h.svc.ScanContainer(id, req.ScanResult)
	if err != nil {
		h.log.Error("scan container failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, container)
}

func (h *PortHandler) SeizeContainer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid container id"})
		return
	}
	var req domain.SeizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	container, err := h.svc.SeizeContainer(id, req.SeizureDescription, req.CaseReference)
	if err != nil {
		h.log.Error("seize container failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, container)
}

func (h *PortHandler) GetSeizureStats(c *gin.Context) {
	stats, err := h.svc.GetSeizureStats()
	if err != nil {
		h.log.Error("get seizure stats failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func SetupRouter(svc *service.PortService, log *zap.Logger) *gin.Engine {
	r := gin.Default()
	handler := NewPortHandler(svc, log)

	api := r.Group("/api/v1/port")
	{
		api.POST("/arrivals", handler.RecordArrival)
		api.GET("/arrivals/:id", handler.GetArrival)
		api.GET("/containers/high-risk", handler.GetHighRiskContainers)
		api.POST("/containers/:id/scan", handler.ScanContainer)
		api.POST("/containers/:id/seize", handler.SeizeContainer)
		api.GET("/stats/seizures", handler.GetSeizureStats)
	}
	return r
}
