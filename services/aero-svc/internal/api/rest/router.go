package rest

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/aero-svc/internal/service"
)

func SetupRouter(svc *service.AeroService, log *zap.Logger) *gin.Engine {
	handler := NewAeroHandler(svc, log)

	r := gin.Default()

	api := r.Group("/api/v1/aero")
	{
		api.GET("/check/:reg", handler.CheckRegistration)
		api.POST("/strips", handler.ReportStrip)
		api.GET("/strips/map", handler.GetStripMap)
		api.POST("/flights/suspicious", handler.ReportSuspiciousFlight)
		api.GET("/stats/strips", handler.GetStripStats)
	}

	return r
}
