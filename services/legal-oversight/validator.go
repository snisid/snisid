package legaloversight

import (
	"context"
	"fmt"

	"github.com/snisid/platform/internal/platform/logger"
)

type Decision struct {
	SubjectID  string
	RiskScore  float64
	Signals    []string
	Confidence float64
}

func ValidateDecision(d Decision) (bool, string) {
	logger.Info(context.Background(), fmt.Sprintf("LEGAL-AI: Validating decision for subject %s", d.SubjectID))

	if d.Confidence < 0.7 {
		return false, "LOW_CONFIDENCE"
	}

	if len(d.Signals) < 2 {
		return false, "INSUFFICIENT_EVIDENCE"
	}

	if d.RiskScore > 0.9 && d.Confidence < 0.85 {
		return false, "HIGH_RISK_REQUIRE_HIGHER_CONFIDENCE"
	}

	return true, "APPROVED"
}
