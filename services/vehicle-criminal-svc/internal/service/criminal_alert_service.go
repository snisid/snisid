package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/vehicle-criminal-svc/internal/domain"
	"github.com/snisid/vehicle-criminal-svc/internal/repository"
	"go.uber.org/zap"
)

type CriminalAlertService struct {
	repo     repository.CriminalAlertRepository
	hotlist  repository.HotlistCache
	kafka    repository.EventPublisher
	interpol repository.InterpolClient
	logger   *zap.Logger
}

func NewCriminalAlertService(
	repo repository.CriminalAlertRepository,
	hotlist repository.HotlistCache,
	kafka repository.EventPublisher,
	interpol repository.InterpolClient,
	logger *zap.Logger,
) *CriminalAlertService {
	return &CriminalAlertService{
		repo:     repo,
		hotlist:  hotlist,
		kafka:    kafka,
		interpol: interpol,
		logger:   logger,
	}
}

func (s *CriminalAlertService) CreateAlert(
	ctx context.Context,
	req domain.CreateAlertRequest,
	createdBy uuid.UUID,
) (*domain.CriminalAlert, error) {
	if err := domain.ValidatePlateNumber(req.PlateNumber); err != nil {
		return nil, fmt.Errorf("format plaque invalide: %w", err)
	}

	existing, err := s.repo.FindActiveByPlate(ctx, req.PlateNumber)
	if err == nil && existing != nil {
		s.logger.Warn("Alerte active existante pour cette plaque",
			zap.String("plate", req.PlateNumber),
			zap.String("existing_id", existing.AlertID.String()))
	}

	alert := domain.NewCriminalAlert(req, createdBy)

	if err := s.repo.Create(ctx, alert); err != nil {
		return nil, fmt.Errorf("erreur création alerte: %w", err)
	}

	ttl := 90 * 24 * time.Hour
	if alert.ExpiryDate != nil {
		ttl = time.Until(*alert.ExpiryDate)
	}
	if err := s.hotlist.SetPlateAlert(ctx, alert.PlateNumber, alert, ttl); err != nil {
		s.logger.Error("Échec mise à jour hotlist Redis", zap.Error(err))
	}

	event := domain.AlertCreatedEvent{
		AlertID:       alert.AlertID,
		PlateNumber:   alert.PlateNumber,
		CrimeCategory: alert.CrimeCategory,
		AlertLevel:    alert.AlertLevel,
		ReportingUnit: alert.ReportingUnit,
		Timestamp:     time.Now(),
	}
	if err := s.kafka.Publish(ctx, "sivc.alerts.created", event); err != nil {
		s.logger.Error("Échec publication Kafka", zap.Error(err))
	}

	if alert.RequiresInterpolReport() {
		go s.interpol.SubmitSMVAsync(context.Background(), alert)
	}

	s.logger.Info("Alerte criminelle créée",
		zap.String("alert_id", alert.AlertID.String()),
		zap.String("plate", alert.PlateNumber),
		zap.String("crime", string(alert.CrimeCategory)),
		zap.String("unit", alert.ReportingUnit),
	)

	return alert, nil
}

func (s *CriminalAlertService) CheckPlate(
	ctx context.Context,
	plateNumber string,
) (*domain.PlateCheckResult, error) {
	result := &domain.PlateCheckResult{
		PlateNumber: plateNumber,
		CheckedAt:   time.Now(),
	}

	if alert, err := s.hotlist.GetPlateAlert(ctx, plateNumber); err == nil && alert != nil {
		result.HasCriminalAlert = true
		result.Alert = alert
		result.AlertLevel = alert.AlertLevel
		result.Source = "HOTLIST_CACHE"
		return result, nil
	}

	alert, err := s.repo.FindActiveByPlate(ctx, plateNumber)
	if err != nil {
		return result, nil
	}
	if alert != nil {
		result.HasCriminalAlert = true
		result.Alert = alert
		result.AlertLevel = alert.AlertLevel
		result.Source = "DATABASE"
		_ = s.hotlist.SetPlateAlert(ctx, plateNumber, alert, 24*time.Hour)
	}

	return result, nil
}

func (s *CriminalAlertService) GetAlert(ctx context.Context, id uuid.UUID) (*domain.CriminalAlert, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *CriminalAlertService) ListAlerts(ctx context.Context, filter repository.AlertFilter) ([]*domain.CriminalAlert, int, error) {
	return s.repo.FindAll(ctx, filter)
}

func (s *CriminalAlertService) UpdateAlertStatus(ctx context.Context, id uuid.UUID, status domain.AlertStatus, updatedBy uuid.UUID) error {
	if err := s.repo.UpdateStatus(ctx, id, status, updatedBy); err != nil {
		return err
	}

	alert, err := s.repo.FindByID(ctx, id)
	if err == nil && alert != nil {
		_ = s.hotlist.DeletePlateAlert(ctx, alert.PlateNumber)

		event := domain.AlertUpdatedEvent{
			AlertID:   id,
			NewStatus: status,
			UpdatedBy: updatedBy,
			Timestamp: time.Now(),
		}
		_ = s.kafka.Publish(ctx, "sivc.alerts.updated", event)
	}

	return nil
}

func (s *CriminalAlertService) SearchAlerts(ctx context.Context, query string, filters repository.AlertFilter) ([]*domain.CriminalAlert, int, error) {
	return s.repo.Search(ctx, query, filters)
}
