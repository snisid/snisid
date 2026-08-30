package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/sipep/internal/domain"
	"github.com/snisid/platform/services/sipep/internal/service"
)

type MovementHandler struct {
	movementService *service.MovementService
	inmateService   *service.InmateService
}

func NewMovementHandler(ms *service.MovementService, is *service.InmateService) *MovementHandler {
	return &MovementHandler{movementService: ms, inmateService: is}
}

func (h *MovementHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/transfers", h.Transfer)
	rg.POST("/health-events", h.HealthEvent)
}

func (h *MovementHandler) Transfer(c *gin.Context) {
	var req domain.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.inmateService.GetInmate(req.InmateID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "inmate not found"})
		return
	}

	movement, err := h.movementService.Transfer(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	inmate, _ := h.inmateService.GetInmate(req.InmateID)
	if inmate != nil {
		inmate.CurrentFacility = req.ToFacility
		inmate.CellBlock = req.ToBlock
		_ = h.inmateService.UpdateInmate(inmate)
	}

	c.JSON(http.StatusCreated, gin.H{"movement": movement})
}

func (h *MovementHandler) HealthEvent(c *gin.Context) {
	var req domain.HealthEvent
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.EventID = uuid.New()
	c.JSON(http.StatusCreated, gin.H{"event": req})
}
