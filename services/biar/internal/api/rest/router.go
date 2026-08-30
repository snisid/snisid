package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/biar/internal/service"
)

func NewRouter(
	weaponSvc *service.WeaponService,
	batchSvc *service.BatchService,
	statsSvc *service.StatsService,
	syncSvc *service.SyncService,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(AuditMiddleware())

	weaponHandler := NewWeaponHandler(weaponSvc)
	batchHandler := NewBatchHandler(batchSvc)
	statsHandler := NewStatsHandler(statsSvc)
	syncHandler := NewSyncHandler(syncSvc)

	v1 := r.Group("/api/v1/biar")
	{
		v1.POST("/weapons", weaponHandler.Declare)
		v1.GET("/weapons/:id", weaponHandler.GetByID)
		v1.GET("/check/serial/:sn", weaponHandler.CheckSerial)

		v1.POST("/batches", batchHandler.Create)
		v1.GET("/batches/:id", batchHandler.GetByID)

		v1.GET("/stats/by-gang", statsHandler.ByGang)
		v1.GET("/stats/by-origin", statsHandler.ByOrigin)
		v1.GET("/stats/routes", statsHandler.Routes)

		v1.POST("/iarms/sync", syncHandler.SyncIARMS)
	}

	return r
}
