package domain

import (
	"time"

	"github.com/google/uuid"
)

type ChangeHistory struct {
	HistoryID    uuid.UUID `json:"history_id"`
	CitizenID    uuid.UUID `json:"citizen_id"`
	FieldChanged string    `json:"field_changed"`
	OldValue     *string   `json:"old_value,omitempty"`
	NewValue     *string   `json:"new_value,omitempty"`
	ChangeReason string    `json:"change_reason"`
	AuthorizedBy uuid.UUID `json:"authorized_by"`
	ChangedAt    time.Time `json:"changed_at"`
}

type PopulationStats struct {
	Total     int `json:"total"`
	Active    int `json:"active"`
	Suspended int `json:"suspended"`
	Deceased  int `json:"deceased"`
	Cancelled int `json:"cancelled"`
	Merged    int `json:"merged"`
}
