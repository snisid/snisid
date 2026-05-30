package events

import (
	"time"
)

type SNISIDEvent struct {
	EventID     string    `json:"event_id"`
	Type        string    `json:"type"`
	User        string    `json:"user"`
	Institution string    `json:"institution"`
	Action      string    `json:"action"`
	RiskScore   float64   `json:"risk_score"`
	Timestamp   time.Time `json:"timestamp"`
	Signature   string    `json:"signature"`
}

func NewEvent(eventType, userID, institution, action string, risk float64) *SNISIDEvent {
	return &SNISIDEvent{
		EventID:     "evt_" + time.Now().Format("20060102150405"),
		Type:        eventType,
		User:        userID,
		Institution: institution,
		Action:      action,
		RiskScore:   risk,
		Timestamp:   time.Now(),
		Signature:   "sig_valid",
	}
}
