package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/snisid/fpr-ht/internal/domain"
	"github.com/snisid/fpr-ht/internal/kafka"
	"github.com/snisid/fpr-ht/internal/repository"
)

type FprService struct {
	repo     repository.Repository
	producer *kafka.Producer
}

func NewFprService(repo repository.Repository, producer *kafka.Producer) *FprService {
	return &FprService{repo: repo, producer: producer}
}

func (s *FprService) CreateWarrant(w *domain.Warrant) error {
	w.ID = uuid.New()
	now := time.Now().UTC()
	w.CreatedAt = now
	w.UpdatedAt = now
	if err := s.repo.SaveWarrant(w); err != nil {
		return err
	}
	if s.producer != nil {
		s.producer.PublishWarrantCreated(w)
	}
	return nil
}

func (s *FprService) CheckCitizen(citizenID string) (*domain.WarrantCheckResult, error) {
	warrants, err := s.repo.FindWarrantsByName(citizenID)
	if err != nil {
		return nil, err
	}

	result := &domain.WarrantCheckResult{
		CheckLog: domain.CheckLog{
			ID:        uuid.New(),
			CitizenID: citizenID,
			Result:    "CLEAR",
			CheckedAt: time.Now().UTC(),
		},
	}
	if len(warrants) > 0 {
		result.WarrantFound = true
		result.Warrant = &warrants[0]
		result.CheckLog.WarrantID = &warrants[0].ID
		result.CheckLog.Result = "WANTED"
	}

	if err := s.repo.SaveCheckLog(&result.CheckLog); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *FprService) CheckByName(name string) (*domain.WarrantCheckResult, error) {
	warrants, err := s.repo.FindWarrantsByName(name)
	if err != nil {
		return nil, err
	}

	result := &domain.WarrantCheckResult{
		CheckLog: domain.CheckLog{
			ID:        uuid.New(),
			CitizenID: name,
			Result:    "CLEAR",
			CheckedAt: time.Now().UTC(),
		},
	}
	if len(warrants) > 0 {
		result.WarrantFound = true
		result.Warrant = &warrants[0]
		result.CheckLog.WarrantID = &warrants[0].ID
		result.CheckLog.Result = "WANTED"
	}

	if err := s.repo.SaveCheckLog(&result.CheckLog); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *FprService) ReportSighting(warrantID uuid.UUID, sighting *domain.Sighting) error {
	sighting.ID = uuid.New()
	sighting.WarrantID = warrantID
	sighting.CreatedAt = time.Now().UTC()
	return s.repo.SaveSighting(sighting)
}

func (s *FprService) ExecuteWarrant(id uuid.UUID) error {
	now := time.Now().UTC()
	return s.repo.UpdateWarrantExecuted(id, now)
}

func (s *FprService) GetArmedDangerous() ([]domain.Warrant, error) {
	return s.repo.GetArmedDangerousWarrants()
}

func (s *FprService) GetDashboardStats() (*domain.DashboardStats, error) {
	return s.repo.GetDashboardStats()
}
