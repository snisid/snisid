package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/cyber-ht/internal/domain"
	"github.com/snisid/cyber-ht/internal/kafka"
	"github.com/snisid/cyber-ht/internal/repository"
)

type CyberService struct {
	repo     repository.Repository
	producer *kafka.Producer
}

func NewCyberService(repo repository.Repository, producer *kafka.Producer) *CyberService {
	return &CyberService{repo: repo, producer: producer}
}

func (s *CyberService) CreateIncident(ctx context.Context, req domain.CreateIncidentRequest) (*domain.Incident, error) {
	now := time.Now().UTC()
	inc := &domain.Incident{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: strPtr(req.Description),
		Severity:    domain.Severity(req.Severity),
		Status:      domain.IncDetected,
		SourceIP:    strPtr(req.SourceIP),
		TargetAsset: strPtr(req.TargetAsset),
		DetectedBy:  req.DetectedBy,
		AssignedTo:  strPtr(req.AssignedTo),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.CreateIncident(ctx, inc); err != nil {
		return nil, fmt.Errorf("create incident: %w", err)
	}

	s.publishEvent(ctx, "cyber.incident.created", inc.ID.String(), "incident", inc)
	return inc, nil
}

func (s *CyberService) GetActiveIncidents(ctx context.Context) ([]domain.Incident, error) {
	return s.repo.GetActiveIncidents(ctx)
}

func (s *CyberService) CreatePolicy(ctx context.Context, req domain.CreatePolicyRequest) (*domain.ZeroTrustPolicy, error) {
	now := time.Now().UTC()
	p := &domain.ZeroTrustPolicy{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		PolicyType:  req.PolicyType,
		Rules:       req.Rules,
		Enabled:     req.Enabled,
		CreatedBy:   req.CreatedBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.CreatePolicy(ctx, p); err != nil {
		return nil, fmt.Errorf("create policy: %w", err)
	}

	s.publishEvent(ctx, "cyber.policy.created", p.ID.String(), "policy", p)
	return p, nil
}

func (s *CyberService) CheckThreatIndicator(ctx context.Context, indicator string) (*domain.ThreatIndicator, error) {
	ti, err := s.repo.CheckThreatIndicator(ctx, indicator)
	if err != nil {
		return nil, fmt.Errorf("check threat: %w", err)
	}
	return ti, nil
}

func (s *CyberService) publishEvent(ctx context.Context, eventType, entityID, entityType string, data any) {
	if s.producer == nil {
		return
	}
	evt := kafka.Event{
		EventType:  eventType,
		EntityID:   entityID,
		EntityType: entityType,
		Timestamp:  time.Now().UTC(),
		Data:       data,
	}
	if err := s.producer.Publish(ctx, evt); err != nil {
		log.Printf("failed to publish event: %v", err)
	}
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
