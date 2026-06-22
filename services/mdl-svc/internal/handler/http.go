package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/snisid/mdl-svc/internal/service"
)

type Handler struct {
	svc *service.MDLService
}

func NewHandler(svc *service.MDLService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/mdl/issue", h.IssueMDL)
	r.GET("/mdl/:identity_id", h.GetMDL)
	r.POST("/mdl/verify", h.VerifyMDL)
	r.POST("/mdl/trust-readers", h.RegisterTrustedReader)
	r.GET("/mdl/trust-registry", h.GetTrustRegistry)
	r.POST("/mdl/reissue", h.ReissueMDL)
}

func (h *Handler) IssueMDL(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		IdentityID string `json:"identity_id"`
		DeviceID   string `json:"device_id"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	identityID, err := uuid.Parse(req.IdentityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid identity_id"})
		return
	}
	result, err := h.svc.IssueMDL(c.Request.Context(), identityID, req.DeviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetMDL(c *gin.Context) {
	identityID, err := uuid.Parse(c.Param("identity_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid identity_id"})
		return
	}
	result, err := h.svc.GetMDLByIdentity(c.Request.Context(), identityID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "mdl not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) VerifyMDL(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		IssuanceID string            `json:"issuance_id"`
		ReaderID   string            `json:"reader_id"`
		Elements   map[string]string `json:"elements"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	issuanceID, err := uuid.Parse(req.IssuanceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid issuance_id"})
		return
	}
	result, err := h.svc.VerifyMDLPresentation(c.Request.Context(), issuanceID, req.ReaderID, req.Elements)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) RegisterTrustedReader(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		ReaderID   string `json:"reader_id"`
		ReaderName string `json:"reader_name"`
		PublicKey  string `json:"public_key"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	result, err := h.svc.RegisterTrustedReader(c.Request.Context(), req.ReaderID, req.ReaderName, req.PublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetTrustRegistry(c *gin.Context) {
	result, err := h.svc.GetTrustRegistry(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *Handler) ReissueMDL(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		IdentityID string `json:"identity_id"`
		DeviceID   string `json:"device_id"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	identityID, err := uuid.Parse(req.IdentityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid identity_id"})
		return
	}
	result, err := h.svc.ReissueMDL(c.Request.Context(), identityID, req.DeviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}
