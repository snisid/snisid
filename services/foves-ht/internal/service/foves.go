package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/foves-ht/internal/domain"
	"github.com/snisid/foves-ht/internal/kafka"
	"github.com/snisid/foves-ht/internal/repository"
)

type FovesService struct {
	repo     repository.Repository
	producer *kafka.Producer
}

func NewFovesService(repo repository.Repository, producer *kafka.Producer) *FovesService {
	return &FovesService{repo: repo, producer: producer}
}

func (s *FovesService) RegisterVehicle(ctx context.Context, v *domain.Vehicle) (*domain.Vehicle, error) {
	v.ID = uuid.New()
	v.IsStolen = false
	v.IsActive = true
	v.RegisteredAt = time.Now().UTC()
	v.UpdatedAt = time.Now().UTC()

	if err := s.repo.CreateVehicle(ctx, v); err != nil {
		return nil, fmt.Errorf("register vehicle: %w", err)
	}

	s.publishEvent(ctx, "foves.vehicle.registered", v)
	return v, nil
}

func (s *FovesService) GetByPlate(ctx context.Context, plate string) (*domain.Vehicle, error) {
	return s.repo.FindByPlate(ctx, plate)
}

func (s *FovesService) GetByVIN(ctx context.Context, vin string) (*domain.Vehicle, error) {
	return s.repo.FindByVIN(ctx, vin)
}

func (s *FovesService) GetByOwner(ctx context.Context, citizenID uuid.UUID) ([]domain.Vehicle, error) {
	return s.repo.FindByOwner(ctx, citizenID)
}

func (s *FovesService) TransferOwnership(ctx context.Context, vehicleID, fromCitizenID, toCitizenID uuid.UUID, contractRef *string) (*domain.OwnershipTransfer, error) {
	transfer := &domain.OwnershipTransfer{
		ID:            uuid.New(),
		VehicleID:     vehicleID,
		FromCitizenID: fromCitizenID,
		ToCitizenID:   toCitizenID,
		TransferDate:  time.Now().UTC(),
		ContractRef:   contractRef,
		CreatedAt:     time.Now().UTC(),
	}

	if err := s.repo.CreateTransfer(ctx, transfer); err != nil {
		return nil, fmt.Errorf("create transfer: %w", err)
	}

	if err := s.repo.UpdateVehicleOwner(ctx, vehicleID, toCitizenID); err != nil {
		return nil, fmt.Errorf("update vehicle owner: %w", err)
	}

	s.publishEvent(ctx, "foves.ownership.transferred", transfer)
	return transfer, nil
}

func (s *FovesService) IssueLicense(ctx context.Context, l *domain.DriverLicense) (*domain.DriverLicense, error) {
	l.ID = uuid.New()
	l.IssuedDate = time.Now().UTC()
	l.CreatedAt = time.Now().UTC()
	l.UpdatedAt = time.Now().UTC()

	if err := s.repo.CreateLicense(ctx, l); err != nil {
		return nil, fmt.Errorf("issue license: %w", err)
	}

	s.publishEvent(ctx, "foves.license.issued", l)
	return l, nil
}

func (s *FovesService) GetLicense(ctx context.Context, citizenID uuid.UUID) (*domain.DriverLicense, error) {
	return s.repo.FindLicenseByCitizen(ctx, citizenID)
}

func (s *FovesService) publishEvent(ctx context.Context, eventType string, data any) {
	if s.producer == nil {
		return
	}
	evt := kafka.Event{
		EventType: eventType,
		Timestamp: time.Now().UTC(),
		Data:      data,
	}
	if err := s.producer.Publish(ctx, evt); err != nil {
		log.Printf("failed to publish event: %v", err)
	}
}
