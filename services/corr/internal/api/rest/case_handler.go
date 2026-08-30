package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/corr/internal/domain"
	"github.com/snisid/platform/services/corr/internal/service"
)

type CaseHandler struct {
	svc *service.CaseService
}

func NewCaseHandler(svc *service.CaseService) *CaseHandler {
	return &CaseHandler{svc: svc}
}

func (h *CaseHandler) Create(c *gin.Context) {
	var req domain.CreateCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "requête invalide"})
		return
	}
	createdBy := uuid.Nil
	if uid, ok := c.Request.Context().Value(ContextKeyUserID).(uuid.UUID); ok {
		createdBy = uid
	}
	cas, err := h.svc.CreateCase(c.Request.Context(), req, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cas)
}

func (h *CaseHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	cas, err := h.svc.GetCase(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cas introuvable"})
		return
	}
	c.JSON(http.StatusOK, cas)
}

func (h *CaseHandler) ListActive(c *gin.Context) {
	cases, err := h.svc.ListActive(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur requête"})
		return
	}
	c.JSON(http.StatusOK, cases)
}
