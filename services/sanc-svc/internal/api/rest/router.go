package rest

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/sanc-svc/internal/service"
)

func SetupRouter(svc *service.SanctionsService, log *zap.Logger) *gin.Engine {
	handler := NewSancHandler(svc, log)
	r := gin.Default()

	api := r.Group("/api/v1/sanc")
	{
		api.GET("/check/:person_id", handler.CheckPerson)
		api.POST("/check/name", handler.CheckByName)
		api.GET("/entries", handler.GetEntries)
		api.GET("/entries/haiti", handler.GetHaitiEntries)
		api.GET("/matches/unconfirmed", handler.GetUnconfirmedMatches)
		api.POST("/matches/:id/confirm", handler.ConfirmMatch)
		api.POST("/sync/trigger", handler.TriggerSync)
		api.GET("/sync/status", handler.GetSyncStatus)
	}

	return r
}
