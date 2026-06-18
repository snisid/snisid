package sandbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReward_LowFraudHighEconomy(t *testing.T) {
	state := EnvState{Fraud: 0.05, Economy: 0.9}
	r := Reward(state)
	// (1.0 - 0.05) + (0.9 * 0.3) = 0.95 + 0.27 = 1.22
	assert.InDelta(t, 1.22, r, 0.001)
}

func TestReward_HighFraudLowEconomy(t *testing.T) {
	state := EnvState{Fraud: 0.8, Economy: 0.2}
	r := Reward(state)
	// (1.0 - 0.8) + (0.2 * 0.3) = 0.2 + 0.06 = 0.26
	assert.InDelta(t, 0.26, r, 0.001)
}

func TestReward_MaxFraudMinEconomy(t *testing.T) {
	state := EnvState{Fraud: 1.0, Economy: 0.0}
	r := Reward(state)
	assert.InDelta(t, 0.0, r, 0.001)
}

func TestReward_ZeroFraudMaxEconomy(t *testing.T) {
	state := EnvState{Fraud: 0.0, Economy: 1.0}
	r := Reward(state)
	assert.InDelta(t, 1.3, r, 0.001)
}

func TestReward_NegativeValues(t *testing.T) {
	state := EnvState{Fraud: -0.1, Economy: -0.2}
	r := Reward(state)
	// (1.0 - (-0.1)) + (-0.2 * 0.3) = 1.1 - 0.06 = 1.04
	assert.InDelta(t, 1.04, r, 0.001)
}

func TestNewDriftMonitor(t *testing.T) {
	m := &DriftMonitor{Baseline: 0.5}
	assert.Equal(t, 0.5, m.Baseline)
}

func TestDetectCollapseRisk_NoDrift(t *testing.T) {
	m := &DriftMonitor{Baseline: 0.5}
	assert.False(t, m.DetectCollapseRisk(0.5))
	assert.False(t, m.DetectCollapseRisk(0.6))
	assert.False(t, m.DetectCollapseRisk(0.4))
}

func TestDetectCollapseRisk_Borderline(t *testing.T) {
	m := &DriftMonitor{Baseline: 0.5}
	assert.False(t, m.DetectCollapseRisk(0.74))
	assert.False(t, m.DetectCollapseRisk(0.26))
}

func TestDetectCollapseRisk_ExceedsThreshold(t *testing.T) {
	m := &DriftMonitor{Baseline: 0.5}
	assert.True(t, m.DetectCollapseRisk(0.76))
	assert.True(t, m.DetectCollapseRisk(0.24))
}

func TestDetectCollapseRisk_ExtremeValues(t *testing.T) {
	m := &DriftMonitor{Baseline: 0.0}
	assert.True(t, m.DetectCollapseRisk(0.26))
	assert.False(t, m.DetectCollapseRisk(0.25))
	assert.True(t, m.DetectCollapseRisk(-0.26))
}

func TestDetectCollapseRisk_NegativeBaseline(t *testing.T) {
	m := &DriftMonitor{Baseline: -0.5}
	assert.False(t, m.DetectCollapseRisk(-0.3))
	assert.True(t, m.DetectCollapseRisk(-0.76))
	assert.True(t, m.DetectCollapseRisk(-0.24))
}

func TestDriftMonitor_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		baseline float64
		current  float64
		expected bool
	}{
		{"no drift", 0.5, 0.5, false},
		{"small positive drift", 0.5, 0.7, false},
		{"small negative drift", 0.5, 0.3, false},
		{"large positive drift", 0.5, 0.8, true},
		{"large negative drift", 0.5, 0.2, true},
		{"from zero baseline large", 0.0, 0.3, true},
		{"from zero baseline small", 0.0, 0.25, false},
		{"at exact threshold", 0.5, 0.75, false},
		{"just over threshold", 0.5, 0.751, true},
		{"negative to positive", -0.5, 0.0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := &DriftMonitor{Baseline: tc.baseline}
			assert.Equal(t, tc.expected, m.DetectCollapseRisk(tc.current))
		})
	}
}

func TestReward_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		fraud    float64
		economy  float64
		expected float64
	}{
		{"ideal state", 0.0, 1.0, 1.30},
		{"high fraud", 0.5, 0.5, 0.65},
		{"catastrophic", 1.0, 0.0, 0.00},
		{"moderate", 0.2, 0.7, 1.01},
		{"recession", 0.3, 0.1, 0.73},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := Reward(EnvState{Fraud: tc.fraud, Economy: tc.economy})
			assert.InDelta(t, tc.expected, r, 0.001)
		})
	}
}
