package rest

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/reso-svc/internal/service"
)

func SetupRouter(svc *service.ResoService, log *zap.Logger) *gin.Engine {
	r := gin.Default()
	handler := NewResoHandler(svc, log)

	api := r.Group("/api/v1/reso")
	{
		api.GET("/network/:person_id", handler.GetPersonNetwork)
		api.GET("/communities", handler.GetCommunities)
		api.GET("/key-actors", handler.GetKeyActors)
		api.GET("/gang-overlap/:g1/:g2", handler.GetGangOverlap)
		api.POST("/analyze/trigger", handler.TriggerAnalysis)
		api.GET("/path/:from_id/:to_id", handler.FindShortestPath)
		api.GET("/centrality-scores", handler.GetCentralityScores)
		api.GET("/emerging-links", handler.GetEmergingLinks)
	}

	return r
}
