package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/afis-svc/internal/api/rest"
	"github.com/snisid/platform/services/afis-svc/internal/service"
)

type Handler struct {
	enroll  *rest.EnrollHandler
	search  *rest.SearchHandler
	latent  *rest.LatentHandler
	quality *rest.QualityHandler
}

func NewHandler(enrollSvc *service.EnrollmentService, searchSvc *service.SearchService, latentSvc *service.LatentService, qualitySvc *service.QualityService) *Handler {
	return &Handler{
		enroll:  rest.NewEnrollHandler(enrollSvc),
		search:  rest.NewSearchHandler(searchSvc),
		latent:  rest.NewLatentHandler(latentSvc),
		quality: rest.NewQualityHandler(qualitySvc),
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/enroll", h.enroll.Enroll)
	r.POST("/search/tenprint", h.search.SearchTenprint)
	r.POST("/search/latent", h.latent.SearchLatent)
	r.GET("/subjects/:id", h.enroll.GetSubject)
	r.POST("/latents", h.latent.Submit)
	r.PATCH("/latents/:id/match", h.latent.ConfirmMatch)
	r.GET("/quality/check", h.quality.CheckQuality)
	r.GET("/stats", h.quality.Stats)
}
