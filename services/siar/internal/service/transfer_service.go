package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/siar/internal/domain"
)

var (
	ErrFirearmNotFound  = errors.New("arme non trouvée")
	ErrLicenseNotFound  = errors.New("licence non trouvée")
	ErrDealerNotFound   = errors.New("marchand d'armes non trouvé")
	ErrTransferNotFound = errors.New("transfert non trouvé")
	ErrSeizureNotFound  = errors.New("saisie non trouvée")
)

type TransferService struct {
	mu        sync.RWMutex
	transfers map[uuid.UUID]*domain.Transfer
	seizures  map[uuid.UUID]*domain.Seizure
}

func NewTransferService() *TransferService {
	return &TransferService{
		transfers: make(map[uuid.UUID]*domain.Transfer),
		seizures:  make(map[uuid.UUID]*domain.Seizure),
	}
}

func (s *TransferService) CreateTransfer(ctx context.Context, req domain.CreateTransferRequest, authorizedBy *uuid.UUID) (*domain.Transfer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := &domain.Transfer{
		TransferID:    uuid.New(),
		FirearmID:     req.FirearmID,
		FromOwnerID:   req.FromOwnerID,
		FromOwnerName: req.FromOwnerName,
		ToOwnerID:     req.ToOwnerID,
		ToOwnerName:   req.ToOwnerName,
		TransferType:  req.TransferType,
		TransferDate:  req.TransferDate,
		PermitRef:     req.PermitRef,
		AuthorizedBy:  authorizedBy,
		Notes:         req.Notes,
		CreatedAt:     time.Now(),
	}

	s.transfers[t.TransferID] = t
	return t, nil
}

func (s *TransferService) GetTransfer(ctx context.Context, id uuid.UUID) (*domain.Transfer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.transfers[id]
	if !ok {
		return nil, ErrTransferNotFound
	}
	return t, nil
}

func (s *TransferService) ListByFirearm(ctx context.Context, firearmID uuid.UUID) ([]*domain.Transfer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*domain.Transfer
	for _, t := range s.transfers {
		if t.FirearmID == firearmID {
			result = append(result, t)
		}
	}
	return result, nil
}

func (s *TransferService) ReportSeizure(ctx context.Context, req domain.CreateSeizureRequest, createdBy *uuid.UUID) (*domain.Seizure, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	seiz := &domain.Seizure{
		SeizureID:      uuid.New(),
		FirearmID:      req.FirearmID,
		SerialNumber:   req.SerialNumber,
		Make:           req.Make,
		Model:          req.Model,
		SeizureDate:    req.SeizureDate,
		SeizingUnit:    req.SeizingUnit,
		SeizingOfficer: req.SeizingOfficer,
		LocationDesc:   req.LocationDesc,
		DeptCode:       req.DeptCode,
		Context:        req.Context,
		FromPersonID:   req.FromPersonID,
		FromPersonName: req.FromPersonName,
		GangID:         req.GangID,
		CaseReference:  req.CaseReference,
		CreatedBy:      createdBy,
		CreatedAt:      time.Now(),
	}

	s.seizures[seiz.SeizureID] = seiz
	return seiz, nil
}

func (s *TransferService) ReportStolen(ctx context.Context, req domain.ReportStolenRequest, createdBy *uuid.UUID) (*domain.Seizure, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	seiz := &domain.Seizure{
		SeizureID:    uuid.New(),
		SerialNumber: req.SerialNumber,
		Make:         req.Make,
		Model:        req.Model,
		SeizureDate:  req.IncidentDate,
		SeizingUnit:  "PNH_STOLEN_REPORT",
		LocationDesc: req.Notes,
		DeptCode:     req.DeptCode,
		FromPersonID: req.OwnerID,
		Context:      "VOL_DECLARE",
		CreatedBy:    createdBy,
		CreatedAt:    time.Now(),
	}

	s.seizures[seiz.SeizureID] = seiz
	return seiz, nil
}

func (s *TransferService) GetSeizure(ctx context.Context, id uuid.UUID) (*domain.Seizure, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	seiz, ok := s.seizures[id]
	if !ok {
		return nil, ErrSeizureNotFound
	}
	return seiz, nil
}

func (s *TransferService) ListSeizures(ctx context.Context) ([]*domain.Seizure, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*domain.Seizure, 0, len(s.seizures))
	for _, sez := range s.seizures {
		result = append(result, sez)
	}
	return result, nil
}
