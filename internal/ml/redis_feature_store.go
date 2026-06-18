package ml

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const (
	keyVelocity  = "snisid:features:%s:velocity"
	keyGraphRisk = "snisid:features:%s:graph_risk"
)

var ErrInvalidPayload = errors.New("feature extractor: invalid payload")

type RedisFeatureStore struct {
	client *redis.Client
}

func NewRedisFeatureStore(client *redis.Client) *RedisFeatureStore {
	return &RedisFeatureStore{client: client}
}

func (s *RedisFeatureStore) GetVelocity(ctx context.Context, userID string) (float64, error) {
	return s.getFloat(ctx, fmt.Sprintf(keyVelocity, userID))
}

func (s *RedisFeatureStore) GetGraphRisk(ctx context.Context, userID string) (float64, error) {
	return s.getFloat(ctx, fmt.Sprintf(keyGraphRisk, userID))
}

func (s *RedisFeatureStore) getFloat(ctx context.Context, key string) (float64, error) {
	val, err := s.client.Get(ctx, key).Float64()
	if errors.Is(err, redis.Nil) {
		return 0.0, nil
	}
	if err != nil {
		return 0.0, fmt.Errorf("redis feature store get %s: %w", key, err)
	}
	return val, nil
}
