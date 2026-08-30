package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/opr-svc/internal/api/rest"
	"github.com/snisid/platform/services/opr-svc/internal/service"
)

type Handler struct {
	opr *rest.OPRHandler
}

func NewHandler(oprSvc *service.OPRService) *Handler {
	return &Handler{
		opr: rest.NewOPRHandler(oprSvc),
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/orders", h.opr.CreateOrder)
	r.GET("/check/:person_id", h.opr.CheckSubject)
	r.POST("/violations", h.opr.RecordViolation)
	r.GET("/expiring-soon", h.opr.GetExpiringSoon)
	r.GET("/orders/by-gang/:id", h.opr.GetByGangID)
}
