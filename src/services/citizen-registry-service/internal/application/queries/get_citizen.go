package queries

import (
	"context"
	"log"

	"github.com/snisid/citizen-registry-service/internal/domain"
)

type GetCitizenQuery struct {
	NIU string
}

type GetCitizenHandler struct {
	// dependencies like OpenSearch client would go here
}

func NewGetCitizenHandler() *GetCitizenHandler {
	return &GetCitizenHandler{}
}

func (h *GetCitizenHandler) Handle(ctx context.Context, query GetCitizenQuery) (*domain.CitizenView, error) {
	if query.NIU == "" {
		return nil, domain.ErrCitizenNotFound
	}

	// 1. Query OpenSearch
	log.Printf("Querying OpenSearch for NIU: %s", query.NIU)

	// Stubbed response for the read model
	view := &domain.CitizenView{
		NIU:          query.NIU,
		FullName:     "Jean Doe",
		DateOfBirth:  "1990-01-01",
		PlaceOfBirth: "Port-au-Prince",
		Status:       "ACTIVE",
		LastUpdated:  "2026-05-23T12:00:00Z",
	}

	return view, nil
}
