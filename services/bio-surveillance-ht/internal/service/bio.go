package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/bio-surveillance-ht/internal/domain"
	"github.com/snisid/bio-surveillance-ht/internal/kafka"
	"github.com/snisid/bio-surveillance-ht/internal/repository"
)

type BioSurveillanceService struct {
	repo     repository.Repository
	producer *kafka.Producer
}

func NewBioSurveillanceService(repo repository.Repository, producer *kafka.Producer) *BioSurveillanceService {
	return &BioSurveillanceService{repo: repo, producer: producer}
}

func (s *BioSurveillanceService) CreateAlert(ctx context.Context, req domain.CreateDiseaseAlertRequest) (*domain.DiseaseAlert, error) {
	firstCase := time.Now().UTC()
	if req.FirstCaseDetected != "" {
		if parsed, err := time.Parse(time.RFC3339, req.FirstCaseDetected); err == nil {
			firstCase = parsed
		}
	}

	alert := &domain.DiseaseAlert{
		ID:                uuid.New(),
		DiseaseName:       req.DiseaseName,
		PathogenType:      domain.PathogenType(req.PathogenType),
		Icd10Code:         req.Icd10Code,
		AlertLevel:        domain.AlertLevel(req.AlertLevel),
		FirstCaseDetected: firstCase,
		TransmissionMode:  domain.TransmissionMode(req.TransmissionMode),
		IncubationDays:    req.IncubationDays,
		FatalityRate:      req.FatalityRate,
		CasesConfirmed:    req.CasesConfirmed,
		CasesSuspected:    req.CasesSuspected,
		CasesDeaths:       req.CasesDeaths,
		AffectedRegions:   req.AffectedRegions,
		CreatedAt:         time.Now().UTC(),
	}

	if req.SymptomsHallmark != "" {
		alert.SymptomsHallmark = &req.SymptomsHallmark
	}
	if req.SourceLab != "" {
		alert.SourceLab = &req.SourceLab
	}
	if req.WhoAlertRef != "" {
		alert.WhoAlertRef = &req.WhoAlertRef
	}
	if req.ContainmentMeasures != "" {
		alert.ContainmentMeasures = &req.ContainmentMeasures
	}

	if err := s.repo.CreateAlert(ctx, alert); err != nil {
		return nil, err
	}

	s.publishEvent(ctx, "biosurv.alert.created", alert)
	return alert, nil
}

func (s *BioSurveillanceService) GetActiveAlerts(ctx context.Context) ([]domain.DiseaseAlert, error) {
	return s.repo.GetActiveAlerts(ctx)
}

func (s *BioSurveillanceService) GetAlertsByRegion(ctx context.Context, region string) ([]domain.DiseaseAlert, error) {
	return s.repo.GetAlertsByRegion(ctx, region)
}

func (s *BioSurveillanceService) CreateCampaign(ctx context.Context, req domain.CreateVaccinationCampaignRequest) (*domain.VaccinationCampaign, error) {
	startDate, _ := time.Parse(time.RFC3339, req.StartDate)
	var endDate *time.Time
	if req.EndDate != "" {
		parsed, err := time.Parse(time.RFC3339, req.EndDate)
		if err == nil {
			endDate = &parsed
		}
	}

	campaign := &domain.VaccinationCampaign{
		ID:                uuid.New(),
		CampaignName:      req.CampaignName,
		TargetDisease:     req.TargetDisease,
		VaccineType:       req.VaccineType,
		TargetPopulation:  req.TargetPopulation,
		DosesAdministered: req.DosesAdministered,
		CoveragePct:       req.CoveragePct,
		RegionsActive:     req.RegionsActive,
		StartDate:         startDate,
		EndDate:           endDate,
		CoordinatorAgency: req.CoordinatorAgency,
		CreatedAt:         time.Now().UTC(),
	}

	if err := s.repo.CreateCampaign(ctx, campaign); err != nil {
		return nil, err
	}

	s.publishEvent(ctx, "biosurv.campaign.created", campaign)
	return campaign, nil
}

func (s *BioSurveillanceService) GetCampaignCoverage(ctx context.Context, id string) (*domain.VaccinationCampaign, error) {
	cid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid campaign id: %w", err)
	}
	return s.repo.GetCampaignCoverage(ctx, cid)
}

func (s *BioSurveillanceService) UpdateFacilityStock(ctx context.Context, id string, req domain.UpdateFacilityStockRequest) (*domain.HealthFacility, error) {
	fid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid facility id: %w", err)
	}
	return s.repo.UpdateFacilityStock(ctx, fid, req)
}

func (s *BioSurveillanceService) GetDashboardNational(ctx context.Context) (*domain.DashboardNational, error) {
	return s.repo.GetDashboardNational(ctx)
}

func (s *BioSurveillanceService) publishEvent(ctx context.Context, eventType string, data any) {
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
