package ml

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type ModelRouter struct {
	models   map[string]Model
	producer EventProducer
	metrics  MetricsCollector
	logger   *zap.Logger
}

type Model interface {
	Predict(fv FeatureVector) (float64, error)
}

type EventProducer interface {
	Publish(ctx context.Context, topic string, event interface{}) error
}

type MetricsCollector interface {
	IncrementModelPrediction(modelID string)
	IncrementRoutingError(modelID string)
	ObservePredictionScore(modelID string, score float64)
}

type ModelRoutingEvent struct {
	UserID    string    `json:"user_id"`
	ModelID   string    `json:"model_id"`
	Score     float64   `json:"score"`
	Timestamp time.Time `json:"timestamp"`
}

func NewModelRouter(models map[string]Model, producer EventProducer, metrics MetricsCollector, logger *zap.Logger) *ModelRouter {
	return &ModelRouter{
		models:   models,
		producer: producer,
		metrics:  metrics,
		logger:   logger,
	}
}

func (r *ModelRouter) selectModel(userID string) string {
	hash := sha256.Sum256([]byte(userID))
	idx := int(hash[0]) % len(r.models)
	for k := range r.models {
		if idx == 0 {
			return k
		}
		idx--
	}
	for k := range r.models {
		return k
	}
	return ""
}

func (r *ModelRouter) RouteAndRecord(ctx context.Context, fv FeatureVector) (string, error) {
	modelID := r.selectModel(fv.UserID)
	score, err := r.models[modelID].Predict(fv)
	if err != nil {
		r.metrics.IncrementRoutingError(modelID)
		return modelID, fmt.Errorf("model %s predict: %w", modelID, err)
	}

	event := ModelRoutingEvent{
		UserID:    fv.UserID,
		ModelID:   modelID,
		Score:     score,
		Timestamp: time.Now().UTC(),
	}
	if err := r.producer.Publish(ctx, "snisid.ml.routing", event); err != nil {
		r.logger.Warn("failed to publish routing event", zap.Error(err))
	}

	r.metrics.IncrementModelPrediction(modelID)
	r.metrics.ObservePredictionScore(modelID, score)

	return modelID, nil
}
