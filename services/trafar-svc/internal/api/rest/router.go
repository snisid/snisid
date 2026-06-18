package rest

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/trafar-svc/internal/service"
)

func RegisterRoutes(r *gin.Engine, svc *service.TrafarService, log *zap.Logger) {
	handler := NewTrafarHandler(svc, log)

	api := r.Group("/api/v1/trafar")
	{
		api.GET("/routes", handler.GetRoutes)
		api.GET("/routes/:id", handler.GetRoute)
		api.POST("/routes", handler.CreateRoute)
		api.POST("/shipments", handler.RecordShipment)
		api.GET("/map", handler.GetRoutesMap)
		api.GET("/stats/by-origin", handler.GetStatsByOrigin)
		api.GET("/suppliers", handler.GetSuppliers)
	}
}
