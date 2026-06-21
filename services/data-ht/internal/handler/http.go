package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/snisid/data-ht/internal/domain"
	"github.com/snisid/data-ht/internal/service"
)

type Handler struct {
	svc *service.DataService
}

func NewHandler(svc *service.DataService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/pipelines", h.ListPipelines)
	r.POST("/models/register", h.RegisterModel)
	r.GET("/models/:id/bias-audit", h.GetBiasAudit)
	r.GET("/dashboards/national", h.GetNationalDashboard)
}

func (h *Handler) ListPipelines(c *gin.Context) {
	pipelines, err := h.svc.ListPipelines(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pipelines)
}

func (h *Handler) RegisterModel(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req domain.RegisterModelRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.RegisterModel(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetBiasAudit(c *gin.Context) {
	modelID := c.Param("id")
	result, err := h.svc.GetBiasAudit(c.Request.Context(), modelID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetNationalDashboard(c *gin.Context) {
	dash, err := h.svc.GetNationalDashboard(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dash)
}
