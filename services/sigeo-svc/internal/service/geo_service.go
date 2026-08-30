package service

import (
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/sigeo-svc/internal/domain"
)

type GeoIntelService struct {
	repo domain.GeoRepository
	log  *zap.Logger
}

func NewGeoIntelService(repo domain.GeoRepository, log *zap.Logger) *GeoIntelService {
	return &GeoIntelService{repo: repo, log: log}
}

func (s *GeoIntelService) IngestIncident(req *domain.IngestIncidentRequest) (*domain.Incident, error) {
	ed, err := time.Parse(time.RFC3339, req.EventDate)
	if err != nil {
		return nil, err
	}

	srid, _ := uuid.Parse(req.SourceRecordID)

	incident := &domain.Incident{
		SourceModule:   req.SourceModule,
		SourceRecordID: srid,
		EventType:      req.EventType,
		EventDate:      ed,
		Lat:            req.Lat,
		Lng:            req.Lng,
		CreatedAt:      time.Now(),
	}

	if req.DeptCode != "" {
		incident.DeptCode = &req.DeptCode
	}
	if req.Commune != "" {
		incident.Commune = &req.Commune
	}
	if req.Severity != nil {
		incident.Severity = req.Severity
	}
	if req.GangID != "" {
		gid, _ := uuid.Parse(req.GangID)
		incident.GangID = &gid
	}
	if req.Description != "" {
		incident.Description = &req.Description
	}

	return s.repo.CreateIncident(incident)
}

func (s *GeoIntelService) ListIncidents(deptCode string, since time.Time) ([]domain.Incident, error) {
	return s.repo.FindIncidents(deptCode, since)
}

func (s *GeoIntelService) ListCheckpoints() ([]domain.Checkpoint, error) {
	return s.repo.FindCheckpoints()
}

func (s *GeoIntelService) GetZoneReport(deptCode string, period time.Duration) (*domain.ZoneSecurityReport, error) {
	since := time.Now().Add(-period)
	count, err := s.repo.CountIncidentsByZone(deptCode, since)
	if err != nil {
		return nil, err
	}

	score := float64(count) * 2.5
	riskLevel := "LOW"
	if score >= 80 {
		riskLevel = "CRITICAL"
	} else if score >= 50 {
		riskLevel = "HIGH"
	} else if score >= 20 {
		riskLevel = "MEDIUM"
	}

	return &domain.ZoneSecurityReport{
		DeptCode:         deptCode,
		Period:           period,
		GeneratedAt:      time.Now(),
		IncidentCount:    count,
		OverallRiskScore: score,
		RiskLevel:        riskLevel,
	}, nil
}
