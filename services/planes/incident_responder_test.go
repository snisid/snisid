package incidentresponder

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type IncidentSeverity string

const (
	SeverityLow      IncidentSeverity = "LOW"
	SeverityMedium   IncidentSeverity = "MEDIUM"
	SeverityHigh     IncidentSeverity = "HIGH"
	SeverityCritical IncidentSeverity = "CRITICAL"
)

type PlaybookAction struct {
	Action       string `json:"action"`
	Target       string `json:"target"`
	Duration     string `json:"duration,omitempty"`
	RequireHuman bool   `json:"require_human"`
}

type IncidentPayload struct {
	IncidentID   string           `json:"incident_id"`
	Severity     IncidentSeverity `json:"severity"`
	TargetType   string           `json:"target_type"`
	TargetID     string           `json:"target_id"`
	Reason       string           `json:"reason"`
	SourceSystem string           `json:"source_system"`
}

func selectPlaybook(severity IncidentSeverity, command string) []PlaybookAction {
	switch severity {
	case SeverityCritical:
		return []PlaybookAction{
			{Action: "ISOLATE", Target: "identity", Duration: "24h", RequireHuman: false},
			{Action: "REVOKE_TOKENS", Target: "all", RequireHuman: false},
			{Action: "ALERT_SOC", Target: "human_operator", RequireHuman: false},
		}
	case SeverityHigh:
		return []PlaybookAction{
			{Action: "RESTRICT", Target: "identity", Duration: "4h", RequireHuman: false},
			{Action: "FLAG_REVIEW", Target: "identity", RequireHuman: true},
		}
	case SeverityMedium:
		return []PlaybookAction{
			{Action: "MONITOR", Target: "identity", Duration: "24h", RequireHuman: false},
		}
	default:
		return []PlaybookAction{
			{Action: "LOG_ONLY", Target: "identity", RequireHuman: false},
		}
	}
}

func executePlaybook(playbook []PlaybookAction, payload IncidentPayload) []map[string]interface{} {
	var results []map[string]interface{}
	for _, action := range playbook {
		result := map[string]interface{}{
			"action":   action.Action,
			"target":   action.Target,
			"status":   "executed",
			"duration": action.Duration,
		}
		results = append(results, result)
	}
	return results
}

func TestSelectPlaybook_CriticalSeverity(t *testing.T) {
	playbook := selectPlaybook(SeverityCritical, "INCIDENT_RESPONDER")
	assert.Len(t, playbook, 3)
	assert.Equal(t, "ISOLATE", playbook[0].Action)
	assert.Equal(t, "REVOKE_TOKENS", playbook[1].Action)
	assert.Equal(t, "ALERT_SOC", playbook[2].Action)
	assert.False(t, playbook[0].RequireHuman)
	assert.Equal(t, "24h", playbook[0].Duration)
}

func TestSelectPlaybook_HighSeverity(t *testing.T) {
	playbook := selectPlaybook(SeverityHigh, "INCIDENT_RESPONDER")
	assert.Len(t, playbook, 2)
	assert.Equal(t, "RESTRICT", playbook[0].Action)
	assert.Equal(t, "FLAG_REVIEW", playbook[1].Action)
	assert.True(t, playbook[1].RequireHuman)
	assert.Equal(t, "4h", playbook[0].Duration)
}

func TestSelectPlaybook_MediumSeverity(t *testing.T) {
	playbook := selectPlaybook(SeverityMedium, "INCIDENT_RESPONDER")
	assert.Len(t, playbook, 1)
	assert.Equal(t, "MONITOR", playbook[0].Action)
	assert.Equal(t, "identity", playbook[0].Target)
	assert.Equal(t, "24h", playbook[0].Duration)
}

func TestSelectPlaybook_LowSeverity(t *testing.T) {
	playbook := selectPlaybook(SeverityLow, "INCIDENT_RESPONDER")
	assert.Len(t, playbook, 1)
	assert.Equal(t, "LOG_ONLY", playbook[0].Action)
	assert.False(t, playbook[0].RequireHuman)
}

func TestSelectPlaybook_DefaultSeverity(t *testing.T) {
	playbook := selectPlaybook("UNKNOWN", "INCIDENT_RESPONDER")
	assert.Len(t, playbook, 1)
	assert.Equal(t, "LOG_ONLY", playbook[0].Action)
}

func TestExecutePlaybook_CriticalFlow(t *testing.T) {
	playbook := selectPlaybook(SeverityCritical, "INCIDENT_RESPONDER")
	payload := IncidentPayload{
		IncidentID:   "INC-001",
		Severity:     SeverityCritical,
		TargetType:   "citizen",
		TargetID:     "CIT-00042",
		Reason:       "Synthetic identity detected",
		SourceSystem: "FRAUD_ENGINE",
	}

	results := executePlaybook(playbook, payload)
	assert.Len(t, results, 3)
	for _, r := range results {
		assert.Equal(t, "executed", r["status"])
	}
	assert.Equal(t, "ISOLATE", results[0]["action"])
	assert.Equal(t, "identity", results[0]["target"])
	assert.Equal(t, "24h", results[0]["duration"])
}

func TestExecutePlaybook_HighFlow(t *testing.T) {
	playbook := selectPlaybook(SeverityHigh, "INCIDENT_RESPONDER")
	payload := IncidentPayload{
		IncidentID: "INC-002", Severity: SeverityHigh,
		TargetType: "identity", TargetID: "CIT-00123",
		Reason: "Unusual access pattern", SourceSystem: "SIEM",
	}

	results := executePlaybook(playbook, payload)
	assert.Len(t, results, 2)
	assert.Equal(t, "RESTRICT", results[0]["action"])
	assert.Equal(t, "4h", results[0]["duration"])
}

func TestExecutePlaybook_EmptyPlaybook(t *testing.T) {
	payload := IncidentPayload{IncidentID: "INC-003"}
	results := executePlaybook([]PlaybookAction{}, payload)
	assert.Empty(t, results)
}

func TestIncidentPayload_Serialization(t *testing.T) {
	payload := IncidentPayload{
		IncidentID: "INC-SER-001", Severity: SeverityCritical,
		TargetType: "citizen", TargetID: "CIT-999",
		Reason: "Test serialization", SourceSystem: "test",
	}

	data, err := json.Marshal(payload)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "INC-SER-001")
	assert.Contains(t, string(data), "CRITICAL")

	var decoded IncidentPayload
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, payload.IncidentID, decoded.IncidentID)
	assert.Equal(t, payload.Severity, decoded.Severity)
}

func TestPlaybookAction_Serialization(t *testing.T) {
	action := PlaybookAction{
		Action: "ISOLATE", Target: "identity",
		Duration: "24h", RequireHuman: false,
	}

	data, err := json.Marshal(action)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "ISOLATE")
	assert.Contains(t, string(data), "24h")
	assert.Contains(t, string(data), "false")

	var decoded PlaybookAction
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, action.Action, decoded.Action)
	assert.Equal(t, action.Duration, decoded.Duration)
}

func TestSeverityLevel_Comparison(t *testing.T) {
	assert.NotEqual(t, SeverityLow, SeverityHigh)
	assert.NotEqual(t, SeverityMedium, SeverityCritical)
	assert.Equal(t, IncidentSeverity("LOW"), SeverityLow)
	assert.Equal(t, IncidentSeverity("HIGH"), SeverityHigh)
}
