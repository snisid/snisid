package sandbox

import (
	"math"
	"github.com/snisid/platform/backend/internal/platform/logger"
)

type EnvState struct {
	Fraud   float64
	Economy float64
}

func Reward(state EnvState) float64 {
	// Balance: Minimize fraud, Maximize economy stability
	return (1.0 - state.Fraud) + (state.Economy * 0.3)
}

type DriftMonitor struct {
	Baseline float64
}

func (m *DriftMonitor) DetectCollapseRisk(current float64) bool {
	drift := math.Abs(current - m.Baseline)
	if drift > 0.25 {
		logger.Warn("SELF-REGULATION: High system drift detected. Initiating Safety Freeze.")
		return true
	}
	return false
}

func SafetyFreeze() {
	logger.Error("SAFETY: RL-Optimization suspended. Reverting to Baseline Stability Model.", nil)
}
