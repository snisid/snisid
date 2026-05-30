package risk

import (
	"context"
	"fmt"
	"sync"

	"github.com/snisid/platform/backend/internal/platform/logger"
	"go.uber.org/zap"
)

type WeightedModel struct {
	Model  RiskModel
	Weight float64
}

type Orchestrator struct {
	models []WeightedModel
}

func NewOrchestrator(models ...WeightedModel) *Orchestrator {
	return &Orchestrator{models: models}
}

func (o *Orchestrator) Assess(ctx context.Context, identity map[string]interface{}) (int, string, map[string]int, []string, error) {
	totalWeightedScore := 0.0
	modelBreakdown := make(map[string]int)
	explanations := make([]string, 0)

	var mu sync.Mutex
	var wg sync.WaitGroup

	logger.Info(ctx, "Starting multi-factor risk assessment", zap.String("identity_id", fmt.Sprintf("%v", identity["identityId"])))

	for _, m := range o.models {
		wg.Add(1)
		go func(wm WeightedModel) {
			defer wg.Done()
			res, err := wm.Model.Evaluate(ctx, identity)
			if err != nil {
				return
			}

			mu.Lock()
			modelBreakdown[wm.Model.Name()] = res.Score
			totalWeightedScore += float64(res.Score) * wm.Weight
			explanations = append(explanations, fmt.Sprintf("[%s] Score %d: %s", wm.Model.Name(), res.Score, res.Reason))
			mu.Unlock()
		}(m)
	}

	wg.Wait()

	score := int(totalWeightedScore)
	level := "LOW"
	if score > 80 {
		level = "CRITICAL"
	} else if score > 50 {
		level = "HIGH"
	} else if score > 20 {
		level = "MEDIUM"
	}

	return score, level, modelBreakdown, explanations, nil
}
