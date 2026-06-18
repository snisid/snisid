package rest

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/sltd-svc/internal/service"
)

func NewRouter(svc *service.SLTDService, log *zap.Logger) *gin.Engine {
	handler := NewSLTDHandler(svc, log)

	r := gin.Default()

	api := r.Group("/api/v1/sltd")
	{
		api.GET("/check/:num", handler.CheckDocument)
		api.POST("/report/lost", handler.ReportLost)
		api.POST("/report/stolen", handler.ReportStolen)
		api.PATCH("/:id/found", handler.MarkFound)
		api.GET("/stats", handler.GetStats)
	}

	return r
}
