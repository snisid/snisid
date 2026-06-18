package events

import (
	"context"
	"errors"
	"testing"
)

type mockHandler struct {
	handleFn func(ctx context.Context, msg []byte) error
}

func (h *mockHandler) Handle(ctx context.Context, msg []byte) error {
	if h.handleFn != nil {
		return h.handleFn(ctx, msg)
	}
	return nil
}

func TestNewConsumer(t *testing.T) {
	c := NewConsumer([]string{"localhost:9092"}, "test-group", "test-topic")
	if c == nil {
		t.Fatal("NewConsumer returned nil")
	}
	if c.reader == nil {
		t.Error("Kafka reader should be initialized")
	}
	if c.dlq == nil {
		t.Error("DLQ producer should be initialized")
	}
}

func TestConsumerConfig_Defaults(t *testing.T) {
	c := NewConsumer([]string{"broker:9092"}, "group-1", "topic-1")
	if c.config.MaxWorkers != 10 {
		t.Errorf("MaxWorkers = %d, want 10", c.config.MaxWorkers)
	}
	if c.config.MaxRetries != 3 {
		t.Errorf("MaxRetries = %d, want 3", c.config.MaxRetries)
	}
	if c.config.RetryDelay.String() != "2s" {
		t.Errorf("RetryDelay = %s, want 2s", c.config.RetryDelay)
	}
	if c.config.Codec == nil {
		t.Error("Codec should not be nil")
	}
}

func TestConsumer_Decode_JSON(t *testing.T) {
	c := NewConsumer([]string{"localhost:9092"}, "test-group", "test-topic")
	var result map[string]interface{}
	err := c.Decode([]byte(`{"key":"value"}`), &result)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if result["key"] != "value" {
		t.Errorf("key = %s, want value", result["key"])
	}
}

func TestConsumer_Decode_InvalidJSON(t *testing.T) {
	c := NewConsumer([]string{"localhost:9092"}, "test-group", "test-topic")
	var result map[string]interface{}
	err := c.Decode([]byte(`{invalid}`), &result)
	if err == nil {
		t.Error("Expected decode error for invalid JSON")
	}
}

func TestConsumer_Close(t *testing.T) {
	c := NewConsumer([]string{"localhost:9092"}, "test-group", "test-topic")
	err := c.Close()
	if err != nil {
		t.Logf("Close returned error (expected if no broker): %v", err)
	}
}

func TestConsumer_processWithRetry_Success(t *testing.T) {
	c := NewConsumer([]string{"localhost:9092"}, "test-group", "test-topic")
	callCount := 0
	handler := func(ctx context.Context, msg []byte) error {
		callCount++
		return nil
	}

	msg := kafka.Message{Value: []byte(`{"test":"data"}`)}
	err := c.processWithRetry(context.Background(), msg, handler)
	if err != nil {
		t.Fatalf("processWithRetry failed: %v", err)
	}
	if callCount != 1 {
		t.Errorf("Handler called %d times, want 1", callCount)
	}
}

func TestConsumer_processWithRetry_RetryThenSuccess(t *testing.T) {
	c := NewConsumer([]string{"localhost:9092"}, "test-group", "test-topic")
	callCount := 0
	handler := func(ctx context.Context, msg []byte) error {
		callCount++
		if callCount < 3 {
			return errors.New("transient error")
		}
		return nil
	}

	msg := kafka.Message{Value: []byte(`{"test":"data"}`)}
	err := c.processWithRetry(context.Background(), msg, handler)
	if err != nil {
		t.Fatalf("processWithRetry failed: %v", err)
	}
	if callCount != 3 {
		t.Errorf("Handler called %d times, want 3", callCount)
	}
}

func TestConsumer_processWithRetry_MaxRetriesExceeded(t *testing.T) {
	c := NewConsumer([]string{"localhost:9092"}, "test-group", "test-topic")
	callCount := 0
	handler := func(ctx context.Context, msg []byte) error {
		callCount++
		return errors.New("persistent error")
	}

	msg := kafka.Message{Value: []byte(`{"test":"data"}`)}
	err := c.processWithRetry(context.Background(), msg, handler)
	if err == nil {
		t.Fatal("Expected error after max retries")
	}
	// MaxRetries=3, so handler should be called 4 times (initial + 3 retries)
	if callCount != 4 {
		t.Errorf("Handler called %d times, want 4", callCount)
	}
}
