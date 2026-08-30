package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/transport-security-ht/internal/domain"
	"github.com/snisid/transport-security-ht/internal/service"
)

type TransportHandler struct {
	svc service.TransportService
}

func NewTransportHandler(svc service.TransportService) *TransportHandler {
	return &TransportHandler{svc: svc}
}

func (h *TransportHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1/transport")
	{
		api.POST("/screenings", h.LogScreening)
		api.GET("/screenings/recent", h.GetRecentScreenings)
		api.POST("/no-fly", h.AddNoFly)
		api.GET("/no-fly/check", h.CheckNoFly)
		api.GET("/zones/:airport", h.GetZonesByAirport)
		api.POST("/zones/:id/breach", h.ReportZoneBreach)
	}
}

func (h *TransportHandler) LogScreening(c *gin.Context) {
	var s domain.PassengerScreening
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.LogScreening(c.Request.Context(), &s); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, s)
}

func (h *TransportHandler) GetRecentScreenings(c *gin.Context) {
	result, err := h.svc.GetRecentScreenings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *TransportHandler) AddNoFly(c *gin.Context) {
	var p domain.NoFlyPassenger
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.AddNoFlyEntry(c.Request.Context(), &p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *TransportHandler) CheckNoFly(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "identity query param required"})
		return
	}
	result, err := h.svc.CheckNoFly(c.Request.Context(), identity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if result == nil {
		c.JSON(http.StatusOK, gin.H{"match": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"match": true, "entry": result})
}

func (h *TransportHandler) GetZonesByAirport(c *gin.Context) {
	airport := c.Param("airport")
	if airport == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "airport code required"})
		return
	}
	result, err := h.svc.GetZonesByAirport(c.Request.Context(), airport)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *TransportHandler) ReportZoneBreach(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid zone id"})
		return
	}
	if err := h.svc.ReportBreach(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "breach reported"})
}
