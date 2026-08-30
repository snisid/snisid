package verified

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSafe_BelowThreshold(t *testing.T) {
	assert.True(t, IsSafe(5, 10))
}

func TestIsSafe_AtThreshold(t *testing.T) {
	assert.True(t, IsSafe(10, 10))
}

func TestIsSafe_AboveThreshold(t *testing.T) {
	assert.False(t, IsSafe(15, 10))
}

func TestIsSafe_ZeroRisk(t *testing.T) {
	assert.True(t, IsSafe(0, 10))
}

func TestIsSafe_ZeroThreshold(t *testing.T) {
	assert.True(t, IsSafe(0, 0))
	assert.False(t, IsSafe(1, 0))
}

func TestIsSafe_NegativeValues(t *testing.T) {
	assert.True(t, IsSafe(-5, 10))
	assert.False(t, IsSafe(5, -10))
}

func TestIsSafe_LargeValues(t *testing.T) {
	assert.True(t, IsSafe(1000000, 1000000))
	assert.False(t, IsSafe(1000001, 1000000))
}

func TestValidatePolicyInvariant_RiskBelowThresholdPolicyAllow(t *testing.T) {
	assert.True(t, ValidatePolicyInvariant(5, 10, "ALLOW"))
}

func TestValidatePolicyInvariant_RiskBelowThresholdPolicyDeny(t *testing.T) {
	assert.False(t, ValidatePolicyInvariant(5, 10, "DENY"))
}

func TestValidatePolicyInvariant_RiskAtThresholdPolicyAllow(t *testing.T) {
	assert.True(t, ValidatePolicyInvariant(10, 10, "ALLOW"))
}

func TestValidatePolicyInvariant_RiskAboveThresholdPolicyAllow(t *testing.T) {
	assert.True(t, ValidatePolicyInvariant(15, 10, "ALLOW"))
}

func TestValidatePolicyInvariant_RiskAboveThresholdPolicyDeny(t *testing.T) {
	assert.True(t, ValidatePolicyInvariant(15, 10, "DENY"))
}

func TestValidatePolicyInvariant_ZeroValues(t *testing.T) {
	assert.True(t, ValidatePolicyInvariant(0, 0, "ALLOW"))
	assert.False(t, ValidatePolicyInvariant(0, 0, "DENY"))
}

func TestValidatePolicyInvariant_NegativeRisk(t *testing.T) {
	assert.False(t, ValidatePolicyInvariant(-1, 0, "DENY"))
}
