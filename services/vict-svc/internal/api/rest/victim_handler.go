package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/vict-svc/internal/domain"
	"github.com/snisid/platform/services/vict-svc/internal/service"
)

type VictimHandler struct {
	svc *service.VictimService
	log *zap.Logger
}

func NewVictimHandler(svc *service.VictimService, log *zap.Logger) *VictimHandler {
	return &VictimHandler{svc: svc, log: log}
}

func (h *VictimHandler) RegisterVictim(c *gin.Context) {
	var req domain.RegisterVictimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	victim, err := h.svc.RegisterVictim(&req)
	if err != nil {
		h.log.Error("register victim failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, victim)
}

func (h *VictimHandler) GetVictim(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	victim, err := h.svc.GetVictim(id)
	if err != nil {
		h.log.Error("get victim failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, victim)
}

func (h *VictimHandler) CreateMassIncident(c *gin.Context) {
	var req domain.CreateMassIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mi, err := h.svc.CreateMassIncident(&req)
	if err != nil {
		h.log.Error("create mass incident failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, mi)
}

func (h *VictimHandler) ListMassIncidents(c *gin.Context) {
	incidents, err := h.svc.ListMassIncidents()
	if err != nil {
		h.log.Error("list mass incidents failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, incidents)
}

func (h *VictimHandler) ListByGang(c *gin.Context) {
	gangID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid gang id"})
		return
	}

	victims, err := h.svc.ListByGang(gangID)
	if err != nil {
		h.log.Error("list by gang failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, victims)
}

func (h *VictimHandler) GetStatsByType(c *gin.Context) {
	stats, err := h.svc.GetStatsByType()
	if err != nil {
		h.log.Error("get stats by type failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *VictimHandler) GetReparationList(c *gin.Context) {
	victims, err := h.svc.GetReparationList()
	if err != nil {
		h.log.Error("get reparation list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, victims)
}
