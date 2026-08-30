package service

import (
	"context"
	"fmt"
	"time"

	"github.com/snisid/humint-ht/internal/domain"
	"github.com/snisid/humint-ht/internal/kafka"
	"github.com/snisid/humint-ht/internal/repository"
)

type EventProducer interface {
	Publish(ctx context.Context, eventType string, payload interface{}) error
	Close() error
}

type HumintService struct {
	repo     repository.HumintRepository
	producer EventProducer
}

func NewHumintService(repo repository.HumintRepository, producer *kafka.Producer) *HumintService {
	return &HumintService{repo: repo, producer: producer}
}

func (s *HumintService) CreateSource(req domain.CreateSourceRequest) (domain.Source, error) {
	source := domain.Source{
		CodeName:          req.CodeName,
		CredibilityRating: req.CredibilityRating,
		ReliabilityRating: req.ReliabilityRating,
		HandlingOfficerID: req.HandlingOfficerID,
		PaymentAmount:     req.PaymentAmount,
		PaymentFrequency:  req.PaymentFrequency,
		RiskLevel:         req.RiskLevel,
		Compartment:       req.Compartment,
		FirstRecruitedAt:  time.Now().UTC(),
		LastContactAt:     time.Now().UTC(),
	}

	result, err := s.repo.CreateSource(source)
	if err != nil {
		return result, err
	}

	if s.producer != nil {
		ctx := context.Background()
		if pubErr := s.producer.Publish(ctx, "source.created", result); pubErr != nil {
			return result, fmt.Errorf("source created but kafka publish failed: %w", pubErr)
		}
	}

	return result, nil
}

func (s *HumintService) UpdateCredibility(code string, req domain.UpdateCredibilityRequest) (domain.Source, error) {
	return s.repo.UpdateCredibility(code, req.CredibilityRating, req.ReliabilityRating)
}

func (s *HumintService) GetReportsBySource(code string) ([]domain.IntelligenceReport, error) {
	return s.repo.GetReportsBySource(code)
}

func (s *HumintService) SubmitReport(req domain.SubmitReportRequest) (domain.IntelligenceReport, error) {
	report := domain.IntelligenceReport{
		SourceCode:      req.SourceCode,
		Classification:  req.Classification,
		ContentHash:     req.ContentHash,
		ThreatActors:    req.ThreatActors,
		SectorsTargeted: req.SectorsTargeted,
		VeracityScore:   req.VeracityScore,
		VerifiedBy:      req.VerifiedBy,
	}

	result, err := s.repo.SubmitReport(report)
	if err != nil {
		return result, err
	}

	if s.producer != nil {
		ctx := context.Background()
		if pubErr := s.producer.Publish(ctx, "report.submitted", result); pubErr != nil {
			return result, fmt.Errorf("report submitted but kafka publish failed: %w", pubErr)
		}
	}

	return result, nil
}

func (s *HumintService) LogDebriefing(req domain.LogDebriefingRequest) (domain.DebriefingSession, error) {
	sessionDate, err := time.Parse(time.RFC3339, req.SessionDate)
	if err != nil {
		return domain.DebriefingSession{}, fmt.Errorf("parse session_date: %w", err)
	}

	var nextMeeting time.Time
	if req.NextMeetingPlannedAt != "" {
		nextMeeting, err = time.Parse(time.RFC3339, req.NextMeetingPlannedAt)
		if err != nil {
			return domain.DebriefingSession{}, fmt.Errorf("parse next_meeting_planned_at: %w", err)
		}
	}

	debriefing := domain.DebriefingSession{
		SourceCode:           req.SourceCode,
		OfficerID:            req.OfficerID,
		SessionDate:          sessionDate,
		LocationMethod:       req.LocationMethod,
		TopicsCovered:        req.TopicsCovered,
		NextMeetingPlannedAt: nextMeeting,
		RiskAssessment:       req.RiskAssessment,
	}

	result, err := s.repo.LogDebriefing(debriefing)
	if err != nil {
		return result, err
	}

	if s.producer != nil {
		ctx := context.Background()
		if pubErr := s.producer.Publish(ctx, "debriefing.logged", result); pubErr != nil {
			return result, fmt.Errorf("debriefing logged but kafka publish failed: %w", pubErr)
		}
	}

	return result, nil
}

func (s *HumintService) GetHighRiskSources() ([]domain.Source, error) {
	return s.repo.GetHighRiskSources()
}

func (s *HumintService) GetSourceNetwork() (domain.SourceNetworkResponse, error) {
	sources, reports, err := s.repo.GetSourceNetwork()
	if err != nil {
		return domain.SourceNetworkResponse{}, err
	}

	var resp domain.SourceNetworkResponse
	for _, src := range sources {
		resp.Nodes = append(resp.Nodes, domain.SourceNetworkNode{
			ID:   src.CodeName,
			Type: "source",
		})
	}

	for _, rep := range reports {
		for _, actor := range rep.ThreatActors {
			resp.Nodes = append(resp.Nodes, domain.SourceNetworkNode{
				ID:   actor,
				Type: "threat_actor",
			})
			resp.Edges = append(resp.Edges, domain.SourceNetworkEdge{
				Source: rep.SourceCode,
				Target: actor,
				Label:  "reported",
			})
		}
	}

	return resp, nil
}
