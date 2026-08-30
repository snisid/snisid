package rest

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/sifr-svc/internal/domain"
	"github.com/snisid/platform/services/sifr-svc/internal/service"
)

type SifrHandler struct {
	svc *service.BorderService
	log *zap.Logger
}

func NewSifrHandler(svc *service.BorderService, log *zap.Logger) *SifrHandler {
	return &SifrHandler{svc: svc, log: log}
}

func (h *SifrHandler) ProcessCrossing(c *gin.Context) {
	var req domain.CrossingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	saved, result, err := h.svc.ProcessCrossing(&req)
	if err != nil {
		h.log.Error("process crossing failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"crossing": saved,
		"result":   result,
	})
}

func (h *SifrHandler) SearchCrossings(c *gin.Context) {
	var postID *uuid.UUID
	if v := c.Query("post_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post_id"})
			return
		}
		postID = &id
	}

	direction := c.Query("direction")
	var dateFrom, dateTo *time.Time
	if v := c.Query("date_from"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date_from"})
			return
		}
		dateFrom = &t
	}
	if v := c.Query("date_to"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date_to"})
			return
		}
		dateTo = &t
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	crossings, total, err := h.svc.SearchCrossings(postID, direction, dateFrom, dateTo, page, pageSize)
	if err != nil {
		h.log.Error("search crossings failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"crossings": crossings,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *SifrHandler) GetPersonHistory(c *gin.Context) {
	idStr := c.Param("id")
	personID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid person id"})
		return
	}

	crossings, err := h.svc.GetPersonHistory(personID)
	if err != nil {
		h.log.Error("get person history failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"crossings": crossings})
}

func (h *SifrHandler) GetActiveAlerts(c *gin.Context) {
	alerts, err := h.svc.GetActiveAlerts()
	if err != nil {
		h.log.Error("get active alerts failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"alerts": alerts})
}

func (h *SifrHandler) ListPosts(c *gin.Context) {
	posts, err := h.svc.ListPosts()
	if err != nil {
		h.log.Error("list posts failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func (h *SifrHandler) GetDailyStats(c *gin.Context) {
	var postID *uuid.UUID
	if v := c.Query("post_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post_id"})
			return
		}
		postID = &id
	}

	stats, err := h.svc.GetDailyStats(postID)
	if err != nil {
		h.log.Error("get daily stats failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *SifrHandler) ReportClandestineCrossing(c *gin.Context) {
	var req domain.ClandestineCrossing
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report, err := h.svc.ReportClandestineCrossing(&req)
	if err != nil {
		h.log.Error("report clandestine crossing failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"report": report})
}
