package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/gang-svc/internal/api/rest"
	"github.com/snisid/platform/services/gang-svc/internal/service"
)

type Handler struct {
	org *rest.OrganizationHandler
}

func NewHandler(orgSvc *service.OrganizationService, incSvc *service.IncidentService, alliSvc *service.AllianceService) *Handler {
	return &Handler{
		org: rest.NewOrganizationHandler(orgSvc, incSvc, alliSvc),
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("", h.org.CreateOrganization)
	r.GET("", h.org.ListOrganizations)
	r.GET("/:id", h.org.GetOrganization)
	r.POST("/:id/incidents", h.org.CreateIncident)
	r.GET("/:id/incidents", h.org.GetIncidentsByGang)
	r.GET("/by-dept/:code", h.org.GetByDeptCode)
	r.GET("/alliances/map", h.org.GetAllianceMap)
	r.GET("/sanctioned", h.org.GetSanctioned)
}
