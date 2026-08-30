package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/sanc-svc/internal/domain"
	"github.com/snisid/platform/services/sanc-svc/internal/service"
)

type SancHandler struct {
	svc *service.SanctionsService
	log *zap.Logger
}

func NewSancHandler(svc *service.SanctionsService, log *zap.Logger) *SancHandler {
	return &SancHandler{svc: svc, log: log}
}

func (h *SancHandler) CheckPerson(c *gin.Context) {
	idStr := c.Param("person_id")
	personID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid person_id"})
		return
	}

	result, err := h.svc.CheckPersonRealTime(c.Request.Context(), personID)
	if err != nil {
		h.log.Error("check person failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SancHandler) CheckByName(c *gin.Context) {
	var req domain.CheckNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entries, err := h.svc.SearchByName(c.Request.Context(), req.Name, req.DateOfBirth)
	if err != nil {
		h.log.Error("search by name failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"entries": entries})
}

func (h *SancHandler) GetEntries(c *gin.Context) {
	entries, total, err := h.svc.GetActiveEntries(c.Request.Context(), 50, 0)
	if err != nil {
		h.log.Error("get entries failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"entries": entries, "total": total})
}

func (h *SancHandler) GetHaitiEntries(c *gin.Context) {
	entries, _, err := h.svc.GetActiveEntries(c.Request.Context(), 500, 0)
	if err != nil {
		h.log.Error("get haiti entries failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	var haitiEntries []domain.SanctionEntry
	for _, e := range entries {
		for _, nat := range e.Nationality {
			if nat == "HT" || nat == "Haiti" || nat == "HAITI" {
				haitiEntries = append(haitiEntries, e)
				break
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"entries": haitiEntries})
}

func (h *SancHandler) GetUnconfirmedMatches(c *gin.Context) {
	matches, err := h.svc.GetUnconfirmedMatches(c.Request.Context())
	if err != nil {
		h.log.Error("get unconfirmed matches failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"matches": matches})
}

func (h *SancHandler) ConfirmMatch(c *gin.Context) {
	idStr := c.Param("id")
	matchID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}

	var req domain.ConfirmMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.ConfirmMatch(c.Request.Context(), matchID, req.ConfirmedBy); err != nil {
		h.log.Error("confirm match failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "match confirmed"})
}

func (h *SancHandler) TriggerSync(c *gin.Context) {
	result, err := h.svc.SyncOFAC(c.Request.Context())
	if err != nil {
		h.log.Error("trigger sync failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SancHandler) GetSyncStatus(c *gin.Context) {
	logs, err := h.svc.GetSyncStatus(c.Request.Context())
	if err != nil {
		h.log.Error("get sync status failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}
