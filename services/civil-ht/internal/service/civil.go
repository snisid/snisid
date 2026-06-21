package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/civil-ht/internal/domain"
	"github.com/snisid/civil-ht/internal/kafka"
	"github.com/snisid/civil-ht/internal/repository"
)

type CivilService struct {
	repo     repository.Repository
	producer *kafka.Producer
}

func NewCivilService(repo repository.Repository, producer *kafka.Producer) *CivilService {
	return &CivilService{repo: repo, producer: producer}
}

func generateActNumber(actType domain.CivilActType, deptCode string, year int, seq int) string {
	return fmt.Sprintf("ACTE-HT-%04d-%s-%s-%06d", year, deptCode, string(actType[:1]), seq)
}

func (s *CivilService) DeclareBirth(ctx context.Context, req domain.BirthDeclaration, registerInfo domain.CivilAct) (*domain.CivilAct, error) {
	registerInfo.ActID = uuid.New()
	registerInfo.ActType = domain.ActBirth
	registerInfo.DeclaredDate = time.Now()
	registerInfo.CreatedAt = time.Now().UTC()
	registerInfo.ActNumber = generateActNumber(domain.ActBirth, registerInfo.DeptCode, registerInfo.EventDate.Year(), 1)

	birth := &req
	birth.ActID = registerInfo.ActID

	if err := s.repo.CreateBirth(ctx, &registerInfo, birth); err != nil {
		return nil, fmt.Errorf("declare birth: %w", err)
	}

	s.publishEvent(ctx, "civil.birth.declared", &registerInfo)
	return &registerInfo, nil
}

func (s *CivilService) DeclareDeath(ctx context.Context, req domain.DeathDeclaration, registerInfo domain.CivilAct) (*domain.CivilAct, error) {
	registerInfo.ActID = uuid.New()
	registerInfo.ActType = domain.ActDeath
	registerInfo.DeclaredDate = time.Now()
	registerInfo.CreatedAt = time.Now().UTC()
	registerInfo.ActNumber = generateActNumber(domain.ActDeath, registerInfo.DeptCode, registerInfo.EventDate.Year(), 1)

	death := &req
	death.ActID = registerInfo.ActID

	if err := s.repo.CreateDeath(ctx, &registerInfo, death); err != nil {
		return nil, fmt.Errorf("declare death: %w", err)
	}

	s.publishEvent(ctx, "civil.death.declared", &registerInfo)
	return &registerInfo, nil
}

func (s *CivilService) RegisterMarriage(ctx context.Context, req domain.MarriageDeclaration, registerInfo domain.CivilAct) (*domain.CivilAct, error) {
	registerInfo.ActID = uuid.New()
	registerInfo.ActType = domain.ActMarriage
	registerInfo.DeclaredDate = time.Now()
	registerInfo.CreatedAt = time.Now().UTC()
	registerInfo.ActNumber = generateActNumber(domain.ActMarriage, registerInfo.DeptCode, registerInfo.EventDate.Year(), 1)

	marriage := &req
	marriage.ActID = registerInfo.ActID

	if err := s.repo.CreateMarriage(ctx, &registerInfo, marriage); err != nil {
		return nil, fmt.Errorf("register marriage: %w", err)
	}

	s.publishEvent(ctx, "civil.marriage.registered", &registerInfo)
	return &registerInfo, nil
}

func (s *CivilService) GetAct(ctx context.Context, actNumber string) (*domain.CivilAct, error) {
	return s.repo.FindByActNumber(ctx, actNumber)
}

func (s *CivilService) GetCitizenActs(ctx context.Context, citizenID string) ([]domain.CivilAct, error) {
	cid, err := uuid.Parse(citizenID)
	if err != nil {
		return nil, fmt.Errorf("invalid citizen id: %w", err)
	}
	return s.repo.FindByCitizenID(ctx, cid)
}

func (s *CivilService) GetBirthDetails(ctx context.Context, actID uuid.UUID) (*domain.BirthDeclaration, error) {
	return s.repo.FindBirthDetails(ctx, actID)
}

func (s *CivilService) GetDeathDetails(ctx context.Context, actID uuid.UUID) (*domain.DeathDeclaration, error) {
	return s.repo.FindDeathDetails(ctx, actID)
}

func (s *CivilService) GetMarriageDetails(ctx context.Context, actID uuid.UUID) (*domain.MarriageDeclaration, error) {
	return s.repo.FindMarriageDetails(ctx, actID)
}

func (s *CivilService) publishEvent(ctx context.Context, eventType string, act *domain.CivilAct) {
	if s.producer == nil {
		return
	}
	evt := kafka.Event{
		EventType: eventType,
		ActID:     act.ActID.String(),
		ActNumber: act.ActNumber,
		Timestamp: time.Now().UTC(),
		Data:      act,
	}
	if err := s.producer.Publish(ctx, evt); err != nil {
		log.Printf("failed to publish event: %v", err)
	}
}
