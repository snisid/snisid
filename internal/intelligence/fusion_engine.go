package intelligence

import (
	"fmt"
	"github.com/snisid/platform/backend/internal/platform/logger"
)

type Signals struct {
	MLScore        float64
	GraphScore     float64
	BiometricMatch bool
	RuleScore      float64
}

type FinalDecision struct {
	CitizenID   string
	RiskScore   float64
	RiskLevel   string
	Explanation []string
}

func EvaluateRisk(citizenID string, s Signals) FinalDecision {
	logger.Info(fmt.Sprintf("NEXUS-INTEL: Fusing signals for citizen %s", citizenID))
	
	explanation := []string{}
	riskScore := 0.0
	riskLevel := "LOW"

	// Priority Signals
	if !s.BiometricMatch {
		explanation = append(explanation, "critical: biometric mismatch")
		return FinalDecision{citizenID, 0.95, "CRITICAL", explanation}
	}

	if s.GraphScore > 0.8 {
		explanation = append(explanation, "high: graph cluster anomaly")
	}

	// Weighted Fusion
	riskScore = (s.MLScore * 0.5) + (s.GraphScore * 0.3) + (s.RuleScore * 0.2)

	if riskScore > 0.85 {
		riskLevel = "CRITICAL"
	} else if riskScore > 0.7 {
		riskLevel = "HIGH"
	} else if riskScore > 0.4 {
		riskLevel = "MEDIUM"
	}

	return FinalDecision{
		CitizenID:   citizenID,
		RiskScore:   riskScore,
		RiskLevel:   riskLevel,
		Explanation: explanation,
	}
}
