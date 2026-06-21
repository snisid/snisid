package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/air-defense-ht/internal/domain"
	"github.com/snisid/air-defense-ht/internal/service"
)

type AirDefenseHandler struct {
	svc service.AirDefenseServiceInterface
}

func NewAirDefenseHandler(svc service.AirDefenseServiceInterface) *AirDefenseHandler {
	return &AirDefenseHandler{svc: svc}
}

type ingestTrackReq struct {
	TrackNumber      string  `json:"track_number" binding:"required"`
	ContactType      string  `json:"contact_type"`
	Latitude         float64 `json:"latitude" binding:"required"`
	Longitude        float64 `json:"longitude" binding:"required"`
	AltitudeM        int     `json:"altitude_m"`
	SpeedKmh         float64 `json:"speed_kmh"`
	HeadingDeg       int     `json:"heading_deg"`
	SourceRadar      string  `json:"source_radar" binding:"required"`
	Identified       bool    `json:"identified"`
	SquawkCode       string  `json:"squawk_code"`
	FlightPlanRef    string  `json:"flight_plan_ref"`
	ThreatAssessment string  `json:"threat_assessment"`
	OperatorNotes    string  `json:"operator_notes"`
}

func (h *AirDefenseHandler) IngestTrack(c *gin.Context) {
	var req ingestTrackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	contact := domain.RadarContact{
		TrackNumber:      req.TrackNumber,
		ContactType:      domain.ContactType(req.ContactType),
		Latitude:         req.Latitude,
		Longitude:        req.Longitude,
		AltitudeM:        req.AltitudeM,
		SpeedKmh:         req.SpeedKmh,
		HeadingDeg:       req.HeadingDeg,
		SourceRadar:      req.SourceRadar,
		Identified:       req.Identified,
		SquawkCode:       req.SquawkCode,
		FlightPlanRef:    req.FlightPlanRef,
		ThreatAssessment: domain.ThreatAssessment(req.ThreatAssessment),
		OperatorNotes:    req.OperatorNotes,
	}
	if err := h.svc.IngestRadarContact(contact); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"contact_id": contact.ContactID})
}

func (h *AirDefenseHandler) GetActiveTracks(c *gin.Context) {
	tracks, err := h.svc.GetActiveTracks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tracks)
}

func (h *AirDefenseHandler) GetTrackByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("track_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid track_id"})
		return
	}
	track, err := h.svc.GetTrackByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if track == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "track not found"})
		return
	}
	c.JSON(http.StatusOK, track)
}

type openIncidentReq struct {
	AircraftID             string `json:"aircraft_id" binding:"required"`
	Severity               string `json:"severity"`
	InterceptionAsset      string `json:"interception_asset"`
	PilotResponse          string `json:"pilot_response"`
	EngagementRulesApplied bool   `json:"engagement_rules_applied"`
	DurationMinutes        int    `json:"duration_minutes"`
}

func (h *AirDefenseHandler) OpenIncident(c *gin.Context) {
	var req openIncidentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	aircraftID, err := uuid.Parse(req.AircraftID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid aircraft_id"})
		return
	}
	incident := domain.AirDefenseIncident{
		AircraftID:             aircraftID,
		Severity:               domain.IncidentSeverity(req.Severity),
		InterceptionAsset:      req.InterceptionAsset,
		PilotResponse:          domain.PilotResponse(req.PilotResponse),
		EngagementRulesApplied: req.EngagementRulesApplied,
		DurationMinutes:        req.DurationMinutes,
	}
	if err := h.svc.OpenIncident(incident); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"incident_id": incident.IncidentID})
}

func (h *AirDefenseHandler) ResolveIncident(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.ResolveIncident(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "resolved"})
}

type addNoFlyReq struct {
	IdentityRef       string `json:"identity_ref" binding:"required"`
	FullName          string `json:"full_name" binding:"required"`
	DocumentNumber    string `json:"document_number"`
	Reason            string `json:"reason" binding:"required"`
	AddedBy           string `json:"added_by" binding:"required"`
	ExpiresAt         string `json:"expires_at" binding:"required"`
	InterpolNoticeRef string `json:"interpol_notice_ref"`
}

func (h *AirDefenseHandler) AddNoFly(c *gin.Context) {
	var req addNoFlyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	expiresAt, err := parseTime(req.ExpiresAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expires_at"})
		return
	}
	entry := domain.NoFlyListEntry{
		IdentityRef:       req.IdentityRef,
		FullName:          req.FullName,
		DocumentNumber:    req.DocumentNumber,
		Reason:            req.Reason,
		AddedBy:           req.AddedBy,
		ExpiresAt:         expiresAt,
		InterpolNoticeRef: req.InterpolNoticeRef,
	}
	if err := h.svc.AddNoFlyEntry(entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"entry_id": entry.EntryID})
}

func (h *AirDefenseHandler) CheckNoFly(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "identity query param required"})
		return
	}
	entry, err := h.svc.CheckNoFly(identity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if entry == nil {
		c.JSON(http.StatusOK, gin.H{"restricted": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"restricted": true, "entry": entry})
}

func parseTime(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05Z", s)
	}
	return t, err
}
