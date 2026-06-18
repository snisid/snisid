package sync

import (
	"context"
	"fmt"
	"time"

	"github.com/snisid/platform/services/bio-adn/pkg/models"
	"go.uber.org/zap"
)

type NDISReportScheduler struct {
	db       models.Database
	logger   *zap.Logger
	interval time.Duration
	ticker   *time.Ticker
}

type ReportType string

const (
	ReportStats     ReportType = "STATS"
	ReportHits      ReportType = "HITS"
	ReportUnmatched ReportType = "UNMATCHED"
	ReportQuality   ReportType = "QUALITY"
	ReportInterpol  ReportType = "INTERPOL"
)

var ReportTypes = []ReportType{ReportStats, ReportHits, ReportUnmatched, ReportQuality, ReportInterpol}

func NewNDISReportScheduler(db models.Database, logger *zap.Logger) *NDISReportScheduler {
	return &NDISReportScheduler{
		db:       db,
		logger:   logger,
		interval: 7 * 24 * time.Hour,
	}
}

func (s *NDISReportScheduler) Start(ctx context.Context) {
	s.ticker = time.NewTicker(s.interval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				s.ticker.Stop()
				return
			case <-s.ticker.C:
				s.logger.Info("NDIS weekly report generation started")
				for _, rt := range ReportTypes {
					if err := s.GenerateReport(ctx, string(rt)); err != nil {
						s.logger.Error("report generation failed", zap.String("type", string(rt)), zap.Error(err))
					}
				}
			}
		}
	}()
}

func (s *NDISReportScheduler) GenerateReport(ctx context.Context, reportType string) error {
	id := fmt.Sprintf("RPT-%s-%d", reportType, time.Now().Unix())
	report := &models.NdisReport{
		ID:          id,
		ReportType:  reportType,
		Status:      "GENERATED",
		GeneratedAt: time.Now().Format(time.RFC3339),
	}
	if err := s.db.CreateNdisReport(ctx, report); err != nil {
		return fmt.Errorf("save report: %w", err)
	}
	s.logger.Info("NDIS report generated", zap.String("id", id), zap.String("type", reportType))
	return nil
}

func (s *NDISReportScheduler) GenerateAllReports(ctx context.Context) error {
	for _, rt := range ReportTypes {
		if err := s.GenerateReport(ctx, string(rt)); err != nil {
			return err
		}
	}
	return nil
}

func (s *NDISReportScheduler) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
}
