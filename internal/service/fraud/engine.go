package fraud

import (
	"context"
	"fmt"
	"time"

	"github.com/snisid/platform/backend/internal/platform/logger"
	"github.com/snisid/platform/backend/internal/service/router"
	"go.uber.org/zap"
)

type ScoringEngine struct {
	rules *router.Engine
	state *StateStore
}

func NewScoringEngine(stateAddr string) (*ScoringEngine, error) {
	engine, err := router.NewEngine()
	if err != nil {
		return nil, err
	}

	state := NewStateStore(stateAddr)

	return &ScoringEngine{
		rules: engine,
		state: state,
	}, nil
}

func (e *ScoringEngine) CalculateScore(ctx context.Context, event map[string]interface{}) (int, string) {
	totalScore := 0
	reasons := ""

	// 1. Heuristic Rules (CEL)
	// We check against pre-defined fraud rules
	if targets := e.rules.Evaluate(ctx, event); len(targets) > 0 {
		totalScore += 30
		reasons += "Matched heuristic fraud rules. "
	}

	// 2. Velocity Checks (Redis)
	identityID := fmt.Sprintf("%v", event["identityId"])
	if identityID != "" {
		key := fmt.Sprintf("velocity:identity:%s", identityID)
		count, err := e.state.IncrementVelocity(ctx, key, 10*time.Minute)
		if err == nil && count > 5 {
			totalScore += 50
			reasons += fmt.Sprintf("High velocity detected: %d events in 10m. ", count)
		}
	}

	// 3. AI Inference Placeholder
	// In a real system, this would call a gRPC model service
	aiScore := e.mockAIInference(event)
	totalScore += aiScore
	if aiScore > 20 {
		reasons += "AI model flagged high risk pattern. "
	}

	return totalScore, reasons
}

func (e *ScoringEngine) mockAIInference(event map[string]interface{}) int {
	// Dummy logic: flag if identityId contains 'test-fraud'
	id := fmt.Sprintf("%v", event["identityId"])
	if id == "test-fraud" {
		return 60
	}
	return 5
}

func (e *ScoringEngine) ReloadRules(rules []router.Rule) error {
	return e.rules.UpdateRules(rules)
}
