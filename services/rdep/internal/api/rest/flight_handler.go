package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/rdep/internal/domain"
	"github.com/snisid/platform/services/rdep/internal/service"
)

type FlightHandler struct {
	flightSvc *service.FlightService
}

func NewFlightHandler(fs *service.FlightService) *FlightHandler {
	return &FlightHandler{flightSvc: fs}
}

func (h *FlightHandler) Create(c *gin.Context) {
	var req domain.CreateFlightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	flight, err := h.flightSvc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, flight)
}

func (h *FlightHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID vol invalide"})
		return
	}

	flight, err := h.flightSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, flight)
}

func (h *FlightHandler) GetByNumber(c *gin.Context) {
	flightNumber := c.Param("flight_number")
	if flightNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "numéro de vol requis"})
		return
	}

	flight, err := h.flightSvc.GetByNumber(c.Request.Context(), flightNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, flight)
}

func (h *FlightHandler) List(c *gin.Context) {
	flights, err := h.flightSvc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"flights": flights})
}
