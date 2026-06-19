package identity

import (
	"context"
	"fmt"

	"github.com/snisid/platform/internal/platform/events"
)

type EventSender struct {
	producer events.ProducerInterface
}

func NewEventSender(brokers []string) *EventSender {
	// Identity service sends to the ingress topic for central routing
	return &EventSender{
		producer: events.NewProducer(brokers, "snisid.ingress"),
	}
}

func (s *EventSender) EmitIdentityCreated(ctx context.Context, identityID string, metadata map[string]string) error {
	return s.emit(ctx, identityID, "identity.created", "active", metadata)
}

func (s *EventSender) EmitIdentityFlagged(ctx context.Context, identityID string, reason string) error {
	return s.emit(ctx, identityID, "identity.flagged", "flagged", map[string]string{"reason": reason})
}

func (s *EventSender) emit(ctx context.Context, id, action, status string, metadata map[string]string) error {
	evt, err := events.NewEvent(action, map[string]interface{}{
		"identityId": id,
		"action":     action,
		"status":     status,
		"metadata":   metadata,
	})
	if err != nil {
		return fmt.Errorf("failed to create event envelope: %w", err)
	}

	return s.producer.Publish(ctx, id, evt)
}
