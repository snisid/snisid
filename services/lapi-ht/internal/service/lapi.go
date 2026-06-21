package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/snisid/lapi-ht/internal/domain"
	"github.com/snisid/lapi-ht/internal/kafka"
	"github.com/snisid/lapi-ht/internal/repository"
)

type LapiService struct {
	repo     repository.Repository
	producer *kafka.Producer
}

func NewLapiService(repo repository.Repository, producer *kafka.Producer) *LapiService {
	return &LapiService{repo: repo, producer: producer}
}

func (s *LapiService) RecordRead(read *domain.PlateRead) error {
	read.ID = uuid.New()
	read.CreatedAt = time.Now().UTC()
	if err := s.repo.SavePlateRead(read); err != nil {
		return err
	}
	if read.AlertTriggered {
		alert := &domain.AlertDispatch{
			ID:           uuid.New(),
			ReadID:       read.ID,
			PlateNumber:  read.PlateNumberNormalized,
			Reason:       "ALPR alert triggered",
			DispatchedAt: time.Now().UTC(),
			IsActive:     true,
		}
		if err := s.repo.SaveAlertDispatch(alert); err != nil {
			return err
		}
		if s.producer != nil {
			s.producer.PublishAlert(alert)
		}
	}
	return nil
}

func (s *LapiService) GetRecentReads(limit int) ([]domain.PlateRead, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.repo.GetRecentReads(limit)
}

func (s *LapiService) GetReadsByPlate(plateNumber string) ([]domain.PlateRead, error) {
	return s.repo.GetReadsByPlate(plateNumber)
}

func (s *LapiService) GetActiveAlerts() ([]domain.AlertDispatch, error) {
	return s.repo.GetActiveAlerts()
}

func (s *LapiService) GetCameraStatus() ([]domain.Camera, error) {
	return s.repo.GetCameras()
}
