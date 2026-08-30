package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/mar-svc/internal/domain"
	"github.com/snisid/platform/services/mar-svc/internal/service"
)

type MARHandler struct {
	svc *service.MaritimeService
	log *zap.Logger
}

func NewMARHandler(svc *service.MaritimeService, log *zap.Logger) *MARHandler {
	return &MARHandler{svc: svc, log: log}
}

func (h *MARHandler) GetVessel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vessel id"})
		return
	}
	vessel, err := h.svc.GetVessel(id)
	if err != nil {
		h.log.Error("get vessel failed", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "vessel not found"})
		return
	}
	c.JSON(http.StatusOK, vessel)
}

func (h *MARHandler) RegisterVessel(c *gin.Context) {
	var vessel domain.Vessel
	if err := c.ShouldBindJSON(&vessel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.RegisterVessel(&vessel); err != nil {
		h.log.Error("register vessel failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, vessel)
}

func (h *MARHandler) ProcessAIS(c *gin.Context) {
	var msg domain.AISMessage
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sighting, err := h.svc.ProcessAISSighting(&msg)
	if err != nil {
		h.log.Error("process ais failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, sighting)
}

func (h *MARHandler) ReportIncident(c *gin.Context) {
	var incident domain.Incident
	if err := c.ShouldBindJSON(&incident); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.ReportIncident(&incident); err != nil {
		h.log.Error("report incident failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, incident)
}

func (h *MARHandler) GetRecentIncidents(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))
	incidents, err := h.svc.GetRecentIncidents(limit)
	if err != nil {
		h.log.Error("get recent incidents failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, incidents)
}

func (h *MARHandler) AddToWatch(c *gin.Context) {
	var watch domain.WatchVessel
	if err := c.ShouldBindJSON(&watch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.AddToWatchList(&watch); err != nil {
		h.log.Error("add to watch failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, watch)
}

func (h *MARHandler) GetActiveWatches(c *gin.Context) {
	watches, err := h.svc.GetActiveWatches()
	if err != nil {
		h.log.Error("get active watches failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, watches)
}

func (h *MARHandler) GetLivePositions(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	positions, err := h.svc.GetLivePositions(limit)
	if err != nil {
		h.log.Error("get live positions failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, positions)
}

func (h *MARHandler) GetZoneActivity(c *gin.Context) {
	zone := c.Param("zone")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	incidents, err := h.svc.GetZoneActivity(zone, limit)
	if err != nil {
		h.log.Error("get zone activity failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, incidents)
}

func (h *MARHandler) GetIncidentStats(c *gin.Context) {
	stats, err := h.svc.GetIncidentStats()
	if err != nil {
		h.log.Error("get incident stats failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func SetupRouter(svc *service.MaritimeService, log *zap.Logger) *gin.Engine {
	r := gin.Default()
	handler := NewMARHandler(svc, log)

	api := r.Group("/api/v1/mar")
	{
		api.GET("/vessels/:id", handler.GetVessel)
		api.POST("/vessels", handler.RegisterVessel)
		api.POST("/ais", handler.ProcessAIS)
		api.POST("/incidents", handler.ReportIncident)
		api.GET("/incidents/recent", handler.GetRecentIncidents)
		api.POST("/watch", handler.AddToWatch)
		api.GET("/watch/active", handler.GetActiveWatches)
		api.GET("/ais/live", handler.GetLivePositions)
		api.GET("/zones/:zone/activity", handler.GetZoneActivity)
		api.GET("/stats/incidents", handler.GetIncidentStats)
	}
	return r
}
