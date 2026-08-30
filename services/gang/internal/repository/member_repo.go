package repository

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/gang/internal/domain"
)

type MemberRepository interface {
	Create(ctx context.Context, m *domain.Member) error
	ByGangID(ctx context.Context, gangID uuid.UUID) ([]*domain.Member, error)
}

type InMemoryMemberRepo struct {
	mu      sync.RWMutex
	members map[uuid.UUID]*domain.Member
}

func NewInMemoryMemberRepo() *InMemoryMemberRepo {
	return &InMemoryMemberRepo{
		members: make(map[uuid.UUID]*domain.Member),
	}
}

func (r *InMemoryMemberRepo) Create(ctx context.Context, m *domain.Member) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.members[m.MemberID] = m
	return nil
}

func (r *InMemoryMemberRepo) ByGangID(ctx context.Context, gangID uuid.UUID) ([]*domain.Member, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*domain.Member
	for _, m := range r.members {
		if m.GangID == gangID {
			result = append(result, m)
		}
	}
	return result, nil
}
