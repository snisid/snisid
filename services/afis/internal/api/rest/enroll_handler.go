package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/afis/internal/domain"
	"github.com/snisid/platform/services/afis/internal/service"
)

type EnrollHandler struct {
	enrollment *service.EnrollmentService
	search     *service.SearchService
}

func NewEnrollHandler(e *service.EnrollmentService, s *service.SearchService) *EnrollHandler {
	return &EnrollHandler{enrollment: e, search: s}
}

func (h *EnrollHandler) Enroll(c *gin.Context) {
	var req domain.EnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	officerIDStr := c.GetHeader("X-User-ID")
	if officerIDStr == "" {
		officerIDStr = c.GetString("user_id")
	}
	if officerIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "en-tête X-User-ID requis"})
		return
	}

	officerID, err := uuid.Parse(officerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-ID invalide"})
		return
	}

	subject, fps, err := h.enrollment.Enroll(c.Request.Context(), req, officerID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	h.search.IndexSubject(subject, fps)

	c.JSON(http.StatusCreated, gin.H{
		"subject":      subject,
		"fingerprints": fps,
	})
}

func (h *EnrollHandler) GetSubject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	subject, err := h.enrollment.GetSubject(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	fps, err := h.enrollment.GetFingerprints(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"subject":      subject,
		"fingerprints": fps,
	})
}

func (h *EnrollHandler) GetSubjectHistory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	subject, err := h.enrollment.GetSubject(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"subject_id": subject.SubjectID,
		"history":    []interface{}{},
	})
}

func (h *EnrollHandler) GetStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"total_subjects":     0,
		"total_fingerprints": 0,
		"total_latents":      0,
		"total_searches":     0,
	})
}
