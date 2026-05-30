package fraud

import (
	"context"
	"fmt"

	"github.com/snisid/platform/backend/internal/service/router"
)

type ModelResult struct {
	Score  int
	Reason string
}

type Model interface {
	Name() string
	Score(ctx context.Context, event map[string]interface{}) (ModelResult, error)
}

type HeuristicModel struct {
	engine *router.Engine
}

func NewHeuristicModel() (*HeuristicModel, error) {
	engine, err := router.NewEngine()
	if err != nil {
		return nil, err
	}
	return &HeuristicModel{engine: engine}, nil
}

func (m *HeuristicModel) Name() string { return "heuristic_rules" }
func (m *HeuristicModel) Score(ctx context.Context, event map[string]interface{}) (ModelResult, error) {
	if targets := m.engine.Evaluate(ctx, event); len(targets) > 0 {
		return ModelResult{Score: 80, Reason: fmt.Sprintf("Matched rules: %v", targets)}, nil
	}
	return ModelResult{Score: 0, Reason: "No rules matched"}, nil
}

type MLModel struct{}

func (m *MLModel) Name() string { return "ai_inference" }
func (m *MLModel) Score(ctx context.Context, event map[string]interface{}) (ModelResult, error) {
	// Placeholder for gRPC call
	id, _ := event["identityId"].(string)
	if id == "test-ai-fraud" {
		return ModelResult{Score: 95, Reason: "AI model detected anomalous behavioral pattern"}, nil
	}
	return ModelResult{Score: 5, Reason: "Low risk profile"}, nil
}
