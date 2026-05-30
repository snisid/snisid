package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/backend/internal/domain/authorization/entity"
	"github.com/snisid/platform/backend/internal/domain/authorization/usecase"
)

type HttpHandler struct {
	engine usecase.AuthorizationEngine
}

func NewHttpHandler(engine usecase.AuthorizationEngine) *HttpHandler {
	return &HttpHandler{engine: engine}
}

func (h *HttpHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/enforce", h.Enforce)
	r.POST("/refresh", h.Refresh)
}

func (h *HttpHandler) Enforce(c *gin.Context) {
	var req entity.AuthorizationRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	decision, err := h.engine.Enforce(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, decision)
}

func (h *HttpHandler) Refresh(c *gin.Context) {
	if err := h.engine.RefreshPolicies(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "policies refreshed"})
}
