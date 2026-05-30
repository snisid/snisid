package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/backend/internal/domain/audit/usecase"
)

type HttpHandler struct {
	forensics usecase.ForensicsService
}

func NewHttpHandler(forensics usecase.ForensicsService) *HttpHandler {
	return &HttpHandler{forensics: forensics}
}

func (h *HttpHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/verify", h.Verify)
	r.GET("/query", h.Query)
}

func (h *HttpHandler) Verify(c *gin.Context) {
	var req struct {
		Start int64 `form:"start"`
		End   int64 `form:"end"`
	}
	if err := c.BindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	if req.End == 0 {
		req.End = 9999999999 // Arbitrary large number for max end
	}

	valid, err := h.forensics.VerifyIntegrity(c.Request.Context(), req.Start, req.End)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"status": "tamper_detected", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "verified", "valid": valid})
}

func (h *HttpHandler) Query(c *gin.Context) {
	corrID := c.Query("correlationId")
	if corrID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "correlationId is required"})
		return
	}

	events, err := h.forensics.QueryByCorrelationID(c.Request.Context(), corrID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}
