package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/bio-ht/internal/domain"
	"github.com/snisid/bio-ht/internal/kafka"
	"github.com/snisid/bio-ht/internal/milvus"
	"github.com/snisid/bio-ht/internal/repository"
)

type BioService struct {
	repo     repository.Repository
	milvus   *milvus.Client
	producer *kafka.Producer
}

func NewBioService(repo repository.Repository, milvusClient *milvus.Client, producer *kafka.Producer) *BioService {
	return &BioService{repo: repo, milvus: milvusClient, producer: producer}
}

func (s *BioService) Enroll(ctx context.Context, req domain.EnrollRequest) (*domain.BioTemplate, error) {
	citizenID, err := uuid.Parse(req.CitizenID)
	if err != nil {
		return nil, fmt.Errorf("invalid citizen_id: %w", err)
	}

	capturedBy, err := uuid.Parse(req.CapturedBy)
	if err != nil {
		return nil, fmt.Errorf("invalid captured_by: %w", err)
	}

	templateID := uuid.New()

	embedding := generateMockEmbedding()

	vectorID, err := s.milvus.StoreVector(ctx, templateID, citizenID, req.Modality, embedding)
	if err != nil {
		return nil, fmt.Errorf("store vector in milvus: %w", err)
	}

	template := &domain.BioTemplate{
		TemplateID:      templateID,
		CitizenID:       citizenID,
		Modality:        domain.Modality(req.Modality),
		MilvusVectorID:  vectorID,
		QualityScore:    0.85,
		CaptureDevice:   strPtr(req.CaptureDevice),
		CaptureLocation: strPtr(req.CaptureLocation),
		CapturedBy:      capturedBy,
		IsActive:        true,
	}

	if err := s.repo.CreateTemplate(ctx, template); err != nil {
		return nil, fmt.Errorf("save template: %w", err)
	}

	s.publishEvent(ctx, "bio.template.enrolled", template)
	return template, nil
}

func (s *BioService) Verify(ctx context.Context, req domain.VerifyRequest) (*domain.VerifyResult, error) {
	citizenID, err := uuid.Parse(req.CitizenID)
	if err != nil {
		return nil, fmt.Errorf("invalid citizen_id: %w", err)
	}

	sampleEmbedding := generateMockEmbedding()

	score, err := s.milvus.Verify(ctx, req.Modality, sampleEmbedding, citizenID)
	if err != nil {
		return nil, fmt.Errorf("milvus verify: %w", err)
	}

	isMatch := score >= 0.85

	logEntry := &domain.VerificationLog{
		VerificationID:   uuid.New(),
		CitizenID:        &citizenID,
		Modality:         domain.Modality(req.Modality),
		RequestingModule: "bio-ht",
		MatchScore:       score,
		IsMatch:          isMatch,
		VerifiedAt:       time.Now(),
	}
	_ = s.repo.LogVerification(ctx, logEntry)

	return &domain.VerifyResult{
		IsMatch: isMatch,
		Score:   score,
	}, nil
}

func (s *BioService) Identify(ctx context.Context, req domain.IdentifyRequest) (*domain.IdentifyResult, error) {
	sampleEmbedding := generateMockEmbedding()

	candidates, err := s.milvus.Identify(ctx, req.Modality, sampleEmbedding, req.Threshold)
	if err != nil {
		return nil, fmt.Errorf("milvus identify: %w", err)
	}

	result := &domain.IdentifyResult{}
	for _, c := range candidates {
		result.Candidates = append(result.Candidates, domain.IdentifyCandidate{
			CitizenID:  c.CitizenID,
			Score:      c.Score,
			TemplateID: c.TemplateID,
		})
	}
	return result, nil
}

func (s *BioService) GetQuality(ctx context.Context, templateID string) (float64, error) {
	tid, err := uuid.Parse(templateID)
	if err != nil {
		return 0, fmt.Errorf("invalid template_id: %w", err)
	}

	t, err := s.repo.GetTemplate(ctx, tid)
	if err != nil {
		return 0, err
	}
	return t.QualityScore, nil
}

func (s *BioService) publishEvent(ctx context.Context, eventType string, template *domain.BioTemplate) {
	if s.producer == nil {
		return
	}
	evt := kafka.Event{
		EventType:  eventType,
		TemplateID: template.TemplateID.String(),
		CitizenID:  template.CitizenID.String(),
		Timestamp:  time.Now().UTC(),
		Data:       template,
	}
	if err := s.producer.Publish(ctx, evt); err != nil {
		log.Printf("failed to publish event: %v", err)
	}
}

func generateMockEmbedding() []float32 {
	return make([]float32, 512)
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
