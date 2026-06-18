package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/corr/internal/service"
)

type AlertHandler struct {
	svc *service.AlertService
}

func NewAlertHandler(svc *service.AlertService) *AlertHandler {
	return &AlertHandler{svc: svc}
}

func (h *AlertHandler) ListBehavioral(c *gin.Context) {
	alerts, err := h.svc.ListAlerts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur requête"})
		return
	}
	c.JSON(http.StatusOK, alerts)
}

func (h *AlertHandler) ListRiskScores(c *gin.Context) {
	scores, err := h.svc.ListRiskScores(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur requête"})
		return
	}
	c.JSON(http.StatusOK, scores)
}
