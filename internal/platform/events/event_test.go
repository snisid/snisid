package events

import (
	"testing"
	"time"
)

type testData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func TestNewEvent_Success(t *testing.T) {
	data := testData{ID: "evt-001", Name: "IdentityCreated"}
	evt, err := NewEvent("identity.created", data)
	if err != nil {
		t.Fatalf("NewEvent failed: %v", err)
	}
	if evt == nil {
		t.Fatal("NewEvent returned nil")
	}
	if evt.EventID == "" {
		t.Error("EventID should not be empty")
	}
	if evt.EventType != "identity.created" {
		t.Errorf("EventType = %s, want identity.created", evt.EventType)
	}
	if evt.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}
	if evt.Data.ID != "evt-001" {
		t.Errorf("Data.ID = %s, want evt-001", evt.Data.ID)
	}
}

func TestNewEvent_EmptyType(t *testing.T) {
	_, err := NewEvent("", testData{ID: "test"})
	if err == nil {
		t.Error("Expected validation error for empty event type")
	}
}

func TestNewEvent_UniqueIDs(t *testing.T) {
	evt1, _ := NewEvent("type.a", testData{ID: "1"})
	evt2, _ := NewEvent("type.a", testData{ID: "2"})
	if evt1.EventID == evt2.EventID {
		t.Error("EventIDs should be unique")
	}
}

func TestNewEvent_TimestampUTC(t *testing.T) {
	evt, _ := NewEvent("test.event", testData{ID: "1"})
	loc := evt.Timestamp.Location()
	if loc.String() != "UTC" {
		t.Errorf("Timestamp location = %s, want UTC", loc.String())
	}
}

func TestEnvelope_JSONStructure(t *testing.T) {
	data := testData{ID: "evt-002", Name: "TestEvent"}
	evt, _ := NewEvent("test.event", data)

	if evt.EventID == "" {
		t.Error("eventId should be present")
	}
	if evt.Timestamp.IsZero() {
		t.Error("timestamp should be set")
	}
}

func TestEnvelope_Generic(t *testing.T) {
	// Test with a different type
	type otherData struct {
		Value int `json:"value"`
	}
	evt, err := NewEvent("score.updated", otherData{Value: 42})
	if err != nil {
		t.Fatalf("NewEvent with otherData failed: %v", err)
	}
	if evt.Data.Value != 42 {
		t.Errorf("Value = %d, want 42", evt.Data.Value)
	}
}

func TestEnvelope_TimePrecision(t *testing.T) {
	before := time.Now().UTC()
	evt, _ := NewEvent("precision.test", testData{ID: "1"})
	after := time.Now().UTC()

	if evt.Timestamp.Before(before) {
		t.Error("Timestamp should be after 'before'")
	}
	if evt.Timestamp.After(after) {
		t.Error("Timestamp should be before 'after'")
	}
}
