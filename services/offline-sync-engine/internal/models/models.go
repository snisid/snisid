package models

import "time"

type OfflineEvent struct {
	ID           string     `gorm:"primaryKey" json:"id"`
	EventType    string     `gorm:"size:100;not null;index" json:"event_type"`
	Payload      string     `gorm:"type:text;not null" json:"payload"`
	Status       string     `gorm:"size:20;default:'pending';index" json:"status"`
	Priority     int        `gorm:"default:0;index" json:"priority"`
	TerminalID   string     `gorm:"size:50;index" json:"terminal_id"`
	AggregateID  string     `gorm:"size:100;index" json:"aggregate_id"`
	VectorClock  string     `gorm:"type:text" json:"vector_clock"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	SyncedAt     *time.Time `json:"synced_at,omitempty"`
	ErrorMessage string     `gorm:"type:text" json:"error_message,omitempty"`
	RetryCount   int        `gorm:"default:0" json:"retry_count"`
	MaxRetries   int        `gorm:"default:3" json:"max_retries"`
}
