package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/snisid/cyber-ht/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCyberRepo struct {
	createIncidentFn       func(ctx context.Context, inc *domain.Incident) error
	getActiveIncidentsFn   func(ctx context.Context) ([]domain.Incident, error)
	createPolicyFn         func(ctx context.Context, p *domain.ZeroTrustPolicy) error
	checkThreatIndicatorFn func(ctx context.Context, indicator string) (*domain.ThreatIndicator, error)
}

func (m *mockCyberRepo) CreateIncident(ctx context.Context, inc *domain.Incident) error {
	return m.createIncidentFn(ctx, inc)
}
func (m *mockCyberRepo) GetActiveIncidents(ctx context.Context) ([]domain.Incident, error) {
	return m.getActiveIncidentsFn(ctx)
}
func (m *mockCyberRepo) CreatePolicy(ctx context.Context, p *domain.ZeroTrustPolicy) error {
	return m.createPolicyFn(ctx, p)
}
func (m *mockCyberRepo) CheckThreatIndicator(ctx context.Context, indicator string) (*domain.ThreatIndicator, error) {
	return m.checkThreatIndicatorFn(ctx, indicator)
}

// producer is nil in tests — service handles nil producer gracefully

func TestCreateIncident(t *testing.T) {
	tests := []struct {
		name    string
		req     domain.CreateIncidentRequest
		repoErr error
		wantErr bool
	}{
		{
			name: "success",
			req: domain.CreateIncidentRequest{
				Title:      "Data breach",
				Severity:   "HIGH",
				DetectedBy: "SOC-01",
			},
			wantErr: false,
		},
		{
			name: "repo error",
			req: domain.CreateIncidentRequest{
				Title:      "Test",
				Severity:   "LOW",
				DetectedBy: "auto",
			},
			repoErr: errors.New("db error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockCyberRepo{
				createIncidentFn: func(ctx context.Context, inc *domain.Incident) error {
					return tt.repoErr
				},
			}
			svc := NewCyberService(repo, nil)
			inc, err := svc.CreateIncident(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, inc.ID)
			assert.Equal(t, tt.req.Title, inc.Title)
			assert.Equal(t, domain.Severity(tt.req.Severity), inc.Severity)
			assert.Equal(t, domain.IncDetected, inc.Status)
		})
	}
}

func TestGetActiveIncidents(t *testing.T) {
	tests := []struct {
		name    string
		repoRes []domain.Incident
		repoErr error
		wantErr bool
		count   int
	}{
		{
			name:    "success with incidents",
			repoRes: []domain.Incident{{ID: uuid.New(), Title: "test"}},
			count:   1,
		},
		{
			name:    "empty list",
			repoRes: []domain.Incident{},
			count:   0,
		},
		{
			name:    "repo error",
			repoErr: errors.New("query error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockCyberRepo{
				getActiveIncidentsFn: func(ctx context.Context) ([]domain.Incident, error) {
					return tt.repoRes, tt.repoErr
				},
			}
			svc := NewCyberService(repo, nil)
			got, err := svc.GetActiveIncidents(context.Background())
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Len(t, got, tt.count)
		})
	}
}

func TestCreatePolicy(t *testing.T) {
	tests := []struct {
		name    string
		req     domain.CreatePolicyRequest
		repoErr error
		wantErr bool
	}{
		{
			name: "success",
			req: domain.CreatePolicyRequest{
				Name:        "MFA Required",
				Description: "Enforce MFA for all users",
				PolicyType:  "ACCESS_CONTROL",
				Rules:       []string{"mfa=true"},
				Enabled:     true,
				CreatedBy:   "admin",
			},
		},
		{
			name: "repo error",
			req: domain.CreatePolicyRequest{
				Name:        "test",
				Description: "test",
				PolicyType:  "test",
				Rules:       []string{"a"},
				CreatedBy:   "admin",
			},
			repoErr: errors.New("insert error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockCyberRepo{
				createPolicyFn: func(ctx context.Context, p *domain.ZeroTrustPolicy) error {
					return tt.repoErr
				},
			}
			svc := NewCyberService(repo, nil)
			p, err := svc.CreatePolicy(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, p.ID)
			assert.Equal(t, tt.req.Name, p.Name)
			assert.True(t, p.Enabled)
		})
	}
}

func TestCheckThreatIndicator(t *testing.T) {
	knownIndicator := &domain.ThreatIndicator{
		ID:          uuid.New(),
		Indicator:   "1.2.3.4",
		Type:        "IP",
		ThreatLevel: "HIGH",
		Source:      "alienvault",
	}
	tests := []struct {
		name      string
		indicator string
		repoRes   *domain.ThreatIndicator
		repoErr   error
		wantErr   bool
	}{
		{
			name:      "found",
			indicator: "1.2.3.4",
			repoRes:   knownIndicator,
		},
		{
			name:      "not found",
			indicator: "unknown",
			repoErr:   errors.New("threat indicator not found"),
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockCyberRepo{
				checkThreatIndicatorFn: func(ctx context.Context, indicator string) (*domain.ThreatIndicator, error) {
					return tt.repoRes, tt.repoErr
				},
			}
			svc := NewCyberService(repo, nil)
			got, err := svc.CheckThreatIndicator(context.Background(), tt.indicator)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.repoRes.Indicator, got.Indicator)
		})
	}
}
