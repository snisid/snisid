package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/service-desk-svc/internal/domain"
	"github.com/snisid/service-desk-svc/internal/kafka"
	"github.com/snisid/service-desk-svc/internal/repository"
)

type ServiceDeskService struct {
	repo     repository.Repository
	producer *kafka.Producer
}

func NewServiceDeskService(repo repository.Repository, producer *kafka.Producer) *ServiceDeskService {
	return &ServiceDeskService{repo: repo, producer: producer}
}

func (s *ServiceDeskService) CreateCase(ctx context.Context, citizenID uuid.UUID, subject, description string) (*domain.SupportCase, error) {
	c := &domain.SupportCase{
		CaseID:      uuid.New(),
		CitizenID:   citizenID,
		Subject:     subject,
		Description: description,
		Status:      domain.CaseStatusOpen,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	if err := s.repo.CreateCase(ctx, c); err != nil {
		return nil, fmt.Errorf("create case: %w", err)
	}
	s.publishEvent(ctx, "service-desk.case.created", c.CaseID.String())
	return c, nil
}

func (s *ServiceDeskService) GetCase(ctx context.Context, caseID uuid.UUID) (*domain.SupportCase, error) {
	return s.repo.GetCase(ctx, caseID)
}

func (s *ServiceDeskService) ListCases(ctx context.Context, status domain.CaseStatus) ([]domain.SupportCase, error) {
	return s.repo.ListCases(ctx, status)
}

func (s *ServiceDeskService) IssueChallenge(ctx context.Context, caseID uuid.UUID, method domain.RecoveryMethod) (*domain.VerificationChallenge, error) {
	ch := &domain.VerificationChallenge{
		ChallengeID: uuid.New(),
		CaseID:      caseID,
		Method:      method,
		Challenge:   fmt.Sprintf("CHALLENGE-%s-%s", caseID[:8], uuid.New().String()[:8]),
		ExpiresAt:   time.Now().UTC().Add(15 * time.Minute),
		IsResolved:  false,
		CreatedAt:   time.Now().UTC(),
	}
	if err := s.repo.CreateChallenge(ctx, ch); err != nil {
		return nil, fmt.Errorf("issue challenge: %w", err)
	}
	s.publishEvent(ctx, "service-desk.challenge.issued", caseID.String())
	return ch, nil
}

func (s *ServiceDeskService) VerifyResponse(ctx context.Context, challengeID uuid.UUID) error {
	if err := s.repo.ResolveChallenge(ctx, challengeID); err != nil {
		return fmt.Errorf("verify response: %w", err)
	}
	s.publishEvent(ctx, "service-desk.challenge.resolved", challengeID.String())
	return nil
}

func (s *ServiceDeskService) ExecuteRecovery(ctx context.Context, caseID uuid.UUID, citizenID uuid.UUID, method domain.RecoveryMethod) (*domain.IdentityRecoveryRequest, error) {
	req := &domain.IdentityRecoveryRequest{
		RequestID:       uuid.New(),
		CaseID:          caseID,
		CitizenID:       citizenID,
		PreferredMethod: method,
		VerifiedMethods: []domain.RecoveryMethod{method},
		IsVerified:      true,
		CreatedAt:       time.Now().UTC(),
	}
	if err := s.repo.CreateRecoveryRequest(ctx, req); err != nil {
		return nil, fmt.Errorf("execute recovery: %w", err)
	}
	if err := s.repo.VerifyRecoveryRequest(ctx, req.RequestID); err != nil {
		return nil, fmt.Errorf("verify recovery: %w", err)
	}
	if err := s.repo.UpdateCaseStatus(ctx, caseID, domain.CaseStatusResolved); err != nil {
		return nil, fmt.Errorf("update case status: %w", err)
	}

	res := &domain.Resolution{
		ResolutionID: uuid.New(),
		CaseID:       caseID,
		Action:       "IDENTITY_RECOVERED",
		Details:      fmt.Sprintf("Identity recovered via %s", method),
		ResolvedBy:   "system",
		CreatedAt:    time.Now().UTC(),
	}
	if err := s.repo.CreateResolution(ctx, res); err != nil {
		return nil, fmt.Errorf("create resolution: %w", err)
	}

	s.publishEvent(ctx, "service-desk.recovery.executed", caseID.String())
	return req, nil
}

func (s *ServiceDeskService) AddNote(ctx context.Context, caseID uuid.UUID, author, content string) (*domain.CaseNote, error) {
	note := &domain.CaseNote{
		NoteID:    uuid.New(),
		CaseID:    caseID,
		Author:    author,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}
	if err := s.repo.AddCaseNote(ctx, note); err != nil {
		return nil, fmt.Errorf("add note: %w", err)
	}
	s.publishEvent(ctx, "service-desk.note.added", caseID.String())
	return note, nil
}

func (s *ServiceDeskService) publishEvent(ctx context.Context, eventType string, caseID string) {
	if s.producer == nil {
		return
	}
	evt := kafka.Event{
		EventType: eventType,
		CaseID:    caseID,
		Timestamp: time.Now().UTC(),
		Data:      map[string]string{"case_id": caseID},
	}
	if err := s.producer.Publish(ctx, evt); err != nil {
		log.Printf("failed to publish event: %v", err)
	}
}
