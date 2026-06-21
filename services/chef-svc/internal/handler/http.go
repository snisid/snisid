package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/chef-svc/internal/api/rest"
	"github.com/snisid/platform/services/chef-svc/internal/service"
)

type Handler struct {
	member *rest.MemberHandler
}

func NewHandler(svc *service.MemberService) *Handler {
	return &Handler{
		member: rest.NewMemberHandler(svc),
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/members", h.member.CreateMember)
	r.GET("/members/:id", h.member.GetMember)
	r.GET("/members/by-gang/:id", h.member.GetByGang)
	r.GET("/members/sanctioned", h.member.GetSanctioned)
	r.GET("/members/leaders", h.member.GetLeaders)
	r.POST("/members/:id/intel", h.member.AddIntelligenceNote)
	r.GET("/members/:id/intel", h.member.GetIntelligenceNotes)
	r.POST("/members/:id/sightings", h.member.RecordSighting)
	r.GET("/members/:id/sightings", h.member.GetSightings)
	r.PATCH("/members/:id/status", h.member.UpdateStatus)
}
