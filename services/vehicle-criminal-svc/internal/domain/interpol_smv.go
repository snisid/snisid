package domain

import (
	"time"

	"github.com/google/uuid"
)

type InterpolSyncLog struct {
	SyncID          uuid.UUID     `json:"sync_id" db:"sync_id"`
	AlertID         *uuid.UUID    `json:"alert_id,omitempty" db:"alert_id"`
	StolenPlateID   *uuid.UUID    `json:"stolen_plate_id,omitempty" db:"stolen_plate_id"`
	InterpolSMVID   *string       `json:"interpol_smv_id,omitempty" db:"interpol_smv_id"`
	InterpolSADID   *string       `json:"interpol_sad_id,omitempty" db:"interpol_sad_id"`
	SyncDirection   SyncDirection `json:"sync_direction" db:"sync_direction"`
	SyncStatus      SyncStatus    `json:"sync_status" db:"sync_status"`
	SyncTimestamp   time.Time     `json:"sync_timestamp" db:"sync_timestamp"`
	RetryCount      int16         `json:"retry_count" db:"retry_count"`
	RequestPayload  interface{}   `json:"request_payload,omitempty" db:"request_payload"`
	ResponsePayload interface{}   `json:"response_payload,omitempty" db:"response_payload"`
	ErrorCode       *string       `json:"error_code,omitempty" db:"error_code"`
	ErrorMessage    *string       `json:"error_message,omitempty" db:"error_message"`
	ProcessedBy     *uuid.UUID    `json:"processed_by,omitempty" db:"processed_by"`
	ProcessedAt     *time.Time    `json:"processed_at,omitempty" db:"processed_at"`
}

type SMVVehicleRecord struct {
	SMVID          string    `json:"smv_id"`
	NCBReference   string    `json:"ncb_reference"`
	OriginCountry  string    `json:"origin_country"`
	PlateNumber    string    `json:"plate_number"`
	VIN            string    `json:"vin,omitempty"`
	Make           string    `json:"make"`
	Model          string    `json:"model"`
	Year           int       `json:"year,omitempty"`
	ColorPrimary   string    `json:"color_primary"`
	StolenDate     time.Time `json:"stolen_date"`
	StolenLocation string    `json:"stolen_location,omitempty"`
	CrimeType      string    `json:"crime_type"`
	AlertLevel     string    `json:"alert_level,omitempty"`
	SubmittedAt    time.Time `json:"submitted_at"`
}
