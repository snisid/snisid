package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/rdep-svc/internal/domain"
	"github.com/snisid/platform/services/rdep-svc/internal/service"
)

type IntakeHandler struct {
	intakeSvc    *service.IntakeService
	screeningSvc *service.ScreeningService
}

func NewIntakeHandler(intakeSvc *service.IntakeService, screeningSvc *service.ScreeningService) *IntakeHandler {
	return &IntakeHandler{
		intakeSvc:    intakeSvc,
		screeningSvc: screeningSvc,
	}
}

func (h *IntakeHandler) ProcessIntake(c *gin.Context) {
	var req domain.DeporteeIntakeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deportee, err := h.intakeSvc.ProcessIntake(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, deportee)
}

func (h *IntakeHandler) GetDeportee(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID invalide"})
		return
	}

	deportee, err := h.intakeSvc.GetDeportee(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Déporté non trouvé"})
		return
	}

	c.JSON(http.StatusOK, deportee)
}

func (h *IntakeHandler) ScreenDeportee(c *gin.Context) {
	idStr := c.Param("id")
	_, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID invalide"})
		return
	}

	var req service.ScreenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.DeporteeID = idStr

	result, err := h.screeningSvc.ScreenDeportee(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *IntakeHandler) GetHighRisk(c *gin.Context) {
	deportees, err := h.intakeSvc.GetHighRisk(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deportees": deportees})
}

func (h *IntakeHandler) GetGangAffiliated(c *gin.Context) {
	deportees, err := h.intakeSvc.GetGangAffiliated(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deportees": deportees})
}

func (h *IntakeHandler) GetStatsByCountry(c *gin.Context) {
	stats, err := h.intakeSvc.GetStatsByCountry(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
