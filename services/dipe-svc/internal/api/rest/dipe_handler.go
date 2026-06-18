package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/dipe-svc/internal/domain"
	"github.com/snisid/platform/services/dipe-svc/internal/service"
)

type DipeHandler struct {
	svc *service.MissingPersonService
	log *zap.Logger
}

func NewDipeHandler(svc *service.MissingPersonService, log *zap.Logger) *DipeHandler {
	return &DipeHandler{svc: svc, log: log}
}

func (h *DipeHandler) ReportDisappearance(c *gin.Context) {
	var req domain.ReportDisappearanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.ReportDisappearance(&req)
	if err != nil {
		h.log.Error("report disappearance failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *DipeHandler) GetCase(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid case id"})
		return
	}

	result, err := h.svc.GetCaseDetail(id)
	if err != nil {
		h.log.Error("get case failed", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "case not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *DipeHandler) GetOpenCases(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	cases, total, err := h.svc.GetOpenCases(limit, offset)
	if err != nil {
		h.log.Error("get open cases failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cases": cases, "total": total, "limit": limit, "offset": offset})
}

func (h *DipeHandler) AddSighting(c *gin.Context) {
	caseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid case id"})
		return
	}

	var req domain.AddSightingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sighting, err := h.svc.AddSighting(caseID, &req)
	if err != nil {
		h.log.Error("add sighting failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, sighting)
}

func (h *DipeHandler) MatchWithRVIN(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid case id"})
		return
	}

	result, err := h.svc.MatchWithRVIN(id)
	if err != nil {
		h.log.Error("match with rvin failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *DipeHandler) ResolveCase(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid case id"})
		return
	}

	var req domain.ResolveCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.ResolveCase(id, &req); err != nil {
		h.log.Error("resolve case failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "case resolved"})
}

func (h *DipeHandler) GetStatsByType(c *gin.Context) {
	stats, err := h.svc.GetStatsByType()
	if err != nil {
		h.log.Error("get stats by type failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

func (h *DipeHandler) GetHotlineTips(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"tips": []interface{}{}, "message": "hotline tips endpoint"})
}
