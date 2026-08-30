package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/sipci-svc/internal/domain"
	"github.com/snisid/platform/services/sipci-svc/internal/service"
)

type InfraHandler struct {
	svc *service.InfraService
	log *zap.Logger
}

func NewInfraHandler(svc *service.InfraService, log *zap.Logger) *InfraHandler {
	return &InfraHandler{svc: svc, log: log}
}

func (h *InfraHandler) RegisterAsset(c *gin.Context) {
	var req domain.RegisterAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	asset, err := h.svc.RegisterAsset(&req)
	if err != nil {
		h.log.Error("register asset failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, asset)
}

func (h *InfraHandler) GetAsset(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	asset, err := h.svc.GetAsset(id)
	if err != nil {
		h.log.Error("get asset failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, asset)
}

func (h *InfraHandler) ListAssets(c *gin.Context) {
	assets, err := h.svc.ListAssets()
	if err != nil {
		h.log.Error("list assets failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, assets)
}

func (h *InfraHandler) ListCritical(c *gin.Context) {
	assets, err := h.svc.ListCritical()
	if err != nil {
		h.log.Error("list critical failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, assets)
}

func (h *InfraHandler) ListUnderThreat(c *gin.Context) {
	assets, err := h.svc.ListUnderThreat()
	if err != nil {
		h.log.Error("list under threat failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, assets)
}

func (h *InfraHandler) ReportIncident(c *gin.Context) {
	var req domain.ReportIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	incident, err := h.svc.ReportIncident(&req)
	if err != nil {
		h.log.Error("report incident failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, incident)
}

func (h *InfraHandler) ListRecentIncidents(c *gin.Context) {
	incidents, err := h.svc.ListRecentIncidents()
	if err != nil {
		h.log.Error("list recent incidents failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, incidents)
}

func (h *InfraHandler) AssessRisk(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	assessment, err := h.svc.AssessRisk(id)
	if err != nil {
		h.log.Error("assess risk failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, assessment)
}
