package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/snisid/age-verification-svc/internal/domain"
	"github.com/snisid/age-verification-svc/internal/service"
)

type Handler struct {
	svc *service.AgeVerificationService
}

func NewHandler(svc *service.AgeVerificationService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/attestations", h.CreateAttestation)
	r.POST("/attestations/verify", h.VerifyAttestation)
	r.GET("/attestations/:attestation_id", h.GetAttestation)
	r.POST("/attestations/selective", h.SelectiveVerification)
	r.POST("/attestations/revoke", h.RevokeAttestation)
}

func (h *Handler) CreateAttestation(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		IdentityID  string `json:"identity_id"`
		DateOfBirth string `json:"date_of_birth"`
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
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date_of_birth, use YYYY-MM-DD"})
		return
	}
	result, err := h.svc.CreateAttestation(c.Request.Context(), identityID, dob)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) VerifyAttestation(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		AttestationID string `json:"attestation_id"`
		VerifierID    string `json:"verifier_id"`
		Bracket       string `json:"bracket"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	attestationID, err := uuid.Parse(req.AttestationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attestation_id"})
		return
	}
	bracket := domain.AgeBracket(req.Bracket)
	switch bracket {
	case domain.AgeBracketOver18, domain.AgeBracketOver21, domain.AgeBracketOver65:
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bracket, use OVER_18, OVER_21, or OVER_65"})
		return
	}
	result, err := h.svc.VerifyAgeClaim(c.Request.Context(), attestationID, req.VerifierID, bracket)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetAttestation(c *gin.Context) {
	attestationID, err := uuid.Parse(c.Param("attestation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attestation_id"})
		return
	}
	result, err := h.svc.GetAttestation(c.Request.Context(), attestationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "attestation not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) SelectiveVerification(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		AttestationID string `json:"attestation_id"`
		VerifierID    string `json:"verifier_id"`
		Bracket       string `json:"bracket"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	attestationID, err := uuid.Parse(req.AttestationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attestation_id"})
		return
	}
	bracket := domain.AgeBracket(req.Bracket)
	switch bracket {
	case domain.AgeBracketOver18, domain.AgeBracketOver21, domain.AgeBracketOver65:
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bracket"})
		return
	}
	result, err := h.svc.SelectiveBracketVerification(c.Request.Context(), attestationID, req.VerifierID, bracket)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) RevokeAttestation(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		AttestationID string `json:"attestation_id"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	attestationID, err := uuid.Parse(req.AttestationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attestation_id"})
		return
	}
	if err := h.svc.RevokeAttestation(c.Request.Context(), attestationID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "revoked"})
}
