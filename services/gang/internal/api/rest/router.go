package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/gang/internal/service"
)

func NewRouter(
	gangSvc *service.GangService,
	memberSvc *service.MemberService,
	incidentSvc *service.IncidentService,
	territorySvc *service.TerritoryService,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(AuditMiddleware())

	gangHandler := NewGangHandler(gangSvc)
	memberHandler := NewMemberHandler(memberSvc)
	incidentHandler := NewIncidentHandler(incidentSvc)
	territoryHandler := NewTerritoryHandler(territorySvc)

	v1 := r.Group("/api/v1")
	{
		gangs := v1.Group("/gangs")
		{
			gangs.POST("", gangHandler.Create)
			gangs.GET("", gangHandler.List)
			gangs.GET("/:id", gangHandler.GetByID)
			gangs.GET("/by-dept/:code", gangHandler.ByDeptCode)
			gangs.GET("/sanctioned", gangHandler.Sanctioned)
		}

		v1.POST("/gangs/:id/incidents", incidentHandler.Create)
		v1.GET("/gangs/:id/incidents", incidentHandler.ListByGang)
		v1.GET("/gangs/:id/members", memberHandler.ListByGang)
		v1.POST("/gangs/:id/members", memberHandler.Create)
		v1.GET("/gangs/:id/territories", territoryHandler.ListByGang)
		v1.POST("/gangs/:id/territories", territoryHandler.Create)
	}

	return r
}
