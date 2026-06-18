package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/gang-svc/internal/service"
)

type OrganizationHandler struct {
	orgSvc      *service.OrganizationService
	incidentSvc *service.IncidentService
	allianceSvc *service.AllianceService
}

func NewOrganizationHandler(
	orgSvc *service.OrganizationService,
	incidentSvc *service.IncidentService,
	allianceSvc *service.AllianceService,
) *OrganizationHandler {
	return &OrganizationHandler{
		orgSvc:      orgSvc,
		incidentSvc: incidentSvc,
		allianceSvc: allianceSvc,
	}
}

func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	var req service.CreateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, err := h.orgSvc.CreateOrganization(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, org)
}

func (h *OrganizationHandler) GetOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID invalide"})
		return
	}

	org, err := h.orgSvc.GetOrganization(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organisation non trouvée"})
		return
	}

	c.JSON(http.StatusOK, org)
}

func (h *OrganizationHandler) ListOrganizations(c *gin.Context) {
	orgs, err := h.orgSvc.ListOrganizations(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"organizations": orgs})
}

func (h *OrganizationHandler) GetByDeptCode(c *gin.Context) {
	code := c.Param("code")
	orgs, err := h.orgSvc.GetByDeptCode(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"organizations": orgs})
}

func (h *OrganizationHandler) GetSanctioned(c *gin.Context) {
	orgs, err := h.orgSvc.GetSanctioned(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"organizations": orgs})
}

func (h *OrganizationHandler) CreateIncident(c *gin.Context) {
	var req service.CreateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	incident, err := h.incidentSvc.CreateIncident(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, incident)
}

func (h *OrganizationHandler) GetIncidentsByGang(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID invalide"})
		return
	}

	incidents, err := h.incidentSvc.GetByGangID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"incidents": incidents})
}

func (h *OrganizationHandler) GetAllianceMap(c *gin.Context) {
	alliances, err := h.allianceSvc.GetAllianceMap(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"alliances": alliances})
}
