package service

import (
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/terr-svc/internal/domain"
)

type TerritoryService struct {
	repo domain.Repository
	log  *zap.Logger
}

func NewTerritoryService(repo domain.Repository, log *zap.Logger) *TerritoryService {
	return &TerritoryService{repo: repo, log: log}
}

func (s *TerritoryService) CheckPointSafety(lat, lng float64) (*domain.SafetyCheckResult, error) {
	result := &domain.SafetyCheckResult{
		Lat:         lat,
		Lng:         lng,
		IsSafe:      true,
		NearbyZones: []domain.ZoneInfo{},
	}

	zones, err := s.repo.FindZonesContainingPoint(lat, lng)
	if err != nil {
		return nil, err
	}

	if len(zones) > 0 {
		result.IsSafe = false
		z := zones[0]
		result.ContainingZone = &domain.ZoneInfo{
			ZoneID:       z.ZoneID,
			ZoneName:     z.ZoneName,
			GangID:       z.GangID,
			ControlLevel: z.ControlLevel,
			DeptCode:     z.DeptCode,
			AreaKm2:      z.AreaKm2,
		}
		for _, z2 := range zones[1:] {
			result.NearbyZones = append(result.NearbyZones, domain.ZoneInfo{
				ZoneID:       z2.ZoneID,
				ZoneName:     z2.ZoneName,
				GangID:       z2.GangID,
				ControlLevel: z2.ControlLevel,
				DeptCode:     z2.DeptCode,
				AreaKm2:      z2.AreaKm2,
			})
		}
	}

	checkpoints, err := s.repo.FindNearbyCheckpoints(lat, lng, 500)
	if err != nil {
		return nil, err
	}

	if len(checkpoints) > 0 {
		result.IsSafe = false
		result.IsInCheckpoint = true
		result.NearbyCheckpoints = checkpoints
	}

	return result, nil
}

func (s *TerritoryService) GetRouteSafety(waypoints []domain.Point) (*domain.RouteSafetyResult, error) {
	result := &domain.RouteSafetyResult{
		TotalPoints: len(waypoints),
		Waypoints:   make([]domain.SafetyCheckResult, 0, len(waypoints)),
	}

	for _, wp := range waypoints {
		check, err := s.CheckPointSafety(wp.Lat, wp.Lng)
		if err != nil {
			return nil, err
		}
		result.Waypoints = append(result.Waypoints, *check)
		if check.IsSafe {
			result.SafePoints++
		} else {
			result.UnsafePoints++
		}
	}

	result.IsSafe = result.UnsafePoints == 0
	return result, nil
}

func (s *TerritoryService) ListZones() ([]domain.TerritoryZone, error) {
	return s.repo.FindAllZones()
}

func (s *TerritoryService) ListZonesByDept(deptCode string) ([]domain.TerritoryZone, error) {
	return s.repo.FindZonesByDept(deptCode)
}

func (s *TerritoryService) ListZonesByGang(gangID uuid.UUID) ([]domain.TerritoryZone, error) {
	return s.repo.FindZonesByGang(gangID)
}

func (s *TerritoryService) CreateZone(req *domain.SeizureRequest) (*domain.TerritoryZone, error) {
	zone := &domain.TerritoryZone{
		GangID:             req.GangID,
		ZoneName:           req.ZoneName,
		DeptCode:           req.DeptCode,
		Commune:            req.Commune,
		SectionCommunale:   req.SectionCommunale,
		Geom:               req.Geom,
		ControlLevel:       req.ControlLevel,
		IntelligenceSource: req.IntelligenceSource,
		ConfidenceLevel:    req.ConfidenceLevel,
		AnalystNotes:       req.AnalystNotes,
		CreatedBy:          req.CreatedBy,
	}

	return s.repo.CreateZone(zone)
}

func (s *TerritoryService) GetZoneHistory(zoneID uuid.UUID) ([]domain.ZoneHistory, error) {
	return s.repo.GetZoneHistory(zoneID)
}

func (s *TerritoryService) ReportCheckpoint(cp *domain.Checkpoint) (*domain.Checkpoint, error) {
	return s.repo.CreateCheckpoint(cp)
}
