package compiler

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAIInfraCompiler(t *testing.T) {
	c := NewAIInfraCompiler()
	assert.NotNil(t, c)
	assert.Equal(t, "infra.scaling.decisions", c.kafkaTopic)
}

func TestExtractDeployment(t *testing.T) {
	tests := []struct {
		name    string
		podName string
		want    string
	}{
		{"single dash", "nginx-7d8b9c", "nginx"},
		{"multi dash", "api-gateway-6f9d2c", "api-gateway"},
		{"no suffix", "myapp", "myapp"},
		{"short", "ab-cd", "ab-cd"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, extractDeployment(tt.podName))
		})
	}
}

func TestCalculateEstimatedSavings(t *testing.T) {
	tests := []struct {
		name      string
		decisions []ScalingDecision
		want      float64
	}{
		{"no savings", []ScalingDecision{{CurrentReplicas: 3, TargetReplicas: 5}}, 0},
		{"some savings", []ScalingDecision{{CurrentReplicas: 10, TargetReplicas: 6}}, 2.0},
		{"multiple savings", []ScalingDecision{
			{CurrentReplicas: 10, TargetReplicas: 6},
			{CurrentReplicas: 8, TargetReplicas: 4},
		}, 4.0},
		{"no decisions", nil, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateEstimatedSavings(tt.decisions)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAnalyzeScaling_EmptyState(t *testing.T) {
	c := NewAIInfraCompiler()
	state := ClusterState{}
	plan := c.analyzeScaling(state)
	assert.Nil(t, plan)
}

func TestAnalyzeScaling_HighUtilization(t *testing.T) {
	c := NewAIInfraCompiler()
	state := ClusterState{
		Nodes: []NodeStatus{
			{Name: "node-1", AllocatableCPU: 8, CPULimit: 10},
		},
		Pods: []PodStatus{
			{Name: "app-abc-1", Namespace: "default", CPURequest: 0.9, CPULimit: 1.0, MemoryRequest: 0.5, Phase: "Running"},
			{Name: "app-abc-2", Namespace: "default", CPURequest: 0.9, CPULimit: 1.0, MemoryRequest: 0.5, Phase: "Running"},
		},
	}
	plan := c.analyzeScaling(state)
	require.NotNil(t, plan)
	require.Len(t, plan.Decisions, 1)
	assert.Equal(t, 1, plan.Decisions[0].Priority)
	assert.Greater(t, plan.Decisions[0].TargetReplicas, plan.Decisions[0].CurrentReplicas)
}

func TestAnalyzeScaling_LowUtilization(t *testing.T) {
	c := NewAIInfraCompiler()
	state := ClusterState{
		Nodes: []NodeStatus{
			{Name: "node-1", AllocatableCPU: 8},
		},
		Pods: []PodStatus{
			{Name: "app-abc-1", Namespace: "default", CPURequest: 0.15, MemoryRequest: 0.1, Phase: "Running"},
			{Name: "app-abc-2", Namespace: "default", CPURequest: 0.15, MemoryRequest: 0.1, Phase: "Running"},
			{Name: "app-abc-3", Namespace: "default", CPURequest: 0.15, MemoryRequest: 0.1, Phase: "Running"},
		},
	}
	plan := c.analyzeScaling(state)
	require.NotNil(t, plan)
	require.Len(t, plan.Decisions, 1)
	assert.Equal(t, 3, plan.Decisions[0].Priority)
	assert.Less(t, plan.Decisions[0].TargetReplicas, plan.Decisions[0].CurrentReplicas)
}

func TestAnalyzeScaling_IgnoresNonRunning(t *testing.T) {
	c := NewAIInfraCompiler()
	state := ClusterState{
		Pods: []PodStatus{
			{Name: "app-abc-1", Namespace: "default", CPURequest: 0.9, Phase: "Pending"},
		},
	}
	plan := c.analyzeScaling(state)
	assert.Nil(t, plan)
}

func TestAnalyzeHealth_HealthyState(t *testing.T) {
	c := NewAIInfraCompiler()
	state := ClusterState{
		Nodes: []NodeStatus{
			{Name: "node-1", HealthScore: 0.95, AllocatableCPU: 8, CPULimit: 10},
		},
	}
	plan := c.analyzeHealth(state)
	assert.Nil(t, plan)
}

func TestAnalyzeHealth_UnhealthyNode(t *testing.T) {
	c := NewAIInfraCompiler()
	state := ClusterState{
		Nodes: []NodeStatus{
			{Name: "node-1", HealthScore: 0.3, AllocatableCPU: 8, CPULimit: 10},
		},
	}
	plan := c.analyzeHealth(state)
	require.NotNil(t, plan)
	require.Len(t, plan.Actions, 2)
	assert.Equal(t, "CORDON", plan.Actions[0].ActionType)
	assert.Equal(t, "DRAIN", plan.Actions[1].ActionType)
}

func TestAnalyzeHealth_CrashLoop(t *testing.T) {
	c := NewAIInfraCompiler()
	state := ClusterState{
		Pods: []PodStatus{
			{Name: "app-abc-1", Namespace: "default", RestartCount: 10, Phase: "Running"},
		},
	}
	plan := c.analyzeHealth(state)
	require.NotNil(t, plan)
	require.Len(t, plan.Actions, 1)
	assert.Equal(t, "RESTART", plan.Actions[0].ActionType)
	assert.Equal(t, 1, plan.Actions[0].Priority)
}

func TestAnalyzeHealth_HighCPUUtilization(t *testing.T) {
	c := NewAIInfraCompiler()
	state := ClusterState{
		Nodes: []NodeStatus{
			{Name: "node-1", HealthScore: 0.9, CPULimit: 10, AllocatableCPU: 0.3},
		},
	}
	plan := c.analyzeHealth(state)
	require.NotNil(t, plan)
	require.Len(t, plan.Actions, 1)
	assert.Equal(t, "RESTART", plan.Actions[0].ActionType)
}

func TestAnalyzePlacement_EmptyNodes(t *testing.T) {
	c := NewAIInfraCompiler()
	state := ClusterState{Nodes: nil, Pods: []PodStatus{{Name: "pod-1", Phase: "Pending"}}}
	plan := c.analyzePlacement(state)
	assert.Nil(t, plan)
}

func TestAnalyzePlacement_PendingPod(t *testing.T) {
	c := NewAIInfraCompiler()
	state := ClusterState{
		Nodes: []NodeStatus{
			{Name: "node-1", HealthScore: 0.9, GPUCount: 0},
			{Name: "node-2", HealthScore: 0.7, GPUCount: 0},
		},
		Pods: []PodStatus{
			{Name: "pod-1", Phase: "Pending", HasGPU: false},
		},
	}
	plan := c.analyzePlacement(state)
	require.NotNil(t, plan)
	require.Len(t, plan.Suggestions, 1)
	assert.Equal(t, "node-1", plan.Suggestions[0].PreferredNode)
}

func TestAnalyzePlacement_GPURequired(t *testing.T) {
	c := NewAIInfraCompiler()
	state := ClusterState{
		Nodes: []NodeStatus{
			{Name: "node-1", HealthScore: 0.9, GPUCount: 0},
			{Name: "node-2", HealthScore: 0.7, GPUCount: 4},
		},
		Pods: []PodStatus{
			{Name: "pod-gpu-1", Phase: "Pending", HasGPU: true},
		},
	}
	plan := c.analyzePlacement(state)
	require.NotNil(t, plan)
	require.Len(t, plan.Suggestions, 1)
	assert.Equal(t, "node-2", plan.Suggestions[0].PreferredNode)
	assert.Equal(t, "nodeAffinity", plan.Suggestions[0].Affinity)
}

func TestAnalyzePlacement_NoPendingPods(t *testing.T) {
	c := NewAIInfraCompiler()
	state := ClusterState{
		Nodes: []NodeStatus{{Name: "node-1", HealthScore: 0.9}},
		Pods:  []PodStatus{{Name: "pod-1", Phase: "Running"}},
	}
	plan := c.analyzePlacement(state)
	assert.Nil(t, plan)
}

func TestCompile_FullPipeline(t *testing.T) {
	c := NewAIInfraCompiler()
	state := ClusterState{
		Nodes: []NodeStatus{
			{Name: "node-1", HealthScore: 0.3, CPULimit: 10, AllocatableCPU: 0.3, GPUCount: 0},
			{Name: "node-2", HealthScore: 0.95, CPULimit: 10, AllocatableCPU: 8, GPUCount: 2},
		},
		Pods: []PodStatus{
			{Name: "app-abc-1", Namespace: "default", CPURequest: 0.9, MemoryRequest: 0.5, Phase: "Running"},
			{Name: "app-abc-2", Namespace: "default", CPURequest: 0.9, MemoryRequest: 0.5, Phase: "Running"},
			{Name: "gpu-pod-1", Namespace: "ml", Phase: "Pending", HasGPU: true},
		},
	}
	plan, err := c.Compile(state)
	require.NoError(t, err)
	require.NotNil(t, plan)
	assert.NotNil(t, plan.Scaling)
	assert.NotNil(t, plan.Healing)
	assert.NotNil(t, plan.Placement)
}

func TestCompile_NilWhenNoChanges(t *testing.T) {
	c := NewAIInfraCompiler()
	state := ClusterState{
		Nodes: []NodeStatus{{Name: "node-1", HealthScore: 0.95, CPULimit: 10, AllocatableCPU: 8}},
		Pods: []PodStatus{
			{Name: "pod-1", Namespace: "default", CPURequest: 0.5, MemoryRequest: 0.3, Phase: "Running"},
		},
	}
	plan, err := c.Compile(state)
	require.NoError(t, err)
	assert.Nil(t, plan.Scaling)
	assert.Nil(t, plan.Healing)
	assert.Nil(t, plan.Placement)
}

func TestEncodeDecisions(t *testing.T) {
	c := NewAIInfraCompiler()
	plan := &InfrastructurePlan{
		Scaling: &ScalingPlan{
			Decisions: []ScalingDecision{
				{Namespace: "default", Deployment: "app", CurrentReplicas: 2, TargetReplicas: 4, Reason: "scale up", Priority: 1, GeneratedAt: time.Now()},
			},
		},
	}
	data, err := c.EncodeDecisions(plan)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	var decoded map[string]interface{}
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, "infra.scaling.decisions", decoded["topic"])
	assert.Equal(t, "1.0", decoded["version"])
}
