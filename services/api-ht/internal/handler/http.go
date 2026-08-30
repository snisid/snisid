package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/snisid/api-ht/internal/service"
)

type Handler struct {
	svc *service.APIService
}

func NewHandler(svc *service.APIService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/register", h.Register)
	r.POST("/keys/request", h.RequestKey)
	r.GET("/catalog", h.GetCatalog)
	r.GET("/usage/:key_id", h.GetUsage)
	r.POST("/keys/:id/revoke", h.RevokeKey)
}

func (h *Handler) Register(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		Email        string  `json:"email"`
		ContactName  string  `json:"contact_name"`
		OrgName      *string `json:"org_name,omitempty"`
		ContactPhone *string `json:"contact_phone,omitempty"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	if req.Email == "" || req.ContactName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and contact_name are required"})
		return
	}

	result, err := h.svc.RegisterDeveloper(c.Request.Context(), req.Email, req.ContactName, req.OrgName, req.ContactPhone)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) RequestKey(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		AccountID   string  `json:"account_id"`
		Description *string `json:"description,omitempty"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	accountID, err := uuid.Parse(req.AccountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	result, err := h.svc.RequestKey(c.Request.Context(), accountID, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetCatalog(c *gin.Context) {
	catalog, err := h.svc.GetCatalog(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": catalog})
}

func (h *Handler) GetUsage(c *gin.Context) {
	keyID, err := uuid.Parse(c.Param("key_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key_id"})
		return
	}

	logs, err := h.svc.GetUsage(c.Request.Context(), keyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": logs})
}

func (h *Handler) RevokeKey(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key id"})
		return
	}

	if err := h.svc.RevokeKey(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "revoked"})
}
