package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/snisid/field-ht/internal/domain"
	"github.com/snisid/field-ht/internal/service"
)

type Handler struct {
	svc *service.FieldService
}

func NewHandler(svc *service.FieldService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/missions", h.CreateMission)
	r.GET("/missions/active", h.GetActiveMissions)
	r.POST("/missions/:id/log", h.CreateMissionLog)
	r.GET("/stats/coverage", h.GetCoverageStats)
}

func (h *Handler) CreateMission(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req domain.CreateMissionRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.CreateMission(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetActiveMissions(c *gin.Context) {
	missions, err := h.svc.GetActiveMissions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, missions)
}

func (h *Handler) CreateMissionLog(c *gin.Context) {
	missionID := c.Param("id")
	body, _ := io.ReadAll(c.Request.Body)
	var req domain.CreateMissionLogRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.CreateMissionLog(c.Request.Context(), missionID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetCoverageStats(c *gin.Context) {
	stats, err := h.svc.GetCoverageStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
