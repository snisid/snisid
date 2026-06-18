package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/afis-svc/internal/domain"
	"github.com/snisid/platform/services/afis-svc/internal/service"
)

type LatentHandler struct {
	svc *service.LatentService
}

func NewLatentHandler(svc *service.LatentService) *LatentHandler {
	return &LatentHandler{svc: svc}
}

func (h *LatentHandler) Submit(c *gin.Context) {
	var req domain.LatentSubmission
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lp, err := h.svc.Submit(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ginH(http.StatusCreated, "Empreinte latente soumise", lp))
}

func (h *LatentHandler) SearchLatent(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID latent invalide"})
		return
	}

	results, err := h.svc.SearchLatent(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ginH(http.StatusOK, "Recherche latente terminée", results))
}

func (h *LatentHandler) ConfirmMatch(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID latent invalide"})
		return
	}

	var req domain.LatentMatchConfirm
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.ConfirmMatch(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ginH(http.StatusOK, "Correspondance confirmée", nil))
}
