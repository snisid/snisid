package repository

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/corr/internal/domain"
)

type WhistleblowerRepository interface {
	Create(ctx context.Context, r *domain.WhistleblowerReport) error
	GetByToken(ctx context.Context, token string) (*domain.WhistleblowerReport, error)
	List(ctx context.Context) ([]*domain.WhistleblowerReport, error)
}

type InMemoryWhistleblowerRepo struct {
	mu      sync.RWMutex
	reports map[uuid.UUID]*domain.WhistleblowerReport
}

func NewInMemoryWhistleblowerRepo() *InMemoryWhistleblowerRepo {
	return &InMemoryWhistleblowerRepo{
		reports: make(map[uuid.UUID]*domain.WhistleblowerReport),
	}
}

func (r *InMemoryWhistleblowerRepo) Create(ctx context.Context, w *domain.WhistleblowerReport) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reports[w.ReportID] = w
	return nil
}

func (r *InMemoryWhistleblowerRepo) GetByToken(ctx context.Context, token string) (*domain.WhistleblowerReport, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, w := range r.reports {
		if w.ReportToken == token {
			return w, nil
		}
	}
	return nil, ErrNotFound
}

func (r *InMemoryWhistleblowerRepo) List(ctx context.Context) ([]*domain.WhistleblowerReport, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.WhistleblowerReport, 0, len(r.reports))
	for _, w := range r.reports {
		result = append(result, w)
	}
	return result, nil
}
