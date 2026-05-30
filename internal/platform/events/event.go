package events

import (
	"time"

	"github.com/google/uuid"
	"github.com/snisid/platform/backend/internal/platform/validation"
)

// Envelope represents the standard SNISID JSON wrapper for all Kafka events
type Envelope[T any] struct {
	EventID   string    `json:"eventId" validate:"required"`
	EventType string    `json:"eventType" validate:"required"`
	Timestamp time.Time `json:"timestamp" validate:"required"`
	Data      T         `json:"data" validate:"required"`
}

// NewEvent creates a strictly validated event envelope.
func NewEvent[T any](eventType string, data T) (*Envelope[T], error) {
	evt := &Envelope[T]{
		EventID:   uuid.NewString(),
		EventType: eventType,
		Timestamp: time.Now().UTC(),
		Data:      data,
	}

	// Schema Validation
	if err := validation.Struct(evt); err != nil {
		return nil, validation.TranslateError(err, "events.NewEvent")
	}

	return evt, nil
}
