package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/afis/internal/service"
)

type QualityHandler struct {
	quality *service.QualityService
}

func NewQualityHandler(q *service.QualityService) *QualityHandler {
	return &QualityHandler{quality: q}
}

func (h *QualityHandler) CheckQuality(c *gin.Context) {
	var req struct {
		ImageBase64 string `json:"image_base64" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	score, err := h.quality.CheckQuality(c.Request.Context(), req.ImageBase64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acceptable := h.quality.IsAcceptable(score)

	c.JSON(http.StatusOK, gin.H{
		"nfiq2_score": score,
		"acceptable":  acceptable,
		"min_score":   h.quality.MinScore(),
	})
}
