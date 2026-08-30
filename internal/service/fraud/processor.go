package fraud

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/snisid/platform/internal/platform/events"
	"github.com/snisid/platform/internal/platform/logger"
	"go.uber.org/zap"
)

type EventConsumer interface {
	Start(ctx context.Context, handler func(ctx context.Context, msg []byte) error) error
}

type EventProducer interface {
	Publish(ctx context.Context, key string, event any) error
}

type IdentityEventProcessor struct {
	consumer EventConsumer
	producer EventProducer
}

func NewIdentityEventProcessor(brokers []string) *IdentityEventProcessor {
	consumer := events.NewConsumer(brokers, "fraud-service-group", "snisid.prod.fraud.v1.events")
	producer := events.NewProducer(brokers, "snisid.prod.risk.v1.updates")

	return &IdentityEventProcessor{
		consumer: consumer,
		producer: producer,
	}
}

func (p *IdentityEventProcessor) Start(ctx context.Context) error {
	return p.consumer.Start(ctx, p.processEvent)
}

func (p *IdentityEventProcessor) processEvent(ctx context.Context, payload []byte) error {
	var event map[string]interface{}
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	logger.Info(ctx, "Processing identity event for fraud detection",
		zap.Any("identity_id", event["identityId"]),
		zap.Any("action", event["action"]),
	)

	// Example Fraud Logic: Flag if rapid updates occur or specific patterns
	// In a real system, this would query Neo4j or a Redis state machine
	isFraudulent := p.analyzeFraudPatterns(event)

	if isFraudulent {
		logger.Warn(ctx, "IDENTITY FRAUD DETECTED", zap.Any("identity_id", event["identityId"]))

		// Publish risk update event
		alert := map[string]interface{}{
			"identityId":      event["identityId"],
			"score_increment": 50,
			"reason":          "Suspicious identity lifecycle activity",
		}
		return p.producer.Publish(ctx, fmt.Sprintf("%v", event["identityId"]), alert)
	}

	return nil
}

func (p *IdentityEventProcessor) analyzeFraudPatterns(event map[string]interface{}) bool {
	// Dummy logic: flag if action is 'deleted' and metadata 'force' is true
	if event["action"] == "deleted" {
		if meta, ok := event["metadata"].(map[string]interface{}); ok {
			return meta["force"] == "true"
		}
	}
	return false
}
