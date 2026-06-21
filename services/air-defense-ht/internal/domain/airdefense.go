package domain

import (
	"time"

	"github.com/google/uuid"
)

type ContactType string

const (
	ContactCommercial ContactType = "COMMERCIAL"
	ContactPrivate    ContactType = "PRIVATE"
	ContactMilitary   ContactType = "MILITARY"
	ContactDrone      ContactType = "DRONE"
	ContactUnknown    ContactType = "UNKNOWN"
)

type ThreatAssessment string

const (
	ThreatSafe       ThreatAssessment = "SAFE"
	ThreatSuspicious ThreatAssessment = "SUSPICIOUS"
	ThreatHostile    ThreatAssessment = "HOSTILE"
	ThreatUnknown    ThreatAssessment = "UNKNOWN"
)

type IncidentSeverity string

const (
	SeverityInfo     IncidentSeverity = "INFO"
	SeverityLow      IncidentSeverity = "LOW"
	SeverityMedium   IncidentSeverity = "MEDIUM"
	SeverityHigh     IncidentSeverity = "HIGH"
	SeverityCritical IncidentSeverity = "CRITICAL"
)

type IncidentStatus string

const (
	StatusDetected     IncidentStatus = "DETECTED"
	StatusIntercepting IncidentStatus = "INTERCEPTING"
	StatusIntercepted  IncidentStatus = "INTERCEPTED"
	StatusEscorted     IncidentStatus = "ESCORTED"
	StatusViolation    IncidentStatus = "VIOLATION"
	StatusClosed       IncidentStatus = "CLOSED"
)

type PilotResponse string

const (
	PilotCompliant    PilotResponse = "COMPLIANT"
	PilotNonCompliant PilotResponse = "NON_COMPLIANT"
	PilotHostile      PilotResponse = "HOSTILE"
)

type RadarContact struct {
	ContactID        uuid.UUID        `json:"contact_id" db:"contact_id"`
	TrackNumber      string           `json:"track_number" db:"track_number"`
	ContactType      ContactType      `json:"contact_type" db:"contact_type"`
	Latitude         float64          `json:"latitude" db:"latitude"`
	Longitude        float64          `json:"longitude" db:"longitude"`
	AltitudeM        int              `json:"altitude_m" db:"altitude_m"`
	SpeedKmh         float64          `json:"speed_kmh" db:"speed_kmh"`
	HeadingDeg       int              `json:"heading_deg" db:"heading_deg"`
	SourceRadar      string           `json:"source_radar" db:"source_radar"`
	Identified       bool             `json:"identified" db:"identified"`
	SquawkCode       string           `json:"squawk_code" db:"squawk_code"`
	FlightPlanRef    string           `json:"flight_plan_ref" db:"flight_plan_ref"`
	ThreatAssessment ThreatAssessment `json:"threat_assessment" db:"threat_assessment"`
	OperatorNotes    string           `json:"operator_notes" db:"operator_notes"`
	FirstDetectedAt  time.Time        `json:"first_detected_at" db:"first_detected_at"`
	LastUpdatedAt    time.Time        `json:"last_updated_at" db:"last_updated_at"`
}

type AirDefenseIncident struct {
	IncidentID             uuid.UUID        `json:"incident_id" db:"incident_id"`
	Severity               IncidentSeverity `json:"severity" db:"severity"`
	Status                 IncidentStatus   `json:"status" db:"status"`
	AircraftID             uuid.UUID        `json:"aircraft_id" db:"aircraft_id"`
	InterceptionAsset      string           `json:"interception_asset" db:"interception_asset"`
	PilotResponse          PilotResponse    `json:"pilot_response" db:"pilot_response"`
	EngagementRulesApplied bool             `json:"engagement_rules_applied" db:"engagement_rules_applied"`
	DurationMinutes        int              `json:"duration_minutes" db:"duration_minutes"`
	CreatedAt              time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time        `json:"updated_at" db:"updated_at"`
}

type NoFlyListEntry struct {
	EntryID           uuid.UUID `json:"entry_id" db:"entry_id"`
	IdentityRef       string    `json:"identity_ref" db:"identity_ref"`
	FullName          string    `json:"full_name" db:"full_name"`
	DocumentNumber    string    `json:"document_number" db:"document_number"`
	Reason            string    `json:"reason" db:"reason"`
	AddedBy           string    `json:"added_by" db:"added_by"`
	ExpiresAt         time.Time `json:"expires_at" db:"expires_at"`
	InterpolNoticeRef string    `json:"interpol_notice_ref" db:"interpol_notice_ref"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}
