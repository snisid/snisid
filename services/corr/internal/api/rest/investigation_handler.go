package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/corr/internal/domain"
	"github.com/snisid/platform/services/corr/internal/service"
)

type InvestigationHandler struct {
	svc *service.InvestigationService
}

func NewInvestigationHandler(svc *service.InvestigationService) *InvestigationHandler {
	return &InvestigationHandler{svc: svc}
}

func (h *InvestigationHandler) Start(c *gin.Context) {
	caseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	investigator := uuid.Nil
	if uid, ok := c.Request.Context().Value(ContextKeyUserID).(uuid.UUID); ok {
		investigator = uid
	}
	cas, err := h.svc.StartInvestigation(c.Request.Context(), caseID, investigator)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cas)
}

type closeRequest struct {
	Status domain.CaseStatus `json:"status" validate:"required"`
	Notes  string            `json:"notes"`
}

func (h *InvestigationHandler) Close(c *gin.Context) {
	caseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	var req closeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "requête invalide"})
		return
	}
	cas, err := h.svc.CloseInvestigation(c.Request.Context(), caseID, req.Status, req.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cas)
}
