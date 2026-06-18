package oversight

import (
	"context"
	"fmt"

	"github.com/snisid/platform/internal/platform/logger"
)

type Decision struct {
	RiskScore  float64
	Signals    []string
	Confidence float64
}

type OversightAI struct {
	PlatformID string
}

func (o *OversightAI) ValidateDecision(d Decision) (bool, string) {
	logger.Info(context.Background(), "OVERSIGHT: Validating automated fraud decision...")

	if d.Confidence < 0.7 {
		return false, "REJECTED: LOW_CONFIDENCE"
	}

	if len(d.Signals) < 2 {
		return false, "REJECTED: INSUFFICIENT_EVIDENCE"
	}

	logger.Info(context.Background(), "OVERSIGHT: Decision APPROVED for investigation.")
	return true, "APPROVED"
}

func (o *OversightAI) AuditLog(d Decision, status string, reason string) {
	fmt.Printf("OVERSIGHT_AUDIT: Decision %v -> %s (%s)\n", d, status, reason)
	// Write to immutable Kafka ledger
}
