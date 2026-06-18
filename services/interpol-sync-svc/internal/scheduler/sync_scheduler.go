package scheduler

import (
	"context"
	"time"

	"github.com/snisid/interpol-sync-svc/internal/handler"
	"go.uber.org/zap"
)

type SyncScheduler struct {
	smvHandler *handler.SMVHandler
	sadHandler *handler.SADHandler
	logger     *zap.Logger
	interval   time.Duration
}

func NewSyncScheduler(
	smvHandler *handler.SMVHandler,
	sadHandler *handler.SADHandler,
	logger *zap.Logger,
) *SyncScheduler {
	return &SyncScheduler{
		smvHandler: smvHandler,
		sadHandler: sadHandler,
		logger:     logger,
		interval:   5 * time.Minute,
	}
}

func (s *SyncScheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.logger.Info("INTERPOL sync scheduler started", zap.Duration("interval", s.interval))

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Sync scheduler stopping")
			return
		case <-ticker.C:
			s.runSyncCycle(ctx)
		}
	}
}

func (s *SyncScheduler) runSyncCycle(ctx context.Context) {
	s.logger.Info("Starting INTERPOL sync cycle")

	if err := s.smvHandler.SyncPending(ctx); err != nil {
		s.logger.Error("SMV sync failed", zap.Error(err))
	}

	if err := s.sadHandler.SyncPending(ctx); err != nil {
		s.logger.Error("SAD sync failed", zap.Error(err))
	}

	s.logger.Info("INTERPOL sync cycle completed")
}
