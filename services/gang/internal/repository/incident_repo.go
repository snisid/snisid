package repository

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/gang/internal/domain"
)

type IncidentRepository interface {
	Create(ctx context.Context, inc *domain.Incident) error
	ByGangID(ctx context.Context, gangID uuid.UUID) ([]*domain.Incident, error)
}

type InMemoryIncidentRepo struct {
	mu        sync.RWMutex
	incidents map[uuid.UUID]*domain.Incident
}

func NewInMemoryIncidentRepo() *InMemoryIncidentRepo {
	return &InMemoryIncidentRepo{
		incidents: make(map[uuid.UUID]*domain.Incident),
	}
}

func (r *InMemoryIncidentRepo) Create(ctx context.Context, inc *domain.Incident) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.incidents[inc.IncidentID] = inc
	return nil
}

func (r *InMemoryIncidentRepo) ByGangID(ctx context.Context, gangID uuid.UUID) ([]*domain.Incident, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*domain.Incident
	for _, inc := range r.incidents {
		if inc.GangID == gangID {
			result = append(result, inc)
		}
	}
	return result, nil
}
