package rest

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(handler *BLANHandler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1/blan")
	{
		api.POST("/cases", handler.OpenCase)
		api.GET("/cases/:id", handler.GetCaseDetail)
		api.POST("/cases/:id/assets", handler.AddSuspiciousAsset)
		api.POST("/cases/:id/chain", handler.DocumentTransactionChain)
		api.GET("/real-estate/flagged", handler.GetFlaggedRealEstate)
		api.GET("/assets/frozen", handler.GetFrozenAssets)
		api.GET("/stats/by-typology", handler.GetStatsByTypology)
	}

	return r
}
