package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/sltd-svc/internal/domain"
	"github.com/snisid/platform/services/sltd-svc/internal/service"
)

type SLTDHandler struct {
	svc *service.SLTDService
	log *zap.Logger
}

func NewSLTDHandler(svc *service.SLTDService, log *zap.Logger) *SLTDHandler {
	return &SLTDHandler{svc: svc, log: log}
}

func (h *SLTDHandler) CheckDocument(c *gin.Context) {
	num := c.Param("num")
	if num == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "document number required"})
		return
	}
	issuingCountry := c.DefaultQuery("country", "HTI")
	source := c.DefaultQuery("source", "LOCAL")
	checkedBy := c.DefaultQuery("checked_by", "system")
	result, err := h.svc.CheckDocument(num, issuingCountry, checkedBy, source)
	if err != nil {
		h.log.Error("check document failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *SLTDHandler) ReportLost(c *gin.Context) {
	var req domain.ReportLostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	doc, err := h.svc.ReportLost(&req)
	if err != nil {
		h.log.Error("report lost failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, doc)
}

func (h *SLTDHandler) ReportStolen(c *gin.Context) {
	var req domain.ReportStolenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	doc, err := h.svc.ReportStolen(&req)
	if err != nil {
		h.log.Error("report stolen failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, doc)
}

func (h *SLTDHandler) MarkFound(c *gin.Context) {
	idStr := c.Param("id")
	docID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document id"})
		return
	}
	var req domain.MarkFoundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	doc, err := h.svc.MarkFound(docID, req.FoundLocation, req.ReportedBy)
	if err != nil {
		h.log.Error("mark found failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, doc)
}

func (h *SLTDHandler) GetStats(c *gin.Context) {
	stats, err := h.svc.GetStats()
	if err != nil {
		h.log.Error("get stats failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func SetupRouter(svc *service.SLTDService, log *zap.Logger) *gin.Engine {
	r := gin.Default()
	handler := NewSLTDHandler(svc, log)

	api := r.Group("/api/v1/sltd")
	{
		api.GET("/check/:num", handler.CheckDocument)
		api.POST("/report/lost", handler.ReportLost)
		api.POST("/report/stolen", handler.ReportStolen)
		api.PATCH("/:id/found", handler.MarkFound)
		api.GET("/stats", handler.GetStats)
	}
	return r
}
