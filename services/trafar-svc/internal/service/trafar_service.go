package service

import (
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/trafar-svc/internal/domain"
	"github.com/snisid/platform/services/trafar-svc/internal/repository/postgres"
)

type TrafarService struct {
	repo *postgres.RouteRepository
	log  *zap.Logger
}

func NewTrafarService(repo *postgres.RouteRepository, log *zap.Logger) *TrafarService {
	return &TrafarService{repo: repo, log: log}
}

func (s *TrafarService) ListRoutes() ([]domain.TrafarRoute, error) {
	routes, err := s.repo.FindAll()
	if err != nil {
		s.log.Error("failed to list routes", zap.Error(err))
		return nil, err
	}
	return routes, nil
}

func (s *TrafarService) GetRoute(id uuid.UUID) (*domain.TrafarRoute, error) {
	route, err := s.repo.FindByID(id)
	if err != nil {
		s.log.Error("failed to get route", zap.Error(err), zap.String("route_id", id.String()))
		return nil, err
	}
	return route, nil
}

func (s *TrafarService) CreateRoute(route *domain.TrafarRoute) error {
	if route.RouteID == uuid.Nil {
		route.RouteID = uuid.New()
	}
	if err := s.repo.CreateRoute(route); err != nil {
		s.log.Error("failed to create route", zap.Error(err))
		return err
	}
	s.log.Info("route created", zap.String("route_id", route.RouteID.String()))
	return nil
}

func (s *TrafarService) RecordShipment(shipment *domain.TrafarShipment) error {
	if shipment.ShipmentID == uuid.Nil {
		shipment.ShipmentID = uuid.New()
	}
	if err := s.repo.CreateShipment(shipment); err != nil {
		s.log.Error("failed to record shipment", zap.Error(err))
		return err
	}
	s.log.Info("shipment recorded", zap.String("shipment_id", shipment.ShipmentID.String()))
	return nil
}

func (s *TrafarService) GetMapGeoJSON() (*domain.GeoJSONFeatureCollection, error) {
	fc, err := s.repo.GetRoutesGeoJSON()
	if err != nil {
		s.log.Error("failed to get map geojson", zap.Error(err))
		return nil, err
	}
	return fc, nil
}

func (s *TrafarService) GetStatsByOrigin() ([]map[string]interface{}, error) {
	stats, err := s.repo.GetStatsByOrigin()
	if err != nil {
		s.log.Error("failed to get stats by origin", zap.Error(err))
		return nil, err
	}
	return stats, nil
}

func (s *TrafarService) ListSuppliers() ([]domain.TrafarSupplier, error) {
	suppliers, err := s.repo.GetSuppliers()
	if err != nil {
		s.log.Error("failed to list suppliers", zap.Error(err))
		return nil, err
	}
	return suppliers, nil
}
