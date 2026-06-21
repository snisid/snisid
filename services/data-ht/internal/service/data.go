package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/data-ht/internal/domain"
	"github.com/snisid/data-ht/internal/kafka"
	"github.com/snisid/data-ht/internal/repository"
)

type DataService struct {
	repo     repository.Repository
	producer *kafka.Producer
}

func NewDataService(repo repository.Repository, producer *kafka.Producer) *DataService {
	return &DataService{repo: repo, producer: producer}
}

func (s *DataService) ListPipelines(ctx context.Context) ([]domain.Pipeline, error) {
	return s.repo.ListPipelines(ctx)
}

func (s *DataService) RegisterModel(ctx context.Context, req domain.RegisterModelRequest) (*domain.MLModel, error) {
	model := &domain.MLModel{
		ID:          uuid.New(),
		Name:        req.Name,
		ModelType:   req.ModelType,
		Version:     req.Version,
		MlflowRunID: req.MlflowRunID,
		IsActive:    true,
		CreatedAt:   time.Now().UTC(),
	}

	if req.BiasMetric != "" {
		model.BiasMetric = &req.BiasMetric
	}
	if req.BiasScore != 0 {
		model.BiasScore = &req.BiasScore
	}
	if req.TrainingDate != "" {
		t, err := time.Parse(time.RFC3339, req.TrainingDate)
		if err == nil {
			model.TrainingDate = &t
		}
	}

	if err := s.repo.CreateModel(ctx, model); err != nil {
		return nil, err
	}

	s.publishEvent(ctx, "data.model.registered", model)
	return model, nil
}

func (s *DataService) GetBiasAudit(ctx context.Context, modelID string) (*domain.BiasAuditResult, error) {
	mid, err := uuid.Parse(modelID)
	if err != nil {
		return nil, fmt.Errorf("invalid model_id: %w", err)
	}

	model, err := s.repo.GetModel(ctx, mid)
	if err != nil {
		return nil, err
	}

	audits, err := s.repo.GetGovernanceAuditsByModel(ctx, mid)
	if err != nil {
		return nil, err
	}

	result := &domain.BiasAuditResult{
		ModelID:    model.ID,
		ModelName:  model.Name,
		BiasMetric: model.BiasMetric,
		BiasScore:  model.BiasScore,
		AuditCount: len(audits),
	}

	if len(audits) > 0 {
		result.LastAudited = &audits[0].ConductedAt
	}

	return result, nil
}

func (s *DataService) GetNationalDashboard(ctx context.Context) (*domain.NationalDashboard, error) {
	return s.repo.GetNationalDashboard(ctx)
}

func (s *DataService) publishEvent(ctx context.Context, eventType string, data any) {
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
