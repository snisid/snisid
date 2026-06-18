package ml

import "context"

type FeatureStore interface {
	GetVelocity(ctx context.Context, userID string) (float64, error)
	GetGraphRisk(ctx context.Context, userID string) (float64, error)
}
