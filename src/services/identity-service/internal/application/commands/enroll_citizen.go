package commands

import (
	"context"
	"log"
	"time"

	"github.com/snisid/identity-service/internal/domain"
)

type EnrollCitizenCommand struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	DateOfBirth  string `json:"dateOfBirth"`
	PlaceOfBirth string `json:"placeOfBirth"`
}

type EnrollCitizenResult struct {
	CitizenID string
	Status    string
}

type EnrollCitizenHandler struct {
	// dependencies like repo and kafka producer would go here
}

func NewEnrollCitizenHandler() *EnrollCitizenHandler {
	return &EnrollCitizenHandler{}
}

func (h *EnrollCitizenHandler) Handle(ctx context.Context, cmd EnrollCitizenCommand) (*EnrollCitizenResult, error) {
	// 1. Create domain entity
	citizen, err := domain.NewCitizen(cmd.FirstName, cmd.LastName, cmd.DateOfBirth, cmd.PlaceOfBirth)
	if err != nil {
		return nil, err
	}

	// 2. Persist to CockroachDB
	log.Printf("Persisting citizen to CockroachDB: %s", citizen.ID)

	// 3. Emit Domain Event to Kafka for Biometric ABIS Deduplication
	log.Printf("Emitting IdentityCreated event to Kafka topic 'snisid.identity.events' for Biometric deduplication at %s", time.Now().UTC())

	return &EnrollCitizenResult{
		CitizenID: citizen.ID.String(),
		Status:    string(citizen.Status),
	}, nil
}
