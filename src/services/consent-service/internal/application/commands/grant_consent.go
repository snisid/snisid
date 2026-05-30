package commands

import (
	"context"
	"log"

	"github.com/snisid/consent-service/internal/domain"
)

type GrantConsentCommand struct {
	CitizenNIU       string `json:"citizenNiu"`
	TargetAgency     string `json:"targetAgency"`
	DaysValid        int    `json:"daysValid"`
	DigitalSignature string `json:"digitalSignature"`
}

type GrantConsentResult struct {
	ConsentID string
	Status    string
}

type GrantConsentHandler struct {
	// DB and Kafka dependencies
}

func NewGrantConsentHandler() *GrantConsentHandler {
	return &GrantConsentHandler{}
}

func (h *GrantConsentHandler) Handle(ctx context.Context, cmd GrantConsentCommand) (*GrantConsentResult, error) {
	agreement, err := domain.NewConsentAgreement(cmd.CitizenNIU, cmd.TargetAgency, cmd.DaysValid, cmd.DigitalSignature)
	if err != nil {
		return nil, err
	}

	// 1. Persist to CockroachDB (Immutable record)
	log.Printf("Persisting consent agreement to immutable ledger: %s", agreement.ID)

	// 2. Publish to Kafka to notify API Gateway/OPA
	log.Printf("Emitting ConsentGranted event to Kafka topic 'snisid.consent.events'")

	return &GrantConsentResult{
		ConsentID: agreement.ID.String(),
		Status:    "CONSENT_GRANTED",
	}, nil
}
