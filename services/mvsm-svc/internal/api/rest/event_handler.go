package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/mvsm-svc/internal/domain"
	"github.com/snisid/platform/services/mvsm-svc/internal/service"
)

type EventHandler struct {
	svc *service.EventService
	log *zap.Logger
}

func NewEventHandler(svc *service.EventService, log *zap.Logger) *EventHandler {
	return &EventHandler{svc: svc, log: log}
}

func (h *EventHandler) CreateEvent(c *gin.Context) {
	var req domain.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := h.svc.CreateEvent(&req)
	if err != nil {
		h.log.Error("create event failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, event)
}

func (h *EventHandler) ListUpcoming(c *gin.Context) {
	events, err := h.svc.ListUpcoming()
	if err != nil {
		h.log.Error("list upcoming failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (h *EventHandler) ListActive(c *gin.Context) {
	events, err := h.svc.ListActive()
	if err != nil {
		h.log.Error("list active failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (h *EventHandler) AddUpdate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req domain.AddUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.AddUpdate(id, &req); err != nil {
		h.log.Error("add update failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "update added"})
}

func (h *EventHandler) UpdateRiskLevel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req domain.UpdateRiskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateRiskLevel(id, &req); err != nil {
		h.log.Error("update risk level failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "risk updated"})
}
