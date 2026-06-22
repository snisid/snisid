package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/snisid/hsm-svc/internal/domain"
	"github.com/snisid/hsm-svc/internal/service"
)

type Handler struct {
	svc *service.HSMService
}

func NewHandler(svc *service.HSMService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/keys/generate", h.GenerateKey)
	r.GET("/keys/:id", h.GetKey)
	r.POST("/keys/wrap", h.WrapKey)
	r.POST("/keys/sign", h.SignData)
	r.GET("/keys", h.ListKeys)
}

func (h *Handler) GenerateKey(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req domain.KeyGenerationRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.GenerateKey(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetKey(c *gin.Context) {
	keyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key id"})
		return
	}

	key, err := h.svc.GetKey(c.Request.Context(), keyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
		return
	}
	c.JSON(http.StatusOK, key)
}

func (h *Handler) WrapKey(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req domain.KeyWrapRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := h.svc.WrapKey(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"wrapped_data": result})
}

func (h *Handler) SignData(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req domain.KeySignRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	signature, err := h.svc.SignData(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"signature": signature})
}

func (h *Handler) ListKeys(c *gin.Context) {
	algorithm := c.Query("algorithm")
	state := c.Query("state")

	keys, err := h.svc.ListKeys(c.Request.Context(), algorithm, state)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if keys == nil {
		keys = []domain.HSMKey{}
	}
	c.JSON(http.StatusOK, gin.H{"data": keys})
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func int64Ptr(i int64) *int64 {
	return &i
}

var _ = strconv.Itoa
