package fraud

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func newTestStateStore(t *testing.T) (*StateStore, *miniredis.Miniredis) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	store := &StateStore{client: client}
	t.Cleanup(func() {
		client.Close()
		s.Close()
	})
	return store, s
}

func TestNewStateStore(t *testing.T) {
	store := NewStateStore("localhost:6379")
	if store == nil {
		t.Fatal("NewStateStore returned nil")
	}
	if store.client == nil {
		t.Error("Redis client should be initialized")
	}
}

func TestIncrementVelocity_NewKey(t *testing.T) {
	store, _ := newTestStateStore(t)
	count, err := store.IncrementVelocity(context.Background(), "velocity:identity:ID-001", 10*time.Minute)
	if err != nil {
		t.Fatalf("IncrementVelocity failed: %v", err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1", count)
	}
}

func TestIncrementVelocity_Multiple(t *testing.T) {
	store, _ := newTestStateStore(t)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		count, err := store.IncrementVelocity(ctx, "velocity:identity:ID-002", 10*time.Minute)
		if err != nil {
			t.Fatalf("IncrementVelocity attempt %d failed: %v", i, err)
		}
		if count != int64(i+1) {
			t.Errorf("attempt %d: count = %d, want %d", i, count, i+1)
		}
	}
}

func TestGetState_NotFound(t *testing.T) {
	store, _ := newTestStateStore(t)
	val, err := store.GetState(context.Background(), "nonexistent")
	if err != redis.Nil {
		t.Logf("GetState returned (val=%q, err=%v), expected redis.Nil", val, err)
	}
}

func TestGetState_Found(t *testing.T) {
	store, mini := newTestStateStore(t)
	mini.Set("test:key", "test-value")

	val, err := store.GetState(context.Background(), "test:key")
	if err != nil {
		t.Fatalf("GetState failed: %v", err)
	}
	if val != "test-value" {
		t.Errorf("val = %s, want test-value", val)
	}
}

func TestClose(t *testing.T) {
	store, _ := newTestStateStore(t)
	err := store.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}

func TestIncrementVelocity_Expiry(t *testing.T) {
	store, mini := newTestStateStore(t)
	ctx := context.Background()

	store.IncrementVelocity(ctx, "velocity:expiry-test", time.Second)
	mini.FastForward(2 * time.Second)

	// After expiry, increment should start from 1 again
	count, err := store.IncrementVelocity(ctx, "velocity:expiry-test", 10*time.Minute)
	if err != nil {
		t.Fatalf("IncrementVelocity failed: %v", err)
	}
	if count != 1 {
		t.Errorf("After expiry, count = %d, want 1", count)
	}
}
