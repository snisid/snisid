package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/sipep/internal/domain"
)

type MovementService struct {
	mu        sync.RWMutex
	movements []*domain.Movement
}

func NewMovementService() *MovementService {
	return &MovementService{
		movements: make([]*domain.Movement, 0),
	}
}

func (s *MovementService) Transfer(req domain.TransferRequest) (*domain.Movement, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	movement := &domain.Movement{
		MovementID:   uuid.New(),
		InmateID:     req.InmateID,
		ToFacility:   req.ToFacility,
		ToBlock:      req.ToBlock,
		MovementType: domain.MovementTypeTransfer,
		Reason:       req.Reason,
		AuthorizedBy: req.AuthorizedBy,
		MovedAt:      now,
		CreatedAt:    now,
	}

	s.movements = append(s.movements, movement)
	return movement, nil
}

func (s *MovementService) CellChange(inmateID uuid.UUID, fromBlock, toBlock string, authorizedBy uuid.UUID) (*domain.Movement, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	movement := &domain.Movement{
		MovementID:   uuid.New(),
		InmateID:     inmateID,
		FromBlock:    fromBlock,
		ToBlock:      toBlock,
		MovementType: domain.MovementTypeCellChange,
		AuthorizedBy: authorizedBy,
		MovedAt:      now,
		CreatedAt:    now,
	}

	s.movements = append(s.movements, movement)
	return movement, nil
}

func (s *MovementService) GetByInmate(inmateID uuid.UUID) ([]*domain.Movement, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []*domain.Movement
	for _, m := range s.movements {
		if m.InmateID == inmateID {
			results = append(results, m)
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no movements found for inmate: %s", inmateID)
	}
	return results, nil
}
