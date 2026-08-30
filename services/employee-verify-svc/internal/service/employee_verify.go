package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/employee-verify-svc/internal/domain"
	"github.com/snisid/employee-verify-svc/internal/kafka"
	"github.com/snisid/employee-verify-svc/internal/repository"
)

type EmployeeVerifyService struct {
	repo     repository.Repository
	producer *kafka.Producer
}

func NewEmployeeVerifyService(repo repository.Repository, producer *kafka.Producer) *EmployeeVerifyService {
	return &EmployeeVerifyService{repo: repo, producer: producer}
}

func (s *EmployeeVerifyService) RegisterEmployer(ctx context.Context, companyName, ein, address, contactEmail, contactPhone string) (*domain.EmployerRegistration, error) {
	employer := &domain.EmployerRegistration{
		EmployerID:   uuid.New(),
		CompanyName:  companyName,
		EIN:          ein,
		Address:      address,
		ContactEmail: contactEmail,
		ContactPhone: contactPhone,
		RegisteredAt: time.Now().UTC(),
		IsActive:     true,
	}
	if err := s.repo.InsertEmployer(ctx, employer); err != nil {
		return nil, fmt.Errorf("register employer: %w", err)
	}
	s.publishEvent(ctx, "employer.registered", employer)
	return employer, nil
}

func (s *EmployeeVerifyService) CreateCase(ctx context.Context, employerID uuid.UUID, employeeName, documentNumber, documentType string) (*domain.VerificationRequest, error) {
	tcn := fmt.Sprintf("TCN-%s-%d", uuid.New().String()[:8], time.Now().Unix())

	vreq := &domain.VerificationRequest{
		TCN:            tcn,
		EmployerID:     employerID,
		EmployeeName:   employeeName,
		DocumentNumber: documentNumber,
		DocumentType:   documentType,
		Status:         domain.StatusPending,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}
	if err := s.repo.InsertCase(ctx, vreq); err != nil {
		return nil, fmt.Errorf("create case: %w", err)
	}

	history := &domain.CaseHistory{
		HistoryID:  uuid.New(),
		TCN:        tcn,
		Action:     "CASE_CREATED",
		ActionedBy: "system",
		ActionedAt: time.Now().UTC(),
		Details:    "Verification case created",
	}
	s.repo.InsertCaseHistory(ctx, history)

	s.publishEvent(ctx, "case.created", vreq)
	return vreq, nil
}

func (s *EmployeeVerifyService) GetCaseByTCN(ctx context.Context, tcn string) (*domain.VerificationRequest, error) {
	return s.repo.FindCaseByTCN(ctx, tcn)
}

func (s *EmployeeVerifyService) SubmitVerificationResponse(ctx context.Context, tcn string, ssaMatch, dhsMatch bool, reason string) (*domain.VerificationResult, error) {
	vreq, err := s.repo.FindCaseByTCN(ctx, tcn)
	if err != nil {
		return nil, fmt.Errorf("case not found: %w", err)
	}

	isEligible := ssaMatch && dhsMatch
	var status domain.VerificationStatus
	if isEligible {
		status = domain.StatusVerified
	} else {
		status = domain.StatusNotVerified
	}

	result := &domain.VerificationResult{
		ResultID:    uuid.New(),
		TCN:         tcn,
		SSAMatch:    ssaMatch,
		DHSMatch:    dhsMatch,
		IsEligible:  isEligible,
		Reason:      reason,
		CompletedAt: time.Now().UTC(),
		Status:      status,
	}
	if err := s.repo.InsertVerificationResult(ctx, result); err != nil {
		return nil, fmt.Errorf("insert result: %w", err)
	}

	if err := s.repo.UpdateCaseStatus(ctx, tcn, status); err != nil {
		return nil, fmt.Errorf("update status: %w", err)
	}

	history := &domain.CaseHistory{
		HistoryID:  uuid.New(),
		TCN:        tcn,
		Action:     fmt.Sprintf("VERIFICATION_%s", status),
		ActionedBy: "system",
		ActionedAt: time.Now().UTC(),
		Details:    fmt.Sprintf("SSA:%v DHS:%v Reason:%s", ssaMatch, dhsMatch, reason),
	}
	s.repo.InsertCaseHistory(ctx, history)

	s.publishEvent(ctx, fmt.Sprintf("case.%s", status), vreq)
	return result, nil
}

func (s *EmployeeVerifyService) ListCasesByEmployer(ctx context.Context, ein string) ([]domain.VerificationRequest, error) {
	return s.repo.FindCasesByEmployer(ctx, ein)
}

func (s *EmployeeVerifyService) GetStats(ctx context.Context) (map[string]int, error) {
	return s.repo.GetStats(ctx)
}

func (s *EmployeeVerifyService) publishEvent(ctx context.Context, eventType string, data any) {
	if s.producer == nil {
		return
	}
	var tcn string
	if v, ok := data.(*domain.VerificationRequest); ok {
		tcn = v.TCN
	}
	evt := kafka.Event{
		EventType: eventType,
		TCN:       tcn,
		Timestamp: time.Now().UTC(),
		Data:      data,
	}
	if err := s.producer.Publish(ctx, evt); err != nil {
		log.Printf("failed to publish event: %v", err)
	}
}
