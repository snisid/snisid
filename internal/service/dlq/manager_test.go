package dlq

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	m := NewManager([]string{"localhost:9092"}, "snisid.dlq.events")
	assert.NotNil(t, m)
	assert.NotNil(t, m.consumer)
	assert.NotNil(t, m.producers)
}

func TestManager_HandleDLQEvent_RetryUnderLimit(t *testing.T) {
	m := NewManager([]string{"localhost:9092"}, "snisid.dlq.events")

	event := map[string]interface{}{
		"retryCount":    0,
		"originalTopic": "snisid.prod.identity.v1.events",
		"originalKey":   "key-123",
		"eventId":       "evt-001",
		"header": map[string]interface{}{
			"correlationId": "corr-123",
		},
	}

	payload, err := json.Marshal(event)
	require.NoError(t, err)

	err = m.handleDLQEvent(context.Background(), payload)
	assert.NoError(t, err)
}

func TestManager_HandleDLQEvent_MaxRetries_Quarantine(t *testing.T) {
	m := NewManager([]string{"localhost:9092"}, "snisid.dlq.events")

	event := map[string]interface{}{
		"retryCount":    3,
		"originalTopic": "snisid.prod.identity.v1.events",
		"originalKey":   "key-456",
		"eventId":       "evt-002",
		"header": map[string]interface{}{
			"correlationId": "corr-456",
		},
	}

	payload, err := json.Marshal(event)
	require.NoError(t, err)

	err = m.handleDLQEvent(context.Background(), payload)
	assert.NoError(t, err) // Quarantine should not error
}

func TestManager_HandleDLQEvent_InvalidPayload(t *testing.T) {
	m := NewManager([]string{"localhost:9092"}, "snisid.dlq.events")
	err := m.handleDLQEvent(context.Background(), []byte("invalid-json"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

func TestManager_HandleDLQEvent_EmptyPayload(t *testing.T) {
	m := NewManager([]string{"localhost:9092"}, "snisid.dlq.events")
	err := m.handleDLQEvent(context.Background(), []byte{})
	assert.Error(t, err)
}

func TestManager_GetOrCreateProducer_ReusesExisting(t *testing.T) {
	m := NewManager([]string{"localhost:9092"}, "snisid.dlq.events")
	p1 := m.getOrCreateProducer("test-topic")
	p2 := m.getOrCreateProducer("test-topic")
	assert.Same(t, p1, p2)
}

func TestManager_GetOrCreateProducer_DifferentTopics(t *testing.T) {
	m := NewManager([]string{"localhost:9092"}, "snisid.dlq.events")
	p1 := m.getOrCreateProducer("topic-a")
	p2 := m.getOrCreateProducer("topic-b")
	assert.NotSame(t, p1, p2)
}

func TestManager_SaveToQuarantine(t *testing.T) {
	m := NewManager([]string{"localhost:9092"}, "snisid.dlq.events")
	err := m.saveToQuarantine(context.Background(), map[string]interface{}{
		"eventId": "quarantine-evt",
		"reason":  "max retries exceeded",
	})
	assert.NoError(t, err)
}

func TestManager_Close(t *testing.T) {
	m := NewManager([]string{"localhost:9092"}, "snisid.dlq.events")
	err := m.Close()
	assert.NoError(t, err)
}
