package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/sipep/internal/domain"
	"github.com/snisid/platform/services/sipep/internal/service"
)

type FacilityHandler struct {
	facilityService *service.FacilityService
	inmateService   *service.InmateService
}

func NewFacilityHandler(fs *service.FacilityService, is *service.InmateService) *FacilityHandler {
	return &FacilityHandler{facilityService: fs, inmateService: is}
}

func (h *FacilityHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/facilities", h.ListFacilities)
	rg.GET("/facilities/occupancy", h.Occupancy)
	rg.GET("/alerts/overcrowding", h.OvercrowdingAlerts)
	rg.GET("/stats/preventive-detention", h.PreventiveDetentionStats)
}

func (h *FacilityHandler) ListFacilities(c *gin.Context) {
	facilities := h.facilityService.GetAll()
	c.JSON(http.StatusOK, gin.H{"facilities": facilities})
}

func (h *FacilityHandler) Occupancy(c *gin.Context) {
	allInmates, _ := h.inmateService.Search("")
	facilityCounts := make(map[string]int)
	for _, inmate := range allInmates {
		if inmate.IsCurrentlyDetained {
			facilityCounts[inmate.CurrentFacility]++
		}
	}
	reports := h.facilityService.GetOccupancy(facilityCounts)
	c.JSON(http.StatusOK, gin.H{"occupancy": reports})
}

func (h *FacilityHandler) OvercrowdingAlerts(c *gin.Context) {
	if h == nil || h.facilityService == nil || h.inmateService == nil {
		c.JSON(http.StatusOK, gin.H{"alerts": []domain.OccupancyReport{}})
		return
	}
	allInmates, _ := h.inmateService.Search("")
	facilityCounts := make(map[string]int)
	for _, inmate := range allInmates {
		if inmate.IsCurrentlyDetained {
			facilityCounts[inmate.CurrentFacility]++
		}
	}
	reports := h.facilityService.GetOccupancy(facilityCounts)
	var alerts []domain.OccupancyReport
	for _, r := range reports {
		if r.OccupancyRate > 1.5 {
			alerts = append(alerts, *r)
		}
	}
	if alerts == nil {
		alerts = []domain.OccupancyReport{}
	}
	c.JSON(http.StatusOK, gin.H{"alerts": alerts})
}

func (h *FacilityHandler) PreventiveDetentionStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "not implemented"})
}
