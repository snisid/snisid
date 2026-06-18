package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/afis-svc/internal/domain"
	"github.com/snisid/platform/services/afis-svc/internal/service"
)

type EnrollHandler struct {
	svc *service.EnrollmentService
}

func NewEnrollHandler(svc *service.EnrollmentService) *EnrollHandler {
	return &EnrollHandler{svc: svc}
}

func (h *EnrollHandler) Enroll(c *gin.Context) {
	var req domain.EnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	officerID, _ := uuid.Parse(c.GetString("user_id"))
	profile, err := h.svc.Enroll(c.Request.Context(), req, officerID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ginH(http.StatusCreated, "Sujet enrôlé avec succès", profile))
}

func (h *EnrollHandler) GetSubject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subject_id": id})
}

func ginH(code int, message string, data interface{}) gin.H {
	return gin.H{"code": code, "message": message, "data": data}
}
