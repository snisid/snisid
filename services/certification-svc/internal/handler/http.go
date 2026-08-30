package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/snisid/certification-svc/internal/domain"
	"github.com/snisid/certification-svc/internal/service"
)

type Handler struct {
	svc *service.CertificationService
}

func NewHandler(svc *service.CertificationService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/profiles", h.CreateProfile)
	r.GET("/profiles/:identity_id", h.GetProfile)
	r.PUT("/profiles/:identity_id/ial", h.UpdateIAL)
	r.PUT("/profiles/:identity_id/aal", h.UpdateAAL)
	r.GET("/profiles/verify", h.VerifyCompliance)
	r.GET("/audit", h.GetAudit)
}

func (h *Handler) CreateProfile(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req domain.AssuranceProfile
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.CreateProfile(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetProfile(c *gin.Context) {
	identityID, err := uuid.Parse(c.Param("identity_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid identity id"})
		return
	}

	profile, err := h.svc.GetProfile(c.Request.Context(), identityID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (h *Handler) UpdateIAL(c *gin.Context) {
	identityID, err := uuid.Parse(c.Param("identity_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid identity id"})
		return
	}

	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		IAL       domain.IALLevel `json:"ial"`
		UpdatedBy string          `json:"updated_by"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.UpdateIAL(c.Request.Context(), identityID, req.IAL, req.UpdatedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) UpdateAAL(c *gin.Context) {
	identityID, err := uuid.Parse(c.Param("identity_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid identity id"})
		return
	}

	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		AAL       domain.AALLevel `json:"aal"`
		UpdatedBy string          `json:"updated_by"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.UpdateAAL(c.Request.Context(), identityID, req.AAL, req.UpdatedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) VerifyCompliance(c *gin.Context) {
	identityIDStr := c.Query("identity_id")
	identityID, err := uuid.Parse(identityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid identity_id query parameter"})
		return
	}

	requiredIAL := domain.IALLevel(c.DefaultQuery("required_ial", "IAL2"))
	requiredAAL := domain.AALLevel(c.DefaultQuery("required_aal", "AAL2"))

	result, err := h.svc.VerifyCompliance(c.Request.Context(), identityID, requiredIAL, requiredAAL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetAudit(c *gin.Context) {
	identityIDStr := c.Query("identity_id")
	var identityID uuid.UUID
	if identityIDStr != "" {
		var err error
		identityID, err = uuid.Parse(identityIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid identity_id"})
			return
		}
	}

	entries, err := h.svc.GetAudit(c.Request.Context(), identityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if entries == nil {
		entries = []domain.CertificationAudit{}
	}
	c.JSON(http.StatusOK, gin.H{"data": entries})
}
