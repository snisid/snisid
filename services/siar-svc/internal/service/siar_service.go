package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/siar-svc/internal/domain"
)

type Repository interface {
	CreateFirearm(f *domain.Firearm) (*domain.Firearm, error)
	FindBySerial(serialNumber string) (*domain.Firearm, error)
	FindByID(id uuid.UUID) (*domain.Firearm, error)
	UpdateFirearm(f *domain.Firearm) error
	CreateSeizure(s *domain.Seizure) error
	CreateLicense(l *domain.License) error
	GetLicensesByPerson(personID uuid.UUID) ([]domain.License, error)
	CreateTransfer(t *domain.Transfer) error
	GetStatsByType() ([]domain.StatsByType, error)
}

type SIARService struct {
	repo Repository
	log  *zap.Logger
}

func NewSIARService(repo Repository, log *zap.Logger) *SIARService {
	return &SIARService{repo: repo, log: log}
}

func (s *SIARService) RegisterFirearm(req *domain.RegisterFirearmRequest) (*domain.Firearm, error) {
	firearm := &domain.Firearm{
		SerialNumber:       req.SerialNumber,
		Make:               req.Make,
		Model:              req.Model,
		Caliber:            req.Caliber,
		WeaponType:         req.WeaponType,
		ManufactureYear:    req.ManufactureYear,
		ManufactureCountry: req.ManufactureCountry,
		Status:             domain.REGISTERED,
		RegType:            req.RegType,
		OwnerSnisidID:      req.OwnerSnisidID,
		OwnerEntityName:    req.OwnerEntityName,
		LicenseNumber:      req.LicenseNumber,
		LicenseExpiry:      req.LicenseExpiry,
		ImportDate:         req.ImportDate,
		ImportCountry:      req.ImportCountry,
		ImportPermitRef:    req.ImportPermitRef,
		ImporterName:       req.ImporterName,
		CustomsEntryRef:    req.CustomsEntryRef,
		CurrentDeptCode:    req.CurrentDeptCode,
		StorageLocation:    req.StorageLocation,
		GangID:             req.GangID,
		Notes:              req.Notes,
		CreatedBy:          req.CreatedBy,
	}

	firearm.NationalSiarID = s.generateSiarID(firearm)

	return s.repo.CreateFirearm(firearm)
}

func (s *SIARService) ReportSeizure(req *domain.SeizureRequest) (*domain.Seizure, error) {
	var firearmID *uuid.UUID

	if req.SerialNumber != nil && *req.SerialNumber != "" {
		existing, err := s.repo.FindBySerial(*req.SerialNumber)
		if err == nil && existing != nil {
			firearmID = &existing.ID
			existing.Status = domain.SEIZED
			if err := s.repo.UpdateFirearm(existing); err != nil {
				s.log.Error("failed to update firearm status", zap.Error(err))
				return nil, err
			}
		} else {
			firearm, err := s.CreateSeizureFirearm(req)
			if err != nil {
				return nil, err
			}
			firearmID = &firearm.ID
		}
	} else {
		firearm, err := s.CreateSeizureFirearm(req)
		if err != nil {
			return nil, err
		}
		firearmID = &firearm.ID
	}

	seizure := domain.NewSeizure(req, firearmID)
	if err := s.repo.CreateSeizure(seizure); err != nil {
		s.log.Error("failed to create seizure record", zap.Error(err))
		return nil, err
	}

	s.log.Info("seizure reported",
		zap.String("seizure_id", seizure.ID.String()),
		zap.String("event", "siar.seizure.reported"))

	return seizure, nil
}

func (s *SIARService) CreateSeizureFirearm(req *domain.SeizureRequest) (*domain.Firearm, error) {
	firearm := &domain.Firearm{
		SerialNumber: req.SerialNumber,
		Make:         req.Make,
		Model:        req.Model,
		Caliber:      req.Caliber,
		WeaponType:   req.WeaponType,
		Status:       domain.SEIZED,
		RegType:      req.RegType,
		CreatedBy:    req.CreatedBy,
	}
	firearm.NationalSiarID = s.generateSiarID(firearm)
	return s.repo.CreateFirearm(firearm)
}

func (s *SIARService) ReportStolen(req *domain.StolenRequest) error {
	firearm, err := s.repo.FindByID(req.FirearmID)
	if err != nil {
		return fmt.Errorf("firearm not found: %w", err)
	}

	firearm.Status = domain.REPORTED_STOLEN
	return s.repo.UpdateFirearm(firearm)
}

func (s *SIARService) CheckSerial(serialNumber string) (*domain.Firearm, error) {
	return s.repo.FindBySerial(serialNumber)
}

func (s *SIARService) GetLicensesByPerson(personID uuid.UUID) ([]domain.License, error) {
	return s.repo.GetLicensesByPerson(personID)
}

func (s *SIARService) CreateLicense(req *domain.CreateLicenseRequest) error {
	license := &domain.License{
		LicenseNumber:      req.LicenseNumber,
		HolderSnisidID:     req.HolderSnisidID,
		LicenseType:        req.LicenseType,
		FirearmsAuthorized: req.FirearmsAuthorized,
		IssueDate:          req.IssueDate,
		ExpiryDate:         req.ExpiryDate,
		IssuingAuthority:   req.IssuingAuthority,
		IsActive:           true,
	}
	if license.FirearmsAuthorized == 0 {
		license.FirearmsAuthorized = 1
	}
	return s.repo.CreateLicense(license)
}

func (s *SIARService) GetStatsByType() ([]domain.StatsByType, error) {
	return s.repo.GetStatsByType()
}

func (s *SIARService) GetFirearmByID(id uuid.UUID) (*domain.Firearm, error) {
	return s.repo.FindByID(id)
}

func (s *SIARService) generateSiarID(f *domain.Firearm) string {
	now := time.Now()
	year := now.Format("2006")
	seq := strings.ToUpper(string(f.WeaponType[:min(3, len(f.WeaponType))]))
	return fmt.Sprintf("SIAR-%s-%s-%s", year, seq, uuid.New().String()[:8])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
