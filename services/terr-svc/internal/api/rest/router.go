package rest

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/terr-svc/internal/service"
)

func SetupRouter(handler *TerrHandler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1/terr")
	{
		api.GET("/check", handler.CheckPointSafety)
		api.GET("/route-safety", handler.GetRouteSafety)
		api.GET("/zones", handler.ListZones)
		api.GET("/zones/dept/:code", handler.ListZonesByDept)
		api.POST("/zones", handler.CreateZone)
		api.GET("/zones/:gang_id", handler.ListZonesByGang)
		api.POST("/checkpoints", handler.CreateCheckpoint)
		api.GET("/history/:zone_id", handler.GetZoneHistory)
	}

	return r
}

func NewRouter(svc *service.TerritoryService, log *zap.Logger) *gin.Engine {
	handler := NewTerrHandler(svc, log)
	return SetupRouter(handler)
}
