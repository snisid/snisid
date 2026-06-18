package fraud

import (
	"context"
	"errors"
	"testing"
)

type testModel struct {
	name    string
	scoreFn func(ctx context.Context, event map[string]interface{}) (ModelResult, error)
}

func (m *testModel) Name() string { return m.name }
func (m *testModel) Score(ctx context.Context, event map[string]interface{}) (ModelResult, error) {
	if m.scoreFn != nil {
		return m.scoreFn(ctx, event)
	}
	return ModelResult{Score: 0, Reason: "noop"}, nil
}

func TestNewOrchestrator_Empty(t *testing.T) {
	o := NewOrchestrator()
	if o == nil {
		t.Fatal("NewOrchestrator returned nil")
	}
	if len(o.models) != 0 {
		t.Errorf("models count = %d, want 0", len(o.models))
	}
}

func TestNewOrchestrator_WithModels(t *testing.T) {
	m1 := &testModel{name: "model_a"}
	m2 := &testModel{name: "model_b"}
	o := NewOrchestrator(
		WeightedModel{Model: m1, Weight: 0.6},
		WeightedModel{Model: m2, Weight: 0.4},
	)
	if len(o.models) != 2 {
		t.Errorf("models count = %d, want 2", len(o.models))
	}
}

func TestCalculate_SingleModel(t *testing.T) {
	m := &testModel{
		name: "heuristic",
		scoreFn: func(ctx context.Context, event map[string]interface{}) (ModelResult, error) {
			return ModelResult{Score: 80, Reason: "Matched rule A"}, nil
		},
	}
	o := NewOrchestrator(WeightedModel{Model: m, Weight: 1.0})

	score, scores, explanations, err := o.Calculate(context.Background(), map[string]interface{}{"id": "1"})
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}
	if score != 80 {
		t.Errorf("total score = %d, want 80", score)
	}
	if len(scores) != 1 {
		t.Errorf("model scores count = %d, want 1", len(scores))
	}
	if len(explanations) != 1 {
		t.Errorf("explanations count = %d, want 1", len(explanations))
	}
}

func TestCalculate_MultiModel(t *testing.T) {
	m1 := &testModel{
		name: "model_a",
		scoreFn: func(ctx context.Context, event map[string]interface{}) (ModelResult, error) {
			return ModelResult{Score: 100, Reason: "Critical"}, nil
		},
	}
	m2 := &testModel{
		name: "model_b",
		scoreFn: func(ctx context.Context, event map[string]interface{}) (ModelResult, error) {
			return ModelResult{Score: 50, Reason: "Medium"}, nil
		},
	}
	o := NewOrchestrator(
		WeightedModel{Model: m1, Weight: 0.5},
		WeightedModel{Model: m2, Weight: 0.5},
	)

	score, scores, _, err := o.Calculate(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}
	expected := int(100*0.5 + 50*0.5) // 75
	if score != expected {
		t.Errorf("total score = %d, want %d", score, expected)
	}
	if len(scores) != 2 {
		t.Errorf("model scores count = %d, want 2", len(scores))
	}
}

func TestCalculate_ModelError(t *testing.T) {
	m := &testModel{
		name: "failing",
		scoreFn: func(ctx context.Context, event map[string]interface{}) (ModelResult, error) {
			return ModelResult{}, errors.New("model failure")
		},
	}
	o := NewOrchestrator(WeightedModel{Model: m, Weight: 1.0})

	// Error model should be skipped, not cause total failure
	score, _, explanations, err := o.Calculate(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}
	if score != 0 {
		t.Errorf("score = %d, want 0 (failed model)", score)
	}
	if len(explanations) != 0 {
		t.Errorf("explanations count = %d, want 0", len(explanations))
	}
}

func TestCalculate_WeightEffects(t *testing.T) {
	m := &testModel{
		name: "weighted",
		scoreFn: func(ctx context.Context, event map[string]interface{}) (ModelResult, error) {
			return ModelResult{Score: 100, Reason: "max"}, nil
		},
	}

	// Weight = 0.3 should give score = 30
	o := NewOrchestrator(WeightedModel{Model: m, Weight: 0.3})
	score, _, _, err := o.Calculate(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}
	if score != 30 {
		t.Errorf("score with 0.3 weight = %d, want 30", score)
	}
}
