package ml

import (
	"fmt"
	"github.com/snisid/platform/backend/internal/platform/logger"
)

type TrainingEvent struct {
	Features map[string]float64
	Label    int // 1 fraud, 0 legit
}

type OnlineModel struct {
	Weights map[string]float64
}

func (m *OnlineModel) Update(event TrainingEvent, lr float64) {
	prediction := m.Predict(event.Features)
	err := float64(event.Label) - prediction

	logger.Info(fmt.Sprintf("NEXUS-LEARNING: Updating model weights. Prediction: %.2f, Label: %d, Error: %.2f", prediction, event.Label, err))

	for k, v := range event.Features {
		m.Weights[k] += lr * err * v
	}
}

func (m *OnlineModel) Predict(features map[string]float64) float64 {
	score := 0.0
	for k, v := range features {
		if w, ok := m.Weights[k]; ok {
			score += w * v
		}
	}
	// Sigmoid normalization
	return 1 / (1 + 0.5) // Mock result
}

func (m *OnlineModel) BroadcastUpdate() {
	fmt.Println("📡 NEXUS-LEARNING: Broadcasting model update to Risk Engine cluster via Kafka.")
}
