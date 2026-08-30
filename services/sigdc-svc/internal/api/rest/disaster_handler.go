package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/sigdc-svc/internal/domain"
	"github.com/snisid/platform/services/sigdc-svc/internal/service"
)

type DisasterHandler struct {
	svc *service.DisasterService
	log *zap.Logger
}

func NewDisasterHandler(svc *service.DisasterService, log *zap.Logger) *DisasterHandler {
	return &DisasterHandler{svc: svc, log: log}
}

func (h *DisasterHandler) DeclareDisaster(c *gin.Context) {
	var req domain.DeclareDisasterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	disaster, err := h.svc.DeclareDisaster(&req)
	if err != nil {
		h.log.Error("declare disaster failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, disaster)
}

func (h *DisasterHandler) ListActiveDisasters(c *gin.Context) {
	disasters, err := h.svc.ListActiveDisasters()
	if err != nil {
		h.log.Error("list active disasters failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, disasters)
}

func (h *DisasterHandler) IssueWarning(c *gin.Context) {
	var req domain.IssueWarningRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.IssueWarning(&req); err != nil {
		h.log.Error("issue warning failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "warning issued"})
}

func (h *DisasterHandler) RegisterVictim(c *gin.Context) {
	var req domain.RegisterVictimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vr, err := h.svc.RegisterVictim(&req)
	if err != nil {
		h.log.Error("register victim failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, vr)
}

func (h *DisasterHandler) ListResources(c *gin.Context) {
	disasterIDStr := c.Query("disaster_id")
	if disasterIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "disaster_id required"})
		return
	}

	disasterID, err := uuid.Parse(disasterIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid disaster_id"})
		return
	}

	resources, err := h.svc.ListResources(disasterID)
	if err != nil {
		h.log.Error("list resources failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, resources)
}
