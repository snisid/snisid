package domain

import (
	"time"

	"github.com/google/uuid"
)

type CameraType string

const (
	CameraTypeFixedIntersection CameraType = "FIXED_INTERSECTION"
	CameraTypeFixedBorder       CameraType = "FIXED_BORDER"
	CameraTypeMobileVehicle     CameraType = "MOBILE_VEHICLE"
	CameraTypeMobileHandheld    CameraType = "MOBILE_HANDHELD"
)

type Camera struct {
	ID        uuid.UUID  `json:"id"`
	Label     string     `json:"label"`
	Type      CameraType `json:"type"`
	Latitude  float64    `json:"latitude"`
	Longitude float64    `json:"longitude"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type PlateRead struct {
	ID                    uuid.UUID `json:"id"`
	CameraID              uuid.UUID `json:"camera_id"`
	PlateNumberRaw        string    `json:"plate_number_raw"`
	PlateNumberNormalized string    `json:"plate_number_normalized"`
	OcrConfidence         float64   `json:"ocr_confidence"`
	Latitude              *float64  `json:"latitude"`
	Longitude             *float64  `json:"longitude"`
	SpeedEstimateKmh      *float64  `json:"speed_estimate_kmh"`
	AlertTriggered        bool      `json:"alert_triggered"`
	CapturedAt            time.Time `json:"captured_at"`
	CreatedAt             time.Time `json:"created_at"`
}

type AlertDispatch struct {
	ID           uuid.UUID  `json:"id"`
	ReadID       uuid.UUID  `json:"read_id"`
	PlateNumber  string     `json:"plate_number"`
	Reason       string     `json:"reason"`
	DispatchedAt time.Time  `json:"dispatched_at"`
	ResolvedAt   *time.Time `json:"resolved_at"`
	IsActive     bool       `json:"is_active"`
}
