package handlers

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/snisid/platform/services/identity-api/internal/kafka"
	"github.com/snisid/platform/services/identity-api/internal/models"
)

type Handler struct {
	db       *gorm.DB
	producer *kafka.Producer
}

func New(db *gorm.DB, producer *kafka.Producer) *Handler {
	return &Handler{db: db, producer: producer}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/identities", h.CreateIdentity)
	r.GET("/identities/:id", h.GetIdentity)
	r.GET("/identities", h.ListIdentities)
	r.PUT("/identities/:id", h.UpdateIdentity)
	r.DELETE("/identities/:id", h.DeleteIdentity)
	r.POST("/identities/flag/:id", h.FlagIdentity)
}

func generateNNU(deptCode string) string {
	if len(deptCode) < 3 {
		deptCode = deptCode + strings.Repeat("X", 3-len(deptCode))
	}
	deptCode = strings.ToUpper(deptCode[:3])
	year := time.Now().Format("06")
	hex := make([]byte, 7)
	for i := range hex {
		n, err := rand.Int(rand.Reader, big.NewInt(16))
		if err != nil {
			hex[i] = '0'
			continue
		}
		hex[i] = "0123456789ABCDEF"[n.Int64()]
	}
	return deptCode + year + string(hex)
}

type createIdentityRequest struct {
	FirstName     string `json:"first_name" binding:"required"`
	LastName      string `json:"last_name" binding:"required"`
	DateOfBirth   string `json:"date_of_birth"`
	Gender        string `json:"gender"`
	Nationality   string `json:"nationality"`
	BiometricHash string `json:"biometric_hash"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
}

func (h *Handler) CreateIdentity(c *gin.Context) {
	var req createIdentityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ident := models.Identity{
		ID:            uuid.NewString(),
		NNU:           generateNNU(req.LastName),
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		DateOfBirth:   req.DateOfBirth,
		Gender:        req.Gender,
		Nationality:   req.Nationality,
		Status:        "pending",
		BiometricHash: req.BiometricHash,
		Email:         req.Email,
		Phone:         req.Phone,
		Address:       req.Address,
		Version:       1,
	}
	ident.CreatedAt = time.Now().UTC()
	ident.UpdatedAt = ident.CreatedAt

	if err := h.db.Create(&ident).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.publishEvent(c.Request.Context(), "identity.created", ident.ID, actorID(c), ident)
	c.JSON(http.StatusCreated, ident)
}

func (h *Handler) GetIdentity(c *gin.Context) {
	id := c.Param("id")
	var ident models.Identity
	if err := h.db.First(&ident, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "identity not found"})
		return
	}
	c.JSON(http.StatusOK, ident)
}

func (h *Handler) ListIdentities(c *gin.Context) {
	page := parseInt(c.DefaultQuery("page", "1"), 1)
	limit := parseInt(c.DefaultQuery("limit", "20"), 20)
	status := c.Query("status")

	var idents []models.Identity
	query := h.db.Model(&models.Identity{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&idents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  idents,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *Handler) UpdateIdentity(c *gin.Context) {
	id := c.Param("id")

	var existing models.Identity
	if err := h.db.First(&existing, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "identity not found"})
		return
	}

	var req createIdentityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actor := actorID(c)
	changes := buildChanges(id, &existing, &req, actor)

	if len(changes) == 0 {
		c.JSON(http.StatusOK, existing)
		return
	}

	existing.Version++
	existing.UpdatedAt = time.Now().UTC()

	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&existing).Error; err != nil {
			return err
		}
		for i := range changes {
			changes[i].ChangedAt = existing.UpdatedAt
			if err := tx.Create(&changes[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.publishEvent(c.Request.Context(), "identity.updated", existing.ID, actor, existing)
	c.JSON(http.StatusOK, existing)
}

func (h *Handler) DeleteIdentity(c *gin.Context) {
	id := c.Param("id")

	var existing models.Identity
	if err := h.db.First(&existing, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "identity not found"})
		return
	}

	existing.Status = "suspended"
	existing.Version++
	existing.UpdatedAt = time.Now().UTC()

	actor := actorID(c)
	history := models.IdentityHistory{
		ID:         uuid.NewString(),
		IdentityID: id,
		FieldName:  "status",
		OldValue:   existing.Status,
		NewValue:   "suspended",
		ChangedBy:  actor,
		ChangedAt:  existing.UpdatedAt,
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&existing).Error; err != nil {
			return err
		}
		return tx.Create(&history).Error
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.publishEvent(c.Request.Context(), "identity.deleted", existing.ID, actor, existing)
	c.JSON(http.StatusOK, gin.H{"status": "suspended"})
}

func (h *Handler) FlagIdentity(c *gin.Context) {
	id := c.Param("id")

	var existing models.Identity
	if err := h.db.First(&existing, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "identity not found"})
		return
	}

	existing.Status = "suspended"
	existing.Version++
	existing.UpdatedAt = time.Now().UTC()

	actor := actorID(c)
	history := models.IdentityHistory{
		ID:         uuid.NewString(),
		IdentityID: id,
		FieldName:  "status",
		OldValue:   existing.Status,
		NewValue:   "suspended",
		ChangedBy:  actor,
		ChangedAt:  existing.UpdatedAt,
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&existing).Error; err != nil {
			return err
		}
		return tx.Create(&history).Error
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.publishEvent(c.Request.Context(), "identity.flagged", existing.ID, actor, existing)
	c.JSON(http.StatusOK, gin.H{"status": "suspended"})
}

func (h *Handler) publishEvent(ctx context.Context, eventType, identityID, actorID string, data any) {
	if h.producer == nil {
		return
	}
	evt := kafka.Event{
		EventType:     eventType,
		IdentityID:    identityID,
		CorrelationID: uuid.NewString(),
		ActorID:       actorID,
		Timestamp:     time.Now().UTC(),
		Data:          data,
	}
	if err := h.producer.Publish(ctx, evt); err != nil {
		fmt.Printf("failed to publish event: %v\n", err)
	}
}

func buildChanges(id string, existing *models.Identity, req *createIdentityRequest, actor string) []models.IdentityHistory {
	var changes []models.IdentityHistory
	now := time.Now().UTC()

	if req.FirstName != "" && req.FirstName != existing.FirstName {
		changes = append(changes, models.IdentityHistory{
			ID: uuid.NewString(), IdentityID: id, FieldName: "first_name",
			OldValue: existing.FirstName, NewValue: req.FirstName, ChangedBy: actor, ChangedAt: now,
		})
		existing.FirstName = req.FirstName
	}
	if req.LastName != "" && req.LastName != existing.LastName {
		changes = append(changes, models.IdentityHistory{
			ID: uuid.NewString(), IdentityID: id, FieldName: "last_name",
			OldValue: existing.LastName, NewValue: req.LastName, ChangedBy: actor, ChangedAt: now,
		})
		existing.LastName = req.LastName
	}
	if req.DateOfBirth != "" && req.DateOfBirth != existing.DateOfBirth {
		changes = append(changes, models.IdentityHistory{
			ID: uuid.NewString(), IdentityID: id, FieldName: "date_of_birth",
			OldValue: existing.DateOfBirth, NewValue: req.DateOfBirth, ChangedBy: actor, ChangedAt: now,
		})
		existing.DateOfBirth = req.DateOfBirth
	}
	if req.Gender != "" && req.Gender != existing.Gender {
		changes = append(changes, models.IdentityHistory{
			ID: uuid.NewString(), IdentityID: id, FieldName: "gender",
			OldValue: existing.Gender, NewValue: req.Gender, ChangedBy: actor, ChangedAt: now,
		})
		existing.Gender = req.Gender
	}
	if req.Nationality != "" && req.Nationality != existing.Nationality {
		changes = append(changes, models.IdentityHistory{
			ID: uuid.NewString(), IdentityID: id, FieldName: "nationality",
			OldValue: existing.Nationality, NewValue: req.Nationality, ChangedBy: actor, ChangedAt: now,
		})
		existing.Nationality = req.Nationality
	}
	if req.Email != "" && req.Email != existing.Email {
		changes = append(changes, models.IdentityHistory{
			ID: uuid.NewString(), IdentityID: id, FieldName: "email",
			OldValue: existing.Email, NewValue: req.Email, ChangedBy: actor, ChangedAt: now,
		})
		existing.Email = req.Email
	}
	if req.Phone != "" && req.Phone != existing.Phone {
		changes = append(changes, models.IdentityHistory{
			ID: uuid.NewString(), IdentityID: id, FieldName: "phone",
			OldValue: existing.Phone, NewValue: req.Phone, ChangedBy: actor, ChangedAt: now,
		})
		existing.Phone = req.Phone
	}
	if req.Address != "" && req.Address != existing.Address {
		changes = append(changes, models.IdentityHistory{
			ID: uuid.NewString(), IdentityID: id, FieldName: "address",
			OldValue: existing.Address, NewValue: req.Address, ChangedBy: actor, ChangedAt: now,
		})
		existing.Address = req.Address
	}
	if req.BiometricHash != "" && req.BiometricHash != existing.BiometricHash {
		changes = append(changes, models.IdentityHistory{
			ID: uuid.NewString(), IdentityID: id, FieldName: "biometric_hash",
			OldValue: existing.BiometricHash, NewValue: req.BiometricHash, ChangedBy: actor, ChangedAt: now,
		})
		existing.BiometricHash = req.BiometricHash
	}

	return changes
}

func actorID(c *gin.Context) string {
	if id, ok := c.Get("actor_id"); ok {
		if s, ok2 := id.(string); ok2 && s != "" {
			return s
		}
	}
	if id := c.GetHeader("X-Actor-ID"); id != "" {
		return id
	}
	return "system"
}

func parseInt(s string, defaultVal int) int {
	val, err := strconv.Atoi(s)
	if err != nil || val < 1 {
		return defaultVal
	}
	return val
}
