package repository

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/biar/internal/domain"
)

type BatchRepository interface {
	Create(ctx context.Context, b *domain.BatchSeizure) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.BatchSeizure, error)
	List(ctx context.Context) ([]*domain.BatchSeizure, error)
}

type InMemoryBatchRepo struct {
	mu      sync.RWMutex
	batches map[uuid.UUID]*domain.BatchSeizure
}

func NewInMemoryBatchRepo() *InMemoryBatchRepo {
	return &InMemoryBatchRepo{
		batches: make(map[uuid.UUID]*domain.BatchSeizure),
	}
}

func (r *InMemoryBatchRepo) Create(ctx context.Context, b *domain.BatchSeizure) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.batches[b.BatchID] = b
	return nil
}

func (r *InMemoryBatchRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.BatchSeizure, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.batches[id]
	if !ok {
		return nil, ErrNotFound
	}
	return b, nil
}

func (r *InMemoryBatchRepo) List(ctx context.Context) ([]*domain.BatchSeizure, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.BatchSeizure, 0, len(r.batches))
	for _, b := range r.batches {
		result = append(result, b)
	}
	return result, nil
}
