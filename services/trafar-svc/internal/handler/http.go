package handler

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

func (h *TrafarHandler) ListRoutes(c *gin.Context) {
	routes, err := h.svc.ListRoutes()
	if err != nil {
		h.log.Error("list routes failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, routes)
}

func (h *TrafarHandler) GetRoute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid route id"})
		return
	}
	route, err := h.svc.GetRoute(id)
	if err != nil {
		h.log.Error("get route failed", zap.Error(err))
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
		h.log.Error("create route failed", zap.Error(err))
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
		h.log.Error("record shipment failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, shipment)
}

func (h *TrafarHandler) GetMap(c *gin.Context) {
	fc, err := h.svc.GetMapGeoJSON()
	if err != nil {
		h.log.Error("get map failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, fc)
}

func (h *TrafarHandler) GetStatsByOrigin(c *gin.Context) {
	stats, err := h.svc.GetStatsByOrigin()
	if err != nil {
		h.log.Error("get stats by origin failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

func (h *TrafarHandler) ListSuppliers(c *gin.Context) {
	suppliers, err := h.svc.ListSuppliers()
	if err != nil {
		h.log.Error("list suppliers failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, suppliers)
}

func SetupRouter(svc *service.TrafarService, log *zap.Logger) *gin.Engine {
	r := gin.Default()
	handler := NewTrafarHandler(svc, log)

	api := r.Group("/api/v1/trafar")
	{
		api.GET("/routes", handler.ListRoutes)
		api.GET("/routes/:id", handler.GetRoute)
		api.POST("/routes", handler.CreateRoute)
		api.POST("/shipments", handler.RecordShipment)
		api.GET("/map", handler.GetMap)
		api.GET("/stats/by-origin", handler.GetStatsByOrigin)
		api.GET("/suppliers", handler.ListSuppliers)
	}
	return r
}
