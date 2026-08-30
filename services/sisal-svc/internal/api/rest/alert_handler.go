package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/sisal-svc/internal/domain"
	"github.com/snisid/platform/services/sisal-svc/internal/service"
)

type AlertHandler struct {
	svc *service.AlertService
	log *zap.Logger
}

func NewAlertHandler(svc *service.AlertService, log *zap.Logger) *AlertHandler {
	return &AlertHandler{svc: svc, log: log}
}

func (h *AlertHandler) IssueAlert(c *gin.Context) {
	var req domain.IssueAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	alert, err := h.svc.IssueAlert(&req)
	if err != nil {
		h.log.Error("issue alert failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, alert)
}

func (h *AlertHandler) ListActiveAlerts(c *gin.Context) {
	alerts, err := h.svc.ListActiveAlerts()
	if err != nil {
		h.log.Error("list active alerts failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, alerts)
}

func (h *AlertHandler) ListHistory(c *gin.Context) {
	alerts, err := h.svc.ListHistory()
	if err != nil {
		h.log.Error("list history failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, alerts)
}

func (h *AlertHandler) CancelAlert(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.CancelAlert(id, req.Reason); err != nil {
		h.log.Error("cancel alert failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "cancelled"})
}

func (h *AlertHandler) Subscribe(c *gin.Context) {
	var req domain.SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.svc.Subscribe(&req)
	if err != nil {
		h.log.Error("subscribe failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, sub)
}
