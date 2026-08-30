package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snisid/pki-ht/internal/domain"
	"github.com/snisid/pki-ht/internal/service"
)

type Handler struct {
	svc *service.PKIService
}

func NewHandler(svc *service.PKIService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/issue", h.Issue)
	r.POST("/revoke", h.Revoke)
	r.GET("/ocsp", h.OCSP)
	r.GET("/crl/:ca_id", h.CRL)
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

func (h *Handler) Revoke(c *gin.Context) {
	var req domain.RevokeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Revoke(c.Request.Context(), req.SerialNumber, req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "revoked"})
}

func (h *Handler) OCSP(c *gin.Context) {
	serial := c.Query("serial")
	if serial == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "serial parameter required"})
		return
	}
	cert, err := h.svc.CheckOCSP(c.Request.Context(), serial)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "certificate not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"serial": cert.SerialNumber, "status": cert.Status})
}

func (h *Handler) CRL(c *gin.Context) {
	caID := c.Param("ca_id")
	crl, err := h.svc.GetCRL(c.Request.Context(), caID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "CRL not found"})
		return
	}
	c.JSON(http.StatusOK, crl)
}
