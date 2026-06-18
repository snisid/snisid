package selfhealing_test

import (
	"testing"
	"time"

	"github.com/snisid/platform/services/governance-os/self-healing"
	"github.com/stretchr/testify/assert"
)

func TestNewMeshHealer(t *testing.T) {
	h := selfhealing.NewMeshHealer("mesh-1", 3)
	assert.NotNil(t, h)
	assert.Equal(t, "mesh-1", h.ID)
}

func TestMonitorTrust_HealthyNode(t *testing.T) {
	h := selfhealing.NewMeshHealer("mesh-1", 2)
	nodes := []selfhealing.NodeState{
		{NodeID: "node-1", TrustScore: 0.95, Status: "ACTIVE", Role: selfhealing.RoleValidator, FailureCount: 0, LastSeen: time.Now()},
		{NodeID: "node-2", TrustScore: 0.88, Status: "ACTIVE", Role: selfhealing.RoleGateway, FailureCount: 1, LastSeen: time.Now()},
	}
	h.MonitorTrust(nodes)

	log := h.GetIsolationLog()
	assert.Empty(t, log, "healthy nodes should not be isolated")
}

func TestMonitorTrust_IsolateLowTrust(t *testing.T) {
	h := selfhealing.NewMeshHealer("mesh-1", 2)
	nodes := []selfhealing.NodeState{
		{NodeID: "node-bad", TrustScore: 0.2, Status: "ACTIVE", Role: selfhealing.RoleExecutor, FailureCount: 0, LastSeen: time.Now()},
		{NodeID: "node-good", TrustScore: 0.95, Status: "ACTIVE", Role: selfhealing.RoleValidator, FailureCount: 0, LastSeen: time.Now()},
	}
	h.MonitorTrust(nodes)

	log := h.GetIsolationLog()
	assert.Len(t, log, 1)
	assert.Equal(t, "node-bad", log[0].NodeID)
	assert.Equal(t, "ISOLATE", log[0].Action)
	assert.Contains(t, log[0].Reason, "trust score below 0.3")
}

func TestMonitorTrust_IsolateDegradedTrust(t *testing.T) {
	h := selfhealing.NewMeshHealer("mesh-1", 2)
	nodes := []selfhealing.NodeState{
		{NodeID: "node-deg", TrustScore: 0.4, Status: "ACTIVE", Role: selfhealing.RoleGateway, FailureCount: 2, LastSeen: time.Now()},
	}
	h.MonitorTrust(nodes)

	log := h.GetIsolationLog()
	assert.Len(t, log, 1)
	assert.Equal(t, "node-deg", log[0].NodeID)
	assert.Contains(t, log[0].Reason, "trust score below 0.5")
}

func TestMonitorTrust_IsolateExcessiveFailures(t *testing.T) {
	h := selfhealing.NewMeshHealer("mesh-1", 2)
	nodes := []selfhealing.NodeState{
		{NodeID: "node-fail", TrustScore: 0.9, Status: "ACTIVE", Role: selfhealing.RoleExecutor, FailureCount: 15, LastSeen: time.Now()},
	}
	h.MonitorTrust(nodes)

	log := h.GetIsolationLog()
	assert.Len(t, log, 1)
	assert.Contains(t, log[0].Reason, "excessive failure count")
}

func TestTriggerRebuild_CreatesPlan(t *testing.T) {
	h := selfhealing.NewMeshHealer("mesh-1", 3)
	h.TriggerRebuild("node-to-rebuild")

	plan := h.GetRebuildPlan("node-to-rebuild")
	assert.NotNil(t, plan)
	assert.False(t, plan.Completed)
	assert.Contains(t, plan.Steps, "drain_node")
	assert.Contains(t, plan.Steps, "reissue_spiffe_cert")
	assert.Contains(t, plan.Steps, "rejoin_mesh")
}

func TestCompleteRebuild(t *testing.T) {
	h := selfhealing.NewMeshHealer("mesh-1", 3)
	h.TriggerRebuild("node-rebuild-complete")

	h.CompleteRebuild("node-rebuild-complete")
	plan := h.GetRebuildPlan("node-rebuild-complete")
	assert.True(t, plan.Completed)
}

func TestGetIsolationLog(t *testing.T) {
	h := selfhealing.NewMeshHealer("mesh-1", 2)
	h.Isolate("node-x", "security threat")
	h.Isolate("node-y", "compromised")

	log := h.GetIsolationLog()
	assert.Len(t, log, 2)
	assert.Equal(t, "node-x", log[0].NodeID)
	assert.Equal(t, "node-y", log[1].NodeID)
}

func TestIsolationLog_Immutability(t *testing.T) {
	h := selfhealing.NewMeshHealer("mesh-1", 2)
	h.Isolate("node-a", "reason-a")
	h.Isolate("node-b", "reason-b")

	log := h.GetIsolationLog()
	log[0] = selfhealing.IsolationAction{} // attempt modification

	original := h.GetIsolationLog()
	assert.Equal(t, "node-a", original[0].NodeID)
	assert.Equal(t, "reason-a", original[0].Reason)
}

func TestQuorumLossDetection(t *testing.T) {
	h := selfhealing.NewMeshHealer("mesh-1", 3)
	nodes := []selfhealing.NodeState{
		{NodeID: "n1", TrustScore: 0.1, Status: "ACTIVE", FailureCount: 0, LastSeen: time.Now()},
		{NodeID: "n2", TrustScore: 0.15, Status: "ACTIVE", FailureCount: 0, LastSeen: time.Now()},
		{NodeID: "n3", TrustScore: 0.95, Status: "ACTIVE", FailureCount: 0, LastSeen: time.Now()},
	}
	h.MonitorTrust(nodes)

	plan := h.GetRebuildPlan("n1")
	assert.NotNil(t, plan)
	plan2 := h.GetRebuildPlan("n2")
	assert.NotNil(t, plan2)
}

func TestRebuildPlan_UnknownNode(t *testing.T) {
	h := selfhealing.NewMeshHealer("mesh-1", 3)
	plan := h.GetRebuildPlan("nonexistent")
	assert.Nil(t, plan)
}

func TestNodeRoleConstants(t *testing.T) {
	assert.Equal(t, selfhealing.NodeRole("VALIDATOR"), selfhealing.RoleValidator)
	assert.Equal(t, selfhealing.NodeRole("GATEWAY"), selfhealing.RoleGateway)
	assert.Equal(t, selfhealing.NodeRole("EXECUTOR"), selfhealing.RoleExecutor)
}
