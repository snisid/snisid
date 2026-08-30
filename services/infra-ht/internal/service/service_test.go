package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/snisid/infra-ht/internal/domain"
	"github.com/snisid/infra-ht/internal/repository"
)

type mockRepo struct {
	getDatacentersFunc func(ctx context.Context) ([]domain.Datacenter, error)
	getClustersFunc    func(ctx context.Context) ([]domain.K8sCluster, error)
	createDRDrillFunc  func(ctx context.Context, d *domain.DRDrill) error
}

func (m *mockRepo) GetDatacenters(ctx context.Context) ([]domain.Datacenter, error) {
	return m.getDatacentersFunc(ctx)
}
func (m *mockRepo) GetClusters(ctx context.Context) ([]domain.K8sCluster, error) {
	return m.getClustersFunc(ctx)
}
func (m *mockRepo) CreateDRDrill(ctx context.Context, d *domain.DRDrill) error {
	return m.createDRDrillFunc(ctx, d)
}

func newTestService(repo repository.Repository) *InfraService {
	return NewInfraService(repo, nil)
}

func TestNewInfraService(t *testing.T) {
	svc := NewInfraService(nil, nil)
	require.NotNil(t, svc)
}

func TestGetHealth(t *testing.T) {
	tests := []struct {
		name     string
		repo     *mockRepo
		expected map[string]any
	}{
		{
			name: "healthy with data",
			repo: &mockRepo{
				getDatacentersFunc: func(ctx context.Context) ([]domain.Datacenter, error) {
					return []domain.Datacenter{{DCID: uuid.New(), DCName: "DC1"}}, nil
				},
				getClustersFunc: func(ctx context.Context) ([]domain.K8sCluster, error) {
					return []domain.K8sCluster{{ClusterID: uuid.New(), ClusterName: "prod-cluster"}}, nil
				},
			},
			expected: map[string]any{"datacenters": 1, "clusters": 1, "status": "healthy"},
		},
		{
			name: "healthy empty",
			repo: &mockRepo{
				getDatacentersFunc: func(ctx context.Context) ([]domain.Datacenter, error) {
					return []domain.Datacenter{}, nil
				},
				getClustersFunc: func(ctx context.Context) ([]domain.K8sCluster, error) {
					return []domain.K8sCluster{}, nil
				},
			},
			expected: map[string]any{"datacenters": 0, "clusters": 0, "status": "healthy"},
		},
		{
			name: "repo errors ignored",
			repo: &mockRepo{
				getDatacentersFunc: func(ctx context.Context) ([]domain.Datacenter, error) {
					return nil, errors.New("db error")
				},
				getClustersFunc: func(ctx context.Context) ([]domain.K8sCluster, error) {
					return nil, errors.New("db error")
				},
			},
			expected: map[string]any{"datacenters": 0, "clusters": 0, "status": "healthy"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.repo)
			result := svc.GetHealth(context.Background())
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetClusters(t *testing.T) {
	clusters := []domain.K8sCluster{
		{ClusterID: uuid.New(), ClusterName: "prod", NodeCount: 10},
	}

	tests := []struct {
		name    string
		repo    *mockRepo
		want    []domain.K8sCluster
		wantErr bool
	}{
		{
			name: "success",
			repo: &mockRepo{
				getClustersFunc: func(ctx context.Context) ([]domain.K8sCluster, error) { return clusters, nil },
			},
			want: clusters,
		},
		{
			name: "repo error",
			repo: &mockRepo{
				getClustersFunc: func(ctx context.Context) ([]domain.K8sCluster, error) {
					return nil, errors.New("db error")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.repo)
			got, err := svc.GetClusters(context.Background())
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRecordDRDrill(t *testing.T) {
	drillDate := time.Now()

	tests := []struct {
		name    string
		drill   domain.DRDrill
		repo    *mockRepo
		wantErr bool
	}{
		{
			name: "success",
			drill: domain.DRDrill{
				DrillDate:    drillDate,
				Scenario:     "DC failover",
				RTOTargetMin: 30,
				RPOTargetMin: 5,
			},
			repo: &mockRepo{
				createDRDrillFunc: func(ctx context.Context, d *domain.DRDrill) error { return nil },
			},
		},
		{
			name: "success with defaults",
			drill: domain.DRDrill{
				DrillDate: drillDate,
				Scenario:  "full region fail",
			},
			repo: &mockRepo{
				createDRDrillFunc: func(ctx context.Context, d *domain.DRDrill) error { return nil },
			},
		},
		{
			name: "repo error",
			drill: domain.DRDrill{
				DrillDate: drillDate,
				Scenario:  "fail",
			},
			repo: &mockRepo{
				createDRDrillFunc: func(ctx context.Context, d *domain.DRDrill) error {
					return errors.New("db error")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.repo)
			result, err := svc.RecordDRDrill(context.Background(), tt.drill)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.NotEmpty(t, result.DrillID)
			assert.False(t, result.CreatedAt.IsZero())
			if tt.drill.RTOTargetMin == 0 {
				assert.Equal(t, 15, result.RTOTargetMin)
			} else {
				assert.Equal(t, tt.drill.RTOTargetMin, result.RTOTargetMin)
			}
			if tt.drill.RPOTargetMin == 0 {
				assert.Equal(t, 1, result.RPOTargetMin)
			} else {
				assert.Equal(t, tt.drill.RPOTargetMin, result.RPOTargetMin)
			}
		})
	}
}
