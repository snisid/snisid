package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/extors-svc/internal/domain"
	"github.com/snisid/platform/services/extors-svc/internal/service"
)

type ExtorsHandler struct {
	svc *service.ExtorsService
	log *zap.Logger
}

func NewExtorsHandler(svc *service.ExtorsService, log *zap.Logger) *ExtorsHandler {
	return &ExtorsHandler{svc: svc, log: log}
}

func (h *ExtorsHandler) OpenCase(c *gin.Context) {
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

func (h *ExtorsHandler) GetCaseDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid case id"})
		return
	}

	result, err := h.svc.GetCaseDetail(id)
	if err != nil {
		h.log.Error("get case detail failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ExtorsHandler) AddNegotiation(c *gin.Context) {
	caseIDStr := c.Param("id")
	caseID, err := uuid.Parse(caseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid case id"})
		return
	}

	var req domain.AddNegotiationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.AddNegotiation(caseID, &req)
	if err != nil {
		h.log.Error("add negotiation failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *ExtorsHandler) CreateTollPoint(c *gin.Context) {
	var req domain.CreateTollPointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.DocumentTollPoint(&req)
	if err != nil {
		h.log.Error("create toll point failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *ExtorsHandler) GetTollsMap(c *gin.Context) {
	geojson, err := h.svc.GetTollsMapGeoJSON()
	if err != nil {
		h.log.Error("get tolls map failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, geojson)
}

func (h *ExtorsHandler) GetGangRevenue(c *gin.Context) {
	gangIDStr := c.Param("id")
	gangID, err := uuid.Parse(gangIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid gang id"})
		return
	}

	gangName := c.Query("name")

	report, err := h.svc.ComputeGangRevenue(gangID, gangName)
	if err != nil {
		h.log.Error("get gang revenue failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, report)
}

func (h *ExtorsHandler) GetStatsByType(c *gin.Context) {
	stats, err := h.svc.GetStatsByType()
	if err != nil {
		h.log.Error("get stats by type failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *ExtorsHandler) GetMoncashPatterns(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "moncash patterns endpoint"})
}
