package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/snisid/bio-surveillance-ht/internal/domain"
	"github.com/snisid/bio-surveillance-ht/internal/service"
)

type Handler struct {
	svc *service.BioSurveillanceService
}

func NewHandler(svc *service.BioSurveillanceService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/alerts", h.CreateAlert)
	r.GET("/alerts/active", h.GetActiveAlerts)
	r.GET("/alerts/region/:region", h.GetAlertsByRegion)
	r.POST("/campaigns", h.CreateCampaign)
	r.GET("/campaigns/:id/coverage", h.GetCampaignCoverage)
	r.PATCH("/facilities/:id/stock", h.UpdateFacilityStock)
	r.GET("/dashboard/national", h.GetDashboardNational)
}

func (h *Handler) CreateAlert(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req domain.CreateDiseaseAlertRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.CreateAlert(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetActiveAlerts(c *gin.Context) {
	alerts, err := h.svc.GetActiveAlerts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, alerts)
}

func (h *Handler) GetAlertsByRegion(c *gin.Context) {
	region := c.Param("region")
	alerts, err := h.svc.GetAlertsByRegion(c.Request.Context(), region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, alerts)
}

func (h *Handler) CreateCampaign(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req domain.CreateVaccinationCampaignRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.CreateCampaign(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetCampaignCoverage(c *gin.Context) {
	id := c.Param("id")
	result, err := h.svc.GetCampaignCoverage(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) UpdateFacilityStock(c *gin.Context) {
	id := c.Param("id")
	body, _ := io.ReadAll(c.Request.Body)
	var req domain.UpdateFacilityStockRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.UpdateFacilityStock(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetDashboardNational(c *gin.Context) {
	dash, err := h.svc.GetDashboardNational(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dash)
}
