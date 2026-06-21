package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snisid/card-ht/internal/domain"
	"github.com/snisid/card-ht/internal/service"
)

type Handler struct {
	svc *service.CardService
}

func NewHandler(svc *service.CardService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/issue", h.Issue)
	r.GET("/verify/:doc_number", h.Verify)
	r.POST("/:id/report-lost", h.ReportLost)
	r.POST("/:id/revoke", h.Revoke)
}

func (h *Handler) Issue(c *gin.Context) {
	var req domain.IssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.Issue(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) Verify(c *gin.Context) {
	docNumber := c.Param("doc_number")
	doc, err := h.svc.Verify(c.Request.Context(), docNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}
	c.JSON(http.StatusOK, doc)
}

func (h *Handler) ReportLost(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.ReportLost(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "reported_lost"})
}

func (h *Handler) Revoke(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Revoke(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "revoked"})
}
