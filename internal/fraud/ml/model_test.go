package ml

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

type mockStateProvider struct {
	client *redis.Client
}

func (m *mockStateProvider) GetState(ctx context.Context, key string) (string, error) {
	return m.client.Get(ctx, key).Result()
}

func (m *mockStateProvider) IncrementVelocity(ctx context.Context, key string, window time.Duration) (int64, error) {
	pipe := m.client.Pipeline()
	count := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return count.Val(), nil
}

func newTestFeatureStore(t *testing.T) (*FeatureStore, *miniredis.Miniredis) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	provider := &mockStateProvider{client: client}
	store := NewFeatureStore(provider)
	t.Cleanup(func() {
		client.Close()
		s.Close()
	})
	return store, s
}

func TestGetTransactionVelocity_NoData(t *testing.T) {
	store, _ := newTestFeatureStore(t)
	ctx := context.Background()

	v := store.GetTransactionVelocity(ctx, "unknown-user")
	if v != 0.0 {
		t.Errorf("velocity = %f, want 0.0", v)
	}
}

func TestGetTransactionVelocity_WithData(t *testing.T) {
	store, mini := newTestFeatureStore(t)
	ctx := context.Background()

	mini.Set("snisid:features:user-001:velocity", "4")

	v := store.GetTransactionVelocity(ctx, "user-001")
	if v != 0.4 {
		t.Errorf("velocity = %f, want 0.4", v)
	}
}

func TestGetTransactionVelocity_Normalized(t *testing.T) {
	store, mini := newTestFeatureStore(t)
	ctx := context.Background()

	mini.Set("snisid:features:high-vel:velocity", "20")

	v := store.GetTransactionVelocity(ctx, "high-vel")
	if v != 1.0 {
		t.Errorf("velocity = %f, want 1.0 (capped at 10)", v)
	}
}

func TestGetTransactionVelocity_EmptyUserID(t *testing.T) {
	store, _ := newTestFeatureStore(t)
	ctx := context.Background()

	v := store.GetTransactionVelocity(ctx, "")
	if v != 0.0 {
		t.Errorf("velocity = %f, want 0.0 for empty userID", v)
	}
}

func TestNewFeatureStoreFromAddr(t *testing.T) {
	store := NewFeatureStoreFromAddr("localhost:6379")
	if store == nil {
		t.Fatal("NewFeatureStoreFromAddr returned nil")
	}
	if store.state == nil {
		t.Error("state provider should be initialized")
	}
}

func TestFraudModel_PredictTraditional(t *testing.T) {
	m := &FraudModel{Version: "1.0"}
	score := m.PredictTraditional(0.5, 100.0)
	expected := (0.5 * 0.7) + (100.0 * 0.3)
	if score != expected {
		t.Errorf("score = %f, want %f", score, expected)
	}
}
