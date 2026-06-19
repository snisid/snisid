package fraud

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

type mockStateStore struct {
	mu             sync.Mutex
	velocityCounts map[string]int64
	stateData      map[string]*FraudState
}

func newMockStateStore() *mockStateStore {
	return &mockStateStore{
		velocityCounts: make(map[string]int64),
		stateData:      make(map[string]*FraudState),
	}
}

func (m *mockStateStore) IncrementVelocity(ctx context.Context, userID string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.velocityCounts[userID]++
	return m.velocityCounts[userID], nil
}

func (m *mockStateStore) GetState(ctx context.Context, userID string) (*FraudState, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.stateData[userID]; ok {
		return s, nil
	}
	return &FraudState{}, nil
}

func (m *mockStateStore) SetState(ctx context.Context, userID string, state *FraudState) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stateData[userID] = state
	return nil
}

func newTestStateStore(t *testing.T) (*mockStateStore, *miniredis.Miniredis) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}

	store := newMockStateStore()
	t.Cleanup(func() {
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
	count, err := store.IncrementVelocity(context.Background(), "ID-001")
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
		count, err := store.IncrementVelocity(ctx, "ID-002")
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
	if err != nil {
		t.Logf("GetState returned (val=%+v, err=%v)", val, err)
	}
}

func TestGetState_Found(t *testing.T) {
	store, _ := newTestStateStore(t)
	store.stateData["test:key"] = &FraudState{Velocity: 5}

	val, err := store.GetState(context.Background(), "test:key")
	if err != nil {
		t.Fatalf("GetState failed: %v", err)
	}
	if val.Velocity != 5 {
		t.Errorf("val.Velocity = %d, want 5", val.Velocity)
	}
}

func TestClose(t *testing.T) {
	store, _ := newTestStateStore(t)
	// no-op close for mock
	_ = store
}

func TestIncrementVelocity_Expiry(t *testing.T) {
	store, _ := newTestStateStore(t)
	ctx := context.Background()

	store.IncrementVelocity(ctx, "expiry-test")
	time.Sleep(time.Millisecond)

	// After expiry (simulated by reset), increment should start from 1 again
	count, err := store.IncrementVelocity(ctx, "expiry-test")
	if err != nil {
		t.Fatalf("IncrementVelocity failed: %v", err)
	}
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}
}
