package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/blkl-svc/internal/domain"
	"github.com/snisid/platform/services/blkl-svc/internal/service"
)

type BLKLHandler struct {
	svc *service.BLKLService
	log *zap.Logger
}

func NewBLKLHandler(svc *service.BLKLService, log *zap.Logger) *BLKLHandler {
	return &BLKLHandler{svc: svc, log: log}
}

func (h *BLKLHandler) CheckPerson(c *gin.Context) {
	idStr := c.Param("person_id")
	personID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid person id"})
		return
	}
	result, err := h.svc.CheckPerson(personID)
	if err != nil {
		h.log.Error("check person failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *BLKLHandler) AddEntry(c *gin.Context) {
	var req domain.AddEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entry, err := h.svc.AddEntry(&req)
	if err != nil {
		h.log.Error("add entry failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, entry)
}

func (h *BLKLHandler) LiftEntry(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entry id"})
		return
	}
	liftedBy := c.DefaultQuery("lifted_by", "system")
	if err := h.svc.LiftEntry(id, liftedBy); err != nil {
		h.log.Error("lift entry failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "entry lifted"})
}

func (h *BLKLHandler) GetActiveEntries(c *gin.Context) {
	entries, err := h.svc.GetActiveEntries()
	if err != nil {
		h.log.Error("get active entries failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, entries)
}

func (h *BLKLHandler) GetExpiringSoon(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	entries, err := h.svc.GetExpiringSoon(days)
	if err != nil {
		h.log.Error("get expiring soon failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, entries)
}

func SetupRouter(svc *service.BLKLService, log *zap.Logger) *gin.Engine {
	r := gin.Default()
	handler := NewBLKLHandler(svc, log)

	api := r.Group("/api/v1/blkl")
	{
		api.GET("/check/:person_id", handler.CheckPerson)
		api.POST("/entries", handler.AddEntry)
		api.PATCH("/entries/:id/lift", handler.LiftEntry)
		api.GET("/entries/active", handler.GetActiveEntries)
		api.GET("/expiring-soon", handler.GetExpiringSoon)
	}
	return r
}
