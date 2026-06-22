package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/snisid/fips-cert-svc/internal/domain"
	"github.com/snisid/fips-cert-svc/internal/service"
)

type Handler struct {
	svc *service.FIPSService
}

func NewHandler(svc *service.FIPSService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/modules", h.RegisterModule)
	r.GET("/modules", h.ListModules)
	r.POST("/modules/:id/validate", h.SubmitValidation)
	r.GET("/modules/:id", h.GetModule)
	r.POST("/modules/:id/cve", h.ReportCVE)
	r.GET("/compliance/:service", h.GetCompliance)
	r.GET("/dashboard", h.GetDashboard)
}

func (h *Handler) RegisterModule(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		Name       string   `json:"name"`
		Version    string   `json:"version"`
		Vendor     string   `json:"vendor"`
		FIPSLevel  string   `json:"fips_level"`
		Algorithms []string `json:"algorithms"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	algos := make([]domain.CertAlgo, len(req.Algorithms))
	for i, a := range req.Algorithms {
		algos[i] = domain.CertAlgo(a)
	}

	mod := domain.CryptoModule{
		Name:       req.Name,
		Version:    req.Version,
		Vendor:     req.Vendor,
		FIPSLevel:  domain.FIPSLevel(req.FIPSLevel),
		Algorithms: algos,
	}

	result, err := h.svc.RegisterModule(c.Request.Context(), mod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) ListModules(c *gin.Context) {
	modules, err := h.svc.ListModules(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": modules})
}

func (h *Handler) GetModule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid module id"})
		return
	}

	mod, err := h.svc.GetModule(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "module not found"})
		return
	}
	c.JSON(http.StatusOK, mod)
}

func (h *Handler) SubmitValidation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid module id"})
		return
	}

	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		CertNumber     string `json:"cert_number"`
		ValidationDate string `json:"validation_date"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	valDate, err := time.Parse("2006-01-02", req.ValidationDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid validation_date, use YYYY-MM-DD"})
		return
	}

	result, err := h.svc.SubmitValidation(c.Request.Context(), id, req.CertNumber, valDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) ReportCVE(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid module id"})
		return
	}

	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		CVEID    string  `json:"cve_id"`
		Severity string  `json:"severity"`
		Notes    *string `json:"notes,omitempty"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.ReportCVE(c.Request.Context(), id, req.CVEID, req.Severity, req.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetCompliance(c *gin.Context) {
	service := c.Param("service")
	report, err := h.svc.GetComplianceByService(c.Request.Context(), service)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "compliance report not available"})
		return
	}
	c.JSON(http.StatusOK, report)
}

func (h *Handler) GetDashboard(c *gin.Context) {
	reports, err := h.svc.GetDashboard(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": reports})
}
