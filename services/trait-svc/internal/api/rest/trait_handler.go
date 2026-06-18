package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/trait-svc/internal/domain"
	"github.com/snisid/platform/services/trait-svc/internal/service"
)

type TraiffickingHandler struct {
	svc *service.TraiffickingService
	log *zap.Logger
}

func NewTraiffickingHandler(svc *service.TraiffickingService, log *zap.Logger) *TraiffickingHandler {
	return &TraiffickingHandler{svc: svc, log: log}
}

func (h *TraiffickingHandler) OpenCase(c *gin.Context) {
	var req domain.OpenCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.OpenCase(&req)
	if err != nil {
		h.log.Error("open case failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *TraiffickingHandler) GetCase(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid case id"})
		return
	}

	result, err := h.svc.GetCaseDetail(id)
	if err != nil {
		h.log.Error("get case failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *TraiffickingHandler) AddVictim(c *gin.Context) {
	idStr := c.Param("id")
	caseID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid case id"})
		return
	}

	var req domain.AddVictimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.AddVictim(caseID, &req)
	if err != nil {
		h.log.Error("add victim failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *TraiffickingHandler) GetMinorVictims(c *gin.Context) {
	victims, err := h.svc.GetMinorVictims()
	if err != nil {
		h.log.Error("get minor victims failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, victims)
}

func (h *TraiffickingHandler) DocumentNetwork(c *gin.Context) {
	var req domain.DocumentNetworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.DocumentNetwork(&req)
	if err != nil {
		h.log.Error("document network failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *TraiffickingHandler) GetActiveNetworks(c *gin.Context) {
	networks, err := h.svc.GetActiveNetworks()
	if err != nil {
		h.log.Error("get active networks failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, networks)
}

func (h *TraiffickingHandler) GetStatsByType(c *gin.Context) {
	stats, err := h.svc.GetStatsByType()
	if err != nil {
		h.log.Error("get stats by type failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *TraiffickingHandler) GetMaritimeCases(c *gin.Context) {
	cases, err := h.svc.GetMaritimeCases()
	if err != nil {
		h.log.Error("get maritime cases failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, cases)
}
