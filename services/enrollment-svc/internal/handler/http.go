package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/snisid/enrollment-svc/internal/domain"
	"github.com/snisid/enrollment-svc/internal/service"
)

type Handler struct {
	svc *service.EnrollmentService
}

func NewHandler(svc *service.EnrollmentService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/requests", h.SubmitRequest)
	r.POST("/requests/:id/documents", h.UploadDocuments)
	r.POST("/requests/:id/biometrics", h.CaptureBiometrics)
	r.POST("/requests/:id/review", h.ReviewRequest)
	r.GET("/requests/:id", h.GetRequest)
	r.GET("/requests", h.ListPending)
}

func (h *Handler) SubmitRequest(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req domain.EnrollmentRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.SubmitRequest(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) UploadDocuments(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		Documents []domain.IdentityDocument `json:"documents"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.UploadDocuments(c.Request.Context(), id, req.Documents)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"documents": result})
}

func (h *Handler) CaptureBiometrics(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		Samples []domain.BiometricSample `json:"samples"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.CaptureBiometrics(c.Request.Context(), id, req.Samples)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"samples": result})
}

func (h *Handler) ReviewRequest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	body, _ := io.ReadAll(c.Request.Body)
	var review domain.EnrollmentReview
	if err := json.Unmarshal(body, &review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.ReviewRequest(c.Request.Context(), id, review)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetRequest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	req, err := h.svc.GetRequest(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h *Handler) ListPending(c *gin.Context) {
	reqs, err := h.svc.ListPending(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if reqs == nil {
		reqs = []domain.EnrollmentRequest{}
	}
	c.JSON(http.StatusOK, gin.H{"data": reqs})
}
