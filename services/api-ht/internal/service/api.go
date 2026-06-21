package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/api-ht/internal/domain"
	"github.com/snisid/api-ht/internal/kafka"
	"github.com/snisid/api-ht/internal/repository"
)

type APIService struct {
	repo     repository.Repository
	producer *kafka.Producer
}

func NewAPIService(repo repository.Repository, producer *kafka.Producer) *APIService {
	return &APIService{repo: repo, producer: producer}
}

func generateAPIKey() string {
	b := make([]byte, 32)
	rand.Read(b)
	return "snisid_" + hex.EncodeToString(b)
}

func (s *APIService) RegisterDeveloper(ctx context.Context, email, contactName string, orgName, contactPhone *string) (*domain.DeveloperAccount, error) {
	existing, err := s.repo.FindAccountByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("account already exists with email: %s", email)
	}

	acc := &domain.DeveloperAccount{
		ID:           uuid.New(),
		Email:        email,
		OrgName:      orgName,
		ContactName:  contactName,
		ContactPhone: contactPhone,
		IsApproved:   false,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	if err := s.repo.CreateAccount(ctx, acc); err != nil {
		return nil, fmt.Errorf("create account: %w", err)
	}

	s.publishEvent(ctx, "api.developer.registered", acc)
	return acc, nil
}

func (s *APIService) RequestKey(ctx context.Context, accountID uuid.UUID, description *string) (*domain.APIKey, error) {
	key := &domain.APIKey{
		ID:          uuid.New(),
		AccountID:   accountID,
		KeyValue:    generateAPIKey(),
		Description: description,
		IsActive:    true,
		CreatedAt:   time.Now().UTC(),
	}
	if err := s.repo.CreateAPIKey(ctx, key); err != nil {
		return nil, fmt.Errorf("create api key: %w", err)
	}

	s.publishEvent(ctx, "api.key.requested", key)
	return key, nil
}

func (s *APIService) GetCatalog(ctx context.Context) ([]domain.APIEndpoint, error) {
	return s.repo.ListCatalog(ctx)
}

func (s *APIService) GetUsage(ctx context.Context, keyID uuid.UUID) ([]domain.UsageLog, error) {
	return s.repo.ListUsageByKey(ctx, keyID)
}

func (s *APIService) RevokeKey(ctx context.Context, id uuid.UUID) error {
	key, err := s.repo.FindKeyByID(ctx, id)
	if err != nil {
		return fmt.Errorf("key not found: %w", err)
	}
	if !key.IsActive {
		return fmt.Errorf("key already revoked")
	}
	if err := s.repo.RevokeKey(ctx, id); err != nil {
		return fmt.Errorf("revoke key: %w", err)
	}

	s.publishEvent(ctx, "api.key.revoked", key)
	return nil
}

func (s *APIService) publishEvent(ctx context.Context, eventType string, data any) {
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
