package monitor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckInvariant_Pass(t *testing.T) {
	c := RuntimeChecker{ID: "checker-1"}
	state := SystemState{
		RiskVector: map[string]int{
			"node-1": 30,
			"node-2": 50,
		},
		Threshold: 100,
		Policies: map[string]string{
			"node-1": "ALLOW",
			"node-2": "ALLOW",
		},
	}
	ok, msg := c.CheckInvariant(state)
	assert.True(t, ok)
	assert.Equal(t, "PASS", msg)
}

func TestCheckInvariant_RiskExceedsThresholdPolicyAllow(t *testing.T) {
	c := RuntimeChecker{ID: "checker-1"}
	state := SystemState{
		RiskVector: map[string]int{
			"node-1": 150,
		},
		Threshold: 100,
		Policies: map[string]string{
			"node-1": "ALLOW",
		},
	}
	ok, msg := c.CheckInvariant(state)
	assert.False(t, ok)
	assert.Contains(t, msg, "INVARIANT_VIOLATION")
}

func TestCheckInvariant_RiskExceedsThresholdPolicyDeny(t *testing.T) {
	c := RuntimeChecker{ID: "checker-1"}
	state := SystemState{
		RiskVector: map[string]int{
			"node-1": 150,
		},
		Threshold: 100,
		Policies: map[string]string{
			"node-1": "DENY",
		},
	}
	ok, msg := c.CheckInvariant(state)
	assert.True(t, ok)
	assert.Equal(t, "PASS", msg)
}

func TestCheckInvariant_RiskAtThreshold(t *testing.T) {
	c := RuntimeChecker{ID: "checker-1"}
	state := SystemState{
		RiskVector: map[string]int{"node-1": 100},
		Threshold:  100,
		Policies:   map[string]string{"node-1": "ALLOW"},
	}
	ok, _ := c.CheckInvariant(state)
	assert.True(t, ok)
}

func TestCheckInvariant_EmptyRiskVector(t *testing.T) {
	c := RuntimeChecker{ID: "checker-1"}
	state := SystemState{
		RiskVector: map[string]int{},
		Threshold:  100,
		Policies:   map[string]string{},
	}
	ok, msg := c.CheckInvariant(state)
	assert.True(t, ok)
	assert.Equal(t, "PASS", msg)
}

func TestCheckInvariant_MultipleNodesOneViolation(t *testing.T) {
	c := RuntimeChecker{ID: "checker-1"}
	state := SystemState{
		RiskVector: map[string]int{
			"node-1": 50,
			"node-2": 200,
			"node-3": 30,
		},
		Threshold: 100,
		Policies: map[string]string{
			"node-1": "ALLOW",
			"node-2": "ALLOW",
			"node-3": "DENY",
		},
	}
	ok, msg := c.CheckInvariant(state)
	assert.False(t, ok)
	assert.Contains(t, msg, "node-2")
}

func TestCheckInvariant_ZeroThreshold(t *testing.T) {
	c := RuntimeChecker{ID: "checker-1"}
	state := SystemState{
		RiskVector: map[string]int{"node-1": 1},
		Threshold:  0,
		Policies:   map[string]string{"node-1": "ALLOW"},
	}
	ok, _ := c.CheckInvariant(state)
	assert.False(t, ok)
}

func TestOnViolation(t *testing.T) {
	c := RuntimeChecker{ID: "checker-1"}
	c.OnViolation("test violation")
}

func TestRuntimeChecker_ID(t *testing.T) {
	c := RuntimeChecker{ID: "checker-prod-01"}
	assert.Equal(t, "checker-prod-01", c.ID)
}
