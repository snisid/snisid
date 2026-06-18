package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/chef/internal/domain"
	"github.com/snisid/platform/services/chef/internal/service"
)

type HTTPHandler struct {
	svc *service.MemberService
}

func NewHTTPHandler(svc *service.MemberService) *HTTPHandler {
	return &HTTPHandler{svc: svc}
}

func (h *HTTPHandler) RegisterRoutes(rg *gin.Engine) {
	v1 := rg.Group("/api/v1/chef")
	{
		v1.POST("/members", h.CreateMember)
		v1.GET("/members/:id", h.GetMember)
		v1.GET("/members/by-gang/:id", h.GetMembersByGang)
		v1.GET("/members/sanctioned", h.GetSanctionedMembers)
		v1.GET("/members/leaders", h.GetActiveLeaders)
		v1.PATCH("/members/:id/status", h.UpdateStatus)
		v1.POST("/members/:id/intel", h.AddIntelNote)
		v1.GET("/members/:id/intel", h.GetIntelNotes)
		v1.POST("/members/:id/sightings", h.AddSighting)
		v1.GET("/members/:id/sightings", h.GetSightings)
		v1.GET("/network/:id", h.GetMemberNetwork)
	}
}

func (h *HTTPHandler) CreateMember(c *gin.Context) {
	var req domain.CreateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member, err := h.svc.CreateMember(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, member)
}

func (h *HTTPHandler) GetMember(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	member, err := h.svc.GetMember(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, member)
}

func (h *HTTPHandler) GetMembersByGang(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	members, err := h.svc.GetMembersByGang(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}

func (h *HTTPHandler) GetSanctionedMembers(c *gin.Context) {
	members, err := h.svc.GetSanctionedMembers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}

func (h *HTTPHandler) GetActiveLeaders(c *gin.Context) {
	members, err := h.svc.GetActiveLeaders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}

func (h *HTTPHandler) UpdateStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var req domain.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedByStr := c.GetHeader("X-User-ID")
	updatedBy, err := uuid.Parse(updatedByStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-ID invalide"})
		return
	}

	if err := h.svc.UpdateStatus(c.Request.Context(), id, req.Status, updatedBy, req.Notes); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *HTTPHandler) AddIntelNote(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var req domain.CreateIntelNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdByStr := c.GetHeader("X-User-ID")
	createdBy, err := uuid.Parse(createdByStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-ID invalide"})
		return
	}

	note, err := h.svc.AddIntelNote(c.Request.Context(), id, req, createdBy)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, note)
}

func (h *HTTPHandler) GetIntelNotes(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	notes, err := h.svc.GetIntelNotes(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notes": notes})
}

func (h *HTTPHandler) AddSighting(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var req domain.CreateSightingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reportedByStr := c.GetHeader("X-User-ID")
	reportedBy, err := uuid.Parse(reportedByStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-ID invalide"})
		return
	}

	sighting, err := h.svc.AddSighting(c.Request.Context(), id, req, reportedBy)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sighting)
}

func (h *HTTPHandler) GetSightings(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	sightings, err := h.svc.GetSightings(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sightings": sightings})
}

func (h *HTTPHandler) GetMemberNetwork(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	links, err := h.svc.GetMemberNetwork(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"links": links})
}
