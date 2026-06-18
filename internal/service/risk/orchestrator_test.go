package risk

import (
	"context"
	"errors"
	"testing"
)

type testRiskModel struct {
	name       string
	evaluateFn func(ctx context.Context, data map[string]interface{}) (RiskResult, error)
}

func (m *testRiskModel) Name() string { return m.name }
func (m *testRiskModel) Evaluate(ctx context.Context, data map[string]interface{}) (RiskResult, error) {
	if m.evaluateFn != nil {
		return m.evaluateFn(ctx, data)
	}
	return RiskResult{Score: 0, Reason: "noop"}, nil
}

func TestNewOrchestrator(t *testing.T) {
	o := NewOrchestrator()
	if o == nil {
		t.Fatal("NewOrchestrator returned nil")
	}
	if len(o.models) != 0 {
		t.Errorf("Models count = %d, want 0", len(o.models))
	}
}

func TestAssess_SingleModel(t *testing.T) {
	m := &testRiskModel{
		name: "simple_check",
		evaluateFn: func(ctx context.Context, data map[string]interface{}) (RiskResult, error) {
			return RiskResult{Score: 50, Reason: "Medium risk"}, nil
		},
	}
	o := NewOrchestrator(WeightedModel{Model: m, Weight: 1.0})

	score, level, breakdown, explanations, err := o.Assess(context.Background(), map[string]interface{}{"identityId": "ID-001"})
	if err != nil {
		t.Fatalf("Assess failed: %v", err)
	}
	if score != 50 {
		t.Errorf("Score = %d, want 50", score)
	}
	if level != "MEDIUM" {
		t.Errorf("Level = %s, want MEDIUM", level)
	}
	if len(breakdown) != 1 {
		t.Errorf("Breakdown count = %d, want 1", len(breakdown))
	}
	if len(explanations) != 1 {
		t.Errorf("Explanations count = %d, want 1", len(explanations))
	}
}

func TestAssess_MultiModelWeighted(t *testing.T) {
	m1 := &testRiskModel{
		name: "sanctions",
		evaluateFn: func(ctx context.Context, data map[string]interface{}) (RiskResult, error) {
			return RiskResult{Score: 100, Reason: "Sanctions hit"}, nil
		},
	}
	m2 := &testRiskModel{
		name: "velocity",
		evaluateFn: func(ctx context.Context, data map[string]interface{}) (RiskResult, error) {
			return RiskResult{Score: 0, Reason: "Normal"}, nil
		},
	}
	o := NewOrchestrator(
		WeightedModel{Model: m1, Weight: 0.7},
		WeightedModel{Model: m2, Weight: 0.3},
	)

	score, level, _, _ := o.Assess(context.Background(), map[string]interface{}{"identityId": "ID-002"})
	expected := int(100*0.7 + 0*0.3) // 70
	if score != expected {
		t.Errorf("Score = %d, want %d", score, expected)
	}
	if level != "HIGH" {
		t.Errorf("Level = %s, want HIGH (score=%d)", level, score)
	}
}

func TestAssess_ModelError(t *testing.T) {
	m := &testRiskModel{
		name: "failing",
		evaluateFn: func(ctx context.Context, data map[string]interface{}) (RiskResult, error) {
			return RiskResult{}, errors.New("model failure")
		},
	}
	o := NewOrchestrator(WeightedModel{Model: m, Weight: 1.0})

	// Error model should be skipped
	score, level, _, _, err := o.Assess(context.Background(), map[string]interface{}{"identityId": "ID-003"})
	if err != nil {
		t.Fatalf("Assess failed: %v", err)
	}
	if score != 0 {
		t.Errorf("Score = %d, want 0", score)
	}
	if level != "LOW" {
		t.Errorf("Level = %s, want LOW", level)
	}
}

func TestAssess_LevelBoundaries(t *testing.T) {
	tests := []struct {
		score int
		level string
	}{
		{0, "LOW"},
		{10, "LOW"},
		{20, "MEDIUM"},
		{21, "MEDIUM"},
		{50, "HIGH"},
		{51, "HIGH"},
		{80, "CRITICAL"},
		{81, "CRITICAL"},
		{100, "CRITICAL"},
	}

	for _, tt := range tests {
		m := &testRiskModel{
			name: "scorer",
			evaluateFn: func(ctx context.Context, data map[string]interface{}) (RiskResult, error) {
				return RiskResult{Score: tt.score, Reason: "test"}, nil
			},
		}
		o := NewOrchestrator(WeightedModel{Model: m, Weight: 1.0})
		_, level, _, _, _ := o.Assess(context.Background(), map[string]interface{}{"identityId": "test"})
		if level != tt.level {
			t.Errorf("Score %d: level = %s, want %s", tt.score, level, tt.level)
		}
	}
}
