package rest

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *UCREFHandler) {
	api := r.Group("/api/v1/ucref")
	{
		api.POST("/str", h.SubmitSTR)
		api.GET("/str/:id", h.GetSTRDetail)
		api.GET("/profile/:person_id", h.GetFinancialProfile)
		api.POST("/moncash/pattern", h.RecordMonCashPattern)
		api.GET("/str/unanalyzed", h.GetUnanalyzedSTRs)
		api.POST("/str/:id/disseminate", h.DisseminateSTR)
		api.GET("/gang-finances/:id", h.GetGangFinances)
	}
}
