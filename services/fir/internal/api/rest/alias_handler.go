package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/fir/internal/domain"
	"github.com/snisid/platform/services/fir/internal/service"
)

type AliasHandler struct {
	aliasSvc  *service.AliasService
	recordSvc *service.RecordService
}

func NewAliasHandler(as *service.AliasService, rs *service.RecordService) *AliasHandler {
	return &AliasHandler{aliasSvc: as, recordSvc: rs}
}

func (h *AliasHandler) AddAlias(c *gin.Context) {
	recordID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID casier invalide"})
		return
	}

	var req struct {
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		BirthDate  string `json:"birth_date"`
		IDDocument string `json:"id_document"`
		Notes      string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	alias := domain.Alias{
		FirstName:  &req.FirstName,
		LastName:   &req.LastName,
		BirthDate:  &req.BirthDate,
		IDDocument: &req.IDDocument,
		Notes:      &req.Notes,
	}

	if req.FirstName == "" {
		alias.FirstName = nil
	}
	if req.LastName == "" {
		alias.LastName = nil
	}
	if req.BirthDate == "" {
		alias.BirthDate = nil
	}
	if req.IDDocument == "" {
		alias.IDDocument = nil
	}
	if req.Notes == "" {
		alias.Notes = nil
	}

	created, err := h.aliasSvc.Add(c.Request.Context(), recordID, alias)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *AliasHandler) RemoveAlias(c *gin.Context) {
	aliasID, err := uuid.Parse(c.Param("alias_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID alias invalide"})
		return
	}

	if err := h.aliasSvc.Remove(c.Request.Context(), aliasID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "supprimé"})
}

func (h *AliasHandler) ListAliases(c *gin.Context) {
	recordID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID casier invalide"})
		return
	}

	aliases, err := h.aliasSvc.ListByRecord(c.Request.Context(), recordID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"aliases": aliases})
}
