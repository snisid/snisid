package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/snisid/audit-svc/internal/domain"
	"github.com/snisid/audit-svc/internal/service"
)

type Handler struct {
	svc *service.AuditService
}

func NewHandler(svc *service.AuditService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/events", h.IngestEvent)
	r.GET("/events/:id", h.GetEvent)
	r.GET("/events", h.SearchEvents)
	r.GET("/events/report", h.GenerateReport)
	r.POST("/events/verify", h.VerifyIntegrity)
	r.GET("/stats", h.GetStats)
}

func (h *Handler) IngestEvent(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		Source     string         `json:"source"`
		EventType  string         `json:"event_type"`
		Category   string         `json:"category"`
		ActorID    string         `json:"actor_id,omitempty"`
		ResourceID string         `json:"resource_id"`
		Action     string         `json:"action"`
		Payload    map[string]any `json:"payload,omitempty"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	var actorID *uuid.UUID
	if req.ActorID != "" {
		if id, err := uuid.Parse(req.ActorID); err == nil {
			actorID = &id
		}
	}

	result, err := h.svc.IngestEvent(c.Request.Context(),
		domain.EventSource(req.Source),
		domain.EventType(req.EventType),
		domain.AuditCategory(req.Category),
		actorID, req.ResourceID, req.Action, req.Payload,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetEvent(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}
	result, err := h.svc.GetEvent(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) SearchEvents(c *gin.Context) {
	q := domain.AuditQuery{Limit: 100}
	if src := c.Query("source"); src != "" {
		s := domain.EventSource(src)
		q.Source = &s
	}
	if et := c.Query("event_type"); et != "" {
		e := domain.EventType(et)
		q.EventType = &e
	}
	if cat := c.Query("category"); cat != "" {
		ca := domain.AuditCategory(cat)
		q.Category = &ca
	}
	if from := c.Query("from"); from != "" {
		t, err := time.Parse(time.RFC3339, from)
		if err == nil {
			q.From = &t
		}
	}
	if to := c.Query("to"); to != "" {
		t, err := time.Parse(time.RFC3339, to)
		if err == nil {
			q.To = &t
		}
	}

	result, err := h.svc.SearchEvents(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *Handler) GenerateReport(c *gin.Context) {
	category := domain.AuditCategory(c.DefaultQuery("category", "OPERATIONAL"))
	fromStr := c.DefaultQuery("from", time.Now().Add(-30*24*time.Hour).Format(time.RFC3339))
	toStr := c.DefaultQuery("to", time.Now().Format(time.RFC3339))

	from, _ := time.Parse(time.RFC3339, fromStr)
	to, _ := time.Parse(time.RFC3339, toStr)

	result, err := h.svc.GenerateReport(c.Request.Context(), category, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) VerifyIntegrity(c *gin.Context) {
	valid, err := h.svc.VerifyIntegrity(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"integrity_valid": valid})
}

func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.svc.GetStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
