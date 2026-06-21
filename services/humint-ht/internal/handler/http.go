package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snisid/humint-ht/internal/domain"
	"github.com/snisid/humint-ht/internal/service"
)

type HumintService interface {
	CreateSource(req domain.CreateSourceRequest) (domain.Source, error)
	UpdateCredibility(code string, req domain.UpdateCredibilityRequest) (domain.Source, error)
	GetReportsBySource(code string) ([]domain.IntelligenceReport, error)
	SubmitReport(req domain.SubmitReportRequest) (domain.IntelligenceReport, error)
	LogDebriefing(req domain.LogDebriefingRequest) (domain.DebriefingSession, error)
	GetHighRiskSources() ([]domain.Source, error)
	GetSourceNetwork() (domain.SourceNetworkResponse, error)
}

type HumintHandler struct {
	svc HumintService
}

func NewHumintHandler(svc *service.HumintService) *HumintHandler {
	return &HumintHandler{svc: svc}
}

func (h *HumintHandler) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1/humint")
	{
		v1.POST("/sources", h.CreateSource)
		v1.PATCH("/sources/:code/credibility", h.UpdateCredibility)
		v1.GET("/sources/:code/reports", h.GetReports)
		v1.POST("/reports", h.SubmitReport)
		v1.POST("/debriefings", h.LogDebriefing)
		v1.GET("/sources/high-risk", h.GetHighRisk)
		v1.GET("/analytics/source-network", h.GetSourceNetwork)
	}
}

func (h *HumintHandler) CreateSource(c *gin.Context) {
	var req domain.CreateSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	source, err := h.svc.CreateSource(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, domain.SourceResponse{Source: source})
}

func (h *HumintHandler) UpdateCredibility(c *gin.Context) {
	code := c.Param("code")

	var req domain.UpdateCredibilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	source, err := h.svc.UpdateCredibility(code, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.SourceResponse{Source: source})
}

func (h *HumintHandler) GetReports(c *gin.Context) {
	code := c.Param("code")

	reports, err := h.svc.GetReportsBySource(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.ReportsResponse{Reports: reports, Total: len(reports)})
}

func (h *HumintHandler) SubmitReport(c *gin.Context) {
	var req domain.SubmitReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report, err := h.svc.SubmitReport(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, domain.ReportResponse{Report: report})
}

func (h *HumintHandler) LogDebriefing(c *gin.Context) {
	var req domain.LogDebriefingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	debriefing, err := h.svc.LogDebriefing(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, domain.DebriefingResponse{Debriefing: debriefing})
}

func (h *HumintHandler) GetHighRisk(c *gin.Context) {
	sources, err := h.svc.GetHighRiskSources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.HighRiskResponse{Sources: sources, Total: len(sources)})
}

func (h *HumintHandler) GetSourceNetwork(c *gin.Context) {
	network, err := h.svc.GetSourceNetwork()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, network)
}
