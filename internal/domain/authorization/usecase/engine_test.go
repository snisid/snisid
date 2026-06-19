package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/snisid/platform/internal/domain/authorization/entity"
)

type mockPolicyRepo struct {
	policiesFn func(ctx context.Context) ([]entity.Policy, error)
	grantsFn   func(ctx context.Context) ([]entity.RoleGrant, error)
}

func (m *mockPolicyRepo) GetActivePolicies(ctx context.Context) ([]entity.Policy, error) {
	if m.policiesFn != nil {
		return m.policiesFn(ctx)
	}
	return []entity.Policy{}, nil
}

func (m *mockPolicyRepo) GetRoleGrants(ctx context.Context) ([]entity.RoleGrant, error) {
	if m.grantsFn != nil {
		return m.grantsFn(ctx)
	}
	return []entity.RoleGrant{}, nil
}

type mockProducer struct {
	publishFn func(ctx context.Context, key string, event interface{}) error
}

func (m *mockProducer) Publish(ctx context.Context, key string, event interface{}) error {
	if m.publishFn != nil {
		return m.publishFn(ctx, key, event)
	}
	return nil
}

func (m *mockProducer) Close() error { return nil }

func TestNewOPAEngine(t *testing.T) {
	repo := &mockPolicyRepo{}
	engine := NewOPAEngine(repo, nil)
	if engine == nil {
		t.Fatal("NewOPAEngine returned nil")
	}
}

func TestRefreshPolicies_Success(t *testing.T) {
	repo := &mockPolicyRepo{
		policiesFn: func(ctx context.Context) ([]entity.Policy, error) {
			return []entity.Policy{
				{ID: "pol-001", Name: "test_policy", Module: `package snisid.abac

allow {
	input.action == "read"
	input.resource == "public"
}`, Enabled: true},
			}, nil
		},
		grantsFn: func(ctx context.Context) ([]entity.RoleGrant, error) {
			return []entity.RoleGrant{
				{ID: "grant-001", Role: "admin", Action: "read", Resource: "public"},
			}, nil
		},
	}
	engine := NewOPAEngine(repo, nil)
	err := engine.RefreshPolicies(context.Background())
	if err != nil {
		t.Fatalf("RefreshPolicies failed: %v", err)
	}
}

func TestRefreshPolicies_RepoError(t *testing.T) {
	repo := &mockPolicyRepo{
		policiesFn: func(ctx context.Context) ([]entity.Policy, error) {
			return nil, errors.New("db connection failed")
		},
	}
	engine := NewOPAEngine(repo, nil)
	err := engine.RefreshPolicies(context.Background())
	if err == nil {
		t.Fatal("Expected error from repo, got nil")
	}
}

func TestEnforce_EngineNotInitialized(t *testing.T) {
	repo := &mockPolicyRepo{
		policiesFn: func(ctx context.Context) ([]entity.Policy, error) {
			return []entity.Policy{}, nil
		},
		grantsFn: func(ctx context.Context) ([]entity.RoleGrant, error) {
			return []entity.RoleGrant{}, nil
		},
	}
	engine := NewOPAEngine(repo, nil)

	req := &entity.AuthorizationRequest{
		Subject:  entity.SubjectData{UserID: "usr-001"},
		Action:   "read",
		Resource: "public",
	}
	_, err := engine.Enforce(context.Background(), req)
	if err == nil {
		t.Fatal("Expected 'policy engine not initialized' error")
	}
}

func TestLogDecision_NoProducer(t *testing.T) {
	engine := &opaEngine{}
	req := &entity.AuthorizationRequest{
		Subject: entity.SubjectData{UserID: "usr-001"},
		Action:  "read",
	}
	dec := &entity.AuthorizationDecision{Allowed: true}
	// Should not panic
	engine.logDecision(req, dec)
}

func TestLogDecision_WithProducer(t *testing.T) {
	published := false
	prod := &mockProducer{
		publishFn: func(ctx context.Context, key string, event interface{}) error {
			published = true
			return nil
		},
	}

	engine := &opaEngine{producer: prod}
	req := &entity.AuthorizationRequest{
		Subject: entity.SubjectData{UserID: "usr-001"},
		Action:  "delete",
	}
	dec := &entity.AuthorizationDecision{Allowed: false, Reason: "Denied"}
	engine.logDecision(req, dec)

	if !published {
		t.Error("Expected event to be published")
	}
}
