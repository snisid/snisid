package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/expl-svc/internal/domain"
	"github.com/snisid/platform/services/expl-svc/internal/service"
)

type ExplHandler struct {
	svc *service.ExplService
	log *zap.Logger
}

func NewExplHandler(svc *service.ExplService, log *zap.Logger) *ExplHandler {
	return &ExplHandler{svc: svc, log: log}
}

func (h *ExplHandler) ReportIncident(c *gin.Context) {
	var incident domain.ExplIncident
	if err := c.ShouldBindJSON(&incident); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.ReportIncident(&incident)
	if err != nil {
		h.log.Error("report incident failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *ExplHandler) GetIncident(c *gin.Context) {
	idStr := c.Param("id")
	_, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid incident id"})
		return
	}
	incident, err := h.svc.GetIncidentsByDept("", 1, 0)
	if err != nil {
		h.log.Error("get incident failed", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "incident not found"})
		return
	}
	c.JSON(http.StatusOK, incident)
}

func (h *ExplHandler) GetIncidentsByDept(c *gin.Context) {
	deptCode := c.Query("dept_code")
	limit := 50
	offset := 0
	incidents, err := h.svc.GetIncidentsByDept(deptCode, limit, offset)
	if err != nil {
		h.log.Error("get incidents by dept failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, incidents)
}

func (h *ExplHandler) GetLegalStocks(c *gin.Context) {
	deptCode := c.Query("dept_code")
	limit := 50
	offset := 0
	stocks, err := h.svc.GetLegalStocks(deptCode, limit, offset)
	if err != nil {
		h.log.Error("get legal stocks failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, stocks)
}

func (h *ExplHandler) ReportLegalStock(c *gin.Context) {
	var stock domain.LegalStock
	if err := c.ShouldBindJSON(&stock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.ReportLegalStock(&stock)
	if err != nil {
		h.log.Error("report legal stock failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func SetupRouter(svc *service.ExplService, log *zap.Logger) *gin.Engine {
	r := gin.Default()
	handler := NewExplHandler(svc, log)

	api := r.Group("/api/v1/expl")
	{
		api.POST("/incidents", handler.ReportIncident)
		api.GET("/incidents/:id", handler.GetIncident)
		api.GET("/incidents/by-dept", handler.GetIncidentsByDept)
		api.GET("/legal-stocks", handler.GetLegalStocks)
		api.POST("/legal-stocks", handler.ReportLegalStock)
	}
	return r
}
