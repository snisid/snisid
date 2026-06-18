package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/siar/internal/domain"
	"github.com/snisid/platform/services/siar/internal/service"
)

type LicenseHandler struct {
	svc *service.LicenseService
}

func NewLicenseHandler(svc *service.LicenseService) *LicenseHandler {
	return &LicenseHandler{svc: svc}
}

func (h *LicenseHandler) Create(c *gin.Context) {
	var req domain.CreateLicenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "requête invalide: " + err.Error()})
		return
	}
	createdBy := uuid.Nil
	if uid, ok := c.Request.Context().Value(ContextKeyUserID).(uuid.UUID); ok {
		createdBy = uid
	}
	lic, err := h.svc.Create(c.Request.Context(), req, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, lic)
}

func (h *LicenseHandler) GetByPerson(c *gin.Context) {
	personID, err := uuid.Parse(c.Param("person"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID personne invalide"})
		return
	}
	licenses, err := h.svc.GetByPerson(c.Request.Context(), personID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur requête"})
		return
	}
	if licenses == nil {
		licenses = []*domain.License{}
	}
	c.JSON(http.StatusOK, gin.H{"licenses": licenses})
}
