package ml

import (
	"context"
	"errors"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

type testRedisStore struct {
	client *redis.Client
}

func (s *testRedisStore) GetFloat(ctx context.Context, key string) (float64, error) {
	val, err := s.client.Get(ctx, key).Float64()
	if errors.Is(err, redis.Nil) {
		return 0.0, nil
	}
	if err != nil {
		return 0.0, err
	}
	return val, nil
}

func TestOnlineModel_Predict_InitialState(t *testing.T) {
	model := NewOnlineModel(0.01, 0.001)

	fv := FeatureVector{
		UserID:    "user1",
		Amount:    100.0,
		Velocity:  0.5,
		GraphRisk: 0.3,
	}

	score := model.Predict(fv)
	assert.Greater(t, score, 0.0)
	assert.Less(t, score, 1.0)
}

func TestOnlineModel_Update_Convergence(t *testing.T) {
	model := NewOnlineModel(0.1, 0.001)

	for i := 0; i < 100; i++ {
		fv := FeatureVector{
			UserID:    "user1",
			Amount:    100.0,
			Velocity:  0.8,
			GraphRisk: 0.6,
		}
		model.Update(fv, 1.0)
	}

	fv := FeatureVector{
		UserID:    "user1",
		Amount:    100.0,
		Velocity:  0.8,
		GraphRisk: 0.6,
	}
	score := model.Predict(fv)
	assert.Greater(t, score, 0.7, "model should predict high fraud for consistent high-risk features")
}

func TestOnlineModel_GetSetWeights(t *testing.T) {
	model := NewOnlineModel(0.01, 0.001)

	weights := model.GetWeights()
	assert.Len(t, weights, 4)

	newWeights := []float64{0.1, 0.2, 0.3, 0.4}
	err := model.SetWeights(newWeights)
	assert.NoError(t, err)

	got := model.GetWeights()
	assert.Equal(t, newWeights, got)
}

func TestOnlineModel_SetWeights_DimensionMismatch(t *testing.T) {
	model := NewOnlineModel(0.01, 0.001)

	err := model.SetWeights([]float64{0.1, 0.2})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "weight dimension mismatch")
}

func TestOnlineModel_UpdateCount(t *testing.T) {
	model := NewOnlineModel(0.01, 0.001)
	assert.Equal(t, int64(0), model.GetUpdateCount())

	fv := FeatureVector{UserID: "u", Amount: 10, Velocity: 0.5, GraphRisk: 0.3}
	model.Update(fv, 1.0)
	model.Update(fv, 0.0)

	assert.Equal(t, int64(2), model.GetUpdateCount())
}
