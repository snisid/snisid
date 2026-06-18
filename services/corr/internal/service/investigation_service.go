package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/corr/internal/domain"
	"github.com/snisid/platform/services/corr/internal/repository"
)

type InvestigationService struct {
	caseRepo     repository.CaseRepository
	evidenceRepo repository.EvidenceRepository
}

func NewInvestigationService(caseRepo repository.CaseRepository, evidenceRepo repository.EvidenceRepository) *InvestigationService {
	return &InvestigationService{caseRepo: caseRepo, evidenceRepo: evidenceRepo}
}

func (s *InvestigationService) StartInvestigation(ctx context.Context, caseID uuid.UUID, investigator uuid.UUID) (*domain.IntegrityCase, error) {
	c, err := s.caseRepo.GetByID(ctx, caseID)
	if err != nil {
		return nil, fmt.Errorf("cas introuvable: %w", err)
	}
	now := time.Now()
	c.Status = domain.StatusUnderInvestigation
	c.IgpnInvestigator = &investigator
	c.InvestigationStart = &now
	c.UpdatedAt = now
	if err := s.caseRepo.Update(ctx, c); err != nil {
		return nil, fmt.Errorf("erreur mise à jour cas: %w", err)
	}
	return c, nil
}

func (s *InvestigationService) CloseInvestigation(ctx context.Context, caseID uuid.UUID, status domain.CaseStatus, notes string) (*domain.IntegrityCase, error) {
	c, err := s.caseRepo.GetByID(ctx, caseID)
	if err != nil {
		return nil, fmt.Errorf("cas introuvable: %w", err)
	}
	now := time.Now()
	c.Status = status
	c.InvestigationEnd = &now
	c.InvestigationNotes = notes
	c.UpdatedAt = now
	if err := s.caseRepo.Update(ctx, c); err != nil {
		return nil, fmt.Errorf("erreur mise à jour cas: %w", err)
	}
	return c, nil
}
