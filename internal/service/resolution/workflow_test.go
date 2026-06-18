package resolution

import (
	"context"
	"testing"
)

type mockWorkflowProducer struct {
	publishFn func(ctx context.Context, key string, event interface{}) error
}

func (m *mockWorkflowProducer) Publish(ctx context.Context, key string, event interface{}) error {
	if m.publishFn != nil {
		return m.publishFn(ctx, key, event)
	}
	return nil
}

func TestNewWorkflow(t *testing.T) {
	arb := NewArbitrator()
	prod := &mockWorkflowProducer{}
	w := NewWorkflow(arb, prod)
	if w == nil {
		t.Fatal("NewWorkflow returned nil")
	}
}

func TestMerge_Success(t *testing.T) {
	published := false
	arb := NewArbitrator()
	prod := &mockWorkflowProducer{
		publishFn: func(ctx context.Context, key string, event interface{}) error {
			published = true
			return nil
		},
	}
	w := NewWorkflow(arb, prod)

	err := w.Merge(context.Background(), "ID-PRIMARY", []string{"ID-SEC-1", "ID-SEC-2"})
	if err != nil {
		t.Fatalf("Merge failed: %v", err)
	}
	if !published {
		t.Error("Merge event was not published")
	}
}

func TestMerge_EmptySecondary(t *testing.T) {
	arb := NewArbitrator()
	prod := &mockWorkflowProducer{}
	w := NewWorkflow(arb, prod)

	err := w.Merge(context.Background(), "ID-PRIMARY", []string{})
	if err != nil {
		t.Fatalf("Merge with empty secondaries failed: %v", err)
	}
}

func TestSplit_Success(t *testing.T) {
	published := false
	arb := NewArbitrator()
	prod := &mockWorkflowProducer{
		publishFn: func(ctx context.Context, key string, event interface{}) error {
			published = true
			return nil
		},
	}
	w := NewWorkflow(arb, prod)

	id1, id2, err := w.Split(context.Background(), "ID-ORIG")
	if err != nil {
		t.Fatalf("Split failed: %v", err)
	}
	if id1 != "ID-ORIG-A" {
		t.Errorf("ID1 = %s, want ID-ORIG-A", id1)
	}
	if id2 != "ID-ORIG-B" {
		t.Errorf("ID2 = %s, want ID-ORIG-B", id2)
	}
	if !published {
		t.Error("Split event was not published")
	}
}

func TestSplit_ProduceError(t *testing.T) {
	arb := NewArbitrator()
	prod := &mockWorkflowProducer{
		publishFn: func(ctx context.Context, key string, event interface{}) error {
			return nil
		},
	}
	w := NewWorkflow(arb, prod)

	_, _, err := w.Split(context.Background(), "ID-001")
	if err != nil {
		t.Fatalf("Split should handle producer errors gracefully: %v", err)
	}
}
