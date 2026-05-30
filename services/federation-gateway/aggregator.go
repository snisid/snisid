package federationgateway

import (
	"fmt"
	"github.com/snisid/platform/backend/internal/platform/logger"
)

type ModelUpdate struct {
	Weights   []float64
	Gradients []float64
	Country   string
	Signature string
}

type Aggregator struct {
	ID string
}

func (a *Aggregator) AggregateWeights(updates []ModelUpdate) []float64 {
	logger.Info(fmt.Sprintf("FEDERATION: Aggregating weights from %d sovereign nodes", len(updates)))
	
	// FedAvg or Differential Privacy aggregation logic
	if len(updates) == 0 {
		return nil
	}

	size := len(updates[0].Weights)
	globalWeights := make([]float64, size)

	for _, u := range updates {
		for i, w := range u.Weights {
			globalWeights[i] += w
		}
	}

	for i := range globalWeights {
		globalWeights[i] /= float64(len(updates))
	}

	return globalWeights
}
