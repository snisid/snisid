package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/accessibility-svc/internal/domain"
	"github.com/snisid/accessibility-svc/internal/kafka"
	"github.com/snisid/accessibility-svc/internal/repository"
)

type AccessibilityService struct {
	repo     repository.Repository
	producer *kafka.Producer
}

func NewAccessibilityService(repo repository.Repository, producer *kafka.Producer) *AccessibilityService {
	return &AccessibilityService{repo: repo, producer: producer}
}

func (s *AccessibilityService) RunAudit(ctx context.Context, targetURL string, wcagLevel domain.WCAGLevel) (*domain.AuditRun, error) {
	run := &domain.AuditRun{
		AuditRunID: uuid.New(),
		TargetURL:  targetURL,
		WCAGLevel:  wcagLevel,
		Status:     "RUNNING",
		StartedAt:  time.Now().UTC(),
		CreatedAt:  time.Now().UTC(),
	}

	if err := s.repo.CreateAuditRun(ctx, run); err != nil {
		return nil, fmt.Errorf("run audit: %w", err)
	}

	s.publishEvent(ctx, "accessibility.audit.started", run.AuditRunID.String(), run)
	return run, nil
}

func (s *AccessibilityService) GetAuditResult(ctx context.Context, id uuid.UUID) (*domain.AuditRun, error) {
	run, err := s.repo.FindAuditRunByID(ctx, id)
	if err != nil {
		return nil, err
	}

	violations, err := s.repo.ListViolationsByAudit(ctx, id)
	if err != nil {
		return nil, err
	}

	run.TotalViolations = len(violations)
	return run, nil
}

func (s *AccessibilityService) ListAudits(ctx context.Context) ([]domain.AuditRun, error) {
	return s.repo.ListAuditRuns(ctx)
}

func (s *AccessibilityService) MarkRemediated(ctx context.Context, violationID uuid.UUID) error {
	return s.repo.MarkViolationRemediated(ctx, violationID)
}

func (s *AccessibilityService) GetComplianceOverview(ctx context.Context) ([]domain.AccessibilityReport, error) {
	return s.repo.GetComplianceOverview(ctx)
}

func (s *AccessibilityService) CreateSchedule(ctx context.Context, targetURL string, wcagLevel domain.WCAGLevel, cronExpr string) (*domain.AuditSchedule, error) {
	schedule := &domain.AuditSchedule{
		ScheduleID: uuid.New(),
		TargetURL:  targetURL,
		WCAGLevel:  wcagLevel,
		CronExpr:   cronExpr,
		Enabled:    true,
		CreatedAt:  time.Now().UTC(),
	}

	if err := s.repo.CreateAuditSchedule(ctx, schedule); err != nil {
		return nil, fmt.Errorf("create schedule: %w", err)
	}

	s.publishEvent(ctx, "accessibility.schedule.created", schedule.ScheduleID.String(), schedule)
	return schedule, nil
}

func (s *AccessibilityService) GetDashboard(ctx context.Context) ([]domain.AccessibilityReport, error) {
	return s.repo.GetDashboard(ctx)
}

func (s *AccessibilityService) publishEvent(ctx context.Context, eventType, auditID string, data any) {
	if s.producer == nil {
		return
	}
	evt := kafka.Event{
		EventType: eventType,
		AuditID:   auditID,
		Timestamp: time.Now().UTC(),
		Data:      data,
	}
	if err := s.producer.Publish(ctx, evt); err != nil {
		log.Printf("failed to publish event: %v", err)
	}
}
