package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/snisid/employee-verify-svc/internal/service"
)

type Handler struct {
	svc *service.EmployeeVerifyService
}

func NewHandler(svc *service.EmployeeVerifyService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/employers/register", h.RegisterEmployer)
	r.POST("/cases", h.CreateCase)
	r.GET("/cases/:tcn", h.GetCase)
	r.POST("/cases/:tcn/verify", h.SubmitVerification)
	r.GET("/cases/employer/:ein", h.ListCasesByEmployer)
	r.GET("/stats", h.GetStats)
}

func (h *Handler) RegisterEmployer(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		CompanyName  string `json:"company_name"`
		EIN          string `json:"ein"`
		Address      string `json:"address"`
		ContactEmail string `json:"contact_email"`
		ContactPhone string `json:"contact_phone"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	result, err := h.svc.RegisterEmployer(c.Request.Context(), req.CompanyName, req.EIN, req.Address, req.ContactEmail, req.ContactPhone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) CreateCase(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		EmployerID     string `json:"employer_id"`
		EmployeeName   string `json:"employee_name"`
		DocumentNumber string `json:"document_number"`
		DocumentType   string `json:"document_type"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	employerID, err := uuid.Parse(req.EmployerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employer_id"})
		return
	}
	result, err := h.svc.CreateCase(c.Request.Context(), employerID, req.EmployeeName, req.DocumentNumber, req.DocumentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetCase(c *gin.Context) {
	tcn := c.Param("tcn")
	result, err := h.svc.GetCaseByTCN(c.Request.Context(), tcn)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "case not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) SubmitVerification(c *gin.Context) {
	tcn := c.Param("tcn")
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		SSAMatch bool   `json:"ssa_match"`
		DHSMatch bool   `json:"dhs_match"`
		Reason   string `json:"reason"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	result, err := h.svc.SubmitVerificationResponse(c.Request.Context(), tcn, req.SSAMatch, req.DHSMatch, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) ListCasesByEmployer(c *gin.Context) {
	ein := c.Param("ein")
	results, err := h.svc.ListCasesByEmployer(c.Request.Context(), ein)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results})
}

func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.svc.GetStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
