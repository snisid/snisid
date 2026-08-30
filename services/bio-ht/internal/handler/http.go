package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/snisid/bio-ht/internal/domain"
	"github.com/snisid/bio-ht/internal/service"
)

type Handler struct {
	svc *service.BioService
}

func NewHandler(svc *service.BioService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/enroll", h.Enroll)
	r.POST("/verify", h.Verify)
	r.POST("/identify", h.Identify)
	r.GET("/quality/:template_id", h.GetQuality)
}

func (h *Handler) Enroll(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req domain.EnrollRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.Enroll(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) Verify(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req domain.VerifyRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.Verify(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) Identify(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req domain.IdentifyRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.Identify(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetQuality(c *gin.Context) {
	templateID := c.Param("template_id")
	quality, err := h.svc.GetQuality(c.Request.Context(), templateID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"template_id": templateID, "quality_score": quality})
}
