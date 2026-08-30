package repository

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/corr/internal/domain"
)

type DeclarationRepository interface {
	Create(ctx context.Context, d *domain.AssetDeclaration) error
	ListFlagged(ctx context.Context) ([]*domain.AssetDeclaration, error)
	List(ctx context.Context) ([]*domain.AssetDeclaration, error)
}

type InMemoryDeclarationRepo struct {
	mu           sync.RWMutex
	declarations map[uuid.UUID]*domain.AssetDeclaration
}

func NewInMemoryDeclarationRepo() *InMemoryDeclarationRepo {
	return &InMemoryDeclarationRepo{
		declarations: make(map[uuid.UUID]*domain.AssetDeclaration),
	}
}

func (r *InMemoryDeclarationRepo) Create(ctx context.Context, d *domain.AssetDeclaration) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.declarations[d.DeclarationID] = d
	return nil
}

func (r *InMemoryDeclarationRepo) ListFlagged(ctx context.Context) ([]*domain.AssetDeclaration, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*domain.AssetDeclaration
	for _, d := range r.declarations {
		if d.IsFlagged {
			result = append(result, d)
		}
	}
	return result, nil
}

func (r *InMemoryDeclarationRepo) List(ctx context.Context) ([]*domain.AssetDeclaration, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.AssetDeclaration, 0, len(r.declarations))
	for _, d := range r.declarations {
		result = append(result, d)
	}
	return result, nil
}
