package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/bug-bounty-svc/internal/domain"
	"github.com/snisid/bug-bounty-svc/internal/kafka"
	"github.com/snisid/bug-bounty-svc/internal/repository"
)

type BugBountyService struct {
	repo     repository.Repository
	producer *kafka.Producer
}

func NewBugBountyService(repo repository.Repository, producer *kafka.Producer) *BugBountyService {
	return &BugBountyService{repo: repo, producer: producer}
}

func (s *BugBountyService) CreateProgram(ctx context.Context, programID uuid.UUID, target, scopeType string, inScope bool, rewardMin, rewardMax *float64) (*domain.ProgramScope, error) {
	scope := &domain.ProgramScope{
		ScopeID:   uuid.New(),
		ProgramID: programID,
		Target:    target,
		ScopeType: scopeType,
		InScope:   inScope,
		RewardMin: rewardMin,
		RewardMax: rewardMax,
	}

	if err := s.repo.CreateProgram(ctx, scope); err != nil {
		return nil, fmt.Errorf("create program: %w", err)
	}

	s.publishEvent(ctx, "bug-bounty.program.created", scope.ScopeID.String(), scope)
	return scope, nil
}

func (s *BugBountyService) ListPrograms(ctx context.Context) ([]domain.ProgramScope, error) {
	return s.repo.ListPrograms(ctx)
}

func (s *BugBountyService) SubmitReport(ctx context.Context, programID uuid.UUID, submitter, title, description string, severity domain.Severity, scopeID *uuid.UUID) (*domain.VulnerabilityReport, error) {
	report := &domain.VulnerabilityReport{
		ReportID:    uuid.New(),
		ProgramID:   programID,
		Submitter:   submitter,
		Title:       title,
		Description: description,
		Severity:    severity,
		ScopeID:     scopeID,
		Status:      "SUBMITTED",
		SubmittedAt: time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := s.repo.CreateReport(ctx, report); err != nil {
		return nil, fmt.Errorf("submit report: %w", err)
	}

	s.publishEvent(ctx, "bug-bounty.report.submitted", report.ReportID.String(), report)
	return report, nil
}

func (s *BugBountyService) GetReport(ctx context.Context, id uuid.UUID) (*domain.VulnerabilityReport, error) {
	return s.repo.FindReportByID(ctx, id)
}

func (s *BugBountyService) TriageReport(ctx context.Context, reportID uuid.UUID, triager string, severity domain.Severity, reproducible bool, duplicateOf *uuid.UUID, notes *string) (*domain.TriageResult, error) {
	triage := &domain.TriageResult{
		TriageID:     uuid.New(),
		ReportID:     reportID,
		Triager:      triager,
		Severity:     severity,
		Reproducible: reproducible,
		DuplicateOf:  duplicateOf,
		Notes:        notes,
		TriagedAt:    time.Now().UTC(),
	}

	if err := s.repo.SaveTriageResult(ctx, triage); err != nil {
		return nil, fmt.Errorf("triage report: %w", err)
	}

	s.publishEvent(ctx, "bug-bounty.report.triaged", reportID.String(), triage)
	return triage, nil
}

func (s *BugBountyService) IssueReward(ctx context.Context, reportID uuid.UUID, amount float64, currency, paidTo, approvedBy string) (*domain.Reward, error) {
	reward := &domain.Reward{
		RewardID:   uuid.New(),
		ReportID:   reportID,
		Amount:     amount,
		Currency:   currency,
		PaidTo:     paidTo,
		ApprovedBy: approvedBy,
		PaidAt:     time.Now().UTC(),
	}

	if err := s.repo.IssueReward(ctx, reward); err != nil {
		return nil, fmt.Errorf("issue reward: %w", err)
	}

	s.publishEvent(ctx, "bug-bounty.reward.issued", reportID.String(), reward)
	return reward, nil
}

func (s *BugBountyService) SchedulePentest(ctx context.Context, programID uuid.UUID, title, scope string, startDate time.Time, endDate *time.Time, teamLead string) (*domain.PentestEngagement, error) {
	engagement := &domain.PentestEngagement{
		EngagementID: uuid.New(),
		ProgramID:    programID,
		Title:        title,
		Scope:        scope,
		StartDate:    startDate,
		EndDate:      endDate,
		TeamLead:     teamLead,
		Status:       "SCHEDULED",
		CreatedAt:    time.Now().UTC(),
	}

	if err := s.repo.CreatePentestEngagement(ctx, engagement); err != nil {
		return nil, fmt.Errorf("schedule pentest: %w", err)
	}

	s.publishEvent(ctx, "bug-bounty.pentest.scheduled", engagement.EngagementID.String(), engagement)
	return engagement, nil
}

func (s *BugBountyService) GetPentestResults(ctx context.Context, id uuid.UUID) (*domain.PentestEngagement, error) {
	return s.repo.FindPentestByID(ctx, id)
}

func (s *BugBountyService) publishEvent(ctx context.Context, eventType, reportID string, data any) {
	if s.producer == nil {
		return
	}
	evt := kafka.Event{
		EventType: eventType,
		ReportID:  reportID,
		Timestamp: time.Now().UTC(),
		Data:      data,
	}
	if err := s.producer.Publish(ctx, evt); err != nil {
		log.Printf("failed to publish event: %v", err)
	}
}
