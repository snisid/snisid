package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/entity-resolution/internal/matching"
	"github.com/snisid/platform/services/entity-resolution/internal/models"
	"gorm.io/gorm"
)

type ResolutionHandler struct {
	engine *matching.CompositeEngine
	db     *gorm.DB
}

func NewResolutionHandler(engine *matching.CompositeEngine, db *gorm.DB) *ResolutionHandler {
	return &ResolutionHandler{engine: engine, db: db}
}

func (h *ResolutionHandler) MatchHandler(c *gin.Context) {
	var req models.MatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	candidates, err := h.engine.Match(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"candidates": candidates,
		"count":      len(candidates),
	})
}

func (h *ResolutionHandler) ReconcileHandler(c *gin.Context) {
	var req struct {
		PrimaryID   string `json:"primary_id" binding:"required"`
		SecondaryID string `json:"secondary_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.engine.Reconcile(req.PrimaryID, req.SecondaryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ResolutionHandler) CandidatesHandler(c *gin.Context) {
	id := c.Param("id")

	var identity models.Identity
	if err := h.db.First(&identity, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "identity not found"})
		return
	}

	req := models.MatchRequest{
		FirstName:     identity.FirstName,
		LastName:      identity.LastName,
		FullName:      identity.FullName,
		DOB:           identity.DOB,
		TaxID:         identity.TaxID,
		NNU:           identity.NNU,
		BiometricHash: identity.BiometricHash,
	}

	candidates, err := h.engine.Match(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filtered := make([]models.MatchCandidate, 0)
	for _, c := range candidates {
		if c.IdentityID != id {
			filtered = append(filtered, c)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"identity_id": id,
		"candidates":  filtered,
		"count":       len(filtered),
	})
}

func (h *ResolutionHandler) MergeHandler(c *gin.Context) {
	var req struct {
		PrimaryID   string `json:"primary_id" binding:"required"`
		SecondaryID string `json:"secondary_id" binding:"required"`
		ResolvedBy  string `json:"resolved_by"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.engine.Merge(req.PrimaryID, req.SecondaryID, req.ResolvedBy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "identities merged successfully",
		"primary_id":   req.PrimaryID,
		"secondary_id": req.SecondaryID,
	})
}

func (h *ResolutionHandler) SplitHandler(c *gin.Context) {
	var req struct {
		IdentityID string `json:"identity_id" binding:"required"`
		ResolvedBy string `json:"resolved_by"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.engine.Split(req.IdentityID, req.ResolvedBy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "identity split successfully",
		"identity_id": req.IdentityID,
	})
}

func (h *ResolutionHandler) StatsHandler(c *gin.Context) {
	stats, err := h.engine.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func init() {
	uuid.New()
}
