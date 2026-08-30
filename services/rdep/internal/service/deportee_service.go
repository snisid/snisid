package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/rdep/internal/domain"
)

var ErrDeporteeNotFound = errors.New("déporté non trouvé")
var ErrPersonNotFound = errors.New("personne SNISID introuvable")

type DeporteeService struct {
	mu        sync.RWMutex
	deportees map[uuid.UUID]*domain.Deportee
	byPerson  map[uuid.UUID]uuid.UUID
	events    map[uuid.UUID][]domain.MonitoringEvent
	records   map[uuid.UUID][]domain.ForeignRecord
	seq       int
}

func NewDeporteeService() *DeporteeService {
	return &DeporteeService{
		deportees: make(map[uuid.UUID]*domain.Deportee),
		byPerson:  make(map[uuid.UUID]uuid.UUID),
		events:    make(map[uuid.UUID][]domain.MonitoringEvent),
		records:   make(map[uuid.UUID][]domain.ForeignRecord),
	}
}

func (s *DeporteeService) generateRDEPID() string {
	s.seq++
	return fmt.Sprintf("RDEP-HT-%d-%06d", time.Now().Year(), s.seq)
}

func (s *DeporteeService) Intake(ctx context.Context, req domain.DeporteeIntakeRequest) (*domain.Deportee, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.byPerson[req.SNISIDPersonID]; exists {
		return nil, fmt.Errorf("déporté existe déjà pour la personne: %s", req.SNISIDPersonID)
	}

	riskLevel := domain.RiskNone
	if req.FBINumber != nil && *req.FBINumber != "" {
		riskLevel = domain.RiskMedium
	}
	if req.GangName != nil && *req.GangName != "" {
		riskLevel = domain.RiskVeryHigh
	}

	now := time.Now()
	deportee := &domain.Deportee{
		DeporteeID:         uuid.New(),
		NationalRDEPID:     s.generateRDEPID(),
		SNISIDPersonID:     req.SNISIDPersonID,
		AFISSubjectID:      req.AFISSubjectID,
		DeportationCountry: req.DeportationCountry,
		DeportationDate:    req.DeportationDate,
		ArrivalPort:        req.ArrivalPort,
		ArrivalDeptCode:    req.ArrivalDeptCode,
		DeportingAgency:    req.DeportingAgency,
		DeportationReason:  req.DeportationReason,
		FlightNumber:       req.FlightNumber,
		ForeignName:        req.ForeignName,
		ForeignAliases:     req.ForeignAliases,
		ForeignIDNumber:    req.ForeignIDNumber,
		ForeignCountryID:   req.ForeignCountryID,
		CriminalRiskLevel:  riskLevel,
		MonitoringStatus:   domain.MonitoringActive,
		MonitoringRequired: req.GangName != nil && *req.GangName != "",
		GangAffiliated:     req.GangName != nil && *req.GangName != "",
		GangName:           req.GangName,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	s.deportees[deportee.DeporteeID] = deportee
	s.byPerson[req.SNISIDPersonID] = deportee.DeporteeID

	return deportee, nil
}

func (s *DeporteeService) GetByID(ctx context.Context, deporteeID uuid.UUID) (*domain.Deportee, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	deportee, ok := s.deportees[deporteeID]
	if !ok {
		return nil, ErrDeporteeNotFound
	}
	return deportee, nil
}

func (s *DeporteeService) GetByPersonID(ctx context.Context, personID uuid.UUID) (*domain.Deportee, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	deporteeID, ok := s.byPerson[personID]
	if !ok {
		return nil, ErrDeporteeNotFound
	}
	deportee, ok := s.deportees[deporteeID]
	if !ok {
		return nil, ErrDeporteeNotFound
	}
	return deportee, nil
}

func (s *DeporteeService) ListHighRisk(ctx context.Context) ([]domain.Deportee, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []domain.Deportee
	for _, d := range s.deportees {
		if d.CriminalRiskLevel == domain.RiskHigh || d.CriminalRiskLevel == domain.RiskVeryHigh {
			results = append(results, *d)
		}
	}
	return results, nil
}

func (s *DeporteeService) ListGangAffiliated(ctx context.Context) ([]domain.Deportee, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []domain.Deportee
	for _, d := range s.deportees {
		if d.GangAffiliated {
			results = append(results, *d)
		}
	}
	return results, nil
}

func (s *DeporteeService) StatsByCountry(ctx context.Context) ([]domain.StatsByCountry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := make(map[domain.DeportationCountry]*domain.StatsByCountry)
	for _, d := range s.deportees {
		stat, ok := stats[d.DeportationCountry]
		if !ok {
			stat = &domain.StatsByCountry{Country: d.DeportationCountry}
			stats[d.DeportationCountry] = stat
		}
		stat.TotalCount++
		if d.CriminalRiskLevel == domain.RiskHigh || d.CriminalRiskLevel == domain.RiskVeryHigh {
			stat.HighRiskCount++
		}
		if d.GangAffiliated {
			stat.GangCount++
		}
	}

	result := make([]domain.StatsByCountry, 0, len(stats))
	for _, s := range stats {
		result = append(result, *s)
	}
	return result, nil
}

func (s *DeporteeService) Screen(ctx context.Context, deporteeID uuid.UUID) (*domain.ScreeningResult, error) {
	s.mu.RLock()
	deportee, ok := s.deportees[deporteeID]
	s.mu.RUnlock()

	if !ok {
		return nil, ErrDeporteeNotFound
	}

	result := &domain.ScreeningResult{
		PersonID:  deportee.SNISIDPersonID,
		RiskLevel: deportee.CriminalRiskLevel,
	}

	if deportee.HasForeignRecord {
		result.HasForeignRecord = true
	}

	if deportee.GangAffiliated {
		result.GangAffiliated = true
	}

	return result, nil
}

func (s *DeporteeService) AddMonitoringEvent(ctx context.Context, deporteeID uuid.UUID, event domain.MonitoringEvent) (*domain.MonitoringEvent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.deportees[deporteeID]; !ok {
		return nil, ErrDeporteeNotFound
	}

	event.EventID = uuid.New()
	event.DeporteeID = deporteeID
	event.CreatedAt = time.Now()
	s.events[deporteeID] = append(s.events[deporteeID], event)

	return &event, nil
}

func (s *DeporteeService) GetMonitoringEvents(ctx context.Context, deporteeID uuid.UUID) ([]domain.MonitoringEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events, ok := s.events[deporteeID]
	if !ok {
		return []domain.MonitoringEvent{}, nil
	}
	return events, nil
}
