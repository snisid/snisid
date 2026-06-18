package rest

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *SifrHandler) {
	sifr := r.Group("/sifr")
	{
		sifr.POST("/crossings", h.ProcessCrossing)
		sifr.GET("/crossings/search", h.SearchCrossings)
		sifr.GET("/crossings/person/:id", h.GetPersonHistory)
		sifr.GET("/alerts/active", h.GetActiveAlerts)
		sifr.GET("/posts", h.ListPosts)
		sifr.GET("/stats/daily", h.GetDailyStats)
		sifr.POST("/clandestine", h.ReportClandestineCrossing)
	}
}
