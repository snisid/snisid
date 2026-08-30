package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/siar/internal/domain"
	"github.com/snisid/platform/services/siar/internal/service"
)

type SeizureHandler struct {
	firearmSvc  *service.FirearmService
	transferSvc *service.TransferService
}

func NewSeizureHandler(firearmSvc *service.FirearmService, transferSvc *service.TransferService) *SeizureHandler {
	return &SeizureHandler{firearmSvc: firearmSvc, transferSvc: transferSvc}
}

func (h *SeizureHandler) ReportSeizure(c *gin.Context) {
	var req domain.CreateSeizureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "requête invalide: " + err.Error()})
		return
	}
	createdBy := uuid.Nil
	if uid, ok := c.Request.Context().Value(ContextKeyUserID).(uuid.UUID); ok {
		createdBy = uid
	}

	if req.SerialNumber != "" {
		existing, _ := h.firearmSvc.FindBySerial(c.Request.Context(), req.SerialNumber)
		if existing != nil {
			req.FirearmID = &existing.FirearmID
			h.firearmSvc.UpdateStatus(c.Request.Context(), existing.FirearmID, domain.StatusSeized)
		}
	}

	seiz, err := h.transferSvc.ReportSeizure(c.Request.Context(), req, &createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, seiz)
}

func (h *SeizureHandler) ReportStolen(c *gin.Context) {
	var req domain.ReportStolenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "requête invalide: " + err.Error()})
		return
	}
	createdBy := uuid.Nil
	if uid, ok := c.Request.Context().Value(ContextKeyUserID).(uuid.UUID); ok {
		createdBy = uid
	}

	if req.FirearmID != nil {
		h.firearmSvc.UpdateStatus(c.Request.Context(), *req.FirearmID, domain.StatusReportedStolen)
	} else if req.SerialNumber != "" {
		existing, _ := h.firearmSvc.FindBySerial(c.Request.Context(), req.SerialNumber)
		if existing != nil {
			uid := existing.FirearmID
			req.FirearmID = &uid
			h.firearmSvc.UpdateStatus(c.Request.Context(), existing.FirearmID, domain.StatusReportedStolen)
		}
	}

	seiz, err := h.transferSvc.ReportStolen(c.Request.Context(), req, &createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, seiz)
}
