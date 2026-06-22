package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/snisid/card-svc/internal/domain"
	"github.com/snisid/card-svc/internal/service"
)

type Handler struct {
	svc *service.CardService
}

func NewHandler(svc *service.CardService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/cards/order", h.OrderCard)
	r.GET("/cards/:card_serial", h.GetCard)
	r.POST("/cards/:card_serial/activate", h.ActivateCard)
	r.POST("/cards/:card_serial/block", h.BlockCard)
	r.GET("/cards/inventory", h.GetInventory)
	r.POST("/cards/shipments", h.RecordShipment)
}

func (h *Handler) OrderCard(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req domain.PersonalizationRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.OrderPersonalization(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetCard(c *gin.Context) {
	cardSerial := c.Param("card_serial")
	card, err := h.svc.GetCard(c.Request.Context(), cardSerial)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "card not found"})
		return
	}
	c.JSON(http.StatusOK, card)
}

func (h *Handler) ActivateCard(c *gin.Context) {
	cardSerial := c.Param("card_serial")
	result, err := h.svc.ActivateCard(c.Request.Context(), cardSerial)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) BlockCard(c *gin.Context) {
	cardSerial := c.Param("card_serial")
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		Reason string `json:"reason"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.BlockCard(c.Request.Context(), cardSerial, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetInventory(c *gin.Context) {
	profileID := c.Query("profile_id")
	result, err := h.svc.GetInventory(c.Request.Context(), profileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	switch v := result.(type) {
	case *domain.CardInventory:
		c.JSON(http.StatusOK, v)
	case []domain.CardInventory:
		if v == nil {
			v = []domain.CardInventory{}
		}
		c.JSON(http.StatusOK, gin.H{"data": v})
	default:
		c.JSON(http.StatusOK, result)
	}
}

func (h *Handler) RecordShipment(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var shipment domain.Shipment
	if err := json.Unmarshal(body, &shipment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.RecordShipment(c.Request.Context(), shipment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}
