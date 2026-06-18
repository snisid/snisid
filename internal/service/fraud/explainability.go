package fraud

import (
	"context"
	"fmt"
	"sync"
)

type WeightedModel struct {
	Model  Model
	Weight float64
}

type Orchestrator struct {
	models []WeightedModel
}

func NewOrchestrator(models ...WeightedModel) *Orchestrator {
	return &Orchestrator{models: models}
}

func (o *Orchestrator) Calculate(ctx context.Context, event map[string]interface{}) (int, map[string]int, []string, error) {
	totalWeightedScore := 0.0
	modelScores := make(map[string]int)
	explanations := make([]string, 0)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, m := range o.models {
		wg.Add(1)
		go func(wm WeightedModel) {
			defer wg.Done()
			res, err := wm.Model.Score(ctx, event)
			if err != nil {
				return
			}

			mu.Lock()
			modelScores[wm.Model.Name()] = res.Score
			totalWeightedScore += float64(res.Score) * wm.Weight
			explanations = append(explanations, fmt.Sprintf("[%s] Score %d: %s", wm.Model.Name(), res.Score, res.Reason))
			mu.Unlock()
		}(m)
	}

	wg.Wait()

	return int(totalWeightedScore), modelScores, explanations, nil
}
