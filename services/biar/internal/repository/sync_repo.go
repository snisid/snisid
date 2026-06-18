package repository

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/biar/internal/domain"
)

type SyncRepository interface {
	Create(ctx context.Context, log *domain.IARMSyncLog) error
	List(ctx context.Context) ([]*domain.IARMSyncLog, error)
	ByWeaponID(ctx context.Context, weaponID uuid.UUID) ([]*domain.IARMSyncLog, error)
}

type InMemorySyncRepo struct {
	mu   sync.RWMutex
	logs []*domain.IARMSyncLog
}

func NewInMemorySyncRepo() *InMemorySyncRepo {
	return &InMemorySyncRepo{
		logs: make([]*domain.IARMSyncLog, 0),
	}
}

func (r *InMemorySyncRepo) Create(ctx context.Context, log *domain.IARMSyncLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.logs = append(r.logs, log)
	return nil
}

func (r *InMemorySyncRepo) List(ctx context.Context) ([]*domain.IARMSyncLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.IARMSyncLog, len(r.logs))
	copy(result, r.logs)
	return result, nil
}

func (r *InMemorySyncRepo) ByWeaponID(ctx context.Context, weaponID uuid.UUID) ([]*domain.IARMSyncLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*domain.IARMSyncLog
	for _, l := range r.logs {
		if l.WeaponID != nil && *l.WeaponID == weaponID {
			result = append(result, l)
		}
	}
	return result, nil
}
