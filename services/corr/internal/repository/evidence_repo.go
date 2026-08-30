package repository

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/corr/internal/domain"
)

type EvidenceRepository interface {
	Create(ctx context.Context, e *domain.Evidence) error
	GetByCaseID(ctx context.Context, caseID uuid.UUID) ([]*domain.Evidence, error)
}

type InMemoryEvidenceRepo struct {
	mu       sync.RWMutex
	evidence map[uuid.UUID]*domain.Evidence
}

func NewInMemoryEvidenceRepo() *InMemoryEvidenceRepo {
	return &InMemoryEvidenceRepo{
		evidence: make(map[uuid.UUID]*domain.Evidence),
	}
}

func (r *InMemoryEvidenceRepo) Create(ctx context.Context, e *domain.Evidence) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.evidence[e.EvidenceID] = e
	return nil
}

func (r *InMemoryEvidenceRepo) GetByCaseID(ctx context.Context, caseID uuid.UUID) ([]*domain.Evidence, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*domain.Evidence
	for _, e := range r.evidence {
		if e.CaseID == caseID {
			result = append(result, e)
		}
	}
	return result, nil
}
