package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/age-verification-svc/internal/domain"
	"github.com/snisid/age-verification-svc/internal/kafka"
	"github.com/snisid/age-verification-svc/internal/repository"
)

type AgeVerificationService struct {
	repo     repository.Repository
	producer *kafka.Producer
}

func NewAgeVerificationService(repo repository.Repository, producer *kafka.Producer) *AgeVerificationService {
	return &AgeVerificationService{repo: repo, producer: producer}
}

func (s *AgeVerificationService) CreateAttestation(ctx context.Context, identityID uuid.UUID, dateOfBirth time.Time) (*domain.AgeAttestation, error) {
	attestation := &domain.AgeAttestation{
		AttestationID: uuid.New(),
		IdentityID:    identityID,
		DateOfBirth:   dateOfBirth,
		IssuedAt:      time.Now().UTC(),
		ExpiresAt:     time.Now().UTC().AddDate(1, 0, 0),
		IsRevoked:     false,
	}
	if err := s.repo.InsertAttestation(ctx, attestation); err != nil {
		return nil, fmt.Errorf("create attestation: %w", err)
	}
	s.publishEvent(ctx, "age.attestation.created", attestation)
	return attestation, nil
}

func (s *AgeVerificationService) VerifyAgeClaim(ctx context.Context, attestationID uuid.UUID, verifierID string, bracket domain.AgeBracket) (*domain.AgeClaim, error) {
	attestation, err := s.repo.FindAttestationByID(ctx, attestationID)
	if err != nil {
		return nil, fmt.Errorf("attestation not found: %w", err)
	}
	if attestation.IsRevoked {
		return nil, fmt.Errorf("attestation is revoked")
	}
	if time.Now().UTC().After(attestation.ExpiresAt) {
		return nil, fmt.Errorf("attestation is expired")
	}

	age := time.Now().UTC().Year() - attestation.DateOfBirth.Year()
	if time.Now().UTC().YearDay() < attestation.DateOfBirth.YearDay() {
		age--
	}

	var requiredAge int
	switch bracket {
	case domain.AgeBracketOver18:
		requiredAge = 18
	case domain.AgeBracketOver21:
		requiredAge = 21
	case domain.AgeBracketOver65:
		requiredAge = 65
	default:
		return nil, fmt.Errorf("unknown age bracket: %s", bracket)
	}

	isSatisfied := age >= requiredAge

	claim := &domain.AgeClaim{
		ClaimID:       uuid.New(),
		AttestationID: attestationID,
		VerifierID:    verifierID,
		Bracket:       bracket,
		IsSatisfied:   isSatisfied,
		ClaimedAt:     time.Now().UTC(),
	}
	if err := s.repo.InsertAgeClaim(ctx, claim); err != nil {
		return nil, fmt.Errorf("insert age claim: %w", err)
	}

	s.publishEvent(ctx, "age.claim.verified", claim)
	return claim, nil
}

func (s *AgeVerificationService) GetAttestation(ctx context.Context, attestationID uuid.UUID) (*domain.AgeAttestation, error) {
	return s.repo.FindAttestationByID(ctx, attestationID)
}

func (s *AgeVerificationService) SelectiveBracketVerification(ctx context.Context, attestationID uuid.UUID, verifierID string, bracket domain.AgeBracket) (*domain.AgeClaim, error) {
	return s.VerifyAgeClaim(ctx, attestationID, verifierID, bracket)
}

func (s *AgeVerificationService) RevokeAttestation(ctx context.Context, attestationID uuid.UUID) error {
	attestation, err := s.repo.FindAttestationByID(ctx, attestationID)
	if err != nil {
		return fmt.Errorf("attestation not found: %w", err)
	}
	if attestation.IsRevoked {
		return fmt.Errorf("attestation already revoked")
	}
	if err := s.repo.UpdateAttestationRevoked(ctx, attestationID); err != nil {
		return fmt.Errorf("revoke attestation: %w", err)
	}
	s.publishEvent(ctx, "age.attestation.revoked", attestation)
	return nil
}

func (s *AgeVerificationService) publishEvent(ctx context.Context, eventType string, data any) {
	if s.producer == nil {
		return
	}
	var attestationID string
	if att, ok := data.(*domain.AgeAttestation); ok {
		attestationID = att.AttestationID.String()
	}
	evt := kafka.Event{
		EventType:     eventType,
		AttestationID: attestationID,
		Timestamp:     time.Now().UTC(),
		Data:          data,
	}
	if err := s.producer.Publish(ctx, evt); err != nil {
		log.Printf("failed to publish event: %v", err)
	}
}
