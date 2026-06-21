package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/ong-svc/internal/domain"
	"github.com/snisid/platform/services/ong-svc/internal/service"
)

type ONGHandler struct {
	svc *service.ONGService
	log *zap.Logger
}

func NewONGHandler(svc *service.ONGService, log *zap.Logger) *ONGHandler {
	return &ONGHandler{svc: svc, log: log}
}

func (h *ONGHandler) RegisterOrg(c *gin.Context) {
	var req domain.RegisterOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	org, err := h.svc.RegisterOrg(&req)
	if err != nil {
		h.log.Error("register org failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, org)
}

func (h *ONGHandler) GetOrg(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	org, err := h.svc.GetOrg(id)
	if err != nil {
		h.log.Error("get org failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, org)
}

func (h *ONGHandler) ListOrgs(c *gin.Context) {
	orgs, err := h.svc.ListOrgs()
	if err != nil {
		h.log.Error("list orgs failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, orgs)
}

func (h *ONGHandler) ListFlagged(c *gin.Context) {
	orgs, err := h.svc.ListFlagged()
	if err != nil {
		h.log.Error("list flagged failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, orgs)
}

func (h *ONGHandler) ListUnregistered(c *gin.Context) {
	orgs, err := h.svc.ListUnregistered()
	if err != nil {
		h.log.Error("list unregistered failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, orgs)
}

func (h *ONGHandler) ScreenOrg(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	result, err := h.svc.ScreenOrg(id)
	if err != nil {
		h.log.Error("screen org failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *ONGHandler) RegisterStaff(c *gin.Context) {
	var req domain.RegisterStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	staff, err := h.svc.RegisterStaff(&req)
	if err != nil {
		h.log.Error("register staff failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, staff)
}

func (h *ONGHandler) RequestAccess(c *gin.Context) {
	var req domain.RequestAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ar, err := h.svc.RequestAccess(&req)
	if err != nil {
		h.log.Error("request access failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, ar)
}

func (h *ONGHandler) ApproveAccess(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.ApproveAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.ApproveAccess(id, &req); err != nil {
		h.log.Error("approve access failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}
