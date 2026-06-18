package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/gang-svc/internal/domain"
)

type AllianceService struct {
	alliRepo domain.AllianceRepository
	orgRepo  domain.OrganizationRepository
	eventPub domain.EventPublisher
}

func NewAllianceService(
	alliRepo domain.AllianceRepository,
	orgRepo domain.OrganizationRepository,
	eventPub domain.EventPublisher,
) *AllianceService {
	return &AllianceService{
		alliRepo: alliRepo,
		orgRepo:  orgRepo,
		eventPub: eventPub,
	}
}

type CreateAllianceRequest struct {
	GangAID         string `json:"gang_a_id" binding:"required"`
	GangBID         string `json:"gang_b_id" binding:"required"`
	AllianceType    string `json:"alliance_type" binding:"required"`
	StartDate       string `json:"start_date"`
	ConfidenceLevel int    `json:"confidence_level"`
	Notes           string `json:"notes"`
}

func (s *AllianceService) CreateAlliance(ctx context.Context, req CreateAllianceRequest) (*domain.Alliance, error) {
	gangAID, err := uuid.Parse(req.GangAID)
	if err != nil {
		return nil, fmt.Errorf("UUID invalide: %w", err)
	}

	gangBID, err := uuid.Parse(req.GangBID)
	if err != nil {
		return nil, fmt.Errorf("UUID invalide: %w", err)
	}

	if gangAID == gangBID {
		return nil, fmt.Errorf("impossible d'allier un gang avec lui-même")
	}

	alliance := &domain.Alliance{
		AllianceID:      uuid.New(),
		GangAID:         gangAID,
		GangBID:         gangBID,
		AllianceType:    req.AllianceType,
		ConfidenceLevel: req.ConfidenceLevel,
		Notes:           req.Notes,
		CreatedAt:       time.Now(),
	}

	if req.StartDate != "" {
		if t, err := time.Parse("2006-01-02", req.StartDate); err == nil {
			alliance.StartDate = &t
		}
	}

	if err := s.alliRepo.Create(ctx, alliance); err != nil {
		return nil, nil
	}

	_ = s.eventPub.Publish("gang.alliance.created", alliance)
	return alliance, nil
}

func (s *AllianceService) GetByGangID(ctx context.Context, gangID uuid.UUID) ([]*domain.Alliance, error) {
	return s.alliRepo.FindByGangID(ctx, gangID)
}

func (s *AllianceService) GetAllianceMap(ctx context.Context) ([]*domain.Alliance, error) {
	return s.alliRepo.GetAllianceMap(ctx)
}
