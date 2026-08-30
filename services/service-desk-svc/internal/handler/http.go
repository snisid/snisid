package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/snisid/service-desk-svc/internal/domain"
	"github.com/snisid/service-desk-svc/internal/service"
)

type Handler struct {
	svc *service.ServiceDeskService
}

func NewHandler(svc *service.ServiceDeskService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/cases", h.CreateCase)
	r.GET("/cases/:id", h.GetCase)
	r.GET("/cases", h.ListCases)
	r.POST("/cases/:id/challenge", h.IssueChallenge)
	r.POST("/cases/:id/verify", h.VerifyResponse)
	r.POST("/cases/:id/recover", h.ExecuteRecovery)
	r.POST("/cases/:id/notes", h.AddNote)
}

func (h *Handler) CreateCase(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		CitizenID   string `json:"citizen_id"`
		Subject     string `json:"subject"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	citizenID, err := uuid.Parse(req.CitizenID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid citizen_id"})
		return
	}

	result, err := h.svc.CreateCase(c.Request.Context(), citizenID, req.Subject, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetCase(c *gin.Context) {
	caseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid case id"})
		return
	}
	result, err := h.svc.GetCase(c.Request.Context(), caseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "case not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) ListCases(c *gin.Context) {
	status := domain.CaseStatus(c.DefaultQuery("status", "OPEN"))
	result, err := h.svc.ListCases(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *Handler) IssueChallenge(c *gin.Context) {
	caseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid case id"})
		return
	}
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		Method string `json:"method"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.IssueChallenge(c.Request.Context(), caseID, domain.RecoveryMethod(req.Method))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) VerifyResponse(c *gin.Context) {
	caseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid case id"})
		return
	}
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		ChallengeID string `json:"challenge_id"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	challengeID, err := uuid.Parse(req.ChallengeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid challenge_id"})
		return
	}

	if err := h.svc.VerifyResponse(c.Request.Context(), challengeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "verified", "case_id": caseID.String()})
}

func (h *Handler) ExecuteRecovery(c *gin.Context) {
	caseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid case id"})
		return
	}
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		CitizenID string `json:"citizen_id"`
		Method    string `json:"method"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	citizenID, err := uuid.Parse(req.CitizenID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid citizen_id"})
		return
	}

	result, err := h.svc.ExecuteRecovery(c.Request.Context(), caseID, citizenID, domain.RecoveryMethod(req.Method))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) AddNote(c *gin.Context) {
	caseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid case id"})
		return
	}
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		Author  string `json:"author"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.AddNote(c.Request.Context(), caseID, req.Author, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}
