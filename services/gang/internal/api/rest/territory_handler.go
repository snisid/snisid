package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/gang/internal/domain"
	"github.com/snisid/platform/services/gang/internal/service"
)

type TerritoryHandler struct {
	svc *service.TerritoryService
}

func NewTerritoryHandler(svc *service.TerritoryService) *TerritoryHandler {
	return &TerritoryHandler{svc: svc}
}

func (h *TerritoryHandler) Create(c *gin.Context) {
	var req domain.CreateTerritoryRequest
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
	t, err := h.svc.CreateTerritory(c.Request.Context(), req, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, t)
}

func (h *TerritoryHandler) ListByGang(c *gin.Context) {
	gangID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	territories, err := h.svc.GetTerritories(c.Request.Context(), gangID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur requête"})
		return
	}
	c.JSON(http.StatusOK, territories)
}
