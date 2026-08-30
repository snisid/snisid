package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/aero-svc/internal/domain"
	"github.com/snisid/platform/services/aero-svc/internal/service"
)

type AeroHandler struct {
	svc *service.AeroService
	log *zap.Logger
}

func NewAeroHandler(svc *service.AeroService, log *zap.Logger) *AeroHandler {
	return &AeroHandler{svc: svc, log: log}
}

func (h *AeroHandler) CheckRegistration(c *gin.Context) {
	reg := c.Param("reg")
	if reg == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "registration mark is required"})
		return
	}

	result, err := h.svc.CheckRegistration(reg)
	if err != nil {
		h.log.Error("check registration failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *AeroHandler) ReportStrip(c *gin.Context) {
	var req domain.ReportStripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	strip, err := h.svc.ReportStrip(&req)
	if err != nil {
		h.log.Error("report strip failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, strip)
}

func (h *AeroHandler) GetStripMap(c *gin.Context) {
	fc, err := h.svc.GetStripMap()
	if err != nil {
		h.log.Error("get strip map failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, fc)
}

func (h *AeroHandler) ReportSuspiciousFlight(c *gin.Context) {
	var req domain.ReportFlightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	flight, err := h.svc.ReportSuspiciousFlight(&req)
	if err != nil {
		h.log.Error("report suspicious flight failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, flight)
}

func (h *AeroHandler) GetStripStats(c *gin.Context) {
	stats, err := h.svc.GetStripStats()
	if err != nil {
		h.log.Error("get strip stats failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
