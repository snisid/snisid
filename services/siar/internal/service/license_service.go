package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/siar/internal/domain"
)

type LicenseService struct {
	mu       sync.RWMutex
	licenses map[uuid.UUID]*domain.License
	byNumber map[string]*domain.License
	seq      int
}

func NewLicenseService() *LicenseService {
	return &LicenseService{
		licenses: make(map[uuid.UUID]*domain.License),
		byNumber: make(map[string]*domain.License),
	}
}

func (s *LicenseService) generateLicenseNumber() string {
	s.seq++
	return fmt.Sprintf("SIAR-LIC-%04d", s.seq)
}

func (s *LicenseService) Create(ctx context.Context, req domain.CreateLicenseRequest, createdBy uuid.UUID) (*domain.License, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	lic := &domain.License{
		LicenseID:          uuid.New(),
		LicenseNumber:      s.generateLicenseNumber(),
		HolderSnisidID:     req.HolderSnisidID,
		HolderName:         req.HolderName,
		LicenseType:        req.LicenseType,
		FirearmsAuthorized: req.FirearmsAuthorized,
		IssueDate:          req.IssueDate,
		ExpiryDate:         req.ExpiryDate,
		IssuingAuthority:   req.IssuingAuthority,
		IsActive:           true,
		CreatedBy:          createdBy,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	if lic.FirearmsAuthorized <= 0 {
		lic.FirearmsAuthorized = 1
	}

	s.licenses[lic.LicenseID] = lic
	s.byNumber[lic.LicenseNumber] = lic
	return lic, nil
}

func (s *LicenseService) GetByPerson(ctx context.Context, personID uuid.UUID) ([]*domain.License, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*domain.License
	for _, lic := range s.licenses {
		if lic.HolderSnisidID == personID {
			result = append(result, lic)
		}
	}
	return result, nil
}

func (s *LicenseService) GetByID(ctx context.Context, id uuid.UUID) (*domain.License, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	lic, ok := s.licenses[id]
	if !ok {
		return nil, ErrLicenseNotFound
	}
	return lic, nil
}

func (s *LicenseService) Revoke(ctx context.Context, id uuid.UUID, reason string) (*domain.License, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	lic, ok := s.licenses[id]
	if !ok {
		return nil, ErrLicenseNotFound
	}
	now := time.Now()
	lic.IsActive = false
	lic.RevocationReason = reason
	lic.RevokedAt = &now
	lic.UpdatedAt = now
	return lic, nil
}
