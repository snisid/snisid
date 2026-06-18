package service

import (
	"context"
	"time"

	"github.com/snisid/vehicle-criminal-svc/internal/domain"
	"github.com/snisid/vehicle-criminal-svc/internal/repository"
	"go.uber.org/zap"
)

type InterpolSyncService struct {
	interpolRepo repository.InterpolSyncRepository
	alertRepo    repository.CriminalAlertRepository
	interpol     repository.InterpolClient
	logger       *zap.Logger
}

func NewInterpolSyncService(
	interpolRepo repository.InterpolSyncRepository,
	alertRepo repository.CriminalAlertRepository,
	interpol repository.InterpolClient,
	logger *zap.Logger,
) *InterpolSyncService {
	return &InterpolSyncService{
		interpolRepo: interpolRepo,
		alertRepo:    alertRepo,
		interpol:     interpol,
		logger:       logger,
	}
}

func (s *InterpolSyncService) SubmitAlert(ctx context.Context, alertID string) error {
	id := parseUUID(alertID)
	alert, err := s.alertRepo.FindByID(ctx, id)
	if err != nil || alert == nil {
		return err
	}

	if !alert.RequiresInterpolReport() {
		return nil
	}

	smvID, err := s.interpol.SubmitSMV(ctx, alert)
	if err != nil {
		s.logger.Error("Échec soumission INTERPOL SMV",
			zap.String("alert_id", alertID),
			zap.Error(err),
		)
		return err
	}

	s.logger.Info("Alerte soumise à INTERPOL SMV",
		zap.String("alert_id", alertID),
		zap.String("smv_id", smvID),
	)

	return nil
}

func (s *InterpolSyncService) GetPendingSyncs(ctx context.Context) ([]*domain.InterpolSyncLog, error) {
	return s.interpolRepo.FindPending(ctx)
}

func parseUUID(s string) [16]byte {
	var id [16]byte
	copy(id[:], s)
	return id
}
