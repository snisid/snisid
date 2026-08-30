package fraud

import (
	"context"
	"encoding/json"
	"time"

	"github.com/snisid/platform/internal/platform/events"
	"github.com/snisid/platform/internal/platform/logger"
)

type IdentityCreatedEvent struct {
	IdentityID string    `json:"identityId"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Agency     string    `json:"agency"`
	Timestamp  time.Time `json:"timestamp"`
}

type FraudScoredEvent struct {
	IdentityID string    `json:"identityId"`
	RiskScore  int       `json:"riskScore"`
	IsFraud    bool      `json:"isFraud"`
	Reason     string    `json:"reason"`
	Timestamp  time.Time `json:"timestamp"`
}

type Service interface {
	Start(ctx context.Context) error
}

type service struct {
	graphRepo GraphRepository
	consumer  events.ConsumerInterface
	producer  events.ProducerInterface
}

func NewService(graphRepo GraphRepository, consumer events.ConsumerInterface, producer events.ProducerInterface) Service {
	return &service{
		graphRepo: graphRepo,
		consumer:  consumer,
		producer:  producer,
	}
}

func (s *service) Start(ctx context.Context) error {
	return s.consumer.Start(ctx, func(ctx context.Context, msg []byte) error {
		var evt IdentityCreatedEvent
		if err := json.Unmarshal(msg, &evt); err != nil {
			logger.Error(ctx, "failed to unmarshal event", err)
			return nil // Skip invalid messages
		}

		// Insert into Graph DB
		if err := s.graphRepo.AddIdentityNode(ctx, evt.IdentityID, evt.Agency); err != nil {
			logger.Error(ctx, "failed to add node to graph", err)
		}

		// Check for Fraud Rings
		isFraudRing, err := s.graphRepo.CheckFraudRing(ctx, evt.IdentityID)
		if err != nil {
			logger.Error(ctx, "failed to check fraud ring", err)
		}

		// Simple scoring logic
		riskScore := 10
		reason := "Normal"
		isFraud := false

		if isFraudRing {
			riskScore = 95
			reason = "Part of a detected fraud ring"
			isFraud = true
		} else if evt.Agency == "suspicious-agency" {
			riskScore = 80
			reason = "High risk agency"
			isFraud = true
		}

		scoreEvt := FraudScoredEvent{
			IdentityID: evt.IdentityID,
			RiskScore:  riskScore,
			IsFraud:    isFraud,
			Reason:     reason,
			Timestamp:  time.Now().UTC(),
		}

		// Publish to fraud.scored
		if err := s.producer.Publish(ctx, evt.IdentityID, scoreEvt); err != nil {
			logger.Error(ctx, "failed to publish fraud score", err)
			return err // Return err to retry message if needed
		}

		logger.Info(ctx, "fraud score evaluated")

		return nil
	})
}
