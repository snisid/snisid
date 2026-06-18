package service

import (
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/cybre-svc/internal/domain"
)

type CybreService struct {
	repo domain.CybreRepository
	log  *zap.Logger
}

func NewCybreService(repo domain.CybreRepository, log *zap.Logger) *CybreService {
	return &CybreService{repo: repo, log: log}
}

func (s *CybreService) DeclareIncident(req *domain.DeclareIncidentRequest) (*domain.CyberIncident, error) {
	id, err := time.Parse(time.RFC3339, req.IncidentDate)
	if err != nil {
		return nil, err
	}

	incident := &domain.CyberIncident{
		CrimeType:    domain.CyberCrimeType(req.CrimeType),
		Severity:     domain.CyberMedium,
		IncidentDate: id,
		ReportedDate: time.Now(),
		SuspectPhone: []string{},
		SuspectEmail: []string{},
		CreatedBy:    uuid.New(),
	}

	if req.Severity != "" {
		incident.Severity = domain.CyberSeverity(req.Severity)
	}
	if req.VictimCount != nil {
		incident.VictimCount = req.VictimCount
	}
	if req.TotalFinancialLossUSD != nil {
		incident.TotalFinancialLossUSD = req.TotalFinancialLossUSD
	}
	if req.AttackVector != "" {
		incident.AttackVector = &req.AttackVector
	}
	if req.TargetedPlatform != "" {
		incident.TargetedPlatform = &req.TargetedPlatform
	}
	if req.CaseReference != "" {
		incident.CaseReference = &req.CaseReference
	}

	return s.repo.CreateIncident(incident)
}

func (s *CybreService) GetIncident(id uuid.UUID) (*domain.CyberIncident, error) {
	return s.repo.FindByID(id)
}

func (s *CybreService) ListRecentIntrusions() ([]domain.IntrusionAttempt, error) {
	return s.repo.FindRecentIntrusions()
}

func (s *CybreService) AddThreatIntel(req *domain.AddThreatIntelRequest) (*domain.ThreatIndicator, error) {
	ti := &domain.ThreatIndicator{
		IndicatorType:  req.IndicatorType,
		IndicatorValue: req.IndicatorValue,
	}

	if req.ThreatCategory != "" {
		cat := domain.CyberCrimeType(req.ThreatCategory)
		ti.ThreatCategory = &cat
	}
	if req.ConfidenceScore != nil {
		ti.ConfidenceScore = req.ConfidenceScore
	}
	if req.Source != "" {
		ti.Source = &req.Source
	}
	isActive := true
	ti.IsActive = &isActive
	now := time.Now()
	ti.FirstSeen = &now
	ti.LastSeen = &now

	return s.repo.CreateThreatIndicator(ti)
}

func (s *CybreService) CheckIndicator(indicatorType string, value string) (*domain.ThreatIndicator, error) {
	return s.repo.FindActiveIndicator(indicatorType, value)
}

func (s *CybreService) GetStatsByType() ([]domain.CyberStats, error) {
	return s.repo.GetStatsByType()
}
