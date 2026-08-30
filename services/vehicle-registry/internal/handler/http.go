package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/vehicle-registry/internal/service"
)

type HTTPHandler struct {
	registry *service.Registry
	foves    *service.FoVesClient
	siv      *service.SIVClient
}

func NewHTTPHandler(r *service.Registry, f *service.FoVesClient, s *service.SIVClient) *HTTPHandler {
	return &HTTPHandler{registry: r, foves: f, siv: s}
}

func (h *HTTPHandler) RegisterRoutes(rg *gin.Engine) {
	rg.POST("/vehicles", h.RegisterVehicle)
	rg.GET("/vehicles/:plate", h.LookupByPlate)
	rg.PUT("/vehicles/:plate/owner", h.TransferOwnership)
	rg.GET("/vehicles/search", h.SearchVehicles)
}

func (h *HTTPHandler) RegisterVehicle(c *gin.Context) {
	var req struct {
		Plate         string `json:"plate" binding:"required"`
		VIN           string `json:"vin" binding:"required"`
		Make          string `json:"make" binding:"required"`
		Model         string `json:"model" binding:"required"`
		Year          int    `json:"year" binding:"required"`
		Color         string `json:"color"`
		OwnerID       string `json:"owner_id" binding:"required"`
		InsuranceData string `json:"insurance_data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	v, err := h.registry.RegisterVehicle(req.Plate, req.VIN, req.Make, req.Model, req.Year, req.Color, req.OwnerID, req.InsuranceData)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, v)
}

func (h *HTTPHandler) LookupByPlate(c *gin.Context) {
	plate := c.Param("plate")
	v, ok := h.registry.LookupByPlate(plate)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "vehicle not found"})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *HTTPHandler) TransferOwnership(c *gin.Context) {
	plate := c.Param("plate")
	var req struct {
		NewOwnerID string `json:"new_owner_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	v, err := h.registry.TransferOwnership(plate, req.NewOwnerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *HTTPHandler) SearchVehicles(c *gin.Context) {
	owner := c.Query("owner_id")
	make := c.Query("make")
	model := c.Query("model")
	results := h.registry.SearchVehicles(owner, make, model)
	c.JSON(http.StatusOK, results)
}
