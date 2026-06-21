package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/siar-svc/internal/domain"
	"github.com/snisid/platform/services/siar-svc/internal/service"
)

type SIARHandler struct {
	svc *service.SIARService
	log *zap.Logger
}

func NewSIARHandler(svc *service.SIARService, log *zap.Logger) *SIARHandler {
	return &SIARHandler{svc: svc, log: log}
}

func (h *SIARHandler) RegisterFirearm(c *gin.Context) {
	var req domain.RegisterFirearmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	firearm, err := h.svc.RegisterFirearm(&req)
	if err != nil {
		h.log.Error("register firearm failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, firearm)
}

func (h *SIARHandler) GetFirearm(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid firearm id"})
		return
	}
	firearm, err := h.svc.GetFirearmByID(id)
	if err != nil {
		h.log.Error("get firearm failed", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "firearm not found"})
		return
	}
	c.JSON(http.StatusOK, firearm)
}

func (h *SIARHandler) CheckSerial(c *gin.Context) {
	sn := c.Param("sn")
	if sn == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "serial number required"})
		return
	}
	firearm, err := h.svc.CheckSerial(sn)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "serial not found"})
		return
	}
	c.JSON(http.StatusOK, firearm)
}

func (h *SIARHandler) ReportSeizure(c *gin.Context) {
	var req domain.SeizureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	seizure, err := h.svc.ReportSeizure(&req)
	if err != nil {
		h.log.Error("report seizure failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, seizure)
}

func (h *SIARHandler) ReportStolen(c *gin.Context) {
	var req domain.StolenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.ReportStolen(&req); err != nil {
		h.log.Error("report stolen failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "firearm reported stolen"})
}

func (h *SIARHandler) GetLicensesByPerson(c *gin.Context) {
	personStr := c.Param("person")
	personID, err := uuid.Parse(personStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid person id"})
		return
	}
	licenses, err := h.svc.GetLicensesByPerson(personID)
	if err != nil {
		h.log.Error("get licenses failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, licenses)
}

func (h *SIARHandler) CreateLicense(c *gin.Context) {
	var req domain.CreateLicenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.CreateLicense(&req); err != nil {
		h.log.Error("create license failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "license created"})
}

func (h *SIARHandler) GetStatsByType(c *gin.Context) {
	stats, err := h.svc.GetStatsByType()
	if err != nil {
		h.log.Error("get stats failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func SetupRouter(svc *service.SIARService, log *zap.Logger) *gin.Engine {
	r := gin.Default()
	handler := NewSIARHandler(svc, log)

	api := r.Group("/api/v1/siar")
	{
		api.POST("/firearms", handler.RegisterFirearm)
		api.GET("/firearms/:id", handler.GetFirearm)
		api.GET("/check/serial/:sn", handler.CheckSerial)
		api.POST("/seizures", handler.ReportSeizure)
		api.POST("/stolen", handler.ReportStolen)
		api.GET("/licenses/:person", handler.GetLicensesByPerson)
		api.POST("/licenses", handler.CreateLicense)
		api.GET("/stats/by-type", handler.GetStatsByType)
	}
	return r
}
