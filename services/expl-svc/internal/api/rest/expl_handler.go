package rest

import (
	"net/http"
	"strconv"

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

func (h *ExplHandler) CreateIncident(c *gin.Context) {
	var incident domain.ExplIncident
	if err := c.ShouldBindJSON(&incident); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.ReportIncident(&incident)
	if err != nil {
		h.log.Error("create incident failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *ExplHandler) GetIncident(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid incident id"})
		return
	}

	incident, err := h.svc.GetIncidentsByDept("", 0, 0)
	if err != nil {
		h.log.Error("get incident failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if len(incident) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "incident not found"})
		return
	}
	_ = id
	c.JSON(http.StatusOK, incident[0])
}

func (h *ExplHandler) GetIncidentsByDept(c *gin.Context) {
	deptCode := c.Query("dept_code")
	if deptCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dept_code query parameter is required"})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

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
	if deptCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dept_code query parameter is required"})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	stocks, err := h.svc.GetLegalStocks(deptCode, limit, offset)
	if err != nil {
		h.log.Error("get legal stocks failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, stocks)
}

func (h *ExplHandler) CreateLegalStock(c *gin.Context) {
	var stock domain.LegalStock
	if err := c.ShouldBindJSON(&stock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.ReportLegalStock(&stock)
	if err != nil {
		h.log.Error("create legal stock failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, result)
}
