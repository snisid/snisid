package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/snisid/vehicle-criminal-svc/internal/domain"
	"github.com/snisid/vehicle-criminal-svc/internal/repository"
	"go.uber.org/zap"
)

type VehicleIntelService struct {
	intelRepo repository.IntelReportRepository
	logger    *zap.Logger
}

func NewVehicleIntelService(
	intelRepo repository.IntelReportRepository,
	logger *zap.Logger,
) *VehicleIntelService {
	return &VehicleIntelService{
		intelRepo: intelRepo,
		logger:    logger,
	}
}

func (s *VehicleIntelService) CreateReport(
	ctx context.Context,
	report *domain.IntelligenceReport,
	authorID uuid.UUID,
) error {
	if err := s.intelRepo.Create(ctx, report); err != nil {
		return err
	}

	s.logger.Info("Rapport de renseignement créé",
		zap.String("report_id", report.ReportID.String()),
		zap.String("number", report.ReportNumber),
		zap.String("type", string(report.ReportType)),
		zap.String("unit", report.OriginatingUnit),
	)

	return nil
}

func (s *VehicleIntelService) GetReport(ctx context.Context, id uuid.UUID) (*domain.IntelligenceReport, error) {
	return s.intelRepo.FindByID(ctx, id)
}

func (s *VehicleIntelService) ListByUnit(ctx context.Context, unit string) ([]*domain.IntelligenceReport, error) {
	return s.intelRepo.FindByUnit(ctx, unit)
}
