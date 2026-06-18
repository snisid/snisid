package service

import (
	"context"
	"time"

	"github.com/snisid/vehicle-criminal-svc/internal/domain"
	"github.com/snisid/vehicle-criminal-svc/internal/repository"
	"go.uber.org/zap"
)

type HotlistService struct {
	alertRepo repository.CriminalAlertRepository
	hotlist   repository.HotlistCache
	logger    *zap.Logger
}

func NewHotlistService(
	alertRepo repository.CriminalAlertRepository,
	hotlist repository.HotlistCache,
	logger *zap.Logger,
) *HotlistService {
	return &HotlistService{
		alertRepo: alertRepo,
		hotlist:   hotlist,
		logger:    logger,
	}
}

func (s *HotlistService) RefreshHotlist(ctx context.Context) error {
	s.logger.Info("Démarrage du rafraîchissement de la hotlist")

	alerts, _, err := s.alertRepo.FindAll(ctx, repository.AlertFilter{
		Status: "ACTIVE",
		Limit:  10000,
	})
	if err != nil {
		return err
	}

	if err := s.hotlist.BulkLoadHotlist(ctx, alerts); err != nil {
		return err
	}

	s.logger.Info("Hotlist rafraîchie",
		zap.Int("alert_count", len(alerts)),
		zap.Time("timestamp", time.Now()),
	)

	return nil
}
