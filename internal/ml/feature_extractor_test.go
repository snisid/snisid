package ml

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockFeatureStore struct {
	mock.Mock
}

func (m *MockFeatureStore) GetVelocity(ctx context.Context, userID string) (float64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockFeatureStore) GetGraphRisk(ctx context.Context, userID string) (float64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(float64), args.Error(1)
}

func TestExtractFeatures_RedisOK_ValeurPresente(t *testing.T) {
	store := &MockFeatureStore{}
	store.On("GetVelocity", mock.Anything, "user123").Return(0.88, nil)
	store.On("GetGraphRisk", mock.Anything, "user123").Return(0.42, nil)

	logger := zap.NewNop()
	extractor := NewFeatureExtractor(store, logger)

	payload := map[string]interface{}{
		"user_id": "user123",
		"amount":  150.0,
	}

	fv, err := extractor.ExtractFeatures(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, "user123", fv.UserID)
	assert.Equal(t, 0.88, fv.Velocity)
	assert.Equal(t, 0.42, fv.GraphRisk)
	assert.Equal(t, 150.0, fv.Amount)
}

func TestExtractFeatures_RedisOK_CleAbsente(t *testing.T) {
	store := &MockFeatureStore{}
	store.On("GetVelocity", mock.Anything, "newuser").Return(0.0, nil)
	store.On("GetGraphRisk", mock.Anything, "newuser").Return(0.0, nil)

	logger := zap.NewNop()
	extractor := NewFeatureExtractor(store, logger)

	payload := map[string]interface{}{
		"user_id": "newuser",
		"amount":  50.0,
	}

	fv, err := extractor.ExtractFeatures(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, 0.0, fv.Velocity)
	assert.Equal(t, 0.0, fv.GraphRisk)
}

func TestExtractFeatures_RedisIndisponible(t *testing.T) {
	store := &MockFeatureStore{}
	store.On("GetVelocity", mock.Anything, "user123").Return(0.0, errors.New("redis: connection refused"))

	logger := zap.NewNop()
	extractor := NewFeatureExtractor(store, logger)

	payload := map[string]interface{}{
		"user_id": "user123",
		"amount":  100.0,
	}

	_, err := extractor.ExtractFeatures(context.Background(), payload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "feature extraction failed")
}

func TestExtractFeatures_UserIDVide(t *testing.T) {
	store := &MockFeatureStore{}
	logger := zap.NewNop()
	extractor := NewFeatureExtractor(store, logger)

	payload := map[string]interface{}{
		"user_id": "",
		"amount":  100.0,
	}

	_, err := extractor.ExtractFeatures(context.Background(), payload)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidPayload))
}

func TestExtractFeatures_UserIDManquant(t *testing.T) {
	store := &MockFeatureStore{}
	logger := zap.NewNop()
	extractor := NewFeatureExtractor(store, logger)

	payload := map[string]interface{}{
		"amount": 100.0,
	}

	_, err := extractor.ExtractFeatures(context.Background(), payload)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidPayload))
}

func TestRedisFeatureStore_KeyFormat(t *testing.T) {
	assert.Equal(t, "snisid:features:user123:velocity", "snisid:features:user123:velocity")
	assert.Equal(t, "snisid:features:user123:graph_risk", "snisid:features:user123:graph_risk")
}
