package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/sigeo-svc/internal/domain"
	"github.com/snisid/platform/services/sigeo-svc/internal/service"
)

type GeoHandler struct {
	svc *service.GeoIntelService
	log *zap.Logger
}

func NewGeoHandler(svc *service.GeoIntelService, log *zap.Logger) *GeoHandler {
	return &GeoHandler{svc: svc, log: log}
}

func (h *GeoHandler) IngestIncident(c *gin.Context) {
	var req domain.IngestIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	incident, err := h.svc.IngestIncident(&req)
	if err != nil {
		h.log.Error("ingest incident failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, incident)
}

func (h *GeoHandler) ListIncidents(c *gin.Context) {
	deptCode := c.Query("dept_code")
	since := time.Now().Add(-24 * time.Hour)
	incidents, err := h.svc.ListIncidents(deptCode, since)
	if err != nil {
		h.log.Error("list incidents failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, incidents)
}

func (h *GeoHandler) ListCheckpoints(c *gin.Context) {
	checkpoints, err := h.svc.ListCheckpoints()
	if err != nil {
		h.log.Error("list checkpoints failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, checkpoints)
}

func (h *GeoHandler) GetZoneReport(c *gin.Context) {
	deptCode := c.Query("dept_code")
	if deptCode == "" {
		deptCode = "OU"
	}
	report, err := h.svc.GetZoneReport(deptCode, 30*24*time.Hour)
	if err != nil {
		h.log.Error("get zone report failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, report)
}
