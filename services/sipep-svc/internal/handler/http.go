package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/sipep-svc/internal/api/rest"
	"github.com/snisid/platform/services/sipep-svc/internal/service"
)

type Handler struct {
	intake *rest.IntakeHandler
	inmate *rest.InmateHandler
	stats  *rest.StatsHandler
}

func NewHandler(intakeSvc *service.IntakeService, releaseSvc *service.ReleaseService, transferSvc *service.TransferService, overcrowdingSvc *service.OvercrowdingService) *Handler {
	return &Handler{
		intake: rest.NewIntakeHandler(intakeSvc),
		inmate: rest.NewInmateHandler(intakeSvc, releaseSvc, transferSvc),
		stats:  rest.NewStatsHandler(overcrowdingSvc),
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/intake", h.intake.ProcessIntake)
	r.GET("/inmates/:id", h.inmate.GetInmate)
	r.GET("/inmates/search", h.inmate.SearchInmates)
	r.POST("/release", h.inmate.ProcessRelease)
	r.POST("/transfers", h.inmate.ProcessTransfer)
	r.GET("/inmates/:id/transfers", h.inmate.GetTransfers)
	r.GET("/facilities/occupancy", h.stats.GetFacilityOccupancy)
	r.GET("/alerts/overcrowding", h.stats.GetOvercrowdingAlerts)
	r.GET("/stats/preventive-detention", h.stats.GetPreventiveDetentionStats)
}
