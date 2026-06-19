package ml

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
)

type StateProvider interface {
	GetState(ctx context.Context, key string) (string, error)
	IncrementVelocity(ctx context.Context, key string, window time.Duration) (int64, error)
}

type FeatureStore struct {
	state StateProvider
}

func NewFeatureStore(state StateProvider) *FeatureStore {
	return &FeatureStore{state: state}
}

func NewFeatureStoreFromAddr(addr string) *FeatureStore {
	client := redis.NewClient(&redis.Options{Addr: addr})
	return NewFeatureStore(NewRedisFeatureStore(client))
}

type RedisFeatureStore struct {
	client *redis.Client
}

func NewRedisFeatureStore(client *redis.Client) *RedisFeatureStore {
	return &RedisFeatureStore{client: client}
}

func (s *RedisFeatureStore) GetState(ctx context.Context, key string) (string, error) {
	return s.client.Get(ctx, key).Result()
}

func (s *RedisFeatureStore) IncrementVelocity(ctx context.Context, key string, window time.Duration) (int64, error) {
	val, err := s.client.Incr(ctx, key).Result()
	if err == nil {
		s.client.Expire(ctx, key, window)
	}
	return val, err
}

func (s *RedisFeatureStore) GetTransactionVelocity(ctx context.Context, userID string) (float64, error) {
	val, err := s.client.Get(ctx, fmt.Sprintf("snisid:features:%s:velocity", userID)).Float64()
	if errors.Is(err, redis.Nil) {
		return 0.0, nil
	}
	if err != nil {
		return 0.0, fmt.Errorf("get velocity for %s: %w", userID, err)
	}
	return val, nil
}

func (fs *FeatureStore) GetTransactionVelocity(ctx context.Context, userID string) float64 {
	if userID == "" || fs.state == nil {
		return 0.0
	}
	key := fmt.Sprintf("snisid:features:%s:velocity", userID)
	count, err := fs.state.GetState(ctx, key)
	if err != nil {
		return 0.0
	}
	var c float64
	fmt.Sscanf(count, "%f", &c)
	return math.Min(c/10.0, 1.0)
}

type FraudModel struct {
	Version         string
	store           *FeatureStore
	lastKnownAmount float64
}

func (m *FraudModel) Predict(ctx context.Context, userID string) (float64, error) {
	if m.store == nil {
		return 0.5, nil
	}
	velocity := m.store.GetTransactionVelocity(ctx, userID)
	score := velocity*0.7 + m.lastKnownAmount*0.3
	return math.Min(score, 1.0), nil
}

func (m *FraudModel) PredictTraditional(velocity, amount float64) float64 {
	return (velocity * 0.7) + (amount * 0.3)
}
