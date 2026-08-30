package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/lapi/internal/service"
)

type HTTPHandler struct {
	syncSvc *service.LAPISyncService
	cache   *service.OfflineCache
}

func NewHTTPHandler(s *service.LAPISyncService, c *service.OfflineCache) *HTTPHandler {
	return &HTTPHandler{syncSvc: s, cache: c}
}

func (h *HTTPHandler) RegisterRoutes(rg *gin.Engine) {
	rg.POST("/lapi/records", h.CreateRecord)
	rg.GET("/lapi/records/:id", h.GetRecord)
	rg.PUT("/lapi/records/:id", h.UpdateRecord)
	rg.GET("/lapi/pending", h.GetPending)
	rg.POST("/lapi/sync", h.TriggerSync)
}

func (h *HTTPHandler) CreateRecord(c *gin.Context) {
	var data map[string]any
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	record, err := h.syncSvc.CreateRecord(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, record)
}

func (h *HTTPHandler) GetRecord(c *gin.Context) {
	id := c.Param("id")
	record, ok := h.cache.Get(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}
	c.JSON(http.StatusOK, record)
}

func (h *HTTPHandler) UpdateRecord(c *gin.Context) {
	id := c.Param("id")
	var data map[string]any
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	record, err := h.syncSvc.UpdateRecord(id, data)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, record)
}

func (h *HTTPHandler) GetPending(c *gin.Context) {
	items := h.cache.GetPending()
	c.JSON(http.StatusOK, items)
}

func (h *HTTPHandler) TriggerSync(c *gin.Context) {
	if err := h.syncSvc.Sync(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "sync triggered"})
}
