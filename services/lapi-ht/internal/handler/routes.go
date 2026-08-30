package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/lapi-ht/internal/domain"
	"github.com/snisid/lapi-ht/internal/service"
)

type Handler struct {
	svc *service.LapiService
}

func NewHandler(svc *service.LapiService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1/lapi")
	{
		v1.POST("/reads", h.CreateRead)
		v1.GET("/reads/recent", h.GetRecentReads)
		v1.GET("/reads/plate/:number", h.GetReadsByPlate)
		v1.GET("/alerts/active", h.GetActiveAlerts)
		v1.GET("/cameras/status", h.GetCameraStatus)
	}
}

type createReadRequest struct {
	CameraID              uuid.UUID  `json:"camera_id" binding:"required"`
	PlateNumberRaw        string     `json:"plate_number_raw" binding:"required"`
	PlateNumberNormalized string     `json:"plate_number_normalized"`
	OcrConfidence         *float64   `json:"ocr_confidence"`
	Latitude              *float64   `json:"latitude"`
	Longitude             *float64   `json:"longitude"`
	SpeedEstimateKmh      *float64   `json:"speed_estimate_kmh"`
	AlertTriggered        *bool      `json:"alert_triggered"`
	CapturedAt            *time.Time `json:"captured_at"`
}

func (h *Handler) CreateRead(c *gin.Context) {
	var req createReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	read := &domain.PlateRead{
		CameraID:              req.CameraID,
		PlateNumberRaw:        req.PlateNumberRaw,
		PlateNumberNormalized: req.PlateNumberNormalized,
		Latitude:              req.Latitude,
		Longitude:             req.Longitude,
		SpeedEstimateKmh:      req.SpeedEstimateKmh,
	}
	if read.PlateNumberNormalized == "" {
		read.PlateNumberNormalized = read.PlateNumberRaw
	}
	if req.OcrConfidence != nil {
		read.OcrConfidence = *req.OcrConfidence
	}
	if req.AlertTriggered != nil {
		read.AlertTriggered = *req.AlertTriggered
	}
	if req.CapturedAt != nil {
		read.CapturedAt = *req.CapturedAt
	} else {
		read.CapturedAt = time.Now().UTC()
	}

	if err := h.svc.RecordRead(read); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, read)
}

func (h *Handler) GetRecentReads(c *gin.Context) {
	limit := 50
	if l := c.Query("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil {
			limit = n
		}
	}
	reads, err := h.svc.GetRecentReads(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reads)
}

func (h *Handler) GetReadsByPlate(c *gin.Context) {
	number := c.Param("number")
	if number == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plate number required"})
		return
	}
	reads, err := h.svc.GetReadsByPlate(number)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reads)
}

func (h *Handler) GetActiveAlerts(c *gin.Context) {
	alerts, err := h.svc.GetActiveAlerts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, alerts)
}

func (h *Handler) GetCameraStatus(c *gin.Context) {
	cameras, err := h.svc.GetCameraStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cameras)
}
