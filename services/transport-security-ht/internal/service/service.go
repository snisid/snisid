package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/transport-security-ht/internal/domain"
	"github.com/snisid/transport-security-ht/internal/kafka"
	"github.com/snisid/transport-security-ht/internal/repository"
)

type TransportService interface {
	LogScreening(ctx context.Context, s *domain.PassengerScreening) error
	GetRecentScreenings(ctx context.Context) ([]domain.PassengerScreening, error)
	AddNoFlyEntry(ctx context.Context, p *domain.NoFlyPassenger) error
	CheckNoFly(ctx context.Context, identityRef string) (*domain.NoFlyPassenger, error)
	GetZonesByAirport(ctx context.Context, airportCode string) ([]domain.AirportSecurityZone, error)
	ReportBreach(ctx context.Context, zoneID uuid.UUID) error
}

type transportService struct {
	repo repository.TransportRepository
	kaf  kafka.Producer
}

func NewTransportService(repo repository.TransportRepository, kaf kafka.Producer) TransportService {
	return &transportService{repo: repo, kaf: kaf}
}

func (s *transportService) LogScreening(ctx context.Context, screening *domain.PassengerScreening) error {
	screening.ScreeningID = uuid.New()
	screening.ScreenedAt = time.Now().UTC()
	if err := s.repo.CreateScreening(ctx, screening); err != nil {
		return err
	}
	_ = s.kaf.Publish(ctx, "screening.created", screening)
	return nil
}

func (s *transportService) GetRecentScreenings(ctx context.Context) ([]domain.PassengerScreening, error) {
	return s.repo.GetRecentScreenings(ctx, 50)
}

func (s *transportService) AddNoFlyEntry(ctx context.Context, p *domain.NoFlyPassenger) error {
	if err := s.repo.AddNoFly(ctx, p); err != nil {
		return err
	}
	_ = s.kaf.Publish(ctx, "no-fly.added", p)
	return nil
}

func (s *transportService) CheckNoFly(ctx context.Context, identityRef string) (*domain.NoFlyPassenger, error) {
	return s.repo.CheckNoFly(ctx, identityRef)
}

func (s *transportService) GetZonesByAirport(ctx context.Context, airportCode string) ([]domain.AirportSecurityZone, error) {
	return s.repo.GetZonesByAirport(ctx, airportCode)
}

func (s *transportService) ReportBreach(ctx context.Context, zoneID uuid.UUID) error {
	if err := s.repo.ReportZoneBreach(ctx, zoneID); err != nil {
		return err
	}
	_ = s.kaf.Publish(ctx, "zone.breach", map[string]interface{}{"zone_id": zoneID.String()})
	return nil
}
