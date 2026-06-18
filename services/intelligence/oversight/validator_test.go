package oversight

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDecision_Approved(t *testing.T) {
	o := &OversightAI{PlatformID: "SNISID-PROD"}
	ok, reason := o.ValidateDecision(Decision{
		RiskScore:  0.85,
		Signals:    []string{"velocity_anomaly", "new_device", "geo_mismatch"},
		Confidence: 0.92,
	})
	assert.True(t, ok)
	assert.Equal(t, "APPROVED", reason)
}

func TestValidateDecision_LowConfidence(t *testing.T) {
	o := &OversightAI{}
	ok, reason := o.ValidateDecision(Decision{
		RiskScore:  0.5,
		Signals:    []string{"signal1", "signal2"},
		Confidence: 0.5,
	})
	assert.False(t, ok)
	assert.Contains(t, reason, "LOW_CONFIDENCE")
}

func TestValidateDecision_InsufficientEvidence(t *testing.T) {
	o := &OversightAI{}
	ok, reason := o.ValidateDecision(Decision{
		RiskScore:  0.8,
		Signals:    []string{"single_signal"},
		Confidence: 0.95,
	})
	assert.False(t, ok)
	assert.Contains(t, reason, "INSUFFICIENT_EVIDENCE")
}

func TestValidateDecision_NoSignals(t *testing.T) {
	o := &OversightAI{}
	ok, reason := o.ValidateDecision(Decision{
		RiskScore:  0.5,
		Signals:    []string{},
		Confidence: 0.9,
	})
	assert.False(t, ok)
	assert.Contains(t, reason, "INSUFFICIENT_EVIDENCE")
}

func TestValidateDecision_NilSignals(t *testing.T) {
	o := &OversightAI{}
	ok, reason := o.ValidateDecision(Decision{
		RiskScore:  0.5,
		Confidence: 0.9,
	})
	assert.False(t, ok)
	assert.Contains(t, reason, "INSUFFICIENT_EVIDENCE")
}

func TestValidateDecision_HighRiskHighConfidence(t *testing.T) {
	o := &OversightAI{}
	ok, reason := o.ValidateDecision(Decision{
		RiskScore:  0.95,
		Signals:    []string{"signal_a", "signal_b", "signal_c"},
		Confidence: 0.99,
	})
	assert.True(t, ok)
	assert.Equal(t, "APPROVED", reason)
}

func TestValidateDecision_BoundaryConfidence(t *testing.T) {
	o := &OversightAI{}

	tests := []struct {
		name       string
		confidence float64
		approved   bool
	}{
		{"below_threshold", 0.69, false},
		{"at_threshold", 0.7, true},
		{"above_threshold", 0.71, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, _ := o.ValidateDecision(Decision{
				RiskScore:  0.5,
				Signals:    []string{"s1", "s2"},
				Confidence: tt.confidence,
			})
			assert.Equal(t, tt.approved, ok)
		})
	}
}

func TestAuditLog(t *testing.T) {
	o := &OversightAI{}
	o.AuditLog(Decision{RiskScore: 0.8, Signals: []string{"s1"}, Confidence: 0.9}, "APPROVED", "all good")
}

func TestPlatformID(t *testing.T) {
	o := &OversightAI{PlatformID: "SNISID-HTI"}
	assert.Equal(t, "SNISID-HTI", o.PlatformID)
}
