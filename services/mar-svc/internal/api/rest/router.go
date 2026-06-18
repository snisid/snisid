package rest

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/mar-svc/internal/service"
)

func SetupRouter(handler *MARHandler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1/mar")
	{
		api.GET("/vessels/:id", handler.GetVessel)
		api.POST("/vessels", handler.CreateVessel)
		api.POST("/incidents", handler.CreateIncident)
		api.GET("/incidents/recent", handler.GetRecentIncidents)
		api.POST("/watch", handler.CreateWatch)
		api.GET("/watch/active", handler.GetActiveWatches)
		api.GET("/ais/live", handler.GetAISLive)
		api.GET("/zones/:zone/activity", handler.GetZoneActivity)
		api.GET("/stats/incidents", handler.GetIncidentStats)
	}

	return r
}

func NewRouter(svc *service.MaritimeService, log *zap.Logger) *gin.Engine {
	handler := NewMARHandler(svc, log)
	return SetupRouter(handler)
}
