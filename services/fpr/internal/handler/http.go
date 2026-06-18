package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/fpr/internal/service"
)

type HTTPHandler struct {
	matcher *service.FPRMatcher
	tmpl    *service.TemplateManager
}

func NewHTTPHandler(m *service.FPRMatcher, t *service.TemplateManager) *HTTPHandler {
	return &HTTPHandler{matcher: m, tmpl: t}
}

func (h *HTTPHandler) RegisterRoutes(rg *gin.Engine) {
	rg.POST("/fpr/enroll", h.Enroll)
	rg.POST("/fpr/verify", h.Verify)
	rg.POST("/fpr/identify", h.Identify)
	rg.GET("/fpr/templates/:id", h.GetTemplate)
}

func (h *HTTPHandler) Enroll(c *gin.Context) {
	var req struct {
		UserID      string `json:"user_id" binding:"required"`
		ImageData   string `json:"image_data" binding:"required"`
		Fingerprint string `json:"fingerprint" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quality, err := h.tmpl.AssessQuality(req.ImageData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quality check failed: " + err.Error()})
		return
	}

	tmpl, err := h.tmpl.Create(req.UserID, req.Fingerprint, quality)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, tmpl)
}

func (h *HTTPHandler) Verify(c *gin.Context) {
	var req struct {
		UserID    string `json:"user_id" binding:"required"`
		ImageData string `json:"image_data" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tmpl, ok := h.tmpl.Get(req.UserID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found for user"})
		return
	}

	score := h.matcher.Verify(req.ImageData, tmpl.Data)
	matched := score >= h.matcher.Threshold
	c.JSON(http.StatusOK, gin.H{
		"matched":   matched,
		"score":     score,
		"threshold": h.matcher.Threshold,
	})
}

func (h *HTTPHandler) Identify(c *gin.Context) {
	var req struct {
		ImageData string `json:"image_data" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	templates := h.tmpl.List()
	results := h.matcher.Identify(req.ImageData, templates)

	c.JSON(http.StatusOK, gin.H{"results": results})
}

func (h *HTTPHandler) GetTemplate(c *gin.Context) {
	id := c.Param("id")
	tmpl, ok := h.tmpl.Get(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}
	c.JSON(http.StatusOK, tmpl)
}
