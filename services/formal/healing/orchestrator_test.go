package healing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectAndHeal_SafeState(t *testing.T) {
	h := &HealingEngine{}
	state := &SystemState{
		Timestamp: 1000,
		RiskData:  map[string]int{"node-1": 50},
		PolicyMap: map[string]string{"node-1": "ALLOW"},
	}
	h.DetectAndHeal(state, true)
	require.NotNil(t, h.LastValidState)
	assert.Equal(t, int64(1000), h.LastValidState.Timestamp)
}

func TestDetectAndHeal_UnsafeState(t *testing.T) {
	h := &HealingEngine{
		LastValidState: &SystemState{Timestamp: 900},
	}
	state := &SystemState{
		Timestamp: 1000,
		RiskData:  map[string]int{"node-1": 200},
	}
	h.DetectAndHeal(state, false)
	require.NotNil(t, h.LastValidState)
	assert.Equal(t, int64(900), h.LastValidState.Timestamp)
}

func TestDetectAndHeal_UnsafeNoSnapshot(t *testing.T) {
	h := &HealingEngine{}
	state := &SystemState{Timestamp: 1000}
	h.DetectAndHeal(state, false)
	assert.Nil(t, h.LastValidState)
}

func TestSnapshot(t *testing.T) {
	h := &HealingEngine{}
	state := &SystemState{Timestamp: 500, RiskData: map[string]int{"a": 1}}
	h.Snapshot(state)
	require.NotNil(t, h.LastValidState)
	assert.Equal(t, int64(500), h.LastValidState.Timestamp)
	assert.Equal(t, 1, h.LastValidState.RiskData["a"])
}

func TestSnapshot_OverwritesPrevious(t *testing.T) {
	h := &HealingEngine{
		LastValidState: &SystemState{Timestamp: 100},
	}
	h.Snapshot(&SystemState{Timestamp: 200})
	assert.Equal(t, int64(200), h.LastValidState.Timestamp)
}

func TestRollbackToLastValid_HasSnapshot(t *testing.T) {
	h := &HealingEngine{
		LastValidState: &SystemState{Timestamp: 500},
	}
	h.RollbackToLastValid()
}

func TestRollbackToLastValid_NoSnapshot(t *testing.T) {
	h := &HealingEngine{}
	h.RollbackToLastValid()
}

func TestSystemState_DeepCopy(t *testing.T) {
	state := &SystemState{
		Timestamp: 100,
		RiskData:  map[string]int{"node-1": 50, "node-2": 30},
		PolicyMap: map[string]string{"node-1": "ALLOW", "node-2": "DENY"},
	}
	require.NotNil(t, state)
	assert.Len(t, state.RiskData, 2)
	assert.Len(t, state.PolicyMap, 2)
}

func TestMultipleSnapshots(t *testing.T) {
	h := &HealingEngine{}

	for i := 0; i < 5; i++ {
		state := &SystemState{
			Timestamp: int64(i * 100),
			RiskData:  map[string]int{"node": i},
		}
		h.DetectAndHeal(state, true)
	}

	require.NotNil(t, h.LastValidState)
	assert.Equal(t, int64(400), h.LastValidState.Timestamp)
	assert.Equal(t, 4, h.LastValidState.RiskData["node"])
}

func TestDetectionAndRollbackSequence(t *testing.T) {
	h := &HealingEngine{}

	h.DetectAndHeal(&SystemState{Timestamp: 100, RiskData: map[string]int{"r": 1}}, true)
	assert.Equal(t, int64(100), h.LastValidState.Timestamp)

	h.DetectAndHeal(&SystemState{Timestamp: 200, RiskData: map[string]int{"r": 999}}, false)
	assert.Equal(t, int64(100), h.LastValidState.Timestamp)
}
