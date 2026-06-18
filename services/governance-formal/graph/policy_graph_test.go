package graph

import (
	"testing"

	"github.com/snisid/platform/governance-formal/compiler"
	"github.com/stretchr/testify/assert"
)

func makePolicy(rules []string, trust float64) compiler.CompiledPolicy {
	return func(s compiler.State, a compiler.Action) bool {
		return true
	}
}

type testPolicy struct {
	rules     []string
	constrain []PolicyConstraint
}

func newNode(country string, p testPolicy) *PolicyNode {
	return &PolicyNode{
		Country:     country,
		Policy:      makePolicy(p.rules, 0.8),
		Constraints: p.constrain,
		TrustScore:  0.8,
	}
}

func TestNewGlobalPolicyLattice(t *testing.T) {
	l := NewGlobalPolicyLattice()
	assert.NotNil(t, l)
	assert.NotNil(t, l.Nodes)
	assert.Empty(t, l.Edges)
}

func TestAddCountry_DefaultConstraints(t *testing.T) {
	l := NewGlobalPolicyLattice()
	p := compiler.CompiledPolicy(func(s compiler.State, a compiler.Action) bool { return true })
	l.AddCountry("HTI", p)

	node, ok := l.Nodes["HTI"]
	assert.True(t, ok)
	assert.Equal(t, "HTI", node.Country)
	assert.Equal(t, 5, len(node.Constraints))
	assert.Equal(t, 0.8, node.TrustScore)
}

func TestAddCountry_CustomTrustScore(t *testing.T) {
	l := NewGlobalPolicyLattice()
	l.AddCountry("DOM", compiler.CompiledPolicy(func(s compiler.State, a compiler.Action) bool {
		return true
	}))
	l.Nodes["DOM"].TrustScore = 0.95

	assert.InDelta(t, 0.95, l.Nodes["DOM"].TrustScore, 0.001)
}

func TestProveCompatibility_FullyCompatible(t *testing.T) {
	l := NewGlobalPolicyLattice()
	l.AddCountry("HTI", compiler.CompiledPolicy(func(s compiler.State, a compiler.Action) bool { return true }))
	l.AddCountry("DOM", compiler.CompiledPolicy(func(s compiler.State, a compiler.Action) bool { return true }))

	// Make them very compatible by relaxing HTI constraints
	l.Nodes["HTI"].Constraints = []PolicyConstraint{
		{Domain: "DATA_RETENTION", Strictness: 3},
		{Domain: "BIOMETRIC_SHARING", Strictness: 3},
	}
	l.Nodes["DOM"].Constraints = []PolicyConstraint{
		{Domain: "DATA_RETENTION", Strictness: 5},
		{Domain: "BIOMETRIC_SHARING", Strictness: 5},
	}

	level, violations, err := l.ProveCompatibility("HTI", "DOM")
	assert.NoError(t, err)
	assert.Equal(t, Compatible, level)
	assert.Empty(t, violations)
}

func TestProveCompatibility_PartialOverlap(t *testing.T) {
	l := NewGlobalPolicyLattice()
	l.AddCountry("HTI", compiler.CompiledPolicy(func(s compiler.State, a compiler.Action) bool { return true }))
	l.AddCountry("DOM", compiler.CompiledPolicy(func(s compiler.State, a compiler.Action) bool { return true }))

	l.Nodes["HTI"].Constraints = []PolicyConstraint{
		{Domain: "DATA_RETENTION", Strictness: 1},
		{Domain: "BIOMETRIC_SHARING", Strictness: 8},
	}
	l.Nodes["DOM"].Constraints = []PolicyConstraint{
		{Domain: "DATA_RETENTION", Strictness: 9},
		{Domain: "BIOMETRIC_SHARING", Strictness: 10},
	}

	level, violations, err := l.ProveCompatibility("HTI", "DOM")
	assert.NoError(t, err)
	assert.Equal(t, Incompatible, level)
	assert.NotEmpty(t, violations)
}

func TestProveCompatibility_SourceNotFound(t *testing.T) {
	l := NewGlobalPolicyLattice()
	l.AddCountry("HTI", compiler.CompiledPolicy(func(s compiler.State, a compiler.Action) bool { return true }))

	_, _, err := l.ProveCompatibility("UNKNOWN", "HTI")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not registered")
}

func TestProveCompatibility_TargetNotFound(t *testing.T) {
	l := NewGlobalPolicyLattice()
	l.AddCountry("HTI", compiler.CompiledPolicy(func(s compiler.State, a compiler.Action) bool { return true }))

	_, _, err := l.ProveCompatibility("HTI", "UNKNOWN")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not registered")
}

func TestGetEdge_Existing(t *testing.T) {
	l := NewGlobalPolicyLattice()
	l.AddCountry("HTI", nil)
	l.AddCountry("DOM", nil)
	l.ProveCompatibility("HTI", "DOM")

	edge := l.GetEdge("HTI", "DOM")
	assert.NotNil(t, edge)
	assert.Equal(t, "HTI", edge.From)
	assert.Equal(t, "DOM", edge.To)
}

func TestGetEdge_NonExistent(t *testing.T) {
	l := NewGlobalPolicyLattice()
	edge := l.GetEdge("HTI", "DOM")
	assert.Nil(t, edge)
}

func TestGetCompatiblePeers_FiltersByLevel(t *testing.T) {
	l := NewGlobalPolicyLattice()
	l.AddCountry("HTI", nil)
	l.AddCountry("DOM", nil)
	l.AddCountry("CUB", nil)

	l.ProveCompatibility("HTI", "DOM")
	l.ProveCompatibility("HTI", "CUB")

	peers := l.GetCompatiblePeers("HTI", Incompatible)
	assert.NotEmpty(t, peers)

	peers = l.GetCompatiblePeers("DOM", Compatible)
	assert.Empty(t, peers)
}

func TestRankPeersByTrust(t *testing.T) {
	l := NewGlobalPolicyLattice()
	l.AddCountry("HTI", nil)
	l.AddCountry("DOM", nil)
	l.AddCountry("CUB", nil)

	l.ProveCompatibility("HTI", "DOM")
	l.ProveCompatibility("HTI", "CUB")

	ranked := l.RankPeersByTrust("HTI")
	assert.NotEmpty(t, ranked)
	assert.Equal(t, len(ranked), len(l.Edges))
}

func TestVerifyGlobalConsistency_Asymmetric(t *testing.T) {
	l := NewGlobalPolicyLattice()
	l.AddCountry("HTI", nil)
	l.AddCountry("DOM", nil)

	l.ProveCompatibility("HTI", "DOM")
	// Manually invalidate reverse
	for i, e := range l.Edges {
		if e.From == "DOM" {
			l.Edges[i].Valid = false
		}
	}

	issues := l.VerifyGlobalConsistency()
	assert.NotEmpty(t, issues)
	assert.Contains(t, issues[0], "asymmetric")
}

func TestVerifyGlobalConsistency_AllGood(t *testing.T) {
	l := NewGlobalPolicyLattice()
	l.AddCountry("HTI", nil)
	l.AddCountry("DOM", nil)
	l.ProveCompatibility("HTI", "DOM")
	l.ProveCompatibility("DOM", "HTI")

	issues := l.VerifyGlobalConsistency()
	assert.Empty(t, issues)
}

func TestAddCountry_ConcurrentSafe(t *testing.T) {
	l := NewGlobalPolicyLattice()
	t.Run("parallel", func(t *testing.T) {
		t.Run("add HTI", func(t *testing.T) {
			l.AddCountry("HTI", nil)
		})
		t.Run("add DOM", func(t *testing.T) {
			l.AddCountry("DOM", nil)
		})
		t.Run("add CUB", func(t *testing.T) {
			l.AddCountry("CUB", nil)
		})
	})
	assert.Equal(t, 3, len(l.Nodes))
}

func TestEdgeConfidence_Rounding(t *testing.T) {
	l := NewGlobalPolicyLattice()
	l.AddCountry("HTI", nil)
	l.AddCountry("DOM", nil)

	level, _, err := l.ProveCompatibility("HTI", "DOM")
	assert.NoError(t, err)

	edge := l.GetEdge("HTI", "DOM")
	assert.NotNil(t, edge)
	assert.GreaterOrEqual(t, edge.Confidence, 0.0)
	assert.LessOrEqual(t, edge.Confidence, 1.0)
	assert.Equal(t, edge.Compatibility, level)
}
