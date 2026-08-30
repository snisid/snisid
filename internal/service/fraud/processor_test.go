package fraud

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
)

type mockProcessorConsumer struct {
	startFn func(ctx context.Context, handler func(ctx context.Context, payload []byte) error) error
}

func (m *mockProcessorConsumer) Start(ctx context.Context, handler func(ctx context.Context, payload []byte) error) error {
	if m.startFn != nil {
		return m.startFn(ctx, handler)
	}
	return nil
}

type mockProcessorProducer struct {
	publishFn func(ctx context.Context, key string, event interface{}) error
}

func (m *mockProcessorProducer) Publish(ctx context.Context, key string, event interface{}) error {
	if m.publishFn != nil {
		return m.publishFn(ctx, key, event)
	}
	return nil
}

func TestNewIdentityEventProcessor(t *testing.T) {
	p := NewIdentityEventProcessor([]string{"localhost:9092"})
	if p == nil {
		t.Fatal("NewIdentityEventProcessor returned nil")
	}
}

func TestProcessEvent_ValidEvent(t *testing.T) {
	published := false
	var wg sync.WaitGroup
	wg.Add(1)

	p := &IdentityEventProcessor{
		consumer: &mockProcessorConsumer{},
		producer: &mockProcessorProducer{
			publishFn: func(ctx context.Context, key string, event interface{}) error {
				published = true
				wg.Done()
				return nil
			},
		},
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"identityId": "ID-001",
		"action":     "deleted",
		"metadata": map[string]interface{}{
			"force": "true",
		},
	})

	err := p.processEvent(context.Background(), payload)
	if err != nil {
		t.Fatalf("processEvent failed: %v", err)
	}
	wg.Wait()

	if !published {
		t.Error("Expected alert to be published for force delete")
	}
}

func TestProcessEvent_InvalidJSON(t *testing.T) {
	p := &IdentityEventProcessor{
		consumer: &mockProcessorConsumer{},
		producer: &mockProcessorProducer{},
	}

	err := p.processEvent(context.Background(), []byte(`{invalid}`))
	if err == nil {
		t.Fatal("Expected error for invalid JSON")
	}
}

func TestProcessEvent_NormalAction(t *testing.T) {
	published := false
	p := &IdentityEventProcessor{
		consumer: &mockProcessorConsumer{},
		producer: &mockProcessorProducer{
			publishFn: func(ctx context.Context, key string, event interface{}) error {
				published = true
				return nil
			},
		},
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"identityId": "ID-002",
		"action":     "created",
	})

	err := p.processEvent(context.Background(), payload)
	if err != nil {
		t.Fatalf("processEvent failed: %v", err)
	}
	if published {
		t.Error("No alert expected for normal action")
	}
}

func TestAnalyzeFraudPatterns_ForceDelete(t *testing.T) {
	p := &IdentityEventProcessor{}
	event := map[string]interface{}{
		"action": "deleted",
		"metadata": map[string]interface{}{
			"force": "true",
		},
	}
	if !p.analyzeFraudPatterns(event) {
		t.Error("Expected fraud pattern detected for force delete")
	}
}

func TestAnalyzeFraudPatterns_NormalDelete(t *testing.T) {
	p := &IdentityEventProcessor{}
	event := map[string]interface{}{
		"action": "deleted",
	}
	if p.analyzeFraudPatterns(event) {
		t.Error("No fraud expected for normal delete without force flag")
	}
}

func TestAnalyzeFraudPatterns_NonDelete(t *testing.T) {
	p := &IdentityEventProcessor{}
	event := map[string]interface{}{
		"action": "created",
	}
	if p.analyzeFraudPatterns(event) {
		t.Error("No fraud expected for non-delete action")
	}
}
