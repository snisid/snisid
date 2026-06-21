package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/fpr-ht/internal/domain"
	"github.com/snisid/fpr-ht/internal/service"
)

type Handler struct {
	svc *service.FprService
}

func NewHandler(svc *service.FprService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1/fpr")
	{
		v1.POST("/warrants", h.CreateWarrant)
		v1.GET("/check/:citizen_id", h.CheckCitizen)
		v1.GET("/check/name", h.CheckByName)
		v1.POST("/warrants/:id/sightings", h.ReportSighting)
		v1.PATCH("/warrants/:id/execute", h.ExecuteWarrant)
		v1.GET("/warrants/armed-dangerous", h.GetArmedDangerous)
		v1.GET("/stats/dashboard", h.GetDashboardStats)
	}
}

type createWarrantRequest struct {
	FullName           string    `json:"full_name" binding:"required"`
	Aliases            []string  `json:"aliases"`
	AfisSubjectID      *string   `json:"afis_subject_id"`
	WarrantType        string    `json:"warrant_type" binding:"required"`
	Charges            []string  `json:"charges"`
	IssuingCourt       string    `json:"issuing_court" binding:"required"`
	DangerLevel        *string   `json:"danger_level"`
	PhotoRefs          []string  `json:"photo_refs"`
	VehiclePlatesKnown []string  `json:"vehicle_plates_known"`
	InterpolNoticeRef  *string   `json:"interpol_notice_ref"`
	IssuedAt           time.Time `json:"issued_at"`
}

func (h *Handler) CreateWarrant(c *gin.Context) {
	var req createWarrantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	w := &domain.Warrant{
		FullName:           req.FullName,
		Aliases:            req.Aliases,
		AfisSubjectID:      req.AfisSubjectID,
		WarrantType:        domain.WarrantType(req.WarrantType),
		Charges:            req.Charges,
		IssuingCourt:       req.IssuingCourt,
		PhotoRefs:          req.PhotoRefs,
		VehiclePlatesKnown: req.VehiclePlatesKnown,
		InterpolNoticeRef:  req.InterpolNoticeRef,
		IssuedAt:           req.IssuedAt,
	}
	if req.DangerLevel != nil {
		d := domain.DangerLevel(*req.DangerLevel)
		w.DangerLevel = &d
	}
	if w.IssuedAt.IsZero() {
		w.IssuedAt = time.Now().UTC()
	}

	if err := h.svc.CreateWarrant(w); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, w)
}

func (h *Handler) CheckCitizen(c *gin.Context) {
	citizenID := c.Param("citizen_id")
	if citizenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "citizen_id required"})
		return
	}
	result, err := h.svc.CheckCitizen(citizenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) CheckByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name query parameter required"})
		return
	}
	result, err := h.svc.CheckByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

type createSightingRequest struct {
	CitizenID   string    `json:"citizen_id" binding:"required"`
	Latitude    *float64  `json:"latitude"`
	Longitude   *float64  `json:"longitude"`
	Description string    `json:"description"`
	ReportedBy  string    `json:"reported_by" binding:"required"`
	SightedAt   time.Time `json:"sighted_at"`
}

func (h *Handler) ReportSighting(c *gin.Context) {
	idStr := c.Param("id")
	warrantID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid warrant id"})
		return
	}

	var req createSightingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sighting := &domain.Sighting{
		CitizenID:   req.CitizenID,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Description: req.Description,
		ReportedBy:  req.ReportedBy,
		SightedAt:   req.SightedAt,
	}
	if sighting.SightedAt.IsZero() {
		sighting.SightedAt = time.Now().UTC()
	}

	if err := h.svc.ReportSighting(warrantID, sighting); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sighting)
}

func (h *Handler) ExecuteWarrant(c *gin.Context) {
	idStr := c.Param("id")
	warrantID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid warrant id"})
		return
	}
	if err := h.svc.ExecuteWarrant(warrantID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "executed"})
}

func (h *Handler) GetArmedDangerous(c *gin.Context) {
	warrants, err := h.svc.GetArmedDangerous()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, warrants)
}

func (h *Handler) GetDashboardStats(c *gin.Context) {
	stats, err := h.svc.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
