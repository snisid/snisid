package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/rvin-svc/internal/domain"
	"github.com/snisid/platform/services/rvin-svc/internal/service"
)

type RemainsHandler struct {
	svc *service.RemainsService
	log *zap.Logger
}

func NewRemainsHandler(svc *service.RemainsService, log *zap.Logger) *RemainsHandler {
	return &RemainsHandler{svc: svc, log: log}
}

func (h *RemainsHandler) RegisterRemains(c *gin.Context) {
	var req domain.RegisterRemainsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	remains, err := h.svc.RegisterRemains(&req)
	if err != nil {
		h.log.Error("register remains failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, remains)
}

func (h *RemainsHandler) GetRemains(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	remains, err := h.svc.GetRemains(id)
	if err != nil {
		h.log.Error("get remains failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, remains)
}

func (h *RemainsHandler) SubmitDNA(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.SubmitDNARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.SubmitDNA(id, &req); err != nil {
		h.log.Error("submit dna failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "dna submitted"})
}

func (h *RemainsHandler) ListUnidentified(c *gin.Context) {
	remains, err := h.svc.ListUnidentified()
	if err != nil {
		h.log.Error("list unidentified failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, remains)
}

func (h *RemainsHandler) GetStatsBySource(c *gin.Context) {
	stats, err := h.svc.GetStatsBySource()
	if err != nil {
		h.log.Error("get stats by source failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, stats)
}
