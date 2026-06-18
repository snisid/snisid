package security_test

import (
	"testing"
	"time"

	"github.com/snisid/platform/services/governance-os/security-mesh"
	"github.com/stretchr/testify/assert"
)

func setupMeshHealer() *security.MeshHealer {
	return security.NewMeshHealer("cluster-prod")
}

func TestNewMeshHealer(t *testing.T) {
	h := setupMeshHealer()
	assert.NotNil(t, h)
	assert.Equal(t, "cluster-prod", h.ClusterName)
	assert.Empty(t, h.GetOperationLog())
	assert.Empty(t, h.GetIsolatedNodes())
}

func TestMonitorMesh_HealthyNodesNoAction(t *testing.T) {
	h := setupMeshHealer()
	nodes := []security.MeshNode{
		{ID: "node-1", TrustScore: 0.95, Status: security.StatusHealthy, LastAttested: time.Now(), SPIFFEID: "spiffe://cluster/ns/default/node-1"},
		{ID: "node-2", TrustScore: 0.91, Status: security.StatusHealthy, LastAttested: time.Now(), SPIFFEID: "spiffe://cluster/ns/default/node-2"},
	}
	h.MonitorMesh(nodes)

	ops := h.GetOperationLog()
	assert.Empty(t, ops, "healthy nodes should trigger no operations")
	assert.Empty(t, h.GetIsolatedNodes())
}

func TestMonitorMesh_CriticalTrustTriggersSelfHealing(t *testing.T) {
	h := setupMeshHealer()
	nodes := []security.MeshNode{
		{ID: "node-comp", TrustScore: 0.15, Status: security.StatusHealthy, LastAttested: time.Now(), FailureCount: 5},
	}
	h.MonitorMesh(nodes)

	ops := h.GetOperationLog()
	assert.Len(t, ops, 1)
	assert.Equal(t, "node-comp", ops[0].NodeID)
	assert.Equal(t, "ISOLATE_NETWORK_POLICY", ops[0].Operation)
	assert.Equal(t, "COMPLETED", ops[0].Status)

	isolated := h.GetIsolatedNodes()
	assert.Contains(t, isolated, "node-comp")
}

func TestMonitorMesh_DegradedTrustDegradesNode(t *testing.T) {
	h := setupMeshHealer()
	nodes := []security.MeshNode{
		{ID: "node-deg", TrustScore: 0.35, Status: security.StatusHealthy, LastAttested: time.Now(), FailureCount: 2},
	}
	h.MonitorMesh(nodes)

	ops := h.GetOperationLog()
	assert.Len(t, ops, 1)
	assert.Equal(t, "node-deg", ops[0].NodeID)
	assert.Equal(t, "DEGRADE", ops[0].Operation)
}

func TestMonitorMesh_RestoresIsolatedNode(t *testing.T) {
	h := setupMeshHealer()

	h.MonitorMesh([]security.MeshNode{
		{ID: "node-iso", TrustScore: 0.15, Status: security.StatusHealthy, LastAttested: time.Now()},
	})

	assert.Len(t, h.GetIsolatedNodes(), 1)

	h.MonitorMesh([]security.MeshNode{
		{ID: "node-iso", TrustScore: 0.85, Status: security.StatusIsolated, LastAttested: time.Now()},
	})

	ops := h.GetOperationLog()
	assert.Len(t, ops, 2)
	assert.Equal(t, "RESTORE", ops[1].Operation)
	assert.Equal(t, "node-iso", ops[1].NodeID)

	assert.Empty(t, h.GetIsolatedNodes())
}

func TestApplyNetworkIsolation(t *testing.T) {
	h := setupMeshHealer()
	nodes := []security.MeshNode{
		{ID: "node-net-iso", TrustScore: 0.1, Status: security.StatusHealthy, LastAttested: time.Now()},
	}
	h.MonitorMesh(nodes)

	isolated := h.GetIsolatedNodes()
	assert.Contains(t, isolated, "node-net-iso")
}

func TestGetIsolatedNodes_Multiple(t *testing.T) {
	h := setupMeshHealer()

	h.MonitorMesh([]security.MeshNode{
		{ID: "n1", TrustScore: 0.1, Status: security.StatusHealthy, LastAttested: time.Now()},
	})
	h.MonitorMesh([]security.MeshNode{
		{ID: "n2", TrustScore: 0.15, Status: security.StatusHealthy, LastAttested: time.Now()},
	})

	isolated := h.GetIsolatedNodes()
	assert.Len(t, isolated, 2)
	assert.Contains(t, isolated, "n1")
	assert.Contains(t, isolated, "n2")
}

func TestOperationLog_Immutability(t *testing.T) {
	h := setupMeshHealer()
	h.MonitorMesh([]security.MeshNode{
		{ID: "node-log", TrustScore: 0.1, Status: security.StatusHealthy, LastAttested: time.Now()},
	})

	ops := h.GetOperationLog()
	ops[0].NodeID = "tampered"

	original := h.GetOperationLog()
	assert.Equal(t, "node-log", original[0].NodeID)
}

func TestNodeStatusConstants(t *testing.T) {
	assert.Equal(t, security.NodeStatus("HEALTHY"), security.StatusHealthy)
	assert.Equal(t, security.NodeStatus("DEGRADED"), security.StatusDegraded)
	assert.Equal(t, security.NodeStatus("ISOLATED"), security.StatusIsolated)
	assert.Equal(t, security.NodeStatus("COMPROMISED"), security.StatusCompromised)
}

func TestMeshHealer_ConcurrentAccess(t *testing.T) {
	h := setupMeshHealer()
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(idx int) {
			nodes := []security.MeshNode{
				{ID: "conc-node", TrustScore: 0.1, Status: security.StatusHealthy, LastAttested: time.Now()},
			}
			h.MonitorMesh(nodes)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	// At least one operation should be recorded
	assert.NotEmpty(t, h.GetOperationLog())
}
