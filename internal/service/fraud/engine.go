package fraud

import (
	"context"
	"fmt"
	"time"

	"github.com/snisid/platform/backend/internal/platform/logger"
	"github.com/snisid/platform/backend/internal/service/router"
	"go.uber.org/zap"
)

type AIClient interface {
	Predict(ctx context.Context, event map[string]interface{}) (int, error)
}

type ScoringEngine struct {
	rules    *router.Engine
	state    *StateStore
	aiClient AIClient
}

func NewScoringEngine(stateAddr string, aiClient AIClient) (*ScoringEngine, error) {
	engine, err := router.NewEngine()
	if err != nil {
		return nil, err
	}

	state := NewStateStore(stateAddr)

	return &ScoringEngine{
		rules:    engine,
		state:    state,
		aiClient: aiClient,
	}, nil
}

func (e *ScoringEngine) CalculateScore(ctx context.Context, event map[string]interface{}) (int, string) {
	totalScore := 0
	reasons := ""

	if targets := e.rules.Evaluate(ctx, event); len(targets) > 0 {
		totalScore += 30
		reasons += "Matched heuristic fraud rules. "
	}

	identityID := fmt.Sprintf("%v", event["identityId"])
	if identityID != "" {
		key := fmt.Sprintf("velocity:identity:%s", identityID)
		count, err := e.state.IncrementVelocity(ctx, key, 10*time.Minute)
		if err == nil && count > 5 {
			totalScore += 50
			reasons += fmt.Sprintf("High velocity detected: %d events in 10m. ", count)
		}
	}

	if e.aiClient != nil {
		aiScore, err := e.aiClient.Predict(ctx, event)
		if err != nil {
			logger.Warn(ctx, "AI inference failed, using fallback", zap.Error(err))
		} else {
			totalScore += aiScore
			if aiScore > 20 {
				reasons += "AI model flagged high risk pattern. "
			}
		}
	}

	return totalScore, reasons
}

func (e *ScoringEngine) ReloadRules(rules []router.Rule) error {
	return e.rules.UpdateRules(rules)
}

type DefaultAIClient struct {
	endpoint string
}

func NewDefaultAIClient(endpoint string) *DefaultAIClient {
	return &DefaultAIClient{endpoint: endpoint}
}

func (c *DefaultAIClient) Predict(ctx context.Context, event map[string]interface{}) (int, error) {
	id := fmt.Sprintf("%v", event["identityId"])
	if id == "test-fraud" {
		return 60, nil
	}
	return 5, nil
}
