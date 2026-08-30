package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/offline-sync-engine/internal/engine"
	"github.com/snisid/platform/services/offline-sync-engine/internal/models"
)

type EnqueueRequest struct {
	EventType   string `json:"event_type" binding:"required"`
	Payload     string `json:"payload" binding:"required"`
	Priority    int    `json:"priority"`
	TerminalID  string `json:"terminal_id" binding:"required"`
	AggregateID string `json:"aggregate_id"`
	MaxRetries  int    `json:"max_retries"`
}

type SyncEngineHandler struct {
	engine *engine.SyncEngine
}

func NewSyncEngineHandler(syncEngine *engine.SyncEngine) *SyncEngineHandler {
	return &SyncEngineHandler{engine: syncEngine}
}

func (h *SyncEngineHandler) EnqueueHandler(c *gin.Context) {
	var req EnqueueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Priority < 0 || req.Priority > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "priority must be between 0 and 3 (0=general, 1=enrollment, 2=judicial, 3=emergency)"})
		return
	}

	event := &models.OfflineEvent{
		ID:          uuid.New().String(),
		EventType:   req.EventType,
		Payload:     req.Payload,
		Priority:    req.Priority,
		TerminalID:  req.TerminalID,
		AggregateID: req.AggregateID,
		MaxRetries:  req.MaxRetries,
	}

	if err := h.engine.QueueEvent(event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         event.ID,
		"status":     event.Status,
		"created_at": event.CreatedAt,
	})
}

func (h *SyncEngineHandler) SyncHandler(c *gin.Context) {
	result, err := h.engine.Sync()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *SyncEngineHandler) ListQueueHandler(c *gin.Context) {
	status := c.Query("status")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	events, total, err := h.engine.ListEvents(status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     events,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *SyncEngineHandler) DeleteHandler(c *gin.Context) {
	id := c.Param("id")
	if err := h.engine.RemoveEvent(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "event removed", "id": id})
}

func (h *SyncEngineHandler) StatusHandler(c *gin.Context) {
	status, err := h.engine.GetQueueStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	lastSyncTime := ""
	var lastSync models.OfflineEvent
	if err := h.engine.LastSynced(&lastSync); err == nil && lastSync.SyncedAt != nil {
		lastSyncTime = lastSync.SyncedAt.Format(time.RFC3339)
	}

	c.JSON(http.StatusOK, gin.H{
		"queue_status":   status,
		"last_sync_time": lastSyncTime,
		"service":        "offline-sync-engine",
		"version":        "1.0.0",
	})
}
