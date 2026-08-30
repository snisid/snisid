package service

import (
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/expl-svc/internal/domain"
)

type ExplService struct {
	repo domain.IncidentRepository
	log  *zap.Logger
}

func NewExplService(repo domain.IncidentRepository, log *zap.Logger) *ExplService {
	return &ExplService{repo: repo, log: log}
}

func (s *ExplService) ReportIncident(incident *domain.ExplIncident) (*domain.ExplIncident, error) {
	count, err := s.repo.CountIncidents()
	if err != nil {
		return nil, fmt.Errorf("report incident: %w", err)
	}
	incident.NationalExplID = fmt.Sprintf("EXPL-HT-%06d", count+1)
	incident.IncidentID = uuid.New()

	if err := s.repo.CreateIncident(incident); err != nil {
		return nil, fmt.Errorf("report incident: %w", err)
	}
	return incident, nil
}

func (s *ExplService) GetIncidentsByDept(deptCode string, limit, offset int) ([]domain.ExplIncident, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.FindByDept(deptCode, limit, offset)
}

func (s *ExplService) ReportLegalStock(stock *domain.LegalStock) (*domain.LegalStock, error) {
	stock.StockID = uuid.New()
	if err := s.repo.CreateLegalStock(stock); err != nil {
		return nil, fmt.Errorf("report legal stock: %w", err)
	}
	return stock, nil
}

func (s *ExplService) GetLegalStocks(deptCode string, limit, offset int) ([]domain.LegalStock, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.GetLegalStocks(deptCode, limit, offset)
}
