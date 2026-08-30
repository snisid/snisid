package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/terr-svc/internal/domain"
	"github.com/snisid/platform/services/terr-svc/internal/service"
)

type TerrHandler struct {
	svc *service.TerritoryService
	log *zap.Logger
}

func NewTerrHandler(svc *service.TerritoryService, log *zap.Logger) *TerrHandler {
	return &TerrHandler{svc: svc, log: log}
}

func (h *TerrHandler) CheckPointSafety(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	if latStr == "" || lngStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lat and lng are required"})
		return
	}

	var lat, lng float64
	if _, err := fmt.Sscanf(latStr, "%f", &lat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lat"})
		return
	}
	if _, err := fmt.Sscanf(lngStr, "%f", &lng); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lng"})
		return
	}

	result, err := h.svc.CheckPointSafety(lat, lng)
	if err != nil {
		h.log.Error("check point safety failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *TerrHandler) GetRouteSafety(c *gin.Context) {
	waypointsJSON := c.Query("waypoints")
	if waypointsJSON == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "waypoints are required"})
		return
	}

	var waypoints []domain.Point
	if err := json.Unmarshal([]byte(waypointsJSON), &waypoints); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid waypoints format"})
		return
	}

	result, err := h.svc.GetRouteSafety(waypoints)
	if err != nil {
		h.log.Error("get route safety failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *TerrHandler) ListZones(c *gin.Context) {
	zones, err := h.svc.ListZones()
	if err != nil {
		h.log.Error("list zones failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, zones)
}

func (h *TerrHandler) ListZonesByDept(c *gin.Context) {
	code := c.Param("code")
	zones, err := h.svc.ListZonesByDept(code)
	if err != nil {
		h.log.Error("list zones by dept failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, zones)
}

func (h *TerrHandler) ListZonesByGang(c *gin.Context) {
	gangIDStr := c.Param("gang_id")
	gangID, err := uuid.Parse(gangIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid gang_id"})
		return
	}

	zones, err := h.svc.ListZonesByGang(gangID)
	if err != nil {
		h.log.Error("list zones by gang failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, zones)
}

func (h *TerrHandler) CreateZone(c *gin.Context) {
	var req domain.SeizureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	zone, err := h.svc.CreateZone(&req)
	if err != nil {
		h.log.Error("create zone failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, zone)
}

func (h *TerrHandler) GetZoneHistory(c *gin.Context) {
	zoneIDStr := c.Param("zone_id")
	zoneID, err := uuid.Parse(zoneIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid zone_id"})
		return
	}

	history, err := h.svc.GetZoneHistory(zoneID)
	if err != nil {
		h.log.Error("get zone history failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, history)
}

func (h *TerrHandler) CreateCheckpoint(c *gin.Context) {
	var cp domain.Checkpoint
	if err := c.ShouldBindJSON(&cp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if cp.ReportedAt.IsZero() {
		cp.ReportedAt = time.Now()
	}

	result, err := h.svc.ReportCheckpoint(&cp)
	if err != nil {
		h.log.Error("create checkpoint failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, result)
}
