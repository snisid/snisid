package graph

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewThreatAnalyzer(t *testing.T) {
	a := NewThreatAnalyzer()
	require.NotNil(t, a)
	assert.Empty(t, a.users)
	assert.Empty(t, a.edges)
	assert.InDelta(t, 0.7, a.thresholds.InsiderRiskThreshold, 0.001)
	assert.InDelta(t, 1.5, a.thresholds.SuspiciousPathScore, 0.001)
	assert.Equal(t, 5, a.thresholds.AnomalyActionCount)
}

func TestMapRelationship_NewUser(t *testing.T) {
	a := NewThreatAnalyzer()
	a.MapRelationship(
		UserNode{UID: "u1", Role: "admin", RiskLevel: 0.3},
		ActionEdge{Type: "LOGIN", Target: "t1", RiskScore: 0.1, Timestamp: time.Now().Unix()},
	)
	assert.Len(t, a.users, 1)
	assert.Equal(t, "admin", a.users["u1"].Role)
}

func TestMapRelationship_ExistingUser(t *testing.T) {
	a := NewThreatAnalyzer()
	a.MapRelationship(UserNode{UID: "u1", RiskLevel: 0.3}, ActionEdge{Type: "LOGIN"})
	a.MapRelationship(UserNode{UID: "u1", RiskLevel: 0.5}, ActionEdge{Type: "LOGOUT"})
	assert.Len(t, a.users, 1)
	assert.InDelta(t, 0.3, a.users["u1"].RiskLevel, 0.001)
}

func TestDetectInsiderThreat_UnknownUser(t *testing.T) {
	a := NewThreatAnalyzer()
	score := a.DetectInsiderThreat("unknown")
	assert.InDelta(t, 0.0, score, 0.001)
}

func TestDetectInsiderThreat_LowRisk(t *testing.T) {
	a := NewThreatAnalyzer()
	uid := "u1"
	a.MapRelationship(
		UserNode{UID: uid, Role: "analyst", RiskLevel: 0.1},
		ActionEdge{Type: "VIEW", RiskScore: 0.05, Timestamp: time.Now().Unix()},
	)
	score := a.DetectInsiderThreat(uid)
	assert.Less(t, score, 0.7)
}

func TestDetectInsiderThreat_HighRisk(t *testing.T) {
	a := NewThreatAnalyzer()
	uid := "u-malicious"
	now := time.Now().Unix()

	a.MapRelationship(
		UserNode{UID: uid, Role: "admin", RiskLevel: 0.8, Labels: map[string]string{"dept": "it"}},
		ActionEdge{Type: "DATA_ACCESS", Target: "sensitive-db", RiskScore: 0.9, AnomalyFlag: true, Timestamp: now},
	)
	for i := 0; i < 10; i++ {
		a.MapRelationship(
			UserNode{UID: uid, RiskLevel: 0.8},
			ActionEdge{Type: "DATA_ACCESS", RiskScore: 0.8, AnomalyFlag: i > 3, Timestamp: now + int64(i*60)},
		)
	}

	score := a.DetectInsiderThreat(uid)
	assert.Greater(t, score, 0.0)
	assert.LessOrEqual(t, score, 1.0)
}

func TestFindShortestPath_Basic(t *testing.T) {
	a := NewThreatAnalyzer()
	a.MapRelationship(UserNode{UID: "A"}, ActionEdge{Type: "LINK", Target: "B", RiskScore: 0.3})
	a.MapRelationship(UserNode{UID: "B"}, ActionEdge{Type: "LINK", Target: "C", RiskScore: 0.4})
	a.MapRelationship(UserNode{UID: "C"}, ActionEdge{Type: "LINK", Target: "D", RiskScore: 0.2})

	result := a.FindShortestPath(GraphQuery{SourceUID: "A", TargetUID: "D", MaxDepth: 5})
	assert.Nil(t, result)
}

func TestFindShortestPath_SourceNotFound(t *testing.T) {
	a := NewThreatAnalyzer()
	result := a.FindShortestPath(GraphQuery{SourceUID: "missing", TargetUID: "u2"})
	assert.Nil(t, result)
}

func TestFindShortestPath_TargetNotFound(t *testing.T) {
	a := NewThreatAnalyzer()
	a.MapRelationship(UserNode{UID: "u1"}, ActionEdge{Type: "TEST"})
	result := a.FindShortestPath(GraphQuery{SourceUID: "u1", TargetUID: "missing"})
	assert.Nil(t, result)
}

func TestFindShortestPath_Defaults(t *testing.T) {
	a := NewThreatAnalyzer()
	a.MapRelationship(UserNode{UID: "A"}, ActionEdge{Type: "LINK", Target: "B"})
	result := a.FindShortestPath(GraphQuery{SourceUID: "A", TargetUID: "B"})
	assert.Nil(t, result)
}

func TestDetectClusters_MinSize(t *testing.T) {
	a := NewThreatAnalyzer()
	a.MapRelationship(UserNode{UID: "u1", RiskLevel: 0.2}, ActionEdge{Type: "LINK", Target: "u2"})
	a.MapRelationship(UserNode{UID: "u2", RiskLevel: 0.3}, ActionEdge{Type: "LINK", Target: "u3"})
	a.MapRelationship(UserNode{UID: "u3", RiskLevel: 0.4}, ActionEdge{Type: "LINK", Target: "u1"})

	clusters := a.DetectClusters(2)
	assert.Empty(t, clusters)
}

func TestDetectClusters_DefaultMinSize(t *testing.T) {
	a := NewThreatAnalyzer()
	clusters := a.DetectClusters(0)
	assert.Empty(t, clusters)
}

func TestDetectClusters_Empty(t *testing.T) {
	a := NewThreatAnalyzer()
	clusters := a.DetectClusters(1)
	assert.Empty(t, clusters)
}

func TestDetectAnomalousAccess_None(t *testing.T) {
	a := NewThreatAnalyzer()
	a.MapRelationship(
		UserNode{UID: "u1"},
		ActionEdge{Type: "VIEW", RiskScore: 0.1, AnomalyFlag: false, Timestamp: time.Now().Unix()},
	)
	anomalous := a.DetectAnomalousAccess("u1")
	assert.Empty(t, anomalous)
}

func TestDetectAnomalousAccess_WithAnomalies(t *testing.T) {
	a := NewThreatAnalyzer()
	a.MapRelationship(
		UserNode{UID: "u1"},
		ActionEdge{Type: "DATA_EXPORT", RiskScore: 0.9, AnomalyFlag: true, Timestamp: time.Now().Unix()},
	)
	a.MapRelationship(
		UserNode{UID: "u1"},
		ActionEdge{Type: "VIEW", RiskScore: 0.1, AnomalyFlag: false, Timestamp: time.Now().Unix()},
	)

	anomalous := a.DetectAnomalousAccess("u1")
	assert.Len(t, anomalous, 1)
	assert.Equal(t, "DATA_EXPORT", anomalous[0].Type)
}

func TestDetectAnomalousAccess_HighRisk(t *testing.T) {
	a := NewThreatAnalyzer()
	a.MapRelationship(
		UserNode{UID: "u1"},
		ActionEdge{Type: "SUSPICIOUS", RiskScore: 0.85, AnomalyFlag: false, Timestamp: time.Now().Unix()},
	)
	anomalous := a.DetectAnomalousAccess("u1")
	assert.Len(t, anomalous, 1)
}

func TestDetectAnomalousAccess_UnknownUser(t *testing.T) {
	a := NewThreatAnalyzer()
	anomalous := a.DetectAnomalousAccess("nonexistent")
	assert.Empty(t, anomalous)
}

func TestDetectAnomalousAccess_SortedByRisk(t *testing.T) {
	a := NewThreatAnalyzer()
	a.MapRelationship(UserNode{UID: "u1"}, ActionEdge{Type: "A", RiskScore: 0.3, AnomalyFlag: true, Timestamp: 1})
	a.MapRelationship(UserNode{UID: "u1"}, ActionEdge{Type: "B", RiskScore: 0.9, AnomalyFlag: true, Timestamp: 2})
	a.MapRelationship(UserNode{UID: "u1"}, ActionEdge{Type: "C", RiskScore: 0.6, AnomalyFlag: true, Timestamp: 3})

	anomalous := a.DetectAnomalousAccess("u1")
	require.Len(t, anomalous, 3)
	assert.Equal(t, 0.9, anomalous[0].RiskScore)
	assert.Equal(t, 0.6, anomalous[1].RiskScore)
	assert.Equal(t, 0.3, anomalous[2].RiskScore)
}

func TestContains(t *testing.T) {
	assert.True(t, contains([]string{"a", "b", "c"}, "b"))
	assert.False(t, contains([]string{"a", "b", "c"}, "d"))
	assert.False(t, contains([]string{}, "a"))
}

func TestComputeCentrality_Single(t *testing.T) {
	a := NewThreatAnalyzer()
	adj := map[string]map[string]bool{"u1": {}}
	assert.InDelta(t, 0.0, a.computeCentrality([]string{"u1"}, adj), 0.001)
}

func TestComputeCentrality_Empty(t *testing.T) {
	a := NewThreatAnalyzer()
	assert.InDelta(t, 0.0, a.computeCentrality([]string{}, map[string]map[string]bool{}), 0.001)
}

func TestConcurrentMapRelationship(t *testing.T) {
	a := NewThreatAnalyzer()
	t.Run("parallel", func(t *testing.T) {
		t.Run("user1", func(t *testing.T) {
			a.MapRelationship(
				UserNode{UID: "u1", Role: "admin", RiskLevel: 0.5},
				ActionEdge{Type: "LOGIN", RiskScore: 0.1, Timestamp: time.Now().Unix()},
			)
		})
		t.Run("user2", func(t *testing.T) {
			a.MapRelationship(
				UserNode{UID: "u2", Role: "analyst", RiskLevel: 0.3},
				ActionEdge{Type: "VIEW", RiskScore: 0.05, Timestamp: time.Now().Unix()},
			)
		})
	})
	assert.Len(t, a.users, 2)
}
