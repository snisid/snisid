package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/dpide-svc/internal/domain"
	"github.com/snisid/platform/services/dpide-svc/internal/service"
)

type IDPHandler struct {
	svc *service.IDPService
	log *zap.Logger
}

func NewIDPHandler(svc *service.IDPService, log *zap.Logger) *IDPHandler {
	return &IDPHandler{svc: svc, log: log}
}

func (h *IDPHandler) RegisterIDP(c *gin.Context) {
	var req domain.RegisterIDPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idp, err := h.svc.RegisterIDP(&req)
	if err != nil {
		h.log.Error("register idp failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, idp)
}

func (h *IDPHandler) GetIDP(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	idp, err := h.svc.GetIDP(id)
	if err != nil {
		h.log.Error("get idp failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, idp)
}

func (h *IDPHandler) ListCamps(c *gin.Context) {
	camps, err := h.svc.ListCamps()
	if err != nil {
		h.log.Error("list camps failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, camps)
}

func (h *IDPHandler) GetStats(c *gin.Context) {
	stats, err := h.svc.GetStats()
	if err != nil {
		h.log.Error("get stats failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *IDPHandler) UpdateStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req domain.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateStatus(id, &req); err != nil {
		h.log.Error("update status failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}
