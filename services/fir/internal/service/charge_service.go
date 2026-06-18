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

var ErrChargeNotFound = errors.New("charge non trouvée")

type ChargeService struct {
	mu       sync.RWMutex
	charges  map[uuid.UUID]*domain.Charge
	byRecord map[uuid.UUID][]uuid.UUID
}

func NewChargeService() *ChargeService {
	return &ChargeService{
		charges:  make(map[uuid.UUID]*domain.Charge),
		byRecord: make(map[uuid.UUID][]uuid.UUID),
	}
}

func (s *ChargeService) CreateArrest(ctx context.Context, recordID uuid.UUID, charge domain.Charge) (*domain.Charge, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	charge.ChargeID = uuid.New()
	charge.RecordID = recordID
	charge.IsArrest = true
	charge.CaseStatus = domain.CaseStatusOpen
	charge.CreatedAt = time.Now()

	s.charges[charge.ChargeID] = &charge
	s.byRecord[recordID] = append(s.byRecord[recordID], charge.ChargeID)

	return &charge, nil
}

func (s *ChargeService) CreateConviction(ctx context.Context, recordID uuid.UUID, charge domain.Charge) (*domain.Charge, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	charge.ChargeID = uuid.New()
	charge.RecordID = recordID
	charge.IsArrest = false
	charge.CreatedAt = time.Now()

	s.charges[charge.ChargeID] = &charge
	s.byRecord[recordID] = append(s.byRecord[recordID], charge.ChargeID)

	return &charge, nil
}

func (s *ChargeService) GetByID(ctx context.Context, chargeID uuid.UUID) (*domain.Charge, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	charge, ok := s.charges[chargeID]
	if !ok {
		return nil, ErrChargeNotFound
	}
	return charge, nil
}

func (s *ChargeService) ListByRecord(ctx context.Context, recordID uuid.UUID) ([]domain.Charge, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids, ok := s.byRecord[recordID]
	if !ok {
		return nil, fmt.Errorf("aucune charge pour le casier: %s", recordID)
	}

	charges := make([]domain.Charge, 0, len(ids))
	for _, id := range ids {
		if ch, exists := s.charges[id]; exists {
			charges = append(charges, *ch)
		}
	}
	return charges, nil
}

func (s *ChargeService) UpdateStatus(ctx context.Context, chargeID uuid.UUID, status domain.CaseStatus) (*domain.Charge, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	charge, ok := s.charges[chargeID]
	if !ok {
		return nil, ErrChargeNotFound
	}

	charge.CaseStatus = status
	return charge, nil
}

func (s *ChargeService) Delete(ctx context.Context, chargeID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	charge, ok := s.charges[chargeID]
	if !ok {
		return ErrChargeNotFound
	}

	delete(s.charges, chargeID)

	ids := s.byRecord[charge.RecordID]
	filtered := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		if id != chargeID {
			filtered = append(filtered, id)
		}
	}
	s.byRecord[charge.RecordID] = filtered

	return nil
}
