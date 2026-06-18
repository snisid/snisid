package domain

import (
	"time"

	"github.com/google/uuid"
)

type VehicleSighting struct {
	SightingID        uuid.UUID   `json:"sighting_id" db:"sighting_id"`
	PlateNumber       string      `json:"plate_number" db:"plate_number"`
	SourceType        string      `json:"source_type" db:"source_type"`
	LAPIUnitID        *string     `json:"lapi_unit_id,omitempty" db:"lapi_unit_id"`
	ReportingAgentID  *uuid.UUID  `json:"reporting_agent_id,omitempty" db:"reporting_agent_id"`
	SightingTimestamp time.Time   `json:"sighting_timestamp" db:"sighting_timestamp"`
	LocationLat       *float64    `json:"location_lat,omitempty" db:"location_lat"`
	LocationLng       *float64    `json:"location_lng,omitempty" db:"location_lng"`
	LocationDesc      *string     `json:"location_desc,omitempty" db:"location_desc"`
	DeptCode          *string     `json:"dept_code,omitempty" db:"dept_code"`
	Commune           *string     `json:"commune,omitempty" db:"commune"`
	CheckpointName    *string     `json:"checkpoint_name,omitempty" db:"checkpoint_name"`
	MatchedAlertID    *uuid.UUID  `json:"matched_alert_id,omitempty" db:"matched_alert_id"`
	MatchedPlateID    *uuid.UUID  `json:"matched_plate_id,omitempty" db:"matched_plate_id"`
	MatchConfidence   *float64    `json:"match_confidence,omitempty" db:"match_confidence"`
	AlertTriggered    bool        `json:"alert_triggered" db:"alert_triggered"`
	AlertLevel        *AlertLevel `json:"alert_level,omitempty" db:"alert_level"`
	AlertSentAt       *time.Time  `json:"alert_sent_at,omitempty" db:"alert_sent_at"`
	AlertRecipients   []string    `json:"alert_recipients" db:"alert_recipients"`
	ImageRef          *string     `json:"image_ref,omitempty" db:"image_ref"`
	VideoClipRef      *string     `json:"video_clip_ref,omitempty" db:"video_clip_ref"`
	IsReviewed        bool        `json:"is_reviewed" db:"is_reviewed"`
	ReviewedBy        *uuid.UUID  `json:"reviewed_by,omitempty" db:"reviewed_by"`
	ReviewedAt        *time.Time  `json:"reviewed_at,omitempty" db:"reviewed_at"`
	ReviewNotes       *string     `json:"review_notes,omitempty" db:"review_notes"`
	FalsePositive     bool        `json:"false_positive" db:"false_positive"`
	CreatedAt         time.Time   `json:"created_at" db:"created_at"`
}
