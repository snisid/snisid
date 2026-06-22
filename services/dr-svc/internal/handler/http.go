package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/snisid/dr-svc/internal/service"
)

type Handler struct {
	svc *service.DRService
}

func NewHandler(svc *service.DRService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/regions", h.ListRegions)
	r.POST("/failover/plan", h.CreateFailoverPlan)
	r.POST("/failover/execute", h.ExecuteFailover)
	r.POST("/failover/test", h.RunDRTest)
	r.GET("/backups", h.ListBackups)
	r.POST("/recovery", h.RestoreFromBackup)
	r.GET("/health", h.GetHealthDashboard)
}

func (h *Handler) ListRegions(c *gin.Context) {
	regions, err := h.svc.ListRegions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": regions})
}

func (h *Handler) CreateFailoverPlan(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		Name         string `json:"name"`
		SourceRegion string `json:"source_region"`
		TargetRegion string `json:"target_region"`
		IsAutomated  bool   `json:"is_automated"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	result, err := h.svc.CreateFailoverPlan(c.Request.Context(), req.Name, req.SourceRegion, req.TargetRegion, req.IsAutomated)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) ExecuteFailover(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		PlanID string `json:"plan_id"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	planID, err := uuid.Parse(req.PlanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan_id"})
		return
	}
	result, err := h.svc.ExecuteFailover(c.Request.Context(), planID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) RunDRTest(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		PlanID   string `json:"plan_id"`
		TestName string `json:"test_name"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	planID, err := uuid.Parse(req.PlanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan_id"})
		return
	}
	result, err := h.svc.RunDRTest(c.Request.Context(), planID, req.TestName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) ListBackups(c *gin.Context) {
	manifests, err := h.svc.ListBackups(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": manifests})
}

func (h *Handler) RestoreFromBackup(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		ManifestID string `json:"manifest_id"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	manifestID, err := uuid.Parse(req.ManifestID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid manifest_id"})
		return
	}
	result, err := h.svc.RestoreFromBackup(c.Request.Context(), manifestID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetHealthDashboard(c *gin.Context) {
	dashboard, err := h.svc.GetHealthDashboard(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dashboard)
}
