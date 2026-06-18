package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/chef-svc/internal/domain"
	"github.com/snisid/platform/services/chef-svc/internal/service"
)

type MemberHandler struct {
	svc *service.MemberService
}

func NewMemberHandler(svc *service.MemberService) *MemberHandler {
	return &MemberHandler{svc: svc}
}

func (h *MemberHandler) CreateMember(c *gin.Context) {
	var member domain.CriminalMember
	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.CreateMember(&member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, member)
}

func (h *MemberHandler) GetMember(c *gin.Context) {
	id := c.Param("id")
	member, err := h.svc.GetMember(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
		return
	}
	c.JSON(http.StatusOK, member)
}

func (h *MemberHandler) GetByGang(c *gin.Context) {
	gangID := c.Param("id")
	members, err := h.svc.GetByGang(gangID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, members)
}

func (h *MemberHandler) GetSanctioned(c *gin.Context) {
	members, err := h.svc.GetSanctioned()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, members)
}

func (h *MemberHandler) GetLeaders(c *gin.Context) {
	members, err := h.svc.GetLeaders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, members)
}

func (h *MemberHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status domain.MemberStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateStatus(id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}

func (h *MemberHandler) AddIntelligenceNote(c *gin.Context) {
	memberID := c.Param("id")
	var note domain.IntelligenceNote
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	note.MemberID = memberID

	if err := h.svc.AddIntelligenceNote(&note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, note)
}

func (h *MemberHandler) GetIntelligenceNotes(c *gin.Context) {
	memberID := c.Param("id")
	notes, err := h.svc.GetIntelligenceNotes(memberID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notes)
}

func (h *MemberHandler) RecordSighting(c *gin.Context) {
	memberID := c.Param("id")
	var sighting domain.Sighting
	if err := c.ShouldBindJSON(&sighting); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sighting.MemberID = memberID

	if err := h.svc.RecordSighting(&sighting); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sighting)
}

func (h *MemberHandler) GetSightings(c *gin.Context) {
	memberID := c.Param("id")
	sightings, err := h.svc.GetSightings(memberID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sightings)
}
