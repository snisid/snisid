package domain

import (
	"time"

	"github.com/google/uuid"
)

type WarrantType string

const (
	WarrantTypeArrest            WarrantType = "ARREST_WARRANT"
	WarrantTypeBench             WarrantType = "BENCH_WARRANT"
	WarrantTypeSearchPerson      WarrantType = "SEARCH_WARRANT_PERSON"
	WarrantTypeProbationViol     WarrantType = "PROBATION_VIOLATION"
	WarrantTypeInterpolRedNotice WarrantType = "INTERPOL_RED_NOTICE"
	WarrantTypeContemptOfCourt   WarrantType = "CONTEMPT_OF_COURT"
)

type DangerLevel string

const (
	DangerLevelArmedAndDangerous DangerLevel = "ARMED_AND_DANGEROUS"
)

type Warrant struct {
	ID                 uuid.UUID    `json:"id"`
	FullName           string       `json:"full_name"`
	Aliases            []string     `json:"aliases"`
	AfisSubjectID      *string      `json:"afis_subject_id"`
	WarrantType        WarrantType  `json:"warrant_type"`
	Charges            []string     `json:"charges"`
	IssuingCourt       string       `json:"issuing_court"`
	DangerLevel        *DangerLevel `json:"danger_level"`
	PhotoRefs          []string     `json:"photo_refs"`
	VehiclePlatesKnown []string     `json:"vehicle_plates_known"`
	InterpolNoticeRef  *string      `json:"interpol_notice_ref"`
	IsExecuted         bool         `json:"is_executed"`
	IssuedAt           time.Time    `json:"issued_at"`
	ExecutedAt         *time.Time   `json:"executed_at"`
	CreatedAt          time.Time    `json:"created_at"`
	UpdatedAt          time.Time    `json:"updated_at"`
}

type Sighting struct {
	ID          uuid.UUID `json:"id"`
	WarrantID   uuid.UUID `json:"warrant_id"`
	CitizenID   string    `json:"citizen_id"`
	Latitude    *float64  `json:"latitude"`
	Longitude   *float64  `json:"longitude"`
	Description string    `json:"description"`
	ReportedBy  string    `json:"reported_by"`
	SightedAt   time.Time `json:"sighted_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type CheckLog struct {
	ID        uuid.UUID  `json:"id"`
	CitizenID string     `json:"citizen_id"`
	WarrantID *uuid.UUID `json:"warrant_id"`
	Result    string     `json:"result"`
	CheckedAt time.Time  `json:"checked_at"`
}

type WarrantCheckResult struct {
	WarrantFound bool     `json:"warrant_found"`
	Warrant      *Warrant `json:"warrant,omitempty"`
	CheckLog     CheckLog `json:"check_log"`
}

type DashboardStats struct {
	TotalWarrants    int `json:"total_warrants"`
	ActiveWarrants   int `json:"active_warrants"`
	ExecutedWarrants int `json:"executed_warrants"`
	ArmedDangerous   int `json:"armed_dangerous"`
	TotalSightings   int `json:"total_sightings"`
	TotalChecks      int `json:"total_checks"`
}
