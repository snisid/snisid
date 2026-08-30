package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/biar-svc/internal/domain"
	"github.com/snisid/platform/services/biar-svc/internal/service"
)

type BIARHandler struct {
	svc *service.BIARService
	log *zap.Logger
}

func NewBIARHandler(svc *service.BIARService, log *zap.Logger) *BIARHandler {
	return &BIARHandler{svc: svc, log: log}
}

func (h *BIARHandler) ReportWeapon(c *gin.Context) {
	var req domain.ReportWeaponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	weapon, err := h.svc.ReportIllicitWeapon(&req)
	if err != nil {
		h.log.Error("report weapon failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, weapon)
}

func (h *BIARHandler) GetWeapon(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid weapon id"})
		return
	}
	weapon, err := h.svc.GetWeapon(id)
	if err != nil {
		h.log.Error("get weapon failed", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "weapon not found"})
		return
	}
	c.JSON(http.StatusOK, weapon)
}

func (h *BIARHandler) CheckSerial(c *gin.Context) {
	sn := c.Param("sn")
	weapons, err := h.svc.CheckSerial(sn)
	if err != nil {
		h.log.Error("check serial failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"weapons": weapons, "count": len(weapons)})
}

func (h *BIARHandler) ReportBatch(c *gin.Context) {
	var req domain.ReportBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	batch, err := h.svc.ReportBatch(&req)
	if err != nil {
		h.log.Error("report batch failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, batch)
}

func (h *BIARHandler) GetStatsByGang(c *gin.Context) {
	stats, err := h.svc.GetStatsByGang()
	if err != nil {
		h.log.Error("get stats by gang failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

func (h *BIARHandler) GetStatsByOrigin(c *gin.Context) {
	stats, err := h.svc.GetStatsByOrigin()
	if err != nil {
		h.log.Error("get stats by origin failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

func (h *BIARHandler) GetRoutes(c *gin.Context) {
	stats, err := h.svc.GetRoutes()
	if err != nil {
		h.log.Error("get routes failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

func (h *BIARHandler) SyncFromIARMS(c *gin.Context) {
	result, err := h.svc.SyncFromIARMS()
	if err != nil {
		h.log.Error("sync from iARMS failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func SetupRouter(svc *service.BIARService, log *zap.Logger) *gin.Engine {
	r := gin.Default()
	handler := NewBIARHandler(svc, log)

	api := r.Group("/api/v1/biar")
	{
		api.POST("/weapons", handler.ReportWeapon)
		api.GET("/weapons/:id", handler.GetWeapon)
		api.GET("/check/serial/:sn", handler.CheckSerial)
		api.POST("/batches", handler.ReportBatch)
		api.GET("/stats/by-gang", handler.GetStatsByGang)
		api.GET("/stats/by-origin", handler.GetStatsByOrigin)
		api.GET("/stats/routes", handler.GetRoutes)
		api.POST("/iarms/sync", handler.SyncFromIARMS)
	}
	return r
}
