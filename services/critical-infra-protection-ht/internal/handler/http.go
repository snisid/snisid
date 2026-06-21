package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/critical-infra-protection-ht/internal/domain"
	"github.com/snisid/critical-infra-protection-ht/internal/service"
)

type Handler struct {
	svc *service.InfraProtService
}

func NewHandler(svc *service.InfraProtService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/assets", h.CreateAsset)
	r.GET("/assets/:sector", h.GetAssetsBySector)
	r.POST("/incidents", h.ReportIncident)
	r.GET("/incidents/active", h.GetActiveIncidents)
	r.GET("/incidents/asset/:asset_id", h.GetIncidentsByAsset)
	r.POST("/assessments", h.CreateAssessment)
	r.GET("/dashboard/national", h.GetNationalDashboard)
}

func (h *Handler) CreateAsset(c *gin.Context) {
	var req domain.CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.CreateAsset(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetAssetsBySector(c *gin.Context) {
	sector := c.Param("sector")
	if sector == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sector is required"})
		return
	}
	assets, err := h.svc.GetAssetsBySector(c.Request.Context(), sector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assets)
}

func (h *Handler) ReportIncident(c *gin.Context) {
	var req domain.ReportIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.ReportIncident(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetActiveIncidents(c *gin.Context) {
	incidents, err := h.svc.GetActiveIncidents(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, incidents)
}

func (h *Handler) GetIncidentsByAsset(c *gin.Context) {
	assetID, err := uuid.Parse(c.Param("asset_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset_id"})
		return
	}
	incidents, err := h.svc.GetIncidentsByAsset(c.Request.Context(), assetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, incidents)
}

func (h *Handler) CreateAssessment(c *gin.Context) {
	var req domain.CreateAssessmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.CreateAssessment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetNationalDashboard(c *gin.Context) {
	data, err := h.svc.GetNationalDashboard(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
