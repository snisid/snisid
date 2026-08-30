package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/reso-svc/internal/service"
)

type ResoHandler struct {
	svc *service.ResoService
	log *zap.Logger
}

func NewResoHandler(svc *service.ResoService, log *zap.Logger) *ResoHandler {
	return &ResoHandler{svc: svc, log: log}
}

func (h *ResoHandler) GetPersonNetwork(c *gin.Context) {
	idStr := c.Param("person_id")
	personID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid person_id"})
		return
	}

	result, err := h.svc.GetPersonNetwork(personID)
	if err != nil {
		h.log.Error("get person network failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ResoHandler) GetCommunities(c *gin.Context) {
	result, err := h.svc.DetectCommunities()
	if err != nil {
		h.log.Error("detect communities failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ResoHandler) GetKeyActors(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	actors, err := h.svc.GetKeyActors(limit)
	if err != nil {
		h.log.Error("get key actors failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, actors)
}

func (h *ResoHandler) GetGangOverlap(c *gin.Context) {
	g1Str := c.Param("g1")
	g2Str := c.Param("g2")

	gangID1, err := uuid.Parse(g1Str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid gang_id g1"})
		return
	}
	gangID2, err := uuid.Parse(g2Str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid gang_id g2"})
		return
	}

	result, err := h.svc.GetGangOverlap(gangID1, gangID2)
	if err != nil {
		h.log.Error("get gang overlap failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ResoHandler) TriggerAnalysis(c *gin.Context) {
	result, err := h.svc.AnalyzeNetwork()
	if err != nil {
		h.log.Error("trigger analysis failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ResoHandler) FindShortestPath(c *gin.Context) {
	fromStr := c.Param("from_id")
	toStr := c.Param("to_id")

	fromID, err := uuid.Parse(fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from_id"})
		return
	}
	toID, err := uuid.Parse(toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to_id"})
		return
	}

	result, err := h.svc.FindShortestPath(fromID, toID)
	if err != nil {
		h.log.Error("find shortest path failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ResoHandler) GetCentralityScores(c *gin.Context) {
	actors, err := h.svc.GetKeyActors(0)
	if err != nil {
		h.log.Error("get centrality scores failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, actors)
}

func (h *ResoHandler) GetEmergingLinks(c *gin.Context) {
	actors, err := h.svc.GetKeyActors(50)
	if err != nil {
		h.log.Error("get emerging links failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"emerging_links": actors, "description": "high-centrality actors with potential new connections"})
}
