package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/gang/internal/domain"
	"github.com/snisid/platform/services/gang/internal/service"
)

type MemberHandler struct {
	svc *service.MemberService
}

func NewMemberHandler(svc *service.MemberService) *MemberHandler {
	return &MemberHandler{svc: svc}
}

func (h *MemberHandler) Create(c *gin.Context) {
	var req domain.CreateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "requête invalide"})
		return
	}
	createdBy := uuid.Nil
	if uid, ok := c.Request.Context().Value(ContextKeyUserID).(uuid.UUID); ok {
		createdBy = uid
	}
	member, err := h.svc.CreateMember(c.Request.Context(), req, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, member)
}

func (h *MemberHandler) ListByGang(c *gin.Context) {
	gangID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	members, err := h.svc.GetMembers(c.Request.Context(), gangID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur requête"})
		return
	}
	c.JSON(http.StatusOK, members)
}
