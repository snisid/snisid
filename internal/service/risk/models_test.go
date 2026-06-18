package risk

import (
	"context"
	"testing"
)

func TestNewSanctionsModel(t *testing.T) {
	m := NewSanctionsModel()
	if m == nil {
		t.Fatal("NewSanctionsModel returned nil")
	}
	if m.Name() != "sanctions_check" {
		t.Errorf("Name = %s, want sanctions_check", m.Name())
	}
}

func TestSanctionsModel_BlockedName(t *testing.T) {
	m := NewSanctionsModel()
	result, err := m.Evaluate(context.Background(), map[string]interface{}{
		"fullName": "TEMPORARY_BLOCKED_CITIZEN",
	})
	if err != nil {
		t.Fatalf("Evaluate failed: %v", err)
	}
	if result.Score != 100 {
		t.Errorf("Score = %d, want 100", result.Score)
	}
}

func TestSanctionsModel_NormalName(t *testing.T) {
	m := NewSanctionsModel()
	result, err := m.Evaluate(context.Background(), map[string]interface{}{
		"fullName": "Jean Dupont",
	})
	if err != nil {
		t.Fatalf("Evaluate failed: %v", err)
	}
	if result.Score != 0 {
		t.Errorf("Score = %d, want 0", result.Score)
	}
}

func TestSanctionsModel_SimilarName(t *testing.T) {
	m := NewSanctionsModel()
	result, err := m.Evaluate(context.Background(), map[string]interface{}{
		"fullName": "TEMPORARY_BLOCKED",
	})
	if err != nil {
		t.Fatalf("Evaluate failed: %v", err)
	}
	// Should have partial similarity but not 100
	t.Logf("Similar name score: %d, reason: %s", result.Score, result.Reason)
}

func TestTravelModel_Normal(t *testing.T) {
	m := &TravelModel{}
	result, err := m.Evaluate(context.Background(), map[string]interface{}{
		"identityId": "ID-001",
		"location":   "Port-au-Prince",
	})
	if err != nil {
		t.Fatalf("Evaluate failed: %v", err)
	}
	if result.Score != 0 {
		t.Errorf("Score = %d, want 0", result.Score)
	}
}

func TestTravelModel_Suspicious(t *testing.T) {
	m := &TravelModel{}
	result, err := m.Evaluate(context.Background(), map[string]interface{}{
		"identityId": "ID-002",
		"location":   "SUSPICIOUS_REMOTE_LOCATION",
	})
	if err != nil {
		t.Fatalf("Evaluate failed: %v", err)
	}
	if result.Score != 70 {
		t.Errorf("Score = %d, want 70", result.Score)
	}
}

func TestTravelModel_Name(t *testing.T) {
	m := &TravelModel{}
	if m.Name() != "travel_velocity" {
		t.Errorf("Name = %s, want travel_velocity", m.Name())
	}
}

func TestRiskModelInterface(t *testing.T) {
	var s RiskModel = NewSanctionsModel()
	var t2 RiskModel = &TravelModel{}

	if s.Name() == t2.Name() {
		t.Error("Models should have different names")
	}
}

func TestRiskResult_Values(t *testing.T) {
	r := RiskResult{Score: 95, Reason: "Critical risk"}
	if r.Score != 95 {
		t.Errorf("Score = %d", r.Score)
	}
}
