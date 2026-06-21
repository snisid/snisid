package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/field-ht/internal/domain"
	"github.com/snisid/field-ht/internal/kafka"
	"github.com/snisid/field-ht/internal/repository"
)

type FieldService struct {
	repo     repository.Repository
	producer *kafka.Producer
}

func NewFieldService(repo repository.Repository, producer *kafka.Producer) *FieldService {
	return &FieldService{repo: repo, producer: producer}
}

func (s *FieldService) CreateMission(ctx context.Context, req domain.CreateMissionRequest) (*domain.Mission, error) {
	mission := &domain.Mission{
		ID:        uuid.New(),
		Title:     req.Title,
		Status:    domain.MissionStatusPlanned,
		DeptCode:  req.DeptCode,
		CreatedAt: time.Now().UTC(),
	}

	if req.Description != "" {
		mission.Description = &req.Description
	}

	if req.AssignedUnitID != "" {
		uid, err := uuid.Parse(req.AssignedUnitID)
		if err != nil {
			return nil, fmt.Errorf("invalid assigned_unit_id: %w", err)
		}
		mission.AssignedUnitID = &uid
	}

	if err := s.repo.CreateMission(ctx, mission); err != nil {
		return nil, err
	}

	s.publishEvent(ctx, "field.mission.created", mission)
	return mission, nil
}

func (s *FieldService) GetActiveMissions(ctx context.Context) ([]domain.Mission, error) {
	return s.repo.GetActiveMissions(ctx)
}

func (s *FieldService) CreateMissionLog(ctx context.Context, missionID string, req domain.CreateMissionLogRequest) (*domain.MissionLog, error) {
	mid, err := uuid.Parse(missionID)
	if err != nil {
		return nil, fmt.Errorf("invalid mission_id: %w", err)
	}

	loggedBy, err := uuid.Parse(req.LoggedBy)
	if err != nil {
		return nil, fmt.Errorf("invalid logged_by: %w", err)
	}

	logEntry := &domain.MissionLog{
		ID:        uuid.New(),
		MissionID: mid,
		LoggedBy:  loggedBy,
		Action:    req.Action,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		CreatedAt: time.Now().UTC(),
	}

	if req.Notes != "" {
		logEntry.Notes = &req.Notes
	}

	if err := s.repo.CreateMissionLog(ctx, logEntry); err != nil {
		return nil, err
	}

	s.publishEvent(ctx, "field.mission.logged", logEntry)
	return logEntry, nil
}

func (s *FieldService) GetCoverageStats(ctx context.Context) (*domain.CoverageStats, error) {
	return s.repo.GetCoverageStats(ctx)
}

func (s *FieldService) publishEvent(ctx context.Context, eventType string, data any) {
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
