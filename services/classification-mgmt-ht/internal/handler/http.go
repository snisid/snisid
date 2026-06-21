package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snisid/classification-mgmt-ht/internal/domain"
	"github.com/snisid/classification-mgmt-ht/internal/service"
)

type Handler struct {
	svc *service.ClassificationService
}

func NewHandler(svc *service.ClassificationService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/rules", h.CreateRule)
	r.GET("/rules/:data_type", h.GetRulesByDataType)
	r.POST("/tags", h.TagResource)
	r.GET("/tags/check", h.GetClassificationByURI)
	r.POST("/audit", h.LogAudit)
	r.GET("/audit/recent", h.GetRecentAuditLogs)
	r.GET("/dashboard", h.GetDashboard)
}

func (h *Handler) CreateRule(c *gin.Context) {
	var req domain.CreateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.CreateRule(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetRulesByDataType(c *gin.Context) {
	dataType := c.Param("data_type")
	if dataType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data_type is required"})
		return
	}
	rules, err := h.svc.GetRulesByDataType(c.Request.Context(), dataType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rules)
}

func (h *Handler) TagResource(c *gin.Context) {
	var req domain.TagResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.TagResource(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetClassificationByURI(c *gin.Context) {
	uri := c.Query("uri")
	if uri == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uri query parameter is required"})
		return
	}
	tag, err := h.svc.GetClassificationByURI(c.Request.Context(), uri)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tag)
}

func (h *Handler) LogAudit(c *gin.Context) {
	var req domain.LogAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.LogAudit(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetRecentAuditLogs(c *gin.Context) {
	logs, err := h.svc.GetRecentAuditLogs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}

func (h *Handler) GetDashboard(c *gin.Context) {
	stats, err := h.svc.GetDashboard(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
