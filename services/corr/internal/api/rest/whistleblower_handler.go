package rest

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/corr/internal/domain"
	"github.com/snisid/platform/services/corr/internal/service"
)

type WhistleblowerHandler struct {
	svc *service.WhistleblowerService
}

func NewWhistleblowerHandler(svc *service.WhistleblowerService) *WhistleblowerHandler {
	return &WhistleblowerHandler{svc: svc}
}

func (h *WhistleblowerHandler) Submit(c *gin.Context) {
	var req domain.CreateWhistleblowerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "requête invalide"})
		return
	}
	ip := c.ClientIP()
	ipHash := fmt.Sprintf("%x", sha256.Sum256([]byte(ip)))
	report, err := h.svc.SubmitReport(c.Request.Context(), req, ipHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"report_id":    report.ReportID,
		"report_token": report.ReportToken,
		"message":      "Signalement soumis anonymement. Conservez votre token pour suivi.",
	})
}

func (h *WhistleblowerHandler) GetByToken(c *gin.Context) {
	token := strings.TrimSpace(c.Param("token"))
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token requis"})
		return
	}
	report, err := h.svc.GetReportByToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "signalement introuvable"})
		return
	}
	c.JSON(http.StatusOK, report)
}
