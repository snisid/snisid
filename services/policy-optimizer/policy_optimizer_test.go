package policyoptimizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPolicyOptimizer(t *testing.T) {
	o := &PolicyOptimizer{LearningRate: 0.1}
	assert.InDelta(t, 0.1, o.LearningRate, 0.001)
}

func TestCalculateReward_FraudReduction(t *testing.T) {
	o := &PolicyOptimizer{}
	state := SystemState{FraudRate: 0.1, FalsePositives: 0.05, Instability: 0.02}
	reward := o.CalculateReward(state, 0.3)
	// fraudReduction = 0.3 - 0.1 = 0.2
	// reward = (0.2 * 100) - (0.05 * 50) - (0.02 * 20) = 20 - 2.5 - 0.4 = 17.1
	assert.InDelta(t, 17.1, reward, 0.001)
}

func TestCalculateReward_NegativeReduction(t *testing.T) {
	o := &PolicyOptimizer{}
	state := SystemState{FraudRate: 0.5, FalsePositives: 0.1, Instability: 0.05}
	reward := o.CalculateReward(state, 0.2)
	// fraudReduction = 0.2 - 0.5 = -0.3
	// reward = (-0.3 * 100) - (0.1 * 50) - (0.05 * 20) = -30 - 5 - 1 = -36
	assert.InDelta(t, -36.0, reward, 0.001)
}

func TestCalculateReward_NoChange(t *testing.T) {
	o := &PolicyOptimizer{}
	state := SystemState{FraudRate: 0.2, FalsePositives: 0, Instability: 0}
	reward := o.CalculateReward(state, 0.2)
	assert.InDelta(t, 0.0, reward, 0.001)
}

func TestCalculateReward_HighPenalties(t *testing.T) {
	o := &PolicyOptimizer{}
	state := SystemState{FraudRate: 0.0, FalsePositives: 1.0, Instability: 1.0}
	reward := o.CalculateReward(state, 0.5)
	// fraudReduction = 0.5 - 0 = 0.5
	// reward = (0.5 * 100) - (1.0 * 50) - (1.0 * 20) = 50 - 50 - 20 = -20
	assert.InDelta(t, -20.0, reward, 0.001)
}

func TestCalculateReward_AllMax(t *testing.T) {
	o := &PolicyOptimizer{}
	state := SystemState{FraudRate: 1.0, FalsePositives: 1.0, Instability: 1.0}
	reward := o.CalculateReward(state, 1.0)
	assert.InDelta(t, -70.0, reward, 0.001)
}

func TestCalculateReward_AllMin(t *testing.T) {
	o := &PolicyOptimizer{}
	state := SystemState{FraudRate: 0, FalsePositives: 0, Instability: 0}
	reward := o.CalculateReward(state, 0)
	assert.InDelta(t, 0.0, reward, 0.001)
}

func TestCalculateReward_TableDriven(t *testing.T) {
	o := &PolicyOptimizer{}

	tests := []struct {
		name          string
		fraudRate     float64
		falsePos      float64
		instability   float64
		previousFraud float64
		expected      float64
	}{
		{"significant improvement", 0.1, 0.02, 0.01, 0.4, 28.8},
		{"slight improvement", 0.25, 0.05, 0.03, 0.3, 1.9},
		{"degradation", 0.4, 0.1, 0.05, 0.3, -16.0},
		{"perfect state", 0.0, 0.0, 0.0, 0.5, 50.0},
		{"high false positives", 0.05, 0.8, 0.1, 0.2, -57.0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			state := SystemState{
				FraudRate:      tc.fraudRate,
				FalsePositives: tc.falsePos,
				Instability:    tc.instability,
			}
			reward := o.CalculateReward(state, tc.previousFraud)
			assert.InDelta(t, tc.expected, reward, 0.001)
		})
	}
}

func TestSuggestAction_HighFraud(t *testing.T) {
	o := &PolicyOptimizer{}
	state := SystemState{FraudRate: 0.5}
	action := o.SuggestAction(state)
	assert.Contains(t, action, "Tighten")
}

func TestSuggestAction_ModerateFraud(t *testing.T) {
	o := &PolicyOptimizer{}
	state := SystemState{FraudRate: 0.4}
	action := o.SuggestAction(state)
	assert.Contains(t, action, "Maintain")
}

func TestSuggestAction_LowFraud(t *testing.T) {
	o := &PolicyOptimizer{}
	state := SystemState{FraudRate: 0.05}
	action := o.SuggestAction(state)
	assert.Contains(t, action, "Maintain")
}

func TestSuggestAction_AtThreshold(t *testing.T) {
	o := &PolicyOptimizer{}
	state := SystemState{FraudRate: 0.4}
	action := o.SuggestAction(state)
	assert.Contains(t, action, "Maintain")

	state.FraudRate = 0.4001
	action = o.SuggestAction(state)
	assert.Contains(t, action, "Tighten")
}

func TestSuggestAction_TableDriven(t *testing.T) {
	o := &PolicyOptimizer{}

	tests := []struct {
		name      string
		fraudRate float64
		contains  string
	}{
		{"critical fraud", 0.9, "Tighten"},
		{"high fraud", 0.75, "Tighten"},
		{"elevated fraud", 0.41, "Tighten"},
		{"boundary fraud", 0.4, "Maintain"},
		{"low fraud", 0.1, "Maintain"},
		{"zero fraud", 0.0, "Maintain"},
		{"negative fraud", -0.1, "Maintain"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			state := SystemState{FraudRate: tc.fraudRate}
			action := o.SuggestAction(state)
			assert.Contains(t, action, tc.contains)
		})
	}
}
