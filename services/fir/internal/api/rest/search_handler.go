package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/fir/internal/domain"
	"github.com/snisid/platform/services/fir/internal/service"
)

type SearchHandler struct {
	recordSvc *service.RecordService
	chargeSvc *service.ChargeService
}

func NewSearchHandler(rs *service.RecordService, cs *service.ChargeService) *SearchHandler {
	return &SearchHandler{recordSvc: rs, chargeSvc: cs}
}

func (h *SearchHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "paramètre q requis"})
		return
	}

	records, err := h.recordSvc.Search(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results":   records,
		"hit_count": len(records),
		"query":     query,
	})
}

func (h *SearchHandler) SearchByPerson(c *gin.Context) {
	personID, err := uuid.Parse(c.Query("person_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "person_id invalide"})
		return
	}

	record, err := h.recordSvc.GetByPersonID(c.Request.Context(), personID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	charges, err := h.chargeSvc.ListByRecord(c.Request.Context(), record.RecordID)
	if err != nil {
		charges = []domain.Charge{}
	}

	c.JSON(http.StatusOK, gin.H{
		"record":  record,
		"charges": charges,
	})
}
