package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/opr-svc/internal/domain"
	"github.com/snisid/platform/services/opr-svc/internal/service"
)

type OPRHandler struct {
	oprSvc *service.OPRService
}

func NewOPRHandler(oprSvc *service.OPRService) *OPRHandler {
	return &OPRHandler{oprSvc: oprSvc}
}

func (h *OPRHandler) CreateOrder(c *gin.Context) {
	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.oprSvc.CreateOrder(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OPRHandler) CheckSubject(c *gin.Context) {
	personIDStr := c.Param("person_id")
	personID, err := uuid.Parse(personIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID invalide"})
		return
	}

	result, err := h.oprSvc.CheckSubject(c.Request.Context(), personID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *OPRHandler) RecordViolation(c *gin.Context) {
	var req domain.ViolationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reportedByStr := c.Query("reported_by")
	reportedBy, err := uuid.Parse(reportedByStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID invalide"})
		return
	}

	if err := h.oprSvc.RecordViolation(c.Request.Context(), req, reportedBy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Violation enregistrée"})
}

func (h *OPRHandler) GetExpiringSoon(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		days = 30
	}

	orders, err := h.oprSvc.GetExpiringSoon(c.Request.Context(), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func (h *OPRHandler) GetByGangID(c *gin.Context) {
	gangIDStr := c.Param("id")
	gangID, err := uuid.Parse(gangIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID invalide"})
		return
	}

	orders, err := h.oprSvc.GetByGangID(c.Request.Context(), gangID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}
