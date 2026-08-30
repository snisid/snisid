package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/enrollment-svc/internal/domain"
	"github.com/snisid/enrollment-svc/internal/kafka"
	"github.com/snisid/enrollment-svc/internal/repository"
)

type EnrollmentService struct {
	repo     repository.Repository
	producer *kafka.Producer
}

func NewEnrollmentService(repo repository.Repository, producer *kafka.Producer) *EnrollmentService {
	return &EnrollmentService{repo: repo, producer: producer}
}

func (s *EnrollmentService) SubmitRequest(ctx context.Context, req domain.EnrollmentRequest) (*domain.EnrollmentRequest, error) {
	req.RequestID = uuid.New()
	req.Status = domain.StatusPendingDocuments
	req.SubmittedAt = time.Now().UTC()
	req.UpdatedAt = time.Now().UTC()

	if req.ProofingLevel == "" {
		req.ProofingLevel = domain.IAL2
	}

	if err := s.repo.CreateRequest(ctx, &req); err != nil {
		return nil, fmt.Errorf("submit enrollment: %w", err)
	}

	s.publishEvent(ctx, "enrollment.request.submitted", &req)
	return &req, nil
}

func (s *EnrollmentService) UploadDocuments(ctx context.Context, requestID uuid.UUID, docs []domain.IdentityDocument) ([]domain.IdentityDocument, error) {
	req, err := s.repo.FindRequestByID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("request not found: %w", err)
	}

	var saved []domain.IdentityDocument
	for i := range docs {
		docs[i].DocID = uuid.New()
		docs[i].RequestID = requestID
		docs[i].UploadedAt = time.Now().UTC()

		if err := s.repo.CreateDocument(ctx, &docs[i]); err != nil {
			return nil, fmt.Errorf("save document: %w", err)
		}
		saved = append(saved, docs[i])
	}

	if req.Status == domain.StatusPendingDocuments || req.Status == domain.StatusDraft {
		if err := s.repo.UpdateRequestStatus(ctx, requestID, domain.StatusDocumentsReceived); err != nil {
			return nil, fmt.Errorf("update status: %w", err)
		}
	}

	s.publishEvent(ctx, "enrollment.documents.uploaded", map[string]any{
		"request_id": requestID.String(),
		"count":      len(docs),
	})
	return saved, nil
}

func (s *EnrollmentService) CaptureBiometrics(ctx context.Context, requestID uuid.UUID, samples []domain.BiometricSample) ([]domain.BiometricSample, error) {
	req, err := s.repo.FindRequestByID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("request not found: %w", err)
	}

	var saved []domain.BiometricSample
	for i := range samples {
		samples[i].SampleID = uuid.New()
		samples[i].RequestID = requestID
		samples[i].CapturedAt = time.Now().UTC()

		if err := s.repo.CreateBiometricSample(ctx, &samples[i]); err != nil {
			return nil, fmt.Errorf("save biometric sample: %w", err)
		}
		saved = append(saved, samples[i])
	}

	status := req.Status
	if status == domain.StatusDocumentsReceived || status == domain.StatusPendingBiometrics {
		if err := s.repo.UpdateRequestStatus(ctx, requestID, domain.StatusBiometricsCaptured); err != nil {
			return nil, fmt.Errorf("update status: %w", err)
		}
	}

	s.publishEvent(ctx, "enrollment.biometrics.captured", map[string]any{
		"request_id": requestID.String(),
		"count":      len(samples),
	})
	return saved, nil
}

func (s *EnrollmentService) ReviewRequest(ctx context.Context, requestID uuid.UUID, review domain.EnrollmentReview) (*domain.EnrollmentRequest, error) {
	req, err := s.repo.FindRequestByID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("request not found: %w", err)
	}

	if req.Status != domain.StatusPendingReview && req.Status != domain.StatusBiometricsCaptured {
		if err := s.repo.UpdateRequestStatus(ctx, requestID, domain.StatusPendingReview); err != nil {
			return nil, fmt.Errorf("mark pending review: %w", err)
		}
	}

	review.ReviewID = uuid.New()
	review.RequestID = requestID
	review.ReviewedAt = time.Now().UTC()

	if err := s.repo.CreateReview(ctx, &review); err != nil {
		return nil, fmt.Errorf("save review: %w", err)
	}

	newStatus := domain.StatusApproved
	if review.Decision == "reject" {
		newStatus = domain.StatusRejected
	}
	if err := s.repo.UpdateRequestStatus(ctx, requestID, newStatus); err != nil {
		return nil, fmt.Errorf("update status after review: %w", err)
	}

	req.Status = newStatus

	eventType := "enrollment.request.approved"
	if review.Decision == "reject" {
		eventType = "enrollment.request.rejected"
	}
	s.publishEvent(ctx, eventType, map[string]any{
		"request_id": requestID.String(),
		"decision":   review.Decision,
		"officer":    review.OfficerName,
	})

	return req, nil
}

func (s *EnrollmentService) GetRequest(ctx context.Context, requestID uuid.UUID) (*domain.EnrollmentRequest, error) {
	return s.repo.FindRequestByID(ctx, requestID)
}

func (s *EnrollmentService) ListPending(ctx context.Context) ([]domain.EnrollmentRequest, error) {
	return s.repo.FindPendingRequests(ctx)
}

func (s *EnrollmentService) publishEvent(ctx context.Context, eventType string, data any) {
	if s.producer == nil {
		return
	}
	evt := kafka.Event{
		EventType: eventType,
		Timestamp: time.Now().UTC(),
		Data:      data,
	}
	if req, ok := data.(*domain.EnrollmentRequest); ok {
		evt.RequestID = req.RequestID.String()
		evt.Status = string(req.Status)
	}
	if m, ok := data.(map[string]any); ok {
		if rid, ok := m["request_id"]; ok {
			evt.RequestID = fmt.Sprintf("%v", rid)
		}
	}
	if err := s.producer.Publish(ctx, evt); err != nil {
		log.Printf("failed to publish event: %v", err)
	}
}
