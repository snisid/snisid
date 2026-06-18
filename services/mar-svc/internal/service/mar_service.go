package service

import (
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/mar-svc/internal/domain"
)

type MaritimeService struct {
	repo domain.MaritimeRepository
	log  *zap.Logger
}

func NewMaritimeService(repo domain.MaritimeRepository, log *zap.Logger) *MaritimeService {
	return &MaritimeService{repo: repo, log: log}
}

func (s *MaritimeService) RegisterVessel(v *domain.Vessel) error {
	return s.repo.CreateVessel(v)
}

func (s *MaritimeService) ProcessAISSighting(msg *domain.AISMessage) (*domain.AISSighting, error) {
	zone := s.detectZone(msg.Lat, msg.Lng)

	lastSighting, _ := s.repo.GetLastSighting(msg.MMSI)

	sighting := &domain.AISSighting{
		MMSI:              msg.MMSI,
		VesselName:        msg.VesselName,
		SightingTimestamp: time.Now(),
		Lat:               msg.Lat,
		Lng:               msg.Lng,
		SpeedKnots:        msg.SpeedKnots,
		HeadingDegrees:    msg.HeadingDegrees,
		Destination:       msg.Destination,
		SourceType:        msg.SourceType,
		ZoneCode:          zone,
	}

	if lastSighting != nil {
		sighting.VesselID = lastSighting.VesselID
	}

	suspicious := false
	if msg.SpeedKnots > 30 {
		suspicious = true
	}

	if lastSighting != nil && lastSighting.ZoneCode != "" && lastSighting.ZoneCode != zone {
		elapsed := msg.SpeedKnots
		_ = elapsed
	}

	if suspicious {
		sighting.AlertTriggered = true
		s.log.Warn("suspicious AIS activity",
			zap.String("mmsi", msg.MMSI),
			zap.Float64("speed", msg.SpeedKnots),
			zap.String("zone", zone))
	}

	if err := s.repo.CreateAISSighting(sighting); err != nil {
		return nil, err
	}

	return sighting, nil
}

func (s *MaritimeService) ReportIncident(i *domain.Incident) error {
	return s.repo.CreateIncident(i)
}

func (s *MaritimeService) AddToWatchList(w *domain.WatchVessel) error {
	return s.repo.CreateWatch(w)
}

func (s *MaritimeService) GetActiveWatches() ([]domain.WatchVessel, error) {
	return s.repo.GetActiveWatches()
}

func (s *MaritimeService) GetLivePositions(limit int) ([]domain.AISSighting, error) {
	if limit <= 0 {
		limit = 100
	}
	return s.repo.GetLiveAIS(limit)
}

func (s *MaritimeService) GetZoneActivity(zone string, limit int) ([]domain.Incident, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.GetIncidentsByZone(zone, limit)
}

func (s *MaritimeService) GetIncidentStats() (map[string]int64, error) {
	return s.repo.GetIncidentStats()
}

func (s *MaritimeService) GetVessel(id uuid.UUID) (*domain.Vessel, error) {
	return s.repo.FindVesselByID(id)
}

func (s *MaritimeService) GetRecentIncidents(limit int) ([]domain.Incident, error) {
	if limit <= 0 {
		limit = 25
	}
	return s.repo.GetRecentIncidents(limit)
}

func (s *MaritimeService) detectZone(lat, lng float64) string {
	if lat >= 19.5 && lat <= 20.5 && lng >= -74.5 && lng <= -73.5 {
		return "WINDWARD_PASS"
	}
	if lat >= 19.8 && lat <= 20.3 && lng >= -74.2 && lng <= -73.5 {
		return "TORTUE"
	}
	if lat >= 18.3 && lat <= 19.2 && lng >= -74.0 && lng <= -72.8 {
		return "GONAVE"
	}
	if lat >= 18.2 && lat <= 18.9 && lng >= -73.0 && lng <= -72.0 {
		return "PAU"
	}
	if lat >= 19.5 && lat <= 20.1 && lng >= -72.6 && lng <= -71.8 {
		return "CAP_HAITIEN"
	}
	if lat >= 17.9 && lat <= 18.5 && lng >= -74.2 && lng <= -73.5 {
		return "LES_CAYES"
	}
	return "UNKNOWN"
}
