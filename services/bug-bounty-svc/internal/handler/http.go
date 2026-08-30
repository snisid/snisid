package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/snisid/bug-bounty-svc/internal/domain"
	"github.com/snisid/bug-bounty-svc/internal/service"
)

type Handler struct {
	svc *service.BugBountyService
}

func NewHandler(svc *service.BugBountyService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/programs", h.CreateProgram)
	r.GET("/programs", h.ListPrograms)
	r.POST("/reports", h.SubmitReport)
	r.GET("/reports/:id", h.GetReport)
	r.POST("/reports/:id/triage", h.TriageReport)
	r.POST("/reports/:id/reward", h.IssueReward)
	r.POST("/pentests", h.SchedulePentest)
	r.GET("/pentests/:id/results", h.GetPentestResults)
}

func (h *Handler) CreateProgram(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		ProgramID string   `json:"program_id"`
		Target    string   `json:"target"`
		ScopeType string   `json:"scope_type"`
		InScope   bool     `json:"in_scope"`
		RewardMin *float64 `json:"reward_min,omitempty"`
		RewardMax *float64 `json:"reward_max,omitempty"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	programID, err := uuid.Parse(req.ProgramID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid program_id"})
		return
	}

	result, err := h.svc.CreateProgram(c.Request.Context(), programID, req.Target, req.ScopeType, req.InScope, req.RewardMin, req.RewardMax)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) ListPrograms(c *gin.Context) {
	programs, err := h.svc.ListPrograms(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": programs})
}

func (h *Handler) SubmitReport(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		ProgramID   string  `json:"program_id"`
		Submitter   string  `json:"submitter"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Severity    string  `json:"severity"`
		ScopeID     *string `json:"scope_id,omitempty"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	programID, err := uuid.Parse(req.ProgramID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid program_id"})
		return
	}

	var scopeID *uuid.UUID
	if req.ScopeID != nil {
		if id, err := uuid.Parse(*req.ScopeID); err == nil {
			scopeID = &id
		}
	}

	result, err := h.svc.SubmitReport(c.Request.Context(), programID, req.Submitter, req.Title, req.Description, domain.Severity(req.Severity), scopeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report id"})
		return
	}

	report, err := h.svc.GetReport(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "report not found"})
		return
	}
	c.JSON(http.StatusOK, report)
}

func (h *Handler) TriageReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report id"})
		return
	}

	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		Triager      string  `json:"triager"`
		Severity     string  `json:"severity"`
		Reproducible bool    `json:"reproducible"`
		DuplicateOf  *string `json:"duplicate_of,omitempty"`
		Notes        *string `json:"notes,omitempty"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	var duplicateOf *uuid.UUID
	if req.DuplicateOf != nil {
		if dup, err := uuid.Parse(*req.DuplicateOf); err == nil {
			duplicateOf = &dup
		}
	}

	result, err := h.svc.TriageReport(c.Request.Context(), id, req.Triager, domain.Severity(req.Severity), req.Reproducible, duplicateOf, req.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) IssueReward(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report id"})
		return
	}

	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		Amount     float64 `json:"amount"`
		Currency   string  `json:"currency"`
		PaidTo     string  `json:"paid_to"`
		ApprovedBy string  `json:"approved_by"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.IssueReward(c.Request.Context(), id, req.Amount, req.Currency, req.PaidTo, req.ApprovedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) SchedulePentest(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		ProgramID string  `json:"program_id"`
		Title     string  `json:"title"`
		Scope     string  `json:"scope"`
		StartDate string  `json:"start_date"`
		EndDate   *string `json:"end_date,omitempty"`
		TeamLead  string  `json:"team_lead"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	programID, err := uuid.Parse(req.ProgramID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid program_id"})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date, use YYYY-MM-DD"})
		return
	}

	var endDate *time.Time
	if req.EndDate != nil {
		t, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date"})
			return
		}
		endDate = &t
	}

	result, err := h.svc.SchedulePentest(c.Request.Context(), programID, req.Title, req.Scope, startDate, endDate, req.TeamLead)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetPentestResults(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pentest id"})
		return
	}

	result, err := h.svc.GetPentestResults(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "pentest not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}
