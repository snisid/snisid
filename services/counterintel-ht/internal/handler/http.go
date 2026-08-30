package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/counterintel-ht/internal/domain"
	"github.com/snisid/counterintel-ht/internal/service"
)

type Handler struct {
	svc *service.CounterintelService
}

func NewHandler(svc *service.CounterintelService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/investigations", h.CreateInvestigation)
	r.GET("/investigations/pending", h.GetPendingInvestigations)
	r.GET("/investigations/:id", h.GetInvestigation)
	r.PATCH("/investigations/:id/adjudicate", h.AdjudicateInvestigation)
	r.POST("/threats", h.ReportThreat)
	r.GET("/threats/active", h.GetActiveThreats)
	r.POST("/contacts", h.ReportContact)
	r.GET("/contacts/:subject_id", h.GetContactsBySubject)
}

func (h *Handler) CreateInvestigation(c *gin.Context) {
	var req domain.CreateInvestigationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.CreateInvestigation(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetInvestigation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	inv, err := h.svc.GetInvestigation(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, inv)
}

func (h *Handler) GetPendingInvestigations(c *gin.Context) {
	invs, err := h.svc.GetPendingInvestigations(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, invs)
}

func (h *Handler) AdjudicateInvestigation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.AdjudicateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.AdjudicateInvestigation(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) ReportThreat(c *gin.Context) {
	var req domain.ReportThreatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.ReportThreat(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetActiveThreats(c *gin.Context) {
	alerts, err := h.svc.GetActiveThreats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, alerts)
}

func (h *Handler) ReportContact(c *gin.Context) {
	var req domain.ReportContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.ReportContact(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetContactsBySubject(c *gin.Context) {
	subjectID := c.Param("subject_id")
	if subjectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subject_id is required"})
		return
	}
	contacts, err := h.svc.GetContactsBySubject(c.Request.Context(), subjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, contacts)
}
