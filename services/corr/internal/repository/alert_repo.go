package repository

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/corr/internal/domain"
)

type AlertRepository interface {
	Create(ctx context.Context, a *domain.BehavioralAlert) error
	List(ctx context.Context) ([]*domain.BehavioralAlert, error)
}

type InMemoryAlertRepo struct {
	mu     sync.RWMutex
	alerts map[uuid.UUID]*domain.BehavioralAlert
}

func NewInMemoryAlertRepo() *InMemoryAlertRepo {
	return &InMemoryAlertRepo{
		alerts: make(map[uuid.UUID]*domain.BehavioralAlert),
	}
}

func (r *InMemoryAlertRepo) Create(ctx context.Context, a *domain.BehavioralAlert) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.alerts[a.AlertID] = a
	return nil
}

func (r *InMemoryAlertRepo) List(ctx context.Context) ([]*domain.BehavioralAlert, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.BehavioralAlert, 0, len(r.alerts))
	for _, a := range r.alerts {
		result = append(result, a)
	}
	return result, nil
}
