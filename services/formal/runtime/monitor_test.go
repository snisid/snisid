package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEvent_WithinThreshold(t *testing.T) {
	m := FormalMonitor{Threshold: 100}
	evt := SystemEvent{ID: "evt-1", Risk: 50, Node: "node-1"}
	assert.True(t, m.ValidateEvent(evt))
}

func TestValidateEvent_AtThreshold(t *testing.T) {
	m := FormalMonitor{Threshold: 100}
	evt := SystemEvent{ID: "evt-2", Risk: 100, Node: "node-2"}
	assert.True(t, m.ValidateEvent(evt))
}

func TestValidateEvent_ExceedsThreshold(t *testing.T) {
	m := FormalMonitor{Threshold: 100}
	evt := SystemEvent{ID: "evt-3", Risk: 150, Node: "node-3"}
	assert.False(t, m.ValidateEvent(evt))
}

func TestValidateEvent_ZeroRisk(t *testing.T) {
	m := FormalMonitor{Threshold: 50}
	evt := SystemEvent{ID: "evt-4", Risk: 0, Node: "node-4"}
	assert.True(t, m.ValidateEvent(evt))
}

func TestValidateEvent_ZeroThreshold(t *testing.T) {
	m := FormalMonitor{Threshold: 0}
	evt := SystemEvent{ID: "evt-5", Risk: 1, Node: "node-5"}
	assert.False(t, m.ValidateEvent(evt))
}

func TestValidateEvent_NegativeRisk(t *testing.T) {
	m := FormalMonitor{Threshold: 100}
	evt := SystemEvent{ID: "evt-6", Risk: -10, Node: "node-6"}
	assert.True(t, m.ValidateEvent(evt))
}

func TestValidateEvent_MultipleEvents(t *testing.T) {
	m := FormalMonitor{Threshold: 50}
	assert.True(t, m.ValidateEvent(SystemEvent{Risk: 30}))
	assert.True(t, m.ValidateEvent(SystemEvent{Risk: 50}))
	assert.False(t, m.ValidateEvent(SystemEvent{Risk: 51}))
}

func TestTriggerEmergencyResponse(t *testing.T) {
	m := FormalMonitor{Threshold: 100}
	evt := SystemEvent{ID: "emergency-1", Risk: 200, Node: "node-crash"}
	m.TriggerEmergencyResponse(evt)
}

func TestFormalMonitor_Struct(t *testing.T) {
	m := FormalMonitor{Threshold: 75}
	assert.Equal(t, 75, m.Threshold)
}
