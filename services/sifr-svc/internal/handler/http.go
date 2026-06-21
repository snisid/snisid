package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/sifr-svc/internal/domain"
	"github.com/snisid/platform/services/sifr-svc/internal/service"
)

type SIFRHandler struct {
	svc *service.BorderService
	log *zap.Logger
}

func NewSIFRHandler(svc *service.BorderService, log *zap.Logger) *SIFRHandler {
	return &SIFRHandler{svc: svc, log: log}
}

func (h *SIFRHandler) ProcessCrossing(c *gin.Context) {
	var req domain.CrossingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	crossing, result, err := h.svc.ProcessCrossing(&req)
	if err != nil {
		h.log.Error("process crossing failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"crossing": crossing, "result": result})
}

func (h *SIFRHandler) SearchCrossings(c *gin.Context) {
	postIDStr := c.Query("post_id")
	var postID *uuid.UUID
	if postIDStr != "" {
		id, err := uuid.Parse(postIDStr)
		if err == nil {
			postID = &id
		}
	}
	direction := c.Query("direction")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	crossings, total, err := h.svc.SearchCrossings(postID, direction, nil, nil, page, pageSize)
	if err != nil {
		h.log.Error("search crossings failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"crossings": crossings, "total": total, "page": page, "page_size": pageSize})
}

func (h *SIFRHandler) GetPersonHistory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid person id"})
		return
	}
	crossings, err := h.svc.GetPersonHistory(id)
	if err != nil {
		h.log.Error("get person history failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, crossings)
}

func (h *SIFRHandler) GetActiveAlerts(c *gin.Context) {
	alerts, err := h.svc.GetActiveAlerts()
	if err != nil {
		h.log.Error("get active alerts failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, alerts)
}

func (h *SIFRHandler) ListPosts(c *gin.Context) {
	posts, err := h.svc.ListPosts()
	if err != nil {
		h.log.Error("list posts failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (h *SIFRHandler) GetDailyStats(c *gin.Context) {
	postIDStr := c.Query("post_id")
	var postID *uuid.UUID
	if postIDStr != "" {
		id, err := uuid.Parse(postIDStr)
		if err == nil {
			postID = &id
		}
	}
	stats, err := h.svc.GetDailyStats(postID)
	if err != nil {
		h.log.Error("get daily stats failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *SIFRHandler) ReportClandestine(c *gin.Context) {
	var req domain.ClandestineCrossing
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.ReportClandestineCrossing(&req)
	if err != nil {
		h.log.Error("report clandestine failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func SetupRouter(svc *service.BorderService, log *zap.Logger) *gin.Engine {
	r := gin.Default()
	handler := NewSIFRHandler(svc, log)

	api := r.Group("/api/v1/sifr")
	{
		api.POST("/crossings", handler.ProcessCrossing)
		api.GET("/crossings/search", handler.SearchCrossings)
		api.GET("/crossings/person/:id", handler.GetPersonHistory)
		api.GET("/alerts/active", handler.GetActiveAlerts)
		api.GET("/posts", handler.ListPosts)
		api.GET("/stats/daily", handler.GetDailyStats)
		api.POST("/clandestine", handler.ReportClandestine)
	}
	return r
}
