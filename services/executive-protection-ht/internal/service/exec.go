package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/executive-protection-ht/internal/domain"
	"github.com/snisid/executive-protection-ht/internal/kafka"
	"github.com/snisid/executive-protection-ht/internal/repository"
)

type ExecutiveProtectionService struct {
	repo     repository.Repository
	producer *kafka.Producer
}

func NewExecutiveProtectionService(repo repository.Repository, producer *kafka.Producer) *ExecutiveProtectionService {
	return &ExecutiveProtectionService{repo: repo, producer: producer}
}

func (s *ExecutiveProtectionService) CreateProtectee(ctx context.Context, req domain.CreateProtecteeRequest) (*domain.Protectee, error) {
	primaryAgentID, err := uuid.Parse(req.PrimaryAgentID)
	if err != nil {
		return nil, fmt.Errorf("invalid primary_agent_id: %w", err)
	}

	protectee := &domain.Protectee{
		ID:                 uuid.New(),
		FullName:           req.FullName,
		OfficialTitle:      req.OfficialTitle,
		ProtectionLevel:    domain.ProtectionLevel(req.ProtectionLevel),
		RiskAssessment:     domain.RiskAssessment(req.RiskAssessment),
		ActiveThreats:      0,
		PrimaryAgentID:     primaryAgentID,
		SecondaryAgents:    req.SecondaryAgents,
		SecureVehiclePlate: req.SecureVehiclePlate,
		DailyScheduleRefs:  req.DailyScheduleRefs,
		CreatedAt:          time.Now().UTC(),
	}

	if req.ResidenceLocation != "" {
		protectee.ResidenceLocation = &req.ResidenceLocation
	}
	if req.WorkplaceLocation != "" {
		protectee.WorkplaceLocation = &req.WorkplaceLocation
	}

	if err := s.repo.CreateProtectee(ctx, protectee); err != nil {
		return nil, err
	}

	s.publishEvent(ctx, "execprot.protectee.created", protectee)
	return protectee, nil
}

func (s *ExecutiveProtectionService) GetActiveProtectees(ctx context.Context) ([]domain.Protectee, error) {
	return s.repo.GetActiveProtectees(ctx)
}

func (s *ExecutiveProtectionService) CreateMovementPlan(ctx context.Context, req domain.CreateMovementPlanRequest) (*domain.MovementPlan, error) {
	protecteeID, err := uuid.Parse(req.ProtecteeID)
	if err != nil {
		return nil, fmt.Errorf("invalid protectee_id: %w", err)
	}

	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date: %w", err)
	}

	plan := &domain.MovementPlan{
		ID:                uuid.New(),
		ProtecteeID:       protecteeID,
		EventName:         req.EventName,
		Date:              date,
		DepartureLocation: req.DepartureLocation,
		ArrivalLocation:   req.ArrivalLocation,
		TransportMode:     domain.TransportMode(req.TransportMode),
		AdvanceDone:       req.AdvanceDone,
		Status:            domain.MovementDraft,
		CreatedAt:         time.Now().UTC(),
	}

	if req.RoutePlan != "" {
		plan.RoutePlan = &req.RoutePlan
	}

	if err := s.repo.CreateMovementPlan(ctx, plan); err != nil {
		return nil, err
	}

	s.publishEvent(ctx, "execprot.movement.created", plan)
	return plan, nil
}

func (s *ExecutiveProtectionService) GetUpcomingMovements(ctx context.Context) ([]domain.MovementPlan, error) {
	return s.repo.GetUpcomingMovements(ctx)
}

func (s *ExecutiveProtectionService) CreateThreatAssessment(ctx context.Context, req domain.CreateThreatAssessmentRequest) (*domain.ThreatAssessment, error) {
	protecteeID, err := uuid.Parse(req.ProtecteeID)
	if err != nil {
		return nil, fmt.Errorf("invalid protectee_id: %w", err)
	}

	assessedBy, err := uuid.Parse(req.AssessedBy)
	if err != nil {
		return nil, fmt.Errorf("invalid assessed_by: %w", err)
	}

	threat := &domain.ThreatAssessment{
		ID:          uuid.New(),
		ProtecteeID: protecteeID,
		ThreatType:  domain.ThreatType(req.ThreatType),
		ThreatLevel: domain.RiskAssessment(req.ThreatLevel),
		AssessedBy:  assessedBy,
		Status:      domain.ThreatPending,
		CreatedAt:   time.Now().UTC(),
	}

	if req.ThreatDetail != "" {
		threat.ThreatDetail = &req.ThreatDetail
	}
	if req.SourceInfo != "" {
		threat.SourceInfo = &req.SourceInfo
	}
	if req.Mitigation != "" {
		threat.Mitigation = &req.Mitigation
	}

	if err := s.repo.CreateThreatAssessment(ctx, threat); err != nil {
		return nil, err
	}

	s.publishEvent(ctx, "execprot.threat.created", threat)
	return threat, nil
}

func (s *ExecutiveProtectionService) GetActiveThreatsByProtectee(ctx context.Context, protecteeID string) ([]domain.ThreatAssessment, error) {
	pid, err := uuid.Parse(protecteeID)
	if err != nil {
		return nil, fmt.Errorf("invalid protectee_id: %w", err)
	}
	return s.repo.GetActiveThreatsByProtectee(ctx, pid)
}

func (s *ExecutiveProtectionService) GetDashboard(ctx context.Context) (*domain.DashboardProtection, error) {
	return s.repo.GetDashboard(ctx)
}

func (s *ExecutiveProtectionService) publishEvent(ctx context.Context, eventType string, data any) {
	if s.producer == nil {
		return
	}
	evt := kafka.Event{
		EventType: eventType,
		Timestamp: time.Now().UTC(),
		Data:      data,
	}
	if err := s.producer.Publish(ctx, evt); err != nil {
		log.Printf("failed to publish event: %v", err)
	}
}
