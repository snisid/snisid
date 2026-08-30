package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/corr-svc/internal/domain"
	"github.com/snisid/platform/services/corr-svc/internal/service"
)

type IntegrityHandler struct {
	svc *service.IntegrityService
	log *zap.Logger
}

func NewIntegrityHandler(svc *service.IntegrityService, log *zap.Logger) *IntegrityHandler {
	return &IntegrityHandler{svc: svc, log: log}
}

func (h *IntegrityHandler) OpenCase(c *gin.Context) {
	var req domain.OpenCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	integrityCase, err := h.svc.OpenCase(&req)
	if err != nil {
		h.log.Error("open case failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, integrityCase)
}

func (h *IntegrityHandler) GetCase(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	integrityCase, err := h.svc.GetCase(id)
	if err != nil {
		h.log.Error("get case failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, integrityCase)
}

func (h *IntegrityHandler) ListActiveCases(c *gin.Context) {
	cases, err := h.svc.ListActiveCases()
	if err != nil {
		h.log.Error("list active cases failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, cases)
}

func (h *IntegrityHandler) SubmitWhistleblower(c *gin.Context) {
	var req domain.SubmitWhistleblowerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	wr, err := h.svc.SubmitWhistleblower(&req)
	if err != nil {
		h.log.Error("submit whistleblower failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"report_token": wr.ReportToken})
}

func (h *IntegrityHandler) TrackWhistleblower(c *gin.Context) {
	token := c.Param("token")
	wr, err := h.svc.TrackWhistleblower(token)
	if err != nil {
		h.log.Error("track whistleblower failed", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, wr)
}

func (h *IntegrityHandler) ListBehavioralAlerts(c *gin.Context) {
	alerts, err := h.svc.ListBehavioralAlerts()
	if err != nil {
		h.log.Error("list behavioral alerts failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, alerts)
}

func (h *IntegrityHandler) SubmitDeclaration(c *gin.Context) {
	var req domain.SubmitDeclarationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ad, err := h.svc.SubmitDeclaration(&req)
	if err != nil {
		h.log.Error("submit declaration failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, ad)
}

func (h *IntegrityHandler) ListFlaggedDeclarations(c *gin.Context) {
	declarations, err := h.svc.ListFlaggedDeclarations()
	if err != nil {
		h.log.Error("list flagged declarations failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, declarations)
}
