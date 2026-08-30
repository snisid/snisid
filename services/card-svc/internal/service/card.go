package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/card-svc/internal/domain"
	"github.com/snisid/card-svc/internal/kafka"
	"github.com/snisid/card-svc/internal/repository"
)

type CardService struct {
	repo     repository.Repository
	producer *kafka.Producer
}

func NewCardService(repo repository.Repository, producer *kafka.Producer) *CardService {
	return &CardService{repo: repo, producer: producer}
}

func (s *CardService) OrderPersonalization(ctx context.Context, req domain.PersonalizationRequest) (*domain.PersonalizationRequest, error) {
	if _, err := s.repo.FindProfileByID(ctx, req.ProfileID); err != nil {
		return nil, fmt.Errorf("profile not found: %w", err)
	}

	req.OrderID = uuid.New()
	req.Status = domain.CardStatusOrdered
	req.OrderedAt = time.Now().UTC()
	req.CreatedAt = time.Now().UTC()
	req.UpdatedAt = time.Now().UTC()

	if req.CardSerial == "" {
		req.CardSerial = fmt.Sprintf("SN-%s-%06d", time.Now().Format("20060102"), uuid.New().ID()%1000000)
	}

	if err := s.repo.CreatePersonalizationOrder(ctx, &req); err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}

	s.publishEvent(ctx, "card.ordered", &req)
	return &req, nil
}

func (s *CardService) GetCard(ctx context.Context, cardSerial string) (*domain.PersonalizationRequest, error) {
	return s.repo.FindOrderBySerial(ctx, cardSerial)
}

func (s *CardService) ActivateCard(ctx context.Context, cardSerial string) (*domain.PersonalizationRequest, error) {
	card, err := s.repo.FindOrderBySerial(ctx, cardSerial)
	if err != nil {
		return nil, fmt.Errorf("card not found: %w", err)
	}

	if card.Status != domain.CardStatusIssued {
		return nil, fmt.Errorf("card must be in ISSUED status to activate, current: %s", card.Status)
	}

	if err := s.repo.UpdateOrderActivated(ctx, cardSerial); err != nil {
		return nil, fmt.Errorf("activate card: %w", err)
	}

	card.Status = domain.CardStatusActive
	now := time.Now().UTC()
	card.ActivatedAt = &now
	card.UpdatedAt = now

	s.publishEvent(ctx, "card.activated", card)
	return card, nil
}

func (s *CardService) BlockCard(ctx context.Context, cardSerial string, reason string) (*domain.PersonalizationRequest, error) {
	card, err := s.repo.FindOrderBySerial(ctx, cardSerial)
	if err != nil {
		return nil, fmt.Errorf("card not found: %w", err)
	}

	if err := s.repo.UpdateOrderBlocked(ctx, cardSerial, reason); err != nil {
		return nil, fmt.Errorf("block card: %w", err)
	}

	card.Status = domain.CardStatusBlocked
	now := time.Now().UTC()
	card.BlockedAt = &now
	card.BlockReason = reason
	card.UpdatedAt = now

	s.publishEvent(ctx, "card.blocked", card)
	return card, nil
}

func (s *CardService) GetInventory(ctx context.Context, profileID string) (any, error) {
	if profileID != "" {
		pid, err := uuid.Parse(profileID)
		if err != nil {
			return nil, fmt.Errorf("invalid profile id: %w", err)
		}
		return s.repo.FindInventoryByProfileID(ctx, pid)
	}
	return s.repo.FindAllInventory(ctx)
}

func (s *CardService) RecordShipment(ctx context.Context, shipment domain.Shipment) (*domain.Shipment, error) {
	shipment.ShipmentID = uuid.New()
	shipment.CreatedAt = time.Now().UTC()

	now := time.Now().UTC()
	shipment.ReceivedAt = &now

	if err := s.repo.CreateShipment(ctx, &shipment); err != nil {
		return nil, fmt.Errorf("record shipment: %w", err)
	}

	stock := domain.CardStock{
		StockID:      uuid.New(),
		ProfileID:    shipment.ProfileID,
		SerialFrom:   shipment.SerialFrom,
		SerialTo:     shipment.SerialTo,
		Quantity:     shipment.Quantity,
		AvailableQty: shipment.Quantity,
		Location:     shipment.Vendor,
		ReceivedAt:   now,
	}

	if err := s.repo.CreateStock(ctx, &stock); err != nil {
		return nil, fmt.Errorf("create stock record: %w", err)
	}

	s.publishEvent(ctx, "card.shipment.received", &shipment)
	return &shipment, nil
}

func (s *CardService) publishEvent(ctx context.Context, eventType string, data any) {
	if s.producer == nil {
		return
	}
	evt := kafka.Event{
		EventType: eventType,
		Timestamp: time.Now().UTC(),
		Data:      data,
	}
	if card, ok := data.(*domain.PersonalizationRequest); ok {
		evt.CardSerial = card.CardSerial
		evt.Status = string(card.Status)
	}
	if sh, ok := data.(*domain.Shipment); ok {
		evt.CardSerial = fmt.Sprintf("%s-%s", sh.SerialFrom, sh.SerialTo)
	}
	if err := s.producer.Publish(ctx, evt); err != nil {
		log.Printf("failed to publish event: %v", err)
	}
}
