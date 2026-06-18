package repository

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/corr/internal/domain"
)

type CaseRepository interface {
	Create(ctx context.Context, c *domain.IntegrityCase) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.IntegrityCase, error)
	ListActive(ctx context.Context) ([]*domain.IntegrityCase, error)
	List(ctx context.Context) ([]*domain.IntegrityCase, error)
	Update(ctx context.Context, c *domain.IntegrityCase) error
}

type InMemoryCaseRepo struct {
	mu    sync.RWMutex
	cases map[uuid.UUID]*domain.IntegrityCase
}

func NewInMemoryCaseRepo() *InMemoryCaseRepo {
	return &InMemoryCaseRepo{
		cases: make(map[uuid.UUID]*domain.IntegrityCase),
	}
}

func (r *InMemoryCaseRepo) Create(ctx context.Context, c *domain.IntegrityCase) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cases[c.CaseID] = c
	return nil
}

func (r *InMemoryCaseRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.IntegrityCase, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.cases[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (r *InMemoryCaseRepo) ListActive(ctx context.Context) ([]*domain.IntegrityCase, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*domain.IntegrityCase
	for _, c := range r.cases {
		if c.Status == domain.StatusReported || c.Status == domain.StatusUnderInvestigation {
			result = append(result, c)
		}
	}
	return result, nil
}

func (r *InMemoryCaseRepo) List(ctx context.Context) ([]*domain.IntegrityCase, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.IntegrityCase, 0, len(r.cases))
	for _, c := range r.cases {
		result = append(result, c)
	}
	return result, nil
}

func (r *InMemoryCaseRepo) Update(ctx context.Context, c *domain.IntegrityCase) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cases[c.CaseID] = c
	return nil
}
