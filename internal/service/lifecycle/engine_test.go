package lifecycle

import (
	"context"
	"testing"
)

type mockLifecycleProducer struct {
	publishFn func(ctx context.Context, key string, event interface{}) error
}

func (m *mockLifecycleProducer) Publish(ctx context.Context, key string, event interface{}) error {
	if m.publishFn != nil {
		return m.publishFn(ctx, key, event)
	}
	return nil
}

func (m *mockLifecycleProducer) Close() error { return nil }

func TestNewEngine(t *testing.T) {
	v := NewValidator()
	p := &mockLifecycleProducer{}
	e := NewEngine(v, p)
	if e == nil {
		t.Fatal("NewEngine returned nil")
	}
}

func TestTransition_Success(t *testing.T) {
	v := NewValidator()
	published := false
	p := &mockLifecycleProducer{
		publishFn: func(ctx context.Context, key string, event interface{}) error {
			published = true
			return nil
		},
	}
	e := NewEngine(v, p)

	err := e.Transition(context.Background(), "ID-001", StateCreated, StateVerified, "Standard enrollment", "agent-001")
	if err != nil {
		t.Fatalf("Transition failed: %v", err)
	}
	if !published {
		t.Error("Lifecycle event was not published")
	}
}

func TestTransition_IllegalMove(t *testing.T) {
	v := NewValidator()
	e := NewEngine(v, &mockLifecycleProducer{})

	// Can't go from CREATED directly to SUSPENDED
	err := e.Transition(context.Background(), "ID-001", StateCreated, StateSuspended, "test", "agent-001")
	if err == nil {
		t.Fatal("Expected error for illegal transition")
	}
}

func TestTransition_SameState(t *testing.T) {
	v := NewValidator()
	e := NewEngine(v, &mockLifecycleProducer{})

	// Transition to the same state should be allowed
	err := e.Transition(context.Background(), "ID-001", StateActive, StateActive, "no change", "agent-001")
	if err != nil {
		t.Fatalf("Same-state transition failed: %v", err)
	}
}

func TestTransition_PublishError(t *testing.T) {
	v := NewValidator()
	p := &mockLifecycleProducer{
		publishFn: func(ctx context.Context, key string, event interface{}) error {
			return nil
		},
	}
	e := NewEngine(v, p)

	err := e.Transition(context.Background(), "ID-001", StateVerified, StateActive, "approve", "admin")
	if err != nil {
		t.Fatalf("Transition failed: %v", err)
	}
}

func TestEngine_NilValidator(t *testing.T) {
	e := NewEngine(nil, &mockLifecycleProducer{})
	if e == nil {
		t.Fatal("NewEngine should handle nil validator")
	}
}
