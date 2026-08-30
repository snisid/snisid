package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/rdep/internal/domain"
	"github.com/snisid/platform/services/rdep/internal/service"
)

type ExtraditionHandler struct {
	extraditionSvc *service.ExtraditionService
}

func NewExtraditionHandler(es *service.ExtraditionService) *ExtraditionHandler {
	return &ExtraditionHandler{extraditionSvc: es}
}

func (h *ExtraditionHandler) Create(c *gin.Context) {
	var req domain.CreateExtraditionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	extradition, err := h.extraditionSvc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, extradition)
}

func (h *ExtraditionHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID extradition invalide"})
		return
	}

	extradition, err := h.extraditionSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, extradition)
}

func (h *ExtraditionHandler) List(c *gin.Context) {
	extraditions, err := h.extraditionSvc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"extraditions": extraditions})
}

func (h *ExtraditionHandler) UpdateStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID extradition invalide"})
		return
	}

	var req domain.UpdateExtraditionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	extradition, err := h.extraditionSvc.UpdateStatus(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, extradition)
}
