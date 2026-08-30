package rest

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
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vessel id"})
		return
	}
	vessel, err := h.svc.GetVessel(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vessel not found"})
		return
	}
	c.JSON(http.StatusOK, vessel)
}

func (h *MARHandler) CreateVessel(c *gin.Context) {
	var v domain.Vessel
	if err := c.ShouldBindJSON(&v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.RegisterVessel(&v); err != nil {
		h.log.Error("create vessel failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create vessel"})
		return
	}
	c.JSON(http.StatusCreated, v)
}

func (h *MARHandler) CreateIncident(c *gin.Context) {
	var i domain.Incident
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.ReportIncident(&i); err != nil {
		h.log.Error("create incident failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create incident"})
		return
	}
	c.JSON(http.StatusCreated, i)
}

func (h *MARHandler) GetRecentIncidents(c *gin.Context) {
	limit := 25
	if l, ok := c.GetQuery("limit"); ok {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}
	incidents, err := h.svc.GetRecentIncidents(limit)
	if err != nil {
		h.log.Error("get recent incidents failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get incidents"})
		return
	}
	c.JSON(http.StatusOK, incidents)
}

func (h *MARHandler) CreateWatch(c *gin.Context) {
	var w domain.WatchVessel
	if err := c.ShouldBindJSON(&w); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.AddToWatchList(&w); err != nil {
		h.log.Error("create watch failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create watch"})
		return
	}
	c.JSON(http.StatusCreated, w)
}

func (h *MARHandler) GetActiveWatches(c *gin.Context) {
	watches, err := h.svc.GetActiveWatches()
	if err != nil {
		h.log.Error("get active watches failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get watches"})
		return
	}
	c.JSON(http.StatusOK, watches)
}

func (h *MARHandler) GetAISLive(c *gin.Context) {
	limit := 100
	if l, ok := c.GetQuery("limit"); ok {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}
	sightings, err := h.svc.GetLivePositions(limit)
	if err != nil {
		h.log.Error("get AIS live failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get AIS data"})
		return
	}
	c.JSON(http.StatusOK, sightings)
}

func (h *MARHandler) GetZoneActivity(c *gin.Context) {
	zone := c.Param("zone")
	limit := 50
	if l, ok := c.GetQuery("limit"); ok {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}
	incidents, err := h.svc.GetZoneActivity(zone, limit)
	if err != nil {
		h.log.Error("get zone activity failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get zone activity"})
		return
	}
	c.JSON(http.StatusOK, incidents)
}

func (h *MARHandler) GetIncidentStats(c *gin.Context) {
	stats, err := h.svc.GetIncidentStats()
	if err != nil {
		h.log.Error("get incident stats failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}
