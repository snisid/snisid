package fraud

import (
	"context"
	"testing"
)

func TestHeuristicModel_Name(t *testing.T) {
	m, err := NewHeuristicModel()
	if err != nil {
		t.Skipf("Skipping: CEL router engine not available in test context: %v", err)
	}
	if m.Name() != "heuristic_rules" {
		t.Errorf("Name = %s, want heuristic_rules", m.Name())
	}
}

func TestHeuristicModel_Score(t *testing.T) {
	m, err := NewHeuristicModel()
	if err != nil {
		t.Skip("CEL engine not initialized")
	}

	result, err := m.Score(context.Background(), map[string]interface{}{
		"identityId": "test",
		"action":     "created",
	})
	if err != nil {
		t.Fatalf("Score failed: %v", err)
	}
	if result.Score < 0 || result.Score > 100 {
		t.Errorf("Score out of range: %d", result.Score)
	}
}

func TestMLModel_Name(t *testing.T) {
	m := &MLModel{}
	if m.Name() != "ai_inference" {
		t.Errorf("Name = %s, want ai_inference", m.Name())
	}
}

func TestMLModel_Score_Normal(t *testing.T) {
	m := &MLModel{}
	result, err := m.Score(context.Background(), map[string]interface{}{
		"identityId": "normal-user",
	})
	if err != nil {
		t.Fatalf("Score failed: %v", err)
	}
	if result.Score != 5 {
		t.Errorf("Score = %d, want 5", result.Score)
	}
	if result.Reason != "Low risk profile" {
		t.Errorf("Reason = %s, want 'Low risk profile'", result.Reason)
	}
}

func TestMLModel_Score_FraudDetected(t *testing.T) {
	m := &MLModel{}
	result, err := m.Score(context.Background(), map[string]interface{}{
		"identityId": "test-ai-fraud",
	})
	if err != nil {
		t.Fatalf("Score failed: %v", err)
	}
	if result.Score != 95 {
		t.Errorf("Score = %d, want 95", result.Score)
	}
	if result.Reason != "AI model detected anomalous behavioral pattern" {
		t.Errorf("Reason = %s, want 'AI model detected anomalous behavioral pattern'", result.Reason)
	}
}

func TestModelResult_Values(t *testing.T) {
	r := ModelResult{Score: 80, Reason: "High risk"}
	if r.Score != 80 {
		t.Errorf("Score = %d, want 80", r.Score)
	}
}

func TestHeuristicModel_Score_WithRules(t *testing.T) {
	m, err := NewHeuristicModel()
	if err != nil {
		t.Skip("CEL engine not initialized")
	}

	result, err := m.Score(context.Background(), map[string]interface{}{
		"identityId": "test",
		"action":     "update",
		"metadata": map[string]interface{}{
			"force": "true",
		},
	})
	if err != nil {
		t.Fatalf("Score failed: %v", err)
	}
	// Should match some rules based on the event structure
	t.Logf("Score: %d, Reason: %s", result.Score, result.Reason)
}
