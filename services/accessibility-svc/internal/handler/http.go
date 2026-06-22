package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/snisid/accessibility-svc/internal/domain"
	"github.com/snisid/accessibility-svc/internal/service"
)

type Handler struct {
	svc *service.AccessibilityService
}

func NewHandler(svc *service.AccessibilityService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/audits", h.RunAudit)
	r.GET("/audits/:id", h.GetAuditResult)
	r.GET("/audits", h.ListAudits)
	r.POST("/violations/:id/remediate", h.MarkRemediated)
	r.GET("/compliance", h.GetCompliance)
	r.POST("/schedules", h.CreateSchedule)
	r.GET("/dashboard", h.GetDashboard)
}

func (h *Handler) RunAudit(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		TargetURL string `json:"target_url"`
		WCAGLevel string `json:"wcag_level"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.RunAudit(c.Request.Context(), req.TargetURL, domain.WCAGLevel(req.WCAGLevel))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetAuditResult(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid audit id"})
		return
	}

	result, err := h.svc.GetAuditResult(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "audit not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) ListAudits(c *gin.Context) {
	audits, err := h.svc.ListAudits(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": audits})
}

func (h *Handler) MarkRemediated(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid violation id"})
		return
	}

	if err := h.svc.MarkRemediated(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "remediated"})
}

func (h *Handler) GetCompliance(c *gin.Context) {
	reports, err := h.svc.GetComplianceOverview(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": reports})
}

func (h *Handler) CreateSchedule(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		TargetURL string `json:"target_url"`
		WCAGLevel string `json:"wcag_level"`
		CronExpr  string `json:"cron_expr"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.CreateSchedule(c.Request.Context(), req.TargetURL, domain.WCAGLevel(req.WCAGLevel), req.CronExpr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetDashboard(c *gin.Context) {
	reports, err := h.svc.GetDashboard(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": reports})
}
