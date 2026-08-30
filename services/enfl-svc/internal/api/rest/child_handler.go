package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/enfl-svc/internal/domain"
	"github.com/snisid/platform/services/enfl-svc/internal/service"
)

type ChildHandler struct {
	svc *service.ChildService
	log *zap.Logger
}

func NewChildHandler(svc *service.ChildService, log *zap.Logger) *ChildHandler {
	return &ChildHandler{svc: svc, log: log}
}

func (h *ChildHandler) RegisterChild(c *gin.Context) {
	var req domain.RegisterChildRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	child, err := h.svc.RegisterChild(&req)
	if err != nil {
		h.log.Error("register child failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, child)
}

func (h *ChildHandler) GetChild(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	child, err := h.svc.GetChild(id)
	if err != nil {
		h.log.Error("get child failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, child)
}

func (h *ChildHandler) ListMissing(c *gin.Context) {
	children, err := h.svc.ListMissing()
	if err != nil {
		h.log.Error("list missing failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, children)
}

func (h *ChildHandler) ListRestaveks(c *gin.Context) {
	restaveks, err := h.svc.ListRestaveks()
	if err != nil {
		h.log.Error("list restaveks failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, restaveks)
}

func (h *ChildHandler) LocateChild(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req domain.LocateChildRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.LocateChild(id, &req); err != nil {
		h.log.Error("locate child failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "located"})
}

func (h *ChildHandler) ListGangRecruited(c *gin.Context) {
	children, err := h.svc.ListGangRecruited()
	if err != nil {
		h.log.Error("list gang recruited failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, children)
}
