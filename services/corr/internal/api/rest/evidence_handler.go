package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/corr/internal/domain"
	"github.com/snisid/platform/services/corr/internal/service"
)

type EvidenceHandler struct {
	svc *service.EvidenceService
}

func NewEvidenceHandler(svc *service.EvidenceService) *EvidenceHandler {
	return &EvidenceHandler{svc: svc}
}

func (h *EvidenceHandler) Create(c *gin.Context) {
	var req domain.CreateEvidenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "requête invalide"})
		return
	}
	collectedBy := uuid.Nil
	if uid, ok := c.Request.Context().Value(ContextKeyUserID).(uuid.UUID); ok {
		collectedBy = uid
	}
	e, err := h.svc.AddEvidence(c.Request.Context(), req, collectedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, e)
}

func (h *EvidenceHandler) ListByCase(c *gin.Context) {
	caseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	evidence, err := h.svc.GetEvidence(c.Request.Context(), caseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur requête"})
		return
	}
	c.JSON(http.StatusOK, evidence)
}
