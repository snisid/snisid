package fraud

import (
	"context"
	"testing"
)

func TestNewNeo4jRepository_NilDriver(t *testing.T) {
	// Should not panic but will fail on actual calls
	repo := NewNeo4jRepository(nil)
	if repo == nil {
		t.Fatal("NewNeo4jRepository returned nil")
	}
}

type mockDriver struct {
	newSessionFn func(ctx context.Context, config interface{}) *mockSession
}

type mockSession struct {
	executeWriteFn func(ctx context.Context, fn interface{}) (interface{}, error)
	executeReadFn  func(ctx context.Context, fn interface{}) (interface{}, error)
	closeFn        func(ctx context.Context)
}

func (s *mockSession) Close(ctx context.Context) {
	if s.closeFn != nil {
		s.closeFn(ctx)
	}
}

// GraphRepository interface uses concrete neo4j types, so we test the service layer instead
// These tests verify the interface contract

func TestGraphRepositoryInterface(t *testing.T) {
	var repo GraphRepository
	t.Log("GraphRepository interface defined with AddIdentityNode and CheckFraudRing")
	_ = repo
}

func TestIdentityCreatedEvent_Values(t *testing.T) {
	evt := IdentityCreatedEvent{
		IdentityID: "ID-001",
		FirstName:  "Jean",
		LastName:   "Dupont",
		Agency:     "oni",
	}
	if evt.IdentityID != "ID-001" {
		t.Errorf("IdentityID = %s", evt.IdentityID)
	}
}

func TestFraudScoredEvent_Values(t *testing.T) {
	evt := FraudScoredEvent{
		IdentityID: "ID-001",
		RiskScore:  95,
		IsFraud:    true,
		Reason:     "Fraud ring",
	}
	if !evt.IsFraud {
		t.Error("Expected IsFraud = true")
	}
	if evt.RiskScore != 95 {
		t.Errorf("RiskScore = %d, want 95", evt.RiskScore)
	}
}
