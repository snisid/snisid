package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snisid/offline-ht/internal/domain"
	"github.com/snisid/offline-ht/internal/service"
)

type Handler struct {
	svc *service.OfflineService
}

func NewHandler(svc *service.OfflineService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/queue/push", h.PushQueue)
	r.POST("/sync/:terminal_id", h.SyncTerminal)
	r.GET("/terminals/status", h.GetTerminalsStatus)
	r.GET("/conflicts", h.GetConflicts)
}

func (h *Handler) PushQueue(c *gin.Context) {
	var req domain.PushQueueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.PushQueue(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) SyncTerminal(c *gin.Context) {
	terminalID := c.Param("terminal_id")
	items, err := h.svc.SyncTerminal(c.Request.Context(), terminalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"terminal_id": terminalID, "items": items})
}

func (h *Handler) GetTerminalsStatus(c *gin.Context) {
	terminals, err := h.svc.GetTerminalsStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, terminals)
}

func (h *Handler) GetConflicts(c *gin.Context) {
	items, err := h.svc.GetConflicts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}
