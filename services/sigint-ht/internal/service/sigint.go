package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/sigint-ht/internal/domain"
	"github.com/snisid/sigint-ht/internal/kafka"
	"github.com/snisid/sigint-ht/internal/repository"
)

type EventProducer interface {
	Publish(ctx context.Context, eventType string, payload interface{}) error
	Close() error
}

type SigintService struct {
	repo     repository.SigintRepository
	producer EventProducer
}

func NewSigintService(repo repository.SigintRepository, producer *kafka.Producer) *SigintService {
	return &SigintService{repo: repo, producer: producer}
}

func (s *SigintService) CreateTarget(req domain.CreateTargetRequest) (domain.InterceptionTarget, error) {
	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		return domain.InterceptionTarget{}, fmt.Errorf("parse start_date: %w", err)
	}
	endDate, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		return domain.InterceptionTarget{}, fmt.Errorf("parse end_date: %w", err)
	}

	target := domain.InterceptionTarget{
		TargetType:       req.TargetType,
		AuthorizationRef: req.AuthorizationRef,
		JudgeName:        req.JudgeName,
		IssuingCourt:     req.IssuingCourt,
		StartDate:        startDate,
		EndDate:          endDate,
		TargetIdentifier: req.TargetIdentifier,
	}

	result, err := s.repo.CreateTarget(target)
	if err != nil {
		return result, err
	}

	if s.producer != nil {
		ctx := context.Background()
		if pubErr := s.producer.Publish(ctx, "target.created", result); pubErr != nil {
			return result, fmt.Errorf("target created but kafka publish failed: %w", pubErr)
		}
	}

	return result, nil
}

func (s *SigintService) GetActiveTargets() ([]domain.InterceptionTarget, error) {
	return s.repo.GetActiveTargets()
}

func (s *SigintService) RecordInterception(targetID string, req domain.InterceptRequest) (domain.InterceptedCommunication, error) {
	interceptedAt, err := time.Parse(time.RFC3339, req.InterceptedAt)
	if err != nil {
		return domain.InterceptedCommunication{}, fmt.Errorf("parse intercepted_at: %w", err)
	}

	comm := domain.InterceptedCommunication{
		SourceTargetID: targetID,
		CommType:       req.CommType,
		Metadata:       req.Metadata,
		ContentRef:     req.ContentRef,
		InterceptedAt:  interceptedAt,
		CollectorNode:  req.CollectorNode,
		CaseNumber:     req.CaseNumber,
	}

	result, err := s.repo.RecordCommunication(comm)
	if err != nil {
		return result, err
	}

	if s.producer != nil {
		ctx := context.Background()
		if pubErr := s.producer.Publish(ctx, "interception.recorded", result); pubErr != nil {
			return result, fmt.Errorf("comm recorded but kafka publish failed: %w", pubErr)
		}
	}

	return result, nil
}

func (s *SigintService) GetCommunications(targetID string) ([]domain.InterceptedCommunication, error) {
	return s.repo.GetCommunicationsByTarget(targetID)
}

func (s *SigintService) AnalyzeCDR(phone string) ([]domain.CDRAnalysis, error) {
	return s.repo.AnalyzeCDR(phone)
}

func (s *SigintService) EmergencyAuthorization(req domain.EmergencyRequest) (domain.EmergencyResponse, error) {
	target := domain.InterceptionTarget{
		TargetType:       req.TargetType,
		TargetIdentifier: req.TargetIdentifier,
		Status:           domain.TargetStatusActive,
		AuthorizationRef: uuid.New().String(),
		JudgeName:        req.AuthorizingOfficer,
		IssuingCourt:     "EMERGENCY_AUTHORIZATION",
		StartDate:        time.Now().UTC(),
		EndDate:          time.Now().UTC().Add(72 * time.Hour),
	}

	result, err := s.repo.CreateEmergencyTarget(target)
	if err != nil {
		return domain.EmergencyResponse{}, err
	}

	if s.producer != nil {
		ctx := context.Background()
		if pubErr := s.producer.Publish(ctx, "emergency.authorized", map[string]interface{}{
			"target_id":           result.ID,
			"authorizing_officer": req.AuthorizingOfficer,
			"reason":              req.Reason,
		}); pubErr != nil {
			return domain.EmergencyResponse{}, fmt.Errorf("emergency target created but kafka publish failed: %w", pubErr)
		}
	}

	return domain.EmergencyResponse{
		Target:   result,
		Approved: true,
		AuthRef:  result.AuthorizationRef,
	}, nil
}
