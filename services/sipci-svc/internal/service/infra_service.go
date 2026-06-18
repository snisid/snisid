package service

import (
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/sipci-svc/internal/domain"
)

type InfraService struct {
	repo domain.AssetRepository
	log  *zap.Logger
}

func NewInfraService(repo domain.AssetRepository, log *zap.Logger) *InfraService {
	return &InfraService{repo: repo, log: log}
}

func (s *InfraService) RegisterAsset(req *domain.RegisterAssetRequest) (*domain.Asset, error) {
	asset := &domain.Asset{
		AssetName:          req.AssetName,
		AssetCategory:      domain.AssetCategory(req.AssetCategory),
		DeptCode:           req.DeptCode,
		Lat:                req.Lat,
		Lng:                req.Lng,
		CurrentThreatLevel: domain.ThreatNormal,
		CreatedBy:          uuid.New(),
	}

	if req.Commune != "" {
		asset.Commune = &req.Commune
	}
	if req.CriticalityScore != nil {
		asset.CriticalityScore = req.CriticalityScore
	}
	if req.PopulationServed != nil {
		asset.PopulationServed = req.PopulationServed
	}
	if req.OwnerEntity != "" {
		asset.OwnerEntity = &req.OwnerEntity
	}
	if req.OperatingOrg != "" {
		asset.OperatingOrg = &req.OperatingOrg
	}
	if req.SiteManagerPhone != "" {
		asset.SiteManagerPhone = &req.SiteManagerPhone
	}

	return s.repo.Create(asset)
}

func (s *InfraService) GetAsset(id uuid.UUID) (*domain.Asset, error) {
	return s.repo.FindByID(id)
}

func (s *InfraService) ListAssets() ([]domain.Asset, error) {
	return s.repo.FindAll()
}

func (s *InfraService) ListCritical() ([]domain.Asset, error) {
	return s.repo.FindCritical()
}

func (s *InfraService) ListUnderThreat() ([]domain.Asset, error) {
	return s.repo.FindUnderThreat()
}

func (s *InfraService) ReportIncident(req *domain.ReportIncidentRequest) (*domain.AssetIncident, error) {
	assetID, err := uuid.Parse(req.AssetID)
	if err != nil {
		return nil, err
	}

	id, _ := time.Parse(time.RFC3339, req.IncidentDate)

	incident := &domain.AssetIncident{
		AssetID:      assetID,
		IncidentType: req.IncidentType,
		IncidentDate: id,
		Description:  req.Description,
		CreatedBy:    uuid.New(),
	}

	if req.PerpetratorType != "" {
		incident.PerpetratorType = &req.PerpetratorType
	}
	if req.GangID != "" {
		gid, _ := uuid.Parse(req.GangID)
		incident.GangID = &gid
	}
	if req.ImpactSeverity != nil {
		incident.ImpactSeverity = req.ImpactSeverity
	}
	if req.EconomicLossUSD != nil {
		incident.EconomicLossUSD = req.EconomicLossUSD
	}

	return s.repo.CreateIncident(incident)
}

func (s *InfraService) ListRecentIncidents() ([]domain.AssetIncident, error) {
	return s.repo.FindRecentIncidents()
}

func (s *InfraService) AssessRisk(id uuid.UUID) (*domain.RiskAssessment, error) {
	asset, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	assessment := &domain.RiskAssessment{
		AssetID: id.String(),
	}

	if asset.CriticalityScore != nil {
		assessment.BaseScore = float64(*asset.CriticalityScore) * 10
	}

	if asset.IsInGangZone != nil && *asset.IsInGangZone {
		assessment.ThreatFactors = append(assessment.ThreatFactors, "IN_GANG_TERRITORY")
		assessment.FinalScore += 40
	}

	if asset.UnderExtortion != nil && *asset.UnderExtortion {
		assessment.ThreatFactors = append(assessment.ThreatFactors, "ACTIVE_EXTORTION")
		assessment.FinalScore += 25
	}

	if asset.IncidentCount12m != nil && *asset.IncidentCount12m > 0 {
		assessment.ThreatFactors = append(assessment.ThreatFactors, "RECENT_INCIDENTS")
		assessment.FinalScore += float64(*asset.IncidentCount12m) * 5
	}

	assessment.FinalScore += assessment.BaseScore

	assessment.ThreatLevel = domain.ThreatNormal
	if assessment.FinalScore >= 80 {
		assessment.ThreatLevel = domain.ThreatCritical
	} else if assessment.FinalScore >= 60 {
		assessment.ThreatLevel = domain.ThreatSevere
	} else if assessment.FinalScore >= 40 {
		assessment.ThreatLevel = domain.ThreatHigh
	} else if assessment.FinalScore >= 20 {
		assessment.ThreatLevel = domain.ThreatElevated
	}

	_ = s.repo.UpdateThreatLevel(id, assessment.ThreatLevel)

	return assessment, nil
}
