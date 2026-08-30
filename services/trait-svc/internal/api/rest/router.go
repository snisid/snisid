package rest

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/trait-svc/internal/service"
)

func SetupRouter(svc *service.TraiffickingService, log *zap.Logger) *gin.Engine {
	handler := NewTraiffickingHandler(svc, log)

	r := gin.Default()

	api := r.Group("/api/v1/trait")
	{
		api.POST("/cases", handler.OpenCase)
		api.GET("/cases/:id", handler.GetCase)
		api.POST("/cases/:id/victims", handler.AddVictim)
		api.GET("/victims/minors", handler.GetMinorVictims)
		api.POST("/networks", handler.DocumentNetwork)
		api.GET("/networks/active", handler.GetActiveNetworks)
		api.GET("/stats/by-type", handler.GetStatsByType)
		api.GET("/cases/maritime", handler.GetMaritimeCases)
	}

	return r
}
