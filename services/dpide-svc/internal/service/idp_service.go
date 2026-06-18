package service

import (
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/dpide-svc/internal/domain"
)

type IDPService struct {
	repo domain.IDPRepository
	log  *zap.Logger
}

func NewIDPService(repo domain.IDPRepository, log *zap.Logger) *IDPService {
	return &IDPService{repo: repo, log: log}
}

func (s *IDPService) RegisterIDP(req *domain.RegisterIDPRequest) (*domain.IDP, error) {
	dd, err := time.Parse(time.RFC3339, req.DisplacementDate)
	if err != nil {
		return nil, err
	}

	idp := &domain.IDP{
		FullName:          req.FullName,
		DisplacementCause: domain.DisplacementCause(req.DisplacementCause),
		DisplacementDate:  dd,
		OriginDeptCode:    req.OriginDeptCode,
		Status:            domain.Displaced,
		MedicalNeeds:      []string{},
	}

	if req.Gender != "" {
		idp.Gender = &req.Gender
	}
	if req.OriginCommune != "" {
		idp.OriginCommune = &req.OriginCommune
	}
	if req.CurrentLocation != "" {
		idp.CurrentLocation = &req.CurrentLocation
	}
	if req.CurrentDeptCode != "" {
		idp.CurrentDeptCode = &req.CurrentDeptCode
	}
	if req.CurrentCommune != "" {
		idp.CurrentCommune = &req.CurrentCommune
	}
	if req.CurrentLat != nil {
		idp.CurrentLat = req.CurrentLat
	}
	if req.CurrentLng != nil {
		idp.CurrentLng = req.CurrentLng
	}
	if req.ShelterType != "" {
		idp.ShelterType = &req.ShelterType
	}
	if req.HouseholdSize != nil {
		idp.HouseholdSize = req.HouseholdSize
	}
	if req.MinorsCount != nil {
		idp.MinorsCount = req.MinorsCount
	}
	if req.DOB != "" {
		dob, _ := time.Parse("2006-01-02", req.DOB)
		idp.DOB = &dob
	}

	return s.repo.Create(idp)
}

func (s *IDPService) GetIDP(id uuid.UUID) (*domain.IDP, error) {
	return s.repo.FindByID(id)
}

func (s *IDPService) ListCamps() ([]domain.Camp, error) {
	return s.repo.FindCamps()
}

func (s *IDPService) GetStats() (*domain.IDPStats, error) {
	return s.repo.GetStats()
}

func (s *IDPService) UpdateStatus(id uuid.UUID, req *domain.UpdateStatusRequest) error {
	return s.repo.UpdateStatus(id, domain.IDPStatus(req.Status))
}
