package causalinference

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEstimateEffect_ExistingFeature(t *testing.T) {
	e := CausalEngine{
		Edges: []CausalEdge{
			{From: "transaction_velocity", To: "fraud_risk", Weight: 0.8},
			{From: "device_trust", To: "fraud_risk", Weight: -0.3},
		},
	}
	effect := e.EstimateEffect("transaction_velocity", 10.0)
	assert.InDelta(t, 8.0, effect, 0.001)
}

func TestEstimateEffect_NegativeWeight(t *testing.T) {
	e := CausalEngine{
		Edges: []CausalEdge{
			{From: "device_trust", To: "fraud_risk", Weight: -0.3},
		},
	}
	effect := e.EstimateEffect("device_trust", 100.0)
	assert.InDelta(t, -30.0, effect, 0.001)
}

func TestEstimateEffect_ZeroDelta(t *testing.T) {
	e := CausalEngine{
		Edges: []CausalEdge{
			{From: "feature_a", To: "outcome", Weight: 0.5},
		},
	}
	effect := e.EstimateEffect("feature_a", 0)
	assert.InDelta(t, 0, effect, 0.001)
}

func TestEstimateEffect_FeatureNotFound(t *testing.T) {
	e := CausalEngine{
		Edges: []CausalEdge{
			{From: "feature_a", To: "outcome", Weight: 0.5},
		},
	}
	effect := e.EstimateEffect("unknown_feature", 100.0)
	assert.InDelta(t, 0, effect, 0.001)
}

func TestEstimateEffect_EmptyEdges(t *testing.T) {
	e := CausalEngine{Edges: nil}
	effect := e.EstimateEffect("anything", 50.0)
	assert.InDelta(t, 0, effect, 0.001)
}

func TestEstimateEffect_MultipleEdgesDifferentFrom(t *testing.T) {
	e := CausalEngine{
		Edges: []CausalEdge{
			{From: "velocity", To: "risk", Weight: 0.7},
			{From: "amount", To: "risk", Weight: 0.3},
			{From: "velocity", To: "flag", Weight: 0.2},
		},
	}
	effect := e.EstimateEffect("velocity", 10.0)
	assert.InDelta(t, 7.0, effect, 0.001)
}

func TestEstimateEffect_LargeDelta(t *testing.T) {
	e := CausalEngine{
		Edges: []CausalEdge{
			{From: "feature_x", To: "outcome", Weight: 2.5},
		},
	}
	effect := e.EstimateEffect("feature_x", 1000.0)
	assert.InDelta(t, 2500.0, effect, 0.001)
}

func TestRecommendIntervention(t *testing.T) {
	e := CausalEngine{}
	recommendation := e.RecommendIntervention("CIT-123456")
	assert.Equal(t, "REDUCE_TRANSACTION_LIMIT", recommendation)
}

func TestRecommendIntervention_AnySubject(t *testing.T) {
	e := CausalEngine{}
	r1 := e.RecommendIntervention("TAX-789")
	r2 := e.RecommendIntervention("BIO-001")
	assert.Equal(t, r1, r2)
}

func TestCausalEdge_Values(t *testing.T) {
	edge := CausalEdge{
		From:   "a",
		To:     "b",
		Weight: 0.75,
	}
	assert.Equal(t, "a", edge.From)
	assert.Equal(t, "b", edge.To)
	assert.Equal(t, 0.75, edge.Weight)
}
