package soc

// No imports needed

type Severity string

const (
	SeverityLow      Severity = "LOW"
	SeverityMedium   Severity = "MEDIUM"
	SeverityHigh     Severity = "HIGH"
	SeverityCritical Severity = "CRITICAL"
)

type SeverityEngine struct{}

func NewSeverityEngine() *SeverityEngine {
	return &SeverityEngine{}
}

func (e *SeverityEngine) Classify(alert map[string]interface{}) Severity {
	score, _ := alert["fraudScore"].(int)
	if score == 0 {
		// Try parsing from generic metadata if present
		if s, ok := alert["score"].(float64); ok {
			score = int(s)
		}
	}

	if score >= 90 {
		return SeverityCritical
	}
	if score >= 70 {
		return SeverityHigh
	}
	if score >= 40 {
		return SeverityMedium
	}

	// Check specific incident types
	status, _ := alert["status"].(string)
	if status == "CRITICAL" {
		return SeverityCritical
	}

	return SeverityLow
}

func (e *SeverityEngine) GetEscalationPath(sev Severity) []string {
	switch sev {
	case SeverityCritical:
		return []string{"pagerduty", "slack", "audit"}
	case SeverityHigh:
		return []string{"slack", "audit"}
	default:
		return []string{"audit"}
	}
}
