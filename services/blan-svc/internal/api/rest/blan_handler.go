package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/blan-svc/internal/domain"
	"github.com/snisid/platform/services/blan-svc/internal/service"
)

type BLANHandler struct {
	svc *service.BLANService
	log *zap.Logger
}

func NewBLANHandler(svc *service.BLANService, log *zap.Logger) *BLANHandler {
	return &BLANHandler{svc: svc, log: log}
}

func (h *BLANHandler) OpenCase(c *gin.Context) {
	var req domain.CreateCaseRequest
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

func (h *BLANHandler) GetCaseDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid case id"})
		return
	}

	result, err := h.svc.GetCaseDetail(id)
	if err != nil {
		h.log.Error("get case detail failed", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "case not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *BLANHandler) AddSuspiciousAsset(c *gin.Context) {
	idStr := c.Param("id")
	caseID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid case id"})
		return
	}

	var req domain.AddAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.AddSuspiciousAsset(caseID, &req)
	if err != nil {
		h.log.Error("add suspicious asset failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *BLANHandler) DocumentTransactionChain(c *gin.Context) {
	idStr := c.Param("id")
	caseID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid case id"})
		return
	}

	var req domain.AddChainStepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.DocumentTransactionChain(caseID, &req)
	if err != nil {
		h.log.Error("document transaction chain failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *BLANHandler) GetFlaggedRealEstate(c *gin.Context) {
	result, err := h.svc.GetFlaggedRealEstate()
	if err != nil {
		h.log.Error("get flagged real estate failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *BLANHandler) GetFrozenAssets(c *gin.Context) {
	result, err := h.svc.GetFrozenAssets()
	if err != nil {
		h.log.Error("get frozen assets failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *BLANHandler) GetStatsByTypology(c *gin.Context) {
	result, err := h.svc.GetStatsByTypology()
	if err != nil {
		h.log.Error("get stats by typology failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, result)
}
