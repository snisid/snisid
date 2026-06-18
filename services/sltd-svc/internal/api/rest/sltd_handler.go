package rest

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
	docNum := c.Param("num")
	issuingCountry := c.DefaultQuery("country", "HTI")
	checkedBy := c.GetHeader("X-User-ID")
	if checkedBy == "" {
		checkedBy = "system"
	}
	source := c.DefaultQuery("source", "API")

	result, err := h.svc.CheckDocument(docNum, issuingCountry, checkedBy, source)
	if err != nil {
		h.log.Error("failed to check document", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check document"})
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
		h.log.Error("failed to report lost document", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to report lost document"})
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
		h.log.Error("failed to report stolen document", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to report stolen document"})
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
		h.log.Error("failed to mark document found", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark document found"})
		return
	}

	c.JSON(http.StatusOK, doc)
}

func (h *SLTDHandler) GetStats(c *gin.Context) {
	stats, err := h.svc.GetStats()
	if err != nil {
		h.log.Error("failed to get stats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
