package ml

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type FeatureVector struct {
	UserID    string
	Amount    float64
	Velocity  float64
	GraphRisk float64
}

type FeatureExtractor struct {
	store  FeatureStore
	logger *zap.Logger
}

func NewFeatureExtractor(store FeatureStore, logger *zap.Logger) *FeatureExtractor {
	return &FeatureExtractor{store: store, logger: logger}
}

func (fe *FeatureExtractor) ExtractFeatures(ctx context.Context, payload map[string]interface{}) (*FeatureVector, error) {
	userID, ok := payload["user_id"].(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("%w: user_id is required", ErrInvalidPayload)
	}

	amount, _ := payload["amount"].(float64)

	velocity, err := fe.store.GetVelocity(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("feature extraction failed for user %s: %w", userID, err)
	}

	graphRisk, err := fe.store.GetGraphRisk(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("feature extraction failed for user %s: %w", userID, err)
	}

	if velocity == 0.0 && graphRisk == 0.0 {
		fe.logger.Debug("feature extraction: no cached features, using zero values (degraded mode)",
			zap.String("user_id", userID))
	}

	return &FeatureVector{
		UserID:    userID,
		Amount:    amount,
		Velocity:  velocity,
		GraphRisk: graphRisk,
	}, nil
}
