package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/rdep-svc/internal/api/rest"
	"github.com/snisid/platform/services/rdep-svc/internal/service"
)

type Handler struct {
	intake     *rest.IntakeHandler
	monitoring *rest.MonitoringHandler
}

func NewHandler(intakeSvc *service.IntakeService, screeningSvc *service.ScreeningService, monitoringSvc *service.MonitoringService) *Handler {
	return &Handler{
		intake:     rest.NewIntakeHandler(intakeSvc, screeningSvc),
		monitoring: rest.NewMonitoringHandler(monitoringSvc),
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/intake", h.intake.ProcessIntake)
	r.POST("/:id/screen", h.intake.ScreenDeportee)
	r.GET("/:id", h.intake.GetDeportee)
	r.GET("/high-risk", h.intake.GetHighRisk)
	r.GET("/gang-affiliated", h.intake.GetGangAffiliated)
	r.GET("/stats/by-country", h.intake.GetStatsByCountry)
	r.POST("/:id/monitoring/events", h.monitoring.RecordEvent)
	r.GET("/:id/monitoring/events", h.monitoring.GetEvents)
	r.PUT("/:id/address", h.monitoring.UpdateAddress)
}
