package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/afis-svc/internal/service"
)

type QualityHandler struct {
	svc *service.QualityService
}

func NewQualityHandler(svc *service.QualityService) *QualityHandler {
	return &QualityHandler{svc: svc}
}

func (h *QualityHandler) CheckQuality(c *gin.Context) {
	var req struct {
		Score int16 `json:"score" binding:"required,min=0,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.ValidateScore(req.Score); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"valid": false,
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":        true,
		"score":        req.Score,
		"high_quality": h.svc.IsHighQuality(req.Score),
	})
}

func (h *QualityHandler) Stats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"pending_checks": h.svc.PendingChecksCount(),
	})
}
