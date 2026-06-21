package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/snisid/all-source-fusion-ht/internal/domain"
	"github.com/snisid/all-source-fusion-ht/internal/kafka"
	"github.com/snisid/all-source-fusion-ht/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockFusionRepo struct {
	repository.FusionRepository
	createProductFn        func(ctx context.Context, p *domain.IntelProduct) error
	getRecentProductsFn    func(ctx context.Context, limit int) ([]domain.IntelProduct, error)
	createThreatActorFn    func(ctx context.Context, a *domain.ThreatActor) error
	getHighRiskActorsFn    func(ctx context.Context) ([]domain.ThreatActor, error)
	createCorrelationFn    func(ctx context.Context, c *domain.CrossDisciplineCorrelation) error
	getSourceMapFn         func(ctx context.Context, productID uuid.UUID) (*domain.IntelProduct, error)
	getNationalEstimatesFn func(ctx context.Context) ([]domain.IntelProduct, error)
}

func (m *mockFusionRepo) CreateProduct(ctx context.Context, p *domain.IntelProduct) error {
	return m.createProductFn(ctx, p)
}
func (m *mockFusionRepo) GetRecentProducts(ctx context.Context, limit int) ([]domain.IntelProduct, error) {
	return m.getRecentProductsFn(ctx, limit)
}
func (m *mockFusionRepo) CreateThreatActor(ctx context.Context, a *domain.ThreatActor) error {
	return m.createThreatActorFn(ctx, a)
}
func (m *mockFusionRepo) GetHighRiskActors(ctx context.Context) ([]domain.ThreatActor, error) {
	return m.getHighRiskActorsFn(ctx)
}
func (m *mockFusionRepo) CreateCorrelation(ctx context.Context, c *domain.CrossDisciplineCorrelation) error {
	return m.createCorrelationFn(ctx, c)
}
func (m *mockFusionRepo) GetSourceMap(ctx context.Context, productID uuid.UUID) (*domain.IntelProduct, error) {
	return m.getSourceMapFn(ctx, productID)
}
func (m *mockFusionRepo) GetNationalEstimates(ctx context.Context) ([]domain.IntelProduct, error) {
	return m.getNationalEstimatesFn(ctx)
}

type mockKafka struct {
	kafka.Producer
	publishFn func(ctx context.Context, key string, msg interface{}) error
}

func (m *mockKafka) Publish(ctx context.Context, key string, msg interface{}) error {
	if m.publishFn != nil {
		return m.publishFn(ctx, key, msg)
	}
	return nil
}

func TestCreateProduct(t *testing.T) {
	tests := []struct {
		name    string
		repoErr error
		wantErr bool
	}{
		{name: "success", wantErr: false},
		{name: "repo error", repoErr: errors.New("db error"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockFusionRepo{
				createProductFn: func(_ context.Context, _ *domain.IntelProduct) error {
					return tt.repoErr
				},
			}
			svc := NewFusionService(repo, &mockKafka{})
			p := &domain.IntelProduct{Title: "Intel Report"}
			err := svc.CreateProduct(context.Background(), p)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, p.ProductID)
		})
	}
}

func TestGetRecentProducts(t *testing.T) {
	products := []domain.IntelProduct{{Title: "Report 1"}}
	tests := []struct {
		name    string
		result  []domain.IntelProduct
		repoErr error
		wantErr bool
	}{
		{name: "success", result: products, wantErr: false},
		{name: "repo error", repoErr: errors.New("db error"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockFusionRepo{
				getRecentProductsFn: func(_ context.Context, _ int) ([]domain.IntelProduct, error) {
					return tt.result, tt.repoErr
				},
			}
			svc := NewFusionService(repo, &mockKafka{})
			result, err := svc.GetRecentProducts(context.Background())
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestCreateThreatActor(t *testing.T) {
	tests := []struct {
		name    string
		repoErr error
		wantErr bool
	}{
		{name: "success", wantErr: false},
		{name: "repo error", repoErr: errors.New("db error"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockFusionRepo{
				createThreatActorFn: func(_ context.Context, _ *domain.ThreatActor) error {
					return tt.repoErr
				},
			}
			svc := NewFusionService(repo, &mockKafka{})
			a := &domain.ThreatActor{Name: "APT-42"}
			err := svc.CreateThreatActor(context.Background(), a)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, a.ActorID)
		})
	}
}

func TestGetHighRiskActors(t *testing.T) {
	actors := []domain.ThreatActor{{Name: "APT-42"}}
	tests := []struct {
		name    string
		result  []domain.ThreatActor
		repoErr error
		wantErr bool
	}{
		{name: "success", result: actors, wantErr: false},
		{name: "repo error", repoErr: errors.New("db error"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockFusionRepo{
				getHighRiskActorsFn: func(_ context.Context) ([]domain.ThreatActor, error) {
					return tt.result, tt.repoErr
				},
			}
			svc := NewFusionService(repo, &mockKafka{})
			result, err := svc.GetHighRiskActors(context.Background())
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestCreateCorrelation(t *testing.T) {
	tests := []struct {
		name    string
		repoErr error
		wantErr bool
	}{
		{name: "success", wantErr: false},
		{name: "repo error", repoErr: errors.New("db error"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockFusionRepo{
				createCorrelationFn: func(_ context.Context, _ *domain.CrossDisciplineCorrelation) error {
					return tt.repoErr
				},
			}
			svc := NewFusionService(repo, &mockKafka{})
			c := &domain.CrossDisciplineCorrelation{DisciplineA: "SIGINT", DisciplineB: "HUMINT"}
			err := svc.CreateCorrelation(context.Background(), c)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, c.CorrelationID)
		})
	}
}

func TestGetSourceMap(t *testing.T) {
	product := &domain.IntelProduct{Title: "Report 1"}
	tests := []struct {
		name    string
		result  *domain.IntelProduct
		repoErr error
		wantErr bool
	}{
		{name: "found", result: product, wantErr: false},
		{name: "not found", result: nil, wantErr: false},
		{name: "repo error", repoErr: errors.New("db error"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockFusionRepo{
				getSourceMapFn: func(_ context.Context, _ uuid.UUID) (*domain.IntelProduct, error) {
					return tt.result, tt.repoErr
				},
			}
			svc := NewFusionService(repo, &mockKafka{})
			result, err := svc.GetSourceMap(context.Background(), uuid.New())
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestGetNationalEstimates(t *testing.T) {
	estimates := []domain.IntelProduct{{Title: "NIE 2025"}}
	tests := []struct {
		name    string
		result  []domain.IntelProduct
		repoErr error
		wantErr bool
	}{
		{name: "success", result: estimates, wantErr: false},
		{name: "repo error", repoErr: errors.New("db error"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockFusionRepo{
				getNationalEstimatesFn: func(_ context.Context) ([]domain.IntelProduct, error) {
					return tt.result, tt.repoErr
				},
			}
			svc := NewFusionService(repo, &mockKafka{})
			result, err := svc.GetNationalEstimates(context.Background())
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestNewFusionService(t *testing.T) {
	svc := NewFusionService(&mockFusionRepo{}, &mockKafka{})
	require.NotNil(t, svc)
}
