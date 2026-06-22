package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/snisid/sla-svc/internal/service"
)

type Handler struct {
	svc *service.SLAService
}

func NewHandler(svc *service.SLAService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/slas", h.DefineSLA)
	r.GET("/slas", h.ListSLAs)
	r.POST("/slas/:id/slis", h.RecordSLI)
	r.GET("/slas/:id/status", h.GetSLAStatus)
	r.GET("/slas/:id/breaches", h.GetBreaches)
	r.GET("/dashboard", h.GetDashboard)
	r.POST("/slas/:id/escalate", h.TriggerEscalation)
}

func (h *Handler) DefineSLA(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Owner       string `json:"owner"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.DefineSLA(c.Request.Context(), req.Name, req.Description, req.Owner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) ListSLAs(c *gin.Context) {
	result, err := h.svc.ListSLAs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *Handler) RecordSLI(c *gin.Context) {
	slaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sla id"})
		return
	}
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		SLOID string  `json:"slo_id"`
		Name  string  `json:"name"`
		Value float64 `json:"value"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	sloID, err := uuid.Parse(req.SLOID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid slo_id"})
		return
	}

	result, err := h.svc.RecordSLI(c.Request.Context(), slaID, sloID, req.Name, req.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetSLAStatus(c *gin.Context) {
	slaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sla id"})
		return
	}
	result, err := h.svc.GetSLAStatus(c.Request.Context(), slaID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "sla not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetBreaches(c *gin.Context) {
	slaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sla id"})
		return
	}
	result, err := h.svc.GetBreaches(c.Request.Context(), slaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *Handler) GetDashboard(c *gin.Context) {
	result, err := h.svc.GetDashboard(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *Handler) TriggerEscalation(c *gin.Context) {
	slaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sla id"})
		return
	}
	if err := h.svc.TriggerEscalation(c.Request.Context(), slaID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "escalation triggered", "sla_id": slaID.String()})
}
