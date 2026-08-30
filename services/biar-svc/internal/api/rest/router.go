package rest

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/biar-svc/internal/service"
)

func SetupRouter(svc *service.BIARService, log *zap.Logger) *gin.Engine {
	handler := NewBIARHandler(svc, log)

	r := gin.Default()

	api := r.Group("/api/v1/biar")
	{
		api.POST("/weapons", handler.ReportWeapon)
		api.GET("/weapons/:id", handler.GetWeapon)
		api.GET("/check/serial/:sn", handler.CheckSerial)
		api.POST("/batches", handler.ReportBatch)
		api.GET("/stats/by-gang", handler.GetStatsByGang)
		api.GET("/stats/by-origin", handler.GetStatsByOrigin)
		api.GET("/stats/routes", handler.GetRoutes)
		api.POST("/iarms/sync", handler.SyncFromIARMS)
	}

	return r
}
