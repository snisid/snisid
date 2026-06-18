package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/gang/internal/domain"
	"github.com/snisid/platform/services/gang/internal/service"
)

type IncidentHandler struct {
	svc *service.IncidentService
}

func NewIncidentHandler(svc *service.IncidentService) *IncidentHandler {
	return &IncidentHandler{svc: svc}
}

func (h *IncidentHandler) Create(c *gin.Context) {
	var req domain.CreateIncidentRequest
	gangID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID gang invalide"})
		return
	}
	req.GangID = gangID
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "requête invalide"})
		return
	}
	createdBy := uuid.Nil
	if uid, ok := c.Request.Context().Value(ContextKeyUserID).(uuid.UUID); ok {
		createdBy = uid
	}
	inc, err := h.svc.CreateIncident(c.Request.Context(), req, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, inc)
}

func (h *IncidentHandler) ListByGang(c *gin.Context) {
	gangID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	incidents, err := h.svc.GetIncidents(c.Request.Context(), gangID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur requête"})
		return
	}
	c.JSON(http.StatusOK, incidents)
}
