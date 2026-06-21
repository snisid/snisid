package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/fir-svc/internal/api/rest"
	"github.com/snisid/platform/services/fir-svc/internal/service"
)

type Handler struct {
	record *rest.RecordHandler
	cert   *rest.CertificateHandler
	search *rest.SearchHandler
}

func NewHandler(recordSvc *service.RecordService, certSvc *service.CertificateService) *Handler {
	return &Handler{
		record: rest.NewRecordHandler(recordSvc),
		cert:   rest.NewCertificateHandler(certSvc),
		search: rest.NewSearchHandler(recordSvc),
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/records", h.record.CreateRecord)
	r.GET("/records/:person_id", h.record.GetRecord)
	r.POST("/records/:id/arrests", h.record.AddArrest)
	r.POST("/records/:id/convictions", h.record.AddConviction)
	r.GET("/records/:id/arrests", h.record.GetArrests)
	r.GET("/records/:id/convictions", h.record.GetConvictions)
	r.POST("/certificates/issue", h.cert.IssueCertificate)
	r.GET("/certificates/verify/:num", h.cert.VerifyCertificate)
	r.GET("/search", h.search.Search)
}
