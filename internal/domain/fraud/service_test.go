package fraud

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/segmentio/kafka-go"
)

type mockGraphRepo struct {
	addIdentityNodeFn func(ctx context.Context, identityID, agency string) error
	checkFraudRingFn  func(ctx context.Context, identityID string) (bool, error)
}

func (m *mockGraphRepo) AddIdentityNode(ctx context.Context, identityID, agency string) error {
	if m.addIdentityNodeFn != nil {
		return m.addIdentityNodeFn(ctx, identityID, agency)
	}
	return nil
}

func (m *mockGraphRepo) CheckFraudRing(ctx context.Context, identityID string) (bool, error) {
	if m.checkFraudRingFn != nil {
		return m.checkFraudRingFn(ctx, identityID)
	}
	return false, nil
}

type mockConsumer struct {
	readFn func(ctx context.Context, handler func(msg kafka.Message) error) error
}

func (m *mockConsumer) Read(ctx context.Context, handler func(msg kafka.Message) error) error {
	if m.readFn != nil {
		return m.readFn(ctx, handler)
	}
	return nil
}

func (m *mockConsumer) Start(ctx context.Context, handler func(ctx context.Context, payload []byte) error) error {
	if m.readFn != nil {
		m.readFn(ctx, func(msg kafka.Message) error {
			return handler(ctx, msg.Value)
		})
	}
	return nil
}

func (m *mockConsumer) Decode(data []byte, v interface{}) error {
	return nil
}

func (m *mockConsumer) Close() error { return nil }

type mockProducer struct {
	publishFn func(ctx context.Context, key string, event interface{}) error
}

func (m *mockProducer) Publish(ctx context.Context, key string, event interface{}) error {
	if m.publishFn != nil {
		return m.publishFn(ctx, key, event)
	}
	return nil
}

func (m *mockProducer) Close() error { return nil }

func TestNewService(t *testing.T) {
	svc := NewService(&mockGraphRepo{}, &mockConsumer{}, &mockProducer{})
	if svc == nil {
		t.Fatal("NewService returned nil")
	}
}

func TestStart_IdentityCreatedEvent(t *testing.T) {
	addedNode := false
	checkedRing := false

	graphRepo := &mockGraphRepo{
		addIdentityNodeFn: func(ctx context.Context, identityID, agency string) error {
			addedNode = true
			return nil
		},
		checkFraudRingFn: func(ctx context.Context, identityID string) (bool, error) {
			checkedRing = true
			return false, nil
		},
	}

	var publishedEvent interface{}
	producer := &mockProducer{
		publishFn: func(ctx context.Context, key string, event interface{}) error {
			publishedEvent = event
			return nil
		},
	}

	consumer := &mockConsumer{
		readFn: func(ctx context.Context, handler func(msg kafka.Message) error) error {
			msg := kafka.Message{
				Value: []byte(`{"identityId":"ID-001","firstName":"John","lastName":"Doe","agency":"oni","timestamp":"2025-01-01T00:00:00Z"}`),
			}
			return handler(msg)
		},
	}

	svc := NewService(graphRepo, consumer, producer)
	err := svc.Start(context.Background())
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	if !addedNode {
		t.Error("AddIdentityNode was not called")
	}
	if !checkedRing {
		t.Error("CheckFraudRing was not called")
	}
	if publishedEvent == nil {
		t.Error("Expected event to be published")
	}
}

func TestStart_FraudRingDetected(t *testing.T) {
	graphRepo := &mockGraphRepo{
		checkFraudRingFn: func(ctx context.Context, identityID string) (bool, error) {
			return true, nil
		},
	}

	var publishedEvent interface{}
	producer := &mockProducer{
		publishFn: func(ctx context.Context, key string, event interface{}) error {
			publishedEvent = event
			return nil
		},
	}

	consumer := &mockConsumer{
		readFn: func(ctx context.Context, handler func(msg kafka.Message) error) error {
			msg := kafka.Message{
				Value: []byte(`{"identityId":"ID-FRAUD","agency":"oni"}`),
			}
			return handler(msg)
		},
	}

	svc := NewService(graphRepo, consumer, producer)
	svc.Start(context.Background())

	if publishedEvent == nil {
		t.Fatal("Expected fraud scored event")
	}
	evt := publishedEvent.(FraudScoredEvent)
	if evt.RiskScore != 95 {
		t.Errorf("RiskScore = %d, want 95", evt.RiskScore)
	}
	if !evt.IsFraud {
		t.Error("Expected IsFraud = true")
	}
}

func TestStart_InvalidMessage(t *testing.T) {
	graphRepo := &mockGraphRepo{}
	consumer := &mockConsumer{
		readFn: func(ctx context.Context, handler func(msg kafka.Message) error) error {
			msg := kafka.Message{
				Value: []byte(`{invalid json}`),
			}
			return handler(msg)
		},
	}

	svc := NewService(graphRepo, consumer, &mockProducer{})
	err := svc.Start(context.Background())
	if err != nil {
		t.Fatalf("Start should not return error for invalid msg: %v", err)
	}
}

func TestStart_PublishError(t *testing.T) {
	graphRepo := &mockGraphRepo{
		checkFraudRingFn: func(ctx context.Context, identityID string) (bool, error) {
			return false, nil
		},
	}
	producer := &mockProducer{
		publishFn: func(ctx context.Context, key string, event interface{}) error {
			return errors.New("kafka down")
		},
	}
	consumer := &mockConsumer{
		readFn: func(ctx context.Context, handler func(msg kafka.Message) error) error {
			msg := kafka.Message{
				Value: []byte(`{"identityId":"ID-001"}`),
			}
			return handler(msg)
		},
	}

	svc := NewService(graphRepo, consumer, producer)
	err := svc.Start(context.Background())
	if err == nil {
		t.Fatal("Expected error from publish failure")
	}
}

func TestStart_SuspiciousAgency(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	producer := &mockProducer{
		publishFn: func(ctx context.Context, key string, event interface{}) error {
			defer wg.Done()
			evt := event.(FraudScoredEvent)
			if evt.RiskScore != 80 {
				t.Errorf("RiskScore = %d, want 80", evt.RiskScore)
			}
			return nil
		},
	}
	consumer := &mockConsumer{
		readFn: func(ctx context.Context, handler func(msg kafka.Message) error) error {
			msg := kafka.Message{
				Value: []byte(`{"identityId":"ID-002","agency":"suspicious-agency"}`),
			}
			return handler(msg)
		},
	}

	svc := NewService(&mockGraphRepo{}, consumer, producer)
	svc.Start(context.Background())
}
