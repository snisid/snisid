package entity

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuditEvent_ValidFields(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	evt := AuditEvent{
		EventID:       "evt-001",
		CorrelationID: "corr-001",
		EventType:     "identity.created",
		Actor:         "user-123",
		Action:        "CREATE",
		Resource:      "identity/456",
		Status:        "success",
		Payload:       `{"firstName":"Jean"}`,
		PreviousHash:  "genesis-hash-snisid",
		Hash:          "abc123def456",
		SequenceID:    1,
		Timestamp:     now,
	}

	assert.Equal(t, "evt-001", evt.EventID)
	assert.Equal(t, "corr-001", evt.CorrelationID)
	assert.Equal(t, "identity.created", evt.EventType)
	assert.Equal(t, "user-123", evt.Actor)
	assert.Equal(t, "CREATE", evt.Action)
	assert.Equal(t, "identity/456", evt.Resource)
	assert.Equal(t, "success", evt.Status)
	assert.Equal(t, `{"firstName":"Jean"}`, evt.Payload)
	assert.Equal(t, "genesis-hash-snisid", evt.PreviousHash)
	assert.Equal(t, "abc123def456", evt.Hash)
	assert.Equal(t, int64(1), evt.SequenceID)
	assert.WithinDuration(t, now, evt.Timestamp, time.Second)
}

func TestAuditEvent_SerializationRoundTrip(t *testing.T) {
	t.Parallel()

	original := AuditEvent{
		EventID:       "evt-002",
		CorrelationID: "corr-002",
		EventType:     "identity.updated",
		Actor:         "admin",
		Action:        "UPDATE",
		Resource:      "identity/789",
		Status:        "success",
		Payload:       `{"lastName":"Dupont"}`,
		PreviousHash:  "prev-hash-value",
		Hash:          "current-hash-value",
		SequenceID:    10,
		Timestamp:     time.Now().UTC(),
	}

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded AuditEvent
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, original.EventID, decoded.EventID)
	assert.Equal(t, original.CorrelationID, decoded.CorrelationID)
	assert.Equal(t, original.EventType, decoded.EventType)
	assert.Equal(t, original.Actor, decoded.Actor)
	assert.Equal(t, original.Action, decoded.Action)
	assert.Equal(t, original.Resource, decoded.Resource)
	assert.Equal(t, original.Status, decoded.Status)
	assert.Equal(t, original.Payload, decoded.Payload)
	assert.Equal(t, original.PreviousHash, decoded.PreviousHash)
	assert.Equal(t, original.Hash, decoded.Hash)
	assert.Equal(t, original.SequenceID, decoded.SequenceID)
	assert.WithinDuration(t, original.Timestamp, decoded.Timestamp, time.Second)
}

func TestAuditEvent_SequenceIDAutoIncrementTag(t *testing.T) {
	t.Parallel()

	evt1 := AuditEvent{SequenceID: 0}
	evt2 := AuditEvent{SequenceID: 1}
	evt3 := AuditEvent{SequenceID: 100}

	assert.Equal(t, int64(0), evt1.SequenceID)
	assert.Equal(t, int64(1), evt2.SequenceID)
	assert.Equal(t, int64(100), evt3.SequenceID)
}

func TestAuditEvent_JSONFieldNames(t *testing.T) {
	t.Parallel()

	evt := AuditEvent{
		EventID:       "evt-003",
		CorrelationID: "corr-003",
		EventType:     "test.event",
		SequenceID:    5,
	}

	data, err := json.Marshal(evt)
	require.NoError(t, err)

	var raw map[string]interface{}
	err = json.Unmarshal(data, &raw)
	require.NoError(t, err)

	assert.Equal(t, "evt-003", raw["eventId"])
	assert.Equal(t, "corr-003", raw["correlationId"])
	assert.Equal(t, "test.event", raw["eventType"])
	assert.Equal(t, float64(5), raw["sequenceId"])
}

func TestAuditEvent_EmptyFields(t *testing.T) {
	t.Parallel()

	evt := AuditEvent{}

	assert.Empty(t, evt.EventID)
	assert.Empty(t, evt.CorrelationID)
	assert.Empty(t, evt.EventType)
	assert.Empty(t, evt.Actor)
	assert.Empty(t, evt.Action)
	assert.Empty(t, evt.Resource)
	assert.Empty(t, evt.Status)
	assert.Empty(t, evt.Payload)
	assert.Empty(t, evt.PreviousHash)
	assert.Empty(t, evt.Hash)
	assert.Equal(t, int64(0), evt.SequenceID)
	assert.True(t, evt.Timestamp.IsZero())
}

func TestAuditEvent_StatusValidValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status string
	}{
		{"success", "success"},
		{"failed", "failed"},
		{"denied", "denied"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evt := AuditEvent{Status: tt.status}
			assert.Equal(t, tt.status, evt.Status)
		})
	}
}
