package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/snisid/foves-ht/internal/domain"
	"github.com/snisid/foves-ht/internal/service"
)

type Handler struct {
	svc *service.FovesService
}

func NewHandler(svc *service.FovesService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/vehicles", h.RegisterVehicle)
	r.GET("/vehicles/plate/:number", h.GetByPlate)
	r.GET("/vehicles/vin/:vin", h.GetByVIN)
	r.GET("/vehicles/owner/:citizen_id", h.GetByOwner)
	r.POST("/transfers", h.TransferOwnership)
	r.POST("/licenses", h.IssueLicense)
	r.GET("/licenses/:citizen_id", h.GetLicense)
}

func (h *Handler) RegisterVehicle(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		PlateNumber    string `json:"plate_number"`
		VIN            string `json:"vin"`
		Make           string `json:"make"`
		Model          string `json:"model"`
		Year           int    `json:"year"`
		Color          string `json:"color,omitempty"`
		Category       string `json:"category"`
		OwnerCitizenID string `json:"owner_citizen_id"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	ownerID, err := uuid.Parse(req.OwnerCitizenID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid owner_citizen_id"})
		return
	}

	vehicle := &domain.Vehicle{
		PlateNumber:    req.PlateNumber,
		VIN:            req.VIN,
		Make:           req.Make,
		Model:          req.Model,
		Year:           req.Year,
		Category:       domain.VehicleCategory(req.Category),
		OwnerCitizenID: ownerID,
	}
	if req.Color != "" {
		vehicle.Color = &req.Color
	}

	result, err := h.svc.RegisterVehicle(c.Request.Context(), vehicle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetByPlate(c *gin.Context) {
	vehicle, err := h.svc.GetByPlate(c.Request.Context(), c.Param("number"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vehicle)
}

func (h *Handler) GetByVIN(c *gin.Context) {
	vehicle, err := h.svc.GetByVIN(c.Request.Context(), c.Param("vin"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vehicle)
}

func (h *Handler) GetByOwner(c *gin.Context) {
	citizenID, err := uuid.Parse(c.Param("citizen_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid citizen_id"})
		return
	}

	vehicles, err := h.svc.GetByOwner(c.Request.Context(), citizenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": vehicles})
}

func (h *Handler) GetLicense(c *gin.Context) {
	citizenID, err := uuid.Parse(c.Param("citizen_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid citizen_id"})
		return
	}

	license, err := h.svc.GetLicense(c.Request.Context(), citizenID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, license)
}

func (h *Handler) TransferOwnership(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		VehicleID     string  `json:"vehicle_id"`
		FromCitizenID string  `json:"from_citizen_id"`
		ToCitizenID   string  `json:"to_citizen_id"`
		ContractRef   *string `json:"contract_ref,omitempty"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	vehicleID, err := uuid.Parse(req.VehicleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vehicle_id"})
		return
	}
	fromID, err := uuid.Parse(req.FromCitizenID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from_citizen_id"})
		return
	}
	toID, err := uuid.Parse(req.ToCitizenID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to_citizen_id"})
		return
	}

	result, err := h.svc.TransferOwnership(c.Request.Context(), vehicleID, fromID, toID, req.ContractRef)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) IssueLicense(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var req struct {
		CitizenID     string `json:"citizen_id"`
		LicenseNumber string `json:"license_number"`
		CategoryA     bool   `json:"category_a"`
		CategoryB     bool   `json:"category_b"`
		CategoryC     bool   `json:"category_c"`
		CategoryD     bool   `json:"category_d"`
		CategoryE     bool   `json:"category_e"`
		CategoryF     bool   `json:"category_f"`
		ExpiryDate    string `json:"expiry_date"`
		PointsBalance int16  `json:"points_balance"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	citizenID, err := uuid.Parse(req.CitizenID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid citizen_id"})
		return
	}

	expiry, err := time.Parse("2006-01-02", req.ExpiryDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expiry_date, use YYYY-MM-DD"})
		return
	}

	license := &domain.DriverLicense{
		CitizenID:     citizenID,
		LicenseNumber: req.LicenseNumber,
		CategoryA:     req.CategoryA,
		CategoryB:     req.CategoryB,
		CategoryC:     req.CategoryC,
		CategoryD:     req.CategoryD,
		CategoryE:     req.CategoryE,
		CategoryF:     req.CategoryF,
		ExpiryDate:    expiry,
		PointsBalance: req.PointsBalance,
		IsSuspended:   false,
	}

	result, err := h.svc.IssueLicense(c.Request.Context(), license)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}
