package rest

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/expl-svc/internal/service"
)

func SetupRouter(svc *service.ExplService, logger *zap.Logger) *gin.Engine {
	h := NewExplHandler(svc, logger)
	r := gin.Default()

	v1 := r.Group("/api/v1/expl")
	{
		v1.POST("/incidents", h.CreateIncident)
		v1.GET("/incidents/:id", h.GetIncident)
		v1.GET("/incidents/by-dept", h.GetIncidentsByDept)
		v1.GET("/legal-stocks", h.GetLegalStocks)
		v1.POST("/legal-stocks", h.CreateLegalStock)
	}

	return r
}
