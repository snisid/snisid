package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/siar/internal/domain"
)

type DealerService struct {
	mu      sync.RWMutex
	dealers map[uuid.UUID]*domain.Dealer
	seq     int
}

func NewDealerService() *DealerService {
	return &DealerService{
		dealers: make(map[uuid.UUID]*domain.Dealer),
	}
}

func (s *DealerService) generateDealerLicense() string {
	s.seq++
	return fmt.Sprintf("SIAR-DLR-%04d", s.seq)
}

func (s *DealerService) Create(ctx context.Context, req domain.CreateDealerRequest, createdBy uuid.UUID) (*domain.Dealer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	lt := req.LicenseType
	if lt == "" {
		lt = domain.LicenseDealer
	}

	d := &domain.Dealer{
		DealerID:          uuid.New(),
		DealerLicenseNo:   s.generateDealerLicense(),
		BusinessName:      req.BusinessName,
		BusinessRegNo:     req.BusinessRegNo,
		OwnerSnisidID:     req.OwnerSnisidID,
		OwnerName:         req.OwnerName,
		Address:           req.Address,
		DeptCode:          req.DeptCode,
		Commune:           req.Commune,
		Phone:             req.Phone,
		Email:             req.Email,
		LicenseType:       lt,
		Status:            domain.DealerActive,
		LicenseIssueDate:  req.LicenseIssueDate,
		LicenseExpiryDate: req.LicenseExpiryDate,
		PremisesInspected: req.PremisesInspected,
		Notes:             req.Notes,
		CreatedBy:         createdBy,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	s.dealers[d.DealerID] = d
	return d, nil
}

func (s *DealerService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Dealer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	d, ok := s.dealers[id]
	if !ok {
		return nil, ErrDealerNotFound
	}
	return d, nil
}

func (s *DealerService) List(ctx context.Context, deptCode string) ([]*domain.Dealer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*domain.Dealer
	for _, d := range s.dealers {
		if deptCode == "" || d.DeptCode == deptCode {
			result = append(result, d)
		}
	}
	return result, nil
}

func (s *DealerService) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.DealerStatus) (*domain.Dealer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	d, ok := s.dealers[id]
	if !ok {
		return nil, ErrDealerNotFound
	}
	d.Status = status
	d.UpdatedAt = time.Now()
	return d, nil
}
