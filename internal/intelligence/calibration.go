package intelligence

import (
	"math"
)

type CalibratedSignals struct {
	MLRisk         float64
	GraphRisk      float64
	BehaviorRisk   float64
	Confidence     float64
}

// Platt Scaling calibration: maps raw scores to true probabilities
func Calibrate(raw float64) float64 {
	// Formula: 1 / (1 + exp(A*raw + B))
	// Example A=-12, B=6 for sigmoid centering at 0.5
	return 1 / (1 + math.Exp(-12*(raw-0.5)))
}

func FusePrecision(s CalibratedSignals) float64 {
	// Calibrate individual signals
	cML := Calibrate(s.MLRisk)
	cGraph := Calibrate(s.GraphRisk)
	cBehavior := Calibrate(s.BehaviorRisk)

	// Weighted fusion
	final := (cML * 0.5) + (cGraph * 0.3) + (cBehavior * 0.2)

	// Confidence Weighting: lower trust in results if confidence is low
	if s.Confidence < 0.6 {
		final *= 0.7
	}

	return final
}
