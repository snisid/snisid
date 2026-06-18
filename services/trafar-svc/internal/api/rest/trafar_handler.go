package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/trafar-svc/internal/domain"
	"github.com/snisid/platform/services/trafar-svc/internal/service"
)

type TrafarHandler struct {
	svc *service.TrafarService
	log *zap.Logger
}

func NewTrafarHandler(svc *service.TrafarService, log *zap.Logger) *TrafarHandler {
	return &TrafarHandler{svc: svc, log: log}
}

func (h *TrafarHandler) GetRoutes(c *gin.Context) {
	routes, err := h.svc.ListRoutes()
	if err != nil {
		h.log.Error("failed to list routes", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, routes)
}

func (h *TrafarHandler) GetRoute(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid route id"})
		return
	}
	route, err := h.svc.GetRoute(id)
	if err != nil {
		h.log.Error("failed to get route", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
		return
	}
	c.JSON(http.StatusOK, route)
}

func (h *TrafarHandler) CreateRoute(c *gin.Context) {
	var route domain.TrafarRoute
	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.CreateRoute(&route); err != nil {
		h.log.Error("failed to create route", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, route)
}

func (h *TrafarHandler) RecordShipment(c *gin.Context) {
	var shipment domain.TrafarShipment
	if err := c.ShouldBindJSON(&shipment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.RecordShipment(&shipment); err != nil {
		h.log.Error("failed to record shipment", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, shipment)
}

func (h *TrafarHandler) GetRoutesMap(c *gin.Context) {
	fc, err := h.svc.GetMapGeoJSON()
	if err != nil {
		h.log.Error("failed to get map", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, fc)
}

func (h *TrafarHandler) GetStatsByOrigin(c *gin.Context) {
	stats, err := h.svc.GetStatsByOrigin()
	if err != nil {
		h.log.Error("failed to get stats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *TrafarHandler) GetSuppliers(c *gin.Context) {
	suppliers, err := h.svc.ListSuppliers()
	if err != nil {
		h.log.Error("failed to list suppliers", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, suppliers)
}
