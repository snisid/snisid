package repository

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/corr/internal/domain"
)

type OfficerRepository interface {
	Upsert(ctx context.Context, o *domain.Officer) error
	GetBySnisidID(ctx context.Context, id uuid.UUID) (*domain.Officer, error)
	List(ctx context.Context) ([]*domain.Officer, error)
}

type InMemoryOfficerRepo struct {
	mu       sync.RWMutex
	officers map[uuid.UUID]*domain.Officer
}

func NewInMemoryOfficerRepo() *InMemoryOfficerRepo {
	return &InMemoryOfficerRepo{
		officers: make(map[uuid.UUID]*domain.Officer),
	}
}

func (r *InMemoryOfficerRepo) Upsert(ctx context.Context, o *domain.Officer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.officers[o.SnisidID] = o
	return nil
}

func (r *InMemoryOfficerRepo) GetBySnisidID(ctx context.Context, id uuid.UUID) (*domain.Officer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.officers[id]
	if !ok {
		return nil, ErrNotFound
	}
	return o, nil
}

func (r *InMemoryOfficerRepo) List(ctx context.Context) ([]*domain.Officer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.Officer, 0, len(r.officers))
	for _, o := range r.officers {
		result = append(result, o)
	}
	return result, nil
}
