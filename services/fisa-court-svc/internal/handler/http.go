package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/fisa-court-svc/internal/domain"
	"github.com/snisid/fisa-court-svc/internal/service"
)

type Handler struct {
	svc *service.FISAService
}

func NewHandler(svc *service.FISAService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/warrants", h.FileWarrant)
	r.PATCH("/warrants/:id/approve", h.ApproveWarrant)
	r.GET("/warrants/active", h.GetActiveWarrants)
	r.POST("/warrants/:id/renew", h.RenewWarrant)
	r.POST("/reports", h.FileReport)
	r.GET("/docket/:term", h.GetDocketByTerm)
	r.POST("/emergency", h.EmergencyAuthorization)
}

func (h *Handler) FileWarrant(c *gin.Context) {
	var req domain.FileWarrantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.FileWarrant(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) ApproveWarrant(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.ApproveWarrantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.ApproveWarrant(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetActiveWarrants(c *gin.Context) {
	warrants, err := h.svc.GetActiveWarrants(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, warrants)
}

func (h *Handler) RenewWarrant(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.RenewWarrantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.RenewWarrant(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) FileReport(c *gin.Context) {
	var req domain.FileReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.FileReport(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetDocketByTerm(c *gin.Context) {
	term := c.Param("term")
	if term == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "term is required"})
		return
	}
	docket, err := h.svc.GetDocketByTerm(c.Request.Context(), term)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, docket)
}

func (h *Handler) EmergencyAuthorization(c *gin.Context) {
	var req domain.EmergencyAuthorizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.EmergencyAuthorization(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}
