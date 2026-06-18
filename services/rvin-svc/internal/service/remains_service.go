package service

import (
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/rvin-svc/internal/domain"
)

type RemainsService struct {
	repo domain.RemainsRepository
	log  *zap.Logger
}

func NewRemainsService(repo domain.RemainsRepository, log *zap.Logger) *RemainsService {
	return &RemainsService{repo: repo, log: log}
}

func (s *RemainsService) RegisterRemains(req *domain.RegisterRemainsRequest) (*domain.UnidentifiedRemains, error) {
	dd, err := time.Parse(time.RFC3339, req.DiscoveryDate)
	if err != nil {
		return nil, err
	}

	remains := &domain.UnidentifiedRemains{
		DiscoveryDate:     dd,
		DiscoveryLocation: req.DiscoveryLocation,
		DeptCode:          req.DeptCode,
		DiscoverySource:   domain.DiscoverySource(req.DiscoverySource),
		Status:            domain.Unidentified,
		ExaminerID:        uuid.New(),
	}

	if req.Commune != "" {
		remains.Commune = &req.Commune
	}
	if req.Lat != nil {
		remains.Lat = req.Lat
	}
	if req.Lng != nil {
		remains.Lng = req.Lng
	}
	if req.EstimatedSex != "" {
		remains.EstimatedSex = &req.EstimatedSex
	}
	if req.EstimatedAgeMin != nil {
		remains.EstimatedAgeMin = req.EstimatedAgeMin
	}
	if req.EstimatedAgeMax != nil {
		remains.EstimatedAgeMax = req.EstimatedAgeMax
	}
	if req.EstimatedHeightCm != nil {
		remains.EstimatedHeightCm = req.EstimatedHeightCm
	}
	if req.DistinguishingMarks != "" {
		remains.DistinguishingMarks = &req.DistinguishingMarks
	}
	if req.MorgueLocation != "" {
		remains.MorgueLocation = &req.MorgueLocation
	}

	return s.repo.Create(remains)
}

func (s *RemainsService) GetRemains(id uuid.UUID) (*domain.UnidentifiedRemains, error) {
	return s.repo.FindByID(id)
}

func (s *RemainsService) SubmitDNA(remainsID uuid.UUID, req *domain.SubmitDNARequest) error {
	dna := &domain.DNAResult{
		RemainsID:        remainsID,
		ReferenceDNARef:  req.ReferenceDNARef,
		MatchProbability: req.MatchProbability,
		IsMatch:          req.IsMatch,
		CreatedAt:        time.Now(),
	}
	if req.LabReference != "" {
		dna.LabReference = &req.LabReference
	}
	return s.repo.SubmitDNA(remainsID, dna)
}

func (s *RemainsService) ListUnidentified() ([]domain.UnidentifiedRemains, error) {
	return s.repo.FindUnidentified()
}

func (s *RemainsService) GetStatsBySource() ([]domain.SourceStats, error) {
	return s.repo.GetStatsBySource()
}
