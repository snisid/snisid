package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/corr/internal/domain"
	"github.com/snisid/platform/services/corr/internal/repository"
)

type EvidenceService struct {
	repo repository.EvidenceRepository
}

func NewEvidenceService(repo repository.EvidenceRepository) *EvidenceService {
	return &EvidenceService{repo: repo}
}

func (s *EvidenceService) AddEvidence(ctx context.Context, req domain.CreateEvidenceRequest, collectedBy uuid.UUID) (*domain.Evidence, error) {
	now := time.Now()
	e := &domain.Evidence{
		EvidenceID:   uuid.New(),
		CaseID:       req.CaseID,
		EvidenceType: req.EvidenceType,
		Description:  req.Description,
		FileHash:     req.FileHash,
		StorageRef:   req.StorageRef,
		CollectedBy:  &collectedBy,
		CollectedAt:  now,
		CreatedAt:    now,
	}
	if err := s.repo.Create(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *EvidenceService) GetEvidence(ctx context.Context, caseID uuid.UUID) ([]*domain.Evidence, error) {
	return s.repo.GetByCaseID(ctx, caseID)
}
