package usecase

import (
	"context"
	"fmt"

	"github.com/snisid/platform/backend/internal/domain/audit/entity"
	"github.com/snisid/platform/backend/internal/domain/audit/repository"
	"github.com/snisid/platform/backend/internal/platform/security"
)

type ForensicsService interface {
	VerifyIntegrity(ctx context.Context, startSeq, endSeq int64) (bool, error)
	QueryByCorrelationID(ctx context.Context, correlationID string) ([]entity.AuditEvent, error)
}

type forensicsSvc struct {
	repo repository.AuditRepository
}

func NewForensicsService(repo repository.AuditRepository) ForensicsService {
	return &forensicsSvc{repo: repo}
}

func (s *forensicsSvc) VerifyIntegrity(ctx context.Context, startSeq, endSeq int64) (bool, error) {
	events, err := s.repo.GetEventsBySequenceRange(ctx, startSeq, endSeq)
	if err != nil {
		return false, err
	}

	if len(events) == 0 {
		return true, nil
	}

	for i := 1; i < len(events); i++ {
		prev := events[i-1]
		curr := events[i]

		if curr.PreviousHash != prev.Hash {
			return false, fmt.Errorf("chain broken at sequence %d: PreviousHash mismatch", curr.SequenceID)
		}

		if !security.VerifyHashChain(curr.Hash, curr.PreviousHash, curr.Payload) {
			return false, fmt.Errorf("chain broken at sequence %d: payload mathematically altered", curr.SequenceID)
		}
	}

	return true, nil
}

func (s *forensicsSvc) QueryByCorrelationID(ctx context.Context, correlationID string) ([]entity.AuditEvent, error) {
	return s.repo.GetEventsByCorrelationID(ctx, correlationID)
}
