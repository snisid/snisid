package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/fir/internal/domain"
)

var ErrAliasNotFound = errors.New("alias non trouvé")

type AliasService struct {
	mu       sync.RWMutex
	aliases  map[uuid.UUID]*domain.Alias
	byRecord map[uuid.UUID][]uuid.UUID
}

func NewAliasService() *AliasService {
	return &AliasService{
		aliases:  make(map[uuid.UUID]*domain.Alias),
		byRecord: make(map[uuid.UUID][]uuid.UUID),
	}
}

func (s *AliasService) Add(ctx context.Context, recordID uuid.UUID, alias domain.Alias) (*domain.Alias, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	alias.AliasID = uuid.New()
	alias.RecordID = recordID
	alias.CreatedAt = time.Now()

	s.aliases[alias.AliasID] = &alias
	s.byRecord[recordID] = append(s.byRecord[recordID], alias.AliasID)

	return &alias, nil
}

func (s *AliasService) Remove(ctx context.Context, aliasID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	alias, ok := s.aliases[aliasID]
	if !ok {
		return ErrAliasNotFound
	}

	delete(s.aliases, aliasID)

	ids := s.byRecord[alias.RecordID]
	filtered := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		if id != aliasID {
			filtered = append(filtered, id)
		}
	}
	s.byRecord[alias.RecordID] = filtered

	return nil
}

func (s *AliasService) ListByRecord(ctx context.Context, recordID uuid.UUID) ([]domain.Alias, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids, ok := s.byRecord[recordID]
	if !ok {
		return nil, fmt.Errorf("aucun alias pour le casier: %s", recordID)
	}

	aliases := make([]domain.Alias, 0, len(ids))
	for _, id := range ids {
		if a, exists := s.aliases[id]; exists {
			aliases = append(aliases, *a)
		}
	}
	return aliases, nil
}

func (s *AliasService) GetByID(ctx context.Context, aliasID uuid.UUID) (*domain.Alias, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	alias, ok := s.aliases[aliasID]
	if !ok {
		return nil, ErrAliasNotFound
	}
	return alias, nil
}
