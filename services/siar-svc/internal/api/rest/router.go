package rest

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/siar-svc/internal/service"
)

func SetupRouter(svc *service.SIARService, log *zap.Logger) *gin.Engine {
	r := gin.Default()
	handler := NewSIARHandler(svc, log)

	api := r.Group("/api/v1/siar")
	{
		api.POST("/firearms", handler.RegisterFirearm)
		api.GET("/firearms/:id", handler.GetFirearm)
		api.GET("/check/serial/:sn", handler.CheckSerial)
		api.POST("/seizures", handler.ReportSeizure)
		api.POST("/stolen", handler.ReportStolen)
		api.GET("/licenses/:person", handler.GetLicensesByPerson)
		api.POST("/licenses", handler.CreateLicense)
		api.GET("/stats/by-type", handler.GetStatsByType)
	}

	return r
}
