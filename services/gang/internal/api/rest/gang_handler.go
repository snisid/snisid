package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/gang/internal/domain"
	"github.com/snisid/platform/services/gang/internal/service"
)

type GangHandler struct {
	svc *service.GangService
}

func NewGangHandler(svc *service.GangService) *GangHandler {
	return &GangHandler{svc: svc}
}

func (h *GangHandler) Create(c *gin.Context) {
	var req domain.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "requête invalide"})
		return
	}
	createdBy := uuid.Nil
	if uid, ok := c.Request.Context().Value(ContextKeyUserID).(uuid.UUID); ok {
		createdBy = uid
	}
	org, err := h.svc.CreateOrganization(c.Request.Context(), req, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, org)
}

func (h *GangHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	org, err := h.svc.GetOrganization(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "organisation introuvable"})
		return
	}
	c.JSON(http.StatusOK, org)
}

func (h *GangHandler) List(c *gin.Context) {
	orgs, err := h.svc.ListOrganizations(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur requête"})
		return
	}
	c.JSON(http.StatusOK, orgs)
}

func (h *GangHandler) ByDeptCode(c *gin.Context) {
	code := c.Param("code")
	orgs, err := h.svc.ByDeptCode(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur requête"})
		return
	}
	c.JSON(http.StatusOK, orgs)
}

func (h *GangHandler) Sanctioned(c *gin.Context) {
	orgs, err := h.svc.Sanctioned(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur requête"})
		return
	}
	c.JSON(http.StatusOK, orgs)
}
