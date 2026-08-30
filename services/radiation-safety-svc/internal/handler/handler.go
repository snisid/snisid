package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/radiation-safety-svc/internal/domain"
	"github.com/snisid/radiation-safety-svc/internal/service"
)

type RadiationHandler struct {
	svc service.RadiationService
}

func NewRadiationHandler(svc service.RadiationService) *RadiationHandler {
	return &RadiationHandler{svc: svc}
}

func (h *RadiationHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1/radiation")
	{
		api.POST("/sources", h.RegisterSource)
		api.PATCH("/sources/:id/status", h.UpdateSourceStatus)
		api.POST("/alerts", h.CreateAlert)
		api.GET("/alerts/unresponded", h.GetUnrespondedAlerts)
		api.POST("/chemicals", h.RegisterChemical)
		api.GET("/chemicals/suspicious", h.GetSuspiciousChemicals)
		api.GET("/dashboard", h.GetDashboard)
	}
}

func (h *RadiationHandler) RegisterSource(c *gin.Context) {
	var s domain.RadioactiveSource
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.RegisterSource(c.Request.Context(), &s); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, s)
}

func (h *RadiationHandler) UpdateSourceStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid source id"})
		return
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.UpdateSourceStatus(c.Request.Context(), id, domain.SourceStatus(body.Status)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *RadiationHandler) CreateAlert(c *gin.Context) {
	var a domain.RadiationAlert
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.CreateAlert(c.Request.Context(), &a); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, a)
}

func (h *RadiationHandler) GetUnrespondedAlerts(c *gin.Context) {
	result, err := h.svc.GetUnrespondedAlerts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *RadiationHandler) RegisterChemical(c *gin.Context) {
	var chem domain.ChemicalPrecursor
	if err := c.ShouldBindJSON(&chem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.RegisterChemical(c.Request.Context(), &chem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, chem)
}

func (h *RadiationHandler) GetSuspiciousChemicals(c *gin.Context) {
	result, err := h.svc.GetSuspiciousChemicals(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *RadiationHandler) GetDashboard(c *gin.Context) {
	stats, err := h.svc.GetDashboard(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
