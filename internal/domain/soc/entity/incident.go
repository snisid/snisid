package entity

import (
	"time"
)

type Severity string

const (
	SeverityLow      Severity = "LOW"
	SeverityMedium   Severity = "MEDIUM"
	SeverityHigh     Severity = "HIGH"
	SeverityCritical Severity = "CRITICAL"
)

type IncidentStatus string

const (
	StatusNew        IncidentStatus = "NEW"
	StatusAssigned   IncidentStatus = "ASSIGNED"
	StatusInAnalysis IncidentStatus = "IN_ANALYSIS"
	StatusContained  IncidentStatus = "CONTAINED"
	StatusResolved   IncidentStatus = "RESOLVED"
)

type Incident struct {
	ID            string         `json:"id" gorm:"primaryKey"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	Severity      Severity       `json:"severity"`
	Status        IncidentStatus `json:"status"`
	Source        string         `json:"source"`
	AffectedUser  string         `json:"affected_user"`
	CorrelationID string         `json:"correlation_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	PlaybookID    string         `json:"playbook_id"`
	ActionsTaken  []Action       `json:"actions_taken" gorm:"serializer:json"`
}

type Action struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Result    string    `json:"result"`
	Success   bool      `json:"success"`
}
