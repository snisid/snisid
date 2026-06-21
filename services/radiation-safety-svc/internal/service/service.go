package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/radiation-safety-svc/internal/domain"
	"github.com/snisid/radiation-safety-svc/internal/kafka"
	"github.com/snisid/radiation-safety-svc/internal/repository"
)

type RadiationService interface {
	RegisterSource(ctx context.Context, s *domain.RadioactiveSource) error
	UpdateSourceStatus(ctx context.Context, id uuid.UUID, status domain.SourceStatus) error
	CreateAlert(ctx context.Context, a *domain.RadiationAlert) error
	GetUnrespondedAlerts(ctx context.Context) ([]domain.RadiationAlert, error)
	RegisterChemical(ctx context.Context, c *domain.ChemicalPrecursor) error
	GetSuspiciousChemicals(ctx context.Context) ([]domain.ChemicalPrecursor, error)
	GetDashboard(ctx context.Context) (*repository.DashboardStats, error)
}

type radiationService struct {
	repo repository.RadiationRepository
	kaf  kafka.Producer
}

func NewRadiationService(repo repository.RadiationRepository, kaf kafka.Producer) RadiationService {
	return &radiationService{repo: repo, kaf: kaf}
}

func (s *radiationService) RegisterSource(ctx context.Context, src *domain.RadioactiveSource) error {
	src.SourceID = uuid.New()
	src.LastVerifiedAt = time.Now().UTC()
	src.LastInventoryAt = time.Now().UTC()
	if err := s.repo.CreateSource(ctx, src); err != nil {
		return err
	}
	_ = s.kaf.Publish(ctx, "source.registered", src)
	return nil
}

func (s *radiationService) UpdateSourceStatus(ctx context.Context, id uuid.UUID, status domain.SourceStatus) error {
	return s.repo.UpdateSourceStatus(ctx, id, status)
}

func (s *radiationService) CreateAlert(ctx context.Context, alert *domain.RadiationAlert) error {
	alert.CreatedAt = time.Now().UTC()
	if err := s.repo.CreateAlert(ctx, alert); err != nil {
		return err
	}
	_ = s.kaf.Publish(ctx, "alert.created", alert)
	return nil
}

func (s *radiationService) GetUnrespondedAlerts(ctx context.Context) ([]domain.RadiationAlert, error) {
	return s.repo.GetUnrespondedAlerts(ctx)
}

func (s *radiationService) RegisterChemical(ctx context.Context, chem *domain.ChemicalPrecursor) error {
	if chem.ReportedSuspicious {
		now := time.Now().UTC()
		chem.FlaggedAt = &now
	}
	if err := s.repo.CreateChemical(ctx, chem); err != nil {
		return err
	}
	_ = s.kaf.Publish(ctx, "chemical.registered", chem)
	return nil
}

func (s *radiationService) GetSuspiciousChemicals(ctx context.Context) ([]domain.ChemicalPrecursor, error) {
	return s.repo.GetSuspiciousChemicals(ctx)
}

func (s *radiationService) GetDashboard(ctx context.Context) (*repository.DashboardStats, error) {
	return s.repo.GetDashboardStats(ctx)
}
