package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/siar/internal/domain"
)

type FirearmService struct {
	mu       sync.RWMutex
	firearms map[uuid.UUID]*domain.Firearm
	bySerial map[string]*domain.Firearm
	seq      int
}

func NewFirearmService() *FirearmService {
	return &FirearmService{
		firearms: make(map[uuid.UUID]*domain.Firearm),
		bySerial: make(map[string]*domain.Firearm),
	}
}

func (s *FirearmService) generateSiarID() string {
	s.seq++
	return fmt.Sprintf("SIAR-HT-%06d", s.seq)
}

func (s *FirearmService) Create(ctx context.Context, req domain.CreateFirearmRequest, createdBy uuid.UUID) (*domain.Firearm, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if req.SerialNumber != "" {
		if existing, ok := s.bySerial[req.SerialNumber]; ok {
			return nil, fmt.Errorf("arme déjà enregistrée avec le numéro de série: %s (SIAR: %s)", req.SerialNumber, existing.NationalSiarID)
		}
	}

	f := &domain.Firearm{
		FirearmID:          uuid.New(),
		NationalSiarID:     s.generateSiarID(),
		SerialNumber:       req.SerialNumber,
		Make:               req.Make,
		Model:              req.Model,
		Caliber:            req.Caliber,
		WeaponType:         req.WeaponType,
		ManufactureYear:    req.ManufactureYear,
		ManufactureCountry: req.ManufactureCountry,
		Status:             domain.StatusRegistered,
		RegType:            req.RegType,
		OwnerSnisidID:      req.OwnerSnisidID,
		OwnerEntity:        req.OwnerEntity,
		LicenseNumber:      req.LicenseNumber,
		ImportDate:         req.ImportDate,
		ImportCountry:      req.ImportCountry,
		ImportPermitRef:    req.ImportPermitRef,
		ImporterName:       req.ImporterName,
		CustomsEntryRef:    req.CustomsEntryRef,
		CurrentDeptCode:    req.CurrentDeptCode,
		StorageLocation:    req.StorageLocation,
		Notes:              req.Notes,
		CreatedBy:          createdBy,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	s.firearms[f.FirearmID] = f
	if f.SerialNumber != "" {
		s.bySerial[f.SerialNumber] = f
	}

	return f, nil
}

func (s *FirearmService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Firearm, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	f, ok := s.firearms[id]
	if !ok {
		return nil, ErrFirearmNotFound
	}
	return f, nil
}

func (s *FirearmService) FindBySerial(ctx context.Context, serial string) (*domain.Firearm, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	f, ok := s.bySerial[serial]
	if !ok {
		return nil, ErrFirearmNotFound
	}
	return f, nil
}

func (s *FirearmService) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.FirearmStatus) (*domain.Firearm, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, ok := s.firearms[id]
	if !ok {
		return nil, ErrFirearmNotFound
	}
	f.Status = status
	f.UpdatedAt = time.Now()
	return f, nil
}

func (s *FirearmService) List(ctx context.Context) ([]*domain.Firearm, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*domain.Firearm, 0, len(s.firearms))
	for _, f := range s.firearms {
		result = append(result, f)
	}
	return result, nil
}

func (s *FirearmService) StatsByType(ctx context.Context) ([]domain.FirearmStatsByType, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	counts := make(map[domain.WeaponType]int)
	for _, f := range s.firearms {
		counts[f.WeaponType]++
	}

	stats := make([]domain.FirearmStatsByType, 0, len(counts))
	for wt, count := range counts {
		stats = append(stats, domain.FirearmStatsByType{WeaponType: wt, Count: count})
	}
	return stats, nil
}
