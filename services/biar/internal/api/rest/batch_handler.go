package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/biar/internal/domain"
	"github.com/snisid/platform/services/biar/internal/service"
)

type BatchHandler struct {
	svc *service.BatchService
}

func NewBatchHandler(svc *service.BatchService) *BatchHandler {
	return &BatchHandler{svc: svc}
}

func (h *BatchHandler) Create(c *gin.Context) {
	var req domain.CreateBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "requête invalide"})
		return
	}
	batch, err := h.svc.CreateBatch(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, batch)
}

func (h *BatchHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	batch, err := h.svc.GetBatch(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "lot introuvable"})
		return
	}
	c.JSON(http.StatusOK, batch)
}
