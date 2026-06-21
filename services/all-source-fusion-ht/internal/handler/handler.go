package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/all-source-fusion-ht/internal/domain"
	"github.com/snisid/all-source-fusion-ht/internal/service"
)

type FusionHandler struct {
	svc service.FusionService
}

func NewFusionHandler(svc service.FusionService) *FusionHandler {
	return &FusionHandler{svc: svc}
}

func (h *FusionHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1/fusion")
	{
		api.POST("/products", h.CreateProduct)
		api.GET("/products/recent", h.GetRecentProducts)
		api.GET("/products/:id/source-map", h.GetSourceMap)
		api.POST("/threat-actors", h.CreateThreatActor)
		api.GET("/threat-actors/high-risk", h.GetHighRiskActors)
		api.POST("/correlations", h.CreateCorrelation)
		api.GET("/estimates/national", h.GetNationalEstimates)
	}
}

func (h *FusionHandler) CreateProduct(c *gin.Context) {
	var p domain.IntelProduct
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.CreateProduct(c.Request.Context(), &p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *FusionHandler) GetRecentProducts(c *gin.Context) {
	result, err := h.svc.GetRecentProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *FusionHandler) GetSourceMap(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}
	result, err := h.svc.GetSourceMap(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if result == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *FusionHandler) CreateThreatActor(c *gin.Context) {
	var a domain.ThreatActor
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.CreateThreatActor(c.Request.Context(), &a); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, a)
}

func (h *FusionHandler) GetHighRiskActors(c *gin.Context) {
	result, err := h.svc.GetHighRiskActors(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *FusionHandler) CreateCorrelation(c *gin.Context) {
	var corr domain.CrossDisciplineCorrelation
	if err := c.ShouldBindJSON(&corr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.CreateCorrelation(c.Request.Context(), &corr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, corr)
}

func (h *FusionHandler) GetNationalEstimates(c *gin.Context) {
	result, err := h.svc.GetNationalEstimates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
