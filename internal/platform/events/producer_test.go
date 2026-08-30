package events

import (
	"context"
	"testing"
)

type mockKafkaWriter struct {
	writeMessagesFn func(ctx context.Context, msgs ...interface{}) error
	closeFn         func() error
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig([]string{"localhost:9092"}, "test-topic")
	if len(cfg.Brokers) != 1 {
		t.Errorf("Brokers count = %d, want 1", len(cfg.Brokers))
	}
	if cfg.Topic != "test-topic" {
		t.Errorf("Topic = %s, want test-topic", cfg.Topic)
	}
	if cfg.BatchSize != 100 {
		t.Errorf("BatchSize = %d, want 100", cfg.BatchSize)
	}
	if cfg.MaxRetries != 3 {
		t.Errorf("MaxRetries = %d, want 3", cfg.MaxRetries)
	}
	if cfg.Codec == nil {
		t.Error("Codec should not be nil")
	}
}

func TestNewProducer(t *testing.T) {
	p := NewProducer([]string{"localhost:9092"}, "test-topic")
	if p == nil {
		t.Fatal("NewProducer returned nil")
	}
	if p.writer == nil {
		t.Error("Kafka writer should be initialized")
	}
	if p.dlq == nil {
		t.Error("DLQ writer should be initialized")
	}
}

func TestNewProducerWithConfig(t *testing.T) {
	cfg := DefaultConfig([]string{"broker:9092"}, "custom-topic")
	cfg.BatchSize = 50
	cfg.MaxRetries = 5
	cfg.Codec = &JSONCodec{}

	p := NewProducerWithConfig(cfg)
	if p == nil {
		t.Fatal("NewProducerWithConfig returned nil")
	}
	if p.config.BatchSize != 50 {
		t.Errorf("BatchSize = %d, want 50", p.config.BatchSize)
	}
	if p.config.MaxRetries != 5 {
		t.Errorf("MaxRetries = %d, want 5", p.config.MaxRetries)
	}
}

func TestProducer_Publish_NoError(t *testing.T) {
	p := NewProducer([]string{"localhost:9092"}, "test-topic")
	err := p.Publish(context.Background(), "key-1", map[string]interface{}{"event": "test"})
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}
}

func TestProducer_Publish_JSONCodec(t *testing.T) {
	cfg := DefaultConfig([]string{"localhost:9092"}, "json-topic")
	cfg.Codec = &JSONCodec{}
	p := NewProducerWithConfig(cfg)

	err := p.Publish(context.Background(), "key-1", map[string]string{"type": "test"})
	if err != nil {
		t.Fatalf("Publish with JSON codec failed: %v", err)
	}
}

func TestProducer_Close(t *testing.T) {
	p := NewProducer([]string{"localhost:9092"}, "test-topic")
	err := p.Close()
	if err != nil {
		t.Logf("Close returned error (expected if no broker): %v", err)
	}
}
