package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/biar/internal/domain"
	"github.com/snisid/platform/services/biar/internal/service"
)

type WeaponHandler struct {
	svc *service.WeaponService
}

func NewWeaponHandler(svc *service.WeaponService) *WeaponHandler {
	return &WeaponHandler{svc: svc}
}

func (h *WeaponHandler) Declare(c *gin.Context) {
	var req domain.CreateWeaponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "requête invalide"})
		return
	}
	createdBy := uuid.Nil
	if uid, ok := c.Request.Context().Value(ContextKeyUserID).(uuid.UUID); ok {
		createdBy = uid
	}
	w, err := h.svc.DeclareIllicitWeapon(c.Request.Context(), req, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, w)
}

func (h *WeaponHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	w, err := h.svc.GetWeapon(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "arme introuvable"})
		return
	}
	c.JSON(http.StatusOK, w)
}

func (h *WeaponHandler) CheckSerial(c *gin.Context) {
	serial := c.Param("sn")
	weapons, err := h.svc.CheckSerial(c.Request.Context(), serial)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur requête"})
		return
	}
	c.JSON(http.StatusOK, weapons)
}
