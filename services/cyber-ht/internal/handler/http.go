package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snisid/cyber-ht/internal/domain"
	"github.com/snisid/cyber-ht/internal/service"
)

type Handler struct {
	svc *service.CyberService
}

func NewHandler(svc *service.CyberService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/incidents", h.CreateIncident)
	r.GET("/incidents/active", h.GetActiveIncidents)
	r.POST("/policies", h.CreatePolicy)
	r.GET("/threat-intel/check", h.CheckThreatIndicator)
}

func (h *Handler) CreateIncident(c *gin.Context) {
	var req domain.CreateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.CreateIncident(c.Request.Context(), req)
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

func (h *Handler) CreatePolicy(c *gin.Context) {
	var req domain.CreatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.CreatePolicy(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) CheckThreatIndicator(c *gin.Context) {
	indicator := c.Query("indicator")
	if indicator == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "indicator query parameter is required"})
		return
	}

	ti, err := h.svc.CheckThreatIndicator(c.Request.Context(), indicator)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "threat indicator not found"})
		return
	}
	c.JSON(http.StatusOK, ti)
}
