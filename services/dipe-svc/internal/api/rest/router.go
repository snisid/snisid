package rest

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/dipe-svc/internal/service"
)

func NewRouter(svc *service.MissingPersonService, log *zap.Logger) *gin.Engine {
	r := gin.Default()
	handler := NewDipeHandler(svc, log)

	api := r.Group("/api/v1/dipe")
	{
		api.POST("/cases", handler.ReportDisappearance)
		api.GET("/cases/:id", handler.GetCase)
		api.GET("/cases/open", handler.GetOpenCases)
		api.POST("/cases/:id/sightings", handler.AddSighting)
		api.GET("/match/rvin/:id", handler.MatchWithRVIN)
		api.PATCH("/cases/:id/resolve", handler.ResolveCase)
		api.GET("/stats/by-type", handler.GetStatsByType)
		api.GET("/hotline/tips", handler.GetHotlineTips)
	}

	return r
}
