package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/afis-svc/internal/domain"
	"github.com/snisid/platform/services/afis-svc/internal/service"
)

type SearchHandler struct {
	svc *service.SearchService
}

func NewSearchHandler(svc *service.SearchService) *SearchHandler {
	return &SearchHandler{svc: svc}
}

func (h *SearchHandler) SearchTenprint(c *gin.Context) {
	var req domain.EnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := h.svc.SearchTenprint(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ginH(http.StatusOK, "Recherche terminée", results))
}
