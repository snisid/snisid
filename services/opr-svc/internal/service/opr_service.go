package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/opr-svc/internal/domain"
)

type OPRService struct {
	repo     domain.ProtectionOrderRepository
	violRepo domain.ViolationRepository
	eventPub domain.EventPublisher
}

func NewOPRService(
	repo domain.ProtectionOrderRepository,
	violRepo domain.ViolationRepository,
	eventPub domain.EventPublisher,
) *OPRService {
	return &OPRService{
		repo:     repo,
		violRepo: violRepo,
		eventPub: eventPub,
	}
}

type CreateOrderRequest struct {
	OrderType          string   `json:"order_type" binding:"required"`
	ProtectedPersonID  string   `json:"protected_person_id" binding:"required"`
	SubjectPersonID    string   `json:"subject_person_id" binding:"required"`
	ExclusionRadiusM   int      `json:"exclusion_radius_m"`
	ExclusionAddresses []string `json:"exclusion_addresses"`
	NoContactModes     []string `json:"no_contact_modes"`
	IssuingCourt       string   `json:"issuing_court" binding:"required"`
	IssuingJudge       string   `json:"issuing_judge"`
	ExpiryDate         string   `json:"expiry_date" binding:"required"`
	CreatedBy          string   `json:"created_by" binding:"required"`
}

func (s *OPRService) CreateOrder(ctx context.Context, req CreateOrderRequest) (*domain.ProtectionOrder, error) {
	protectedID, err := uuid.Parse(req.ProtectedPersonID)
	if err != nil {
		return nil, fmt.Errorf("UUID invalide: %w", err)
	}

	subjectID, err := uuid.Parse(req.SubjectPersonID)
	if err != nil {
		return nil, fmt.Errorf("UUID invalide: %w", err)
	}

	createdBy, err := uuid.Parse(req.CreatedBy)
	if err != nil {
		return nil, fmt.Errorf("UUID invalide: %w", err)
	}

	expiryDate, err := time.Parse("2006-01-02", req.ExpiryDate)
	if err != nil {
		return nil, fmt.Errorf("date invalide: %w", err)
	}

	order := &domain.ProtectionOrder{
		OrderID:            uuid.New(),
		OrderNumber:        fmt.Sprintf("OPR-HT-%d-%s", time.Now().Year(), uuid.New().String()[:8]),
		OrderType:          domain.OrderType(req.OrderType),
		Status:             domain.StatusActive,
		ProtectedPersonID:  protectedID,
		SubjectPersonID:    subjectID,
		ExclusionRadiusM:   req.ExclusionRadiusM,
		ExclusionAddresses: req.ExclusionAddresses,
		NoContactModes:     req.NoContactModes,
		IssuingCourt:       req.IssuingCourt,
		IssuingJudge:       req.IssuingJudge,
		IssueDate:          time.Now(),
		ExpiryDate:         expiryDate,
		IsRenewable:        true,
		CreatedBy:          createdBy,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, nil
	}

	_ = s.eventPub.Publish("opr.order.created", order)
	return order, nil
}

func (s *OPRService) CheckSubject(ctx context.Context, personID uuid.UUID) (*domain.OPRCheckResult, error) {
	orders, err := s.repo.FindActiveBySubject(ctx, personID)
	if err != nil || len(orders) == 0 {
		return &domain.OPRCheckResult{HasActiveOrder: false}, nil
	}
	return &domain.OPRCheckResult{
		HasActiveOrder: true,
		Orders:         orders,
		HighestType:    s.getHighestSeverity(orders),
	}, nil
}

func (s *OPRService) RecordViolation(ctx context.Context, req domain.ViolationRequest, reportedBy uuid.UUID) error {
	order, err := s.repo.FindByID(ctx, req.OrderID)
	if err != nil {
		return fmt.Errorf("ordonnance introuvable: %w", err)
	}

	violation := &domain.Violation{
		ViolationID:   uuid.New(),
		OrderID:       req.OrderID,
		ViolationDate: time.Now(),
		ViolationType: req.ViolationType,
		LocationDesc:  req.LocationDesc,
		DeptCode:      req.DeptCode,
		ReportedBy:    reportedBy,
		CreatedAt:     time.Now(),
	}

	if err := s.violRepo.Create(ctx, violation); err != nil {
		return fmt.Errorf("création violation: %w", err)
	}

	now := time.Now()
	order.ViolationCount++
	order.LastViolationAt = &now
	order.Status = domain.StatusViolated
	_ = s.repo.Update(ctx, order)

	_ = s.eventPub.Publish("opr.violation.reported", domain.ViolationEvent{
		OrderID:    req.OrderID,
		PersonID:   order.SubjectPersonID,
		ViolType:   req.ViolationType,
		ReportedBy: reportedBy,
	})

	if order.ViolationCount >= 3 {
		_ = s.eventPub.Publish("opr.warrant.request", domain.WarrantRequestEvent{
			PersonID: order.SubjectPersonID,
			Reason:   "OPR violations repetees >= 3",
		})
	}

	return nil
}

func (s *OPRService) GetExpiringSoon(ctx context.Context, days int) ([]*domain.ProtectionOrder, error) {
	return s.repo.FindExpiringSoon(ctx, days)
}

func (s *OPRService) GetByGangID(ctx context.Context, gangID uuid.UUID) ([]*domain.ProtectionOrder, error) {
	return s.repo.FindByGangID(ctx, gangID)
}

func (s *OPRService) getHighestSeverity(orders []*domain.ProtectionOrder) domain.OrderType {
	severity := map[domain.OrderType]int{
		domain.OrderTypeNoContact:         1,
		domain.OrderTypeStayAway:          2,
		domain.OrderTypeRestraintingOrder: 3,
		domain.OrderTypeProtective:        4,
		domain.OrderTypeTravelRestriction: 5,
		domain.OrderTypeGangExclusionZone: 6,
		domain.OrderTypeWitnessProtection: 7,
	}
	highest := orders[0].OrderType
	for _, o := range orders {
		if severity[o.OrderType] > severity[highest] {
			highest = o.OrderType
		}
	}
	return highest
}
