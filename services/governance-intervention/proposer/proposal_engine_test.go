package proposer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPropose_HighRisk(t *testing.T) {
	e := &ProposalEngine{}
	i := e.Propose("ACCT-001", 0.95, []string{"suspicious_transaction", "high_velocity"})

	assert.Equal(t, "HIGH_RISK", i.Type)
	assert.Equal(t, "FREEZE_ACCOUNT", i.Action)
	assert.Equal(t, StatusProposed, i.Status)
	assert.Greater(t, i.CreatedAt, int64(0))
}

func TestPropose_MediumRisk(t *testing.T) {
	e := &ProposalEngine{}
	i := e.Propose("ACCT-002", 0.75, []string{"unusual_pattern"})

	assert.Equal(t, "MEDIUM_RISK", i.Type)
	assert.Equal(t, "INVESTIGATE", i.Action)
}

func TestPropose_LowRisk(t *testing.T) {
	e := &ProposalEngine{}
	i := e.Propose("ACCT-003", 0.3, []string{})

	assert.Equal(t, "LOW_RISK", i.Type)
	assert.Equal(t, "MONITOR", i.Action)
}

func TestPropose_BoundaryRiskScores(t *testing.T) {
	e := &ProposalEngine{}

	tests := []struct {
		name       string
		score      float64
		wantType   string
		wantAction string
	}{
		{"exactly_high_threshold", 0.9, "HIGH_RISK", "FREEZE_ACCOUNT"},
		{"just_above_medium", 0.71, "MEDIUM_RISK", "INVESTIGATE"},
		{"exactly_medium_threshold", 0.7, "MEDIUM_RISK", "INVESTIGATE"},
		{"just_below_medium", 0.69, "LOW_RISK", "MONITOR"},
		{"max_risk", 1.0, "HIGH_RISK", "FREEZE_ACCOUNT"},
		{"zero_risk", 0.0, "LOW_RISK", "MONITOR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := e.Propose("ACCT-TEST", tt.score, []string{"test"})
			assert.Equal(t, tt.wantType, i.Type)
			assert.Equal(t, tt.wantAction, i.Action)
		})
	}
}

func TestPropose_IDUniqueness(t *testing.T) {
	e := &ProposalEngine{}
	i1 := e.Propose("ACCT-001", 0.5, nil)
	i2 := e.Propose("ACCT-001", 0.5, nil)

	assert.NotEqual(t, i1.ID, i2.ID)
}

func TestPropose_TargetPreserved(t *testing.T) {
	e := &ProposalEngine{}
	i := e.Propose("ACCT-999", 0.85, []string{"alert_triggered"})

	assert.Equal(t, "ACCT-999", i.Target)
	assert.Equal(t, []string{"alert_triggered"}, i.Reason)
}

func TestPropose_EmptyReasons(t *testing.T) {
	e := &ProposalEngine{}
	i := e.Propose("ACCT-001", 0.5, []string{})

	assert.Empty(t, i.Reason)
	assert.NotEmpty(t, i.ID)
}

func TestPropose_NilReasons(t *testing.T) {
	e := &ProposalEngine{}
	i := e.Propose("ACCT-001", 0.5, nil)

	assert.Nil(t, i.Reason)
}
