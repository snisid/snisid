package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/siar/internal/domain"
	"github.com/snisid/platform/services/siar/internal/service"
)

type DealerHandler struct {
	svc *service.DealerService
}

func NewDealerHandler(svc *service.DealerService) *DealerHandler {
	return &DealerHandler{svc: svc}
}

func (h *DealerHandler) Create(c *gin.Context) {
	var req domain.CreateDealerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "requête invalide: " + err.Error()})
		return
	}
	createdBy := uuid.Nil
	if uid, ok := c.Request.Context().Value(ContextKeyUserID).(uuid.UUID); ok {
		createdBy = uid
	}
	d, err := h.svc.Create(c.Request.Context(), req, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, d)
}

func (h *DealerHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	d, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "marchand d'armes non trouvé"})
		return
	}
	c.JSON(http.StatusOK, d)
}

func (h *DealerHandler) List(c *gin.Context) {
	deptCode := c.Query("dept_code")
	dealers, err := h.svc.List(c.Request.Context(), deptCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur requête"})
		return
	}
	if dealers == nil {
		dealers = []*domain.Dealer{}
	}
	c.JSON(http.StatusOK, gin.H{"dealers": dealers})
}
