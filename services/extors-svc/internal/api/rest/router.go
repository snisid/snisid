package rest

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/extors-svc/internal/service"
)

func SetupRouter(handler *ExtorsHandler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1/extors")
	{
		api.POST("/cases", handler.OpenCase)
		api.GET("/cases/:id", handler.GetCaseDetail)
		api.POST("/cases/:id/negotiations", handler.AddNegotiation)
		api.POST("/toll-points", handler.CreateTollPoint)
		api.GET("/toll-points/map", handler.GetTollsMap)
		api.GET("/gang/:id/revenue", handler.GetGangRevenue)
		api.GET("/stats/by-type", handler.GetStatsByType)
		api.GET("/moncash/patterns", handler.GetMoncashPatterns)
	}

	return r
}

func NewRouter(svc *service.ExtorsService, log *zap.Logger) *gin.Engine {
	handler := NewExtorsHandler(svc, log)
	return SetupRouter(handler)
}
