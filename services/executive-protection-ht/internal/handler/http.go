package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/snisid/executive-protection-ht/internal/domain"
	"github.com/snisid/executive-protection-ht/internal/service"
)

type Handler struct {
	svc *service.ExecutiveProtectionService
}

func NewHandler(svc *service.ExecutiveProtectionService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/protectees", h.CreateProtectee)
	r.GET("/protectees/active", h.GetActiveProtectees)
	r.POST("/movements", h.CreateMovementPlan)
	r.GET("/movements/upcoming", h.GetUpcomingMovements)
	r.POST("/threats", h.CreateThreatAssessment)
	r.GET("/threats/active/:protectee_id", h.GetActiveThreatsByProtectee)
	r.GET("/dashboard", h.GetDashboard)
}

func (h *Handler) CreateProtectee(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req domain.CreateProtecteeRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.CreateProtectee(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetActiveProtectees(c *gin.Context) {
	protectees, err := h.svc.GetActiveProtectees(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, protectees)
}

func (h *Handler) CreateMovementPlan(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req domain.CreateMovementPlanRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.CreateMovementPlan(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetUpcomingMovements(c *gin.Context) {
	plans, err := h.svc.GetUpcomingMovements(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plans)
}

func (h *Handler) CreateThreatAssessment(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req domain.CreateThreatAssessmentRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.CreateThreatAssessment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetActiveThreatsByProtectee(c *gin.Context) {
	protecteeID := c.Param("protectee_id")
	threats, err := h.svc.GetActiveThreatsByProtectee(c.Request.Context(), protecteeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, threats)
}

func (h *Handler) GetDashboard(c *gin.Context) {
	dash, err := h.svc.GetDashboard(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dash)
}
