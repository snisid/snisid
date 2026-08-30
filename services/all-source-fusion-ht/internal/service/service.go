package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/all-source-fusion-ht/internal/domain"
	"github.com/snisid/all-source-fusion-ht/internal/kafka"
	"github.com/snisid/all-source-fusion-ht/internal/repository"
)

type FusionService interface {
	CreateProduct(ctx context.Context, p *domain.IntelProduct) error
	GetRecentProducts(ctx context.Context) ([]domain.IntelProduct, error)
	CreateThreatActor(ctx context.Context, a *domain.ThreatActor) error
	GetHighRiskActors(ctx context.Context) ([]domain.ThreatActor, error)
	CreateCorrelation(ctx context.Context, c *domain.CrossDisciplineCorrelation) error
	GetSourceMap(ctx context.Context, productID uuid.UUID) (*domain.IntelProduct, error)
	GetNationalEstimates(ctx context.Context) ([]domain.IntelProduct, error)
}

type fusionService struct {
	repo repository.FusionRepository
	kaf  kafka.Producer
}

func NewFusionService(repo repository.FusionRepository, kaf kafka.Producer) FusionService {
	return &fusionService{repo: repo, kaf: kaf}
}

func (s *fusionService) CreateProduct(ctx context.Context, p *domain.IntelProduct) error {
	p.ProductID = uuid.New()
	p.ValidFrom = time.Now().UTC()
	if err := s.repo.CreateProduct(ctx, p); err != nil {
		return err
	}
	_ = s.kaf.Publish(ctx, "product.created", p)
	return nil
}

func (s *fusionService) GetRecentProducts(ctx context.Context) ([]domain.IntelProduct, error) {
	return s.repo.GetRecentProducts(ctx, 50)
}

func (s *fusionService) CreateThreatActor(ctx context.Context, a *domain.ThreatActor) error {
	a.ActorID = uuid.New()
	return s.repo.CreateThreatActor(ctx, a)
}

func (s *fusionService) GetHighRiskActors(ctx context.Context) ([]domain.ThreatActor, error) {
	return s.repo.GetHighRiskActors(ctx)
}

func (s *fusionService) CreateCorrelation(ctx context.Context, c *domain.CrossDisciplineCorrelation) error {
	c.CorrelationID = uuid.New()
	c.CreatedAt = time.Now().UTC()
	if err := s.repo.CreateCorrelation(ctx, c); err != nil {
		return err
	}
	_ = s.kaf.Publish(ctx, "correlation.created", c)
	return nil
}

func (s *fusionService) GetSourceMap(ctx context.Context, productID uuid.UUID) (*domain.IntelProduct, error) {
	return s.repo.GetSourceMap(ctx, productID)
}

func (s *fusionService) GetNationalEstimates(ctx context.Context) ([]domain.IntelProduct, error) {
	return s.repo.GetNationalEstimates(ctx)
}
