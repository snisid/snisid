package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snisid/sigint-ht/internal/domain"
	"github.com/snisid/sigint-ht/internal/service"
)

type SigintService interface {
	CreateTarget(req domain.CreateTargetRequest) (domain.InterceptionTarget, error)
	GetActiveTargets() ([]domain.InterceptionTarget, error)
	RecordInterception(targetID string, req domain.InterceptRequest) (domain.InterceptedCommunication, error)
	GetCommunications(targetID string) ([]domain.InterceptedCommunication, error)
	AnalyzeCDR(phone string) ([]domain.CDRAnalysis, error)
	EmergencyAuthorization(req domain.EmergencyRequest) (domain.EmergencyResponse, error)
}

type SigintHandler struct {
	svc SigintService
}

func NewSigintHandler(svc *service.SigintService) *SigintHandler {
	return &SigintHandler{svc: svc}
}

func (h *SigintHandler) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1/sigint")
	{
		v1.POST("/targets", h.CreateTarget)
		v1.GET("/targets/active", h.GetActiveTargets)
		v1.POST("/targets/:id/intercept", h.RecordInterception)
		v1.GET("/targets/:id/communications", h.GetCommunications)
		v1.GET("/cdr/analysis", h.AnalyzeCDR)
		v1.POST("/emergency", h.EmergencyAuthorization)
	}
}

func (h *SigintHandler) CreateTarget(c *gin.Context) {
	var req domain.CreateTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	target, err := h.svc.CreateTarget(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, domain.TargetResponse{Target: target})
}

func (h *SigintHandler) GetActiveTargets(c *gin.Context) {
	targets, err := h.svc.GetActiveTargets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.TargetsResponse{Targets: targets, Total: len(targets)})
}

func (h *SigintHandler) RecordInterception(c *gin.Context) {
	targetID := c.Param("id")

	var req domain.InterceptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comm, err := h.svc.RecordInterception(targetID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, domain.CommunicationResponse{Communication: comm})
}

func (h *SigintHandler) GetCommunications(c *gin.Context) {
	targetID := c.Param("id")

	comms, err := h.svc.GetCommunications(targetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.CommunicationsResponse{Communications: comms, Total: len(comms)})
}

func (h *SigintHandler) AnalyzeCDR(c *gin.Context) {
	phone := c.Query("phone")
	if phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone query parameter is required"})
		return
	}

	records, err := h.svc.AnalyzeCDR(phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.CDRAnalysisResponse{Records: records, Total: len(records)})
}

func (h *SigintHandler) EmergencyAuthorization(c *gin.Context) {
	var req domain.EmergencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.svc.EmergencyAuthorization(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}
