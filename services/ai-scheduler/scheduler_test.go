package scheduler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptimizePlacement_EmptyInputs(t *testing.T) {
	s := AIScheduler{}
	placement := s.OptimizePlacement(nil, nil)
	assert.Empty(t, placement)
}

func TestOptimizePlacement_NoNodes(t *testing.T) {
	s := AIScheduler{}
	pods := []Pod{
		{Name: "pod-1", CPUReq: 2, MemReq: 4},
	}
	placement := s.OptimizePlacement(pods, nil)
	assert.Empty(t, placement)
}

func TestOptimizePlacement_SimpleAssignment(t *testing.T) {
	s := AIScheduler{}
	pods := []Pod{
		{Name: "pod-1", CPUReq: 2, MemReq: 4},
	}
	nodes := []Node{
		{ID: "node-1", Region: "us-east", Capacity: 10, Free: 8, Health: 0.95},
	}
	placement := s.OptimizePlacement(pods, nodes)
	assert.Equal(t, "node-1", placement["pod-1"])
}

func TestOptimizePlacement_PrefersHealthiest(t *testing.T) {
	s := AIScheduler{}
	pods := []Pod{
		{Name: "pod-1", CPUReq: 2, MemReq: 4},
	}
	nodes := []Node{
		{ID: "node-1", Region: "us-east", Capacity: 10, Free: 8, Health: 0.7},
		{ID: "node-2", Region: "us-east", Capacity: 10, Free: 8, Health: 0.95},
		{ID: "node-3", Region: "us-east", Capacity: 10, Free: 8, Health: 0.5},
	}
	placement := s.OptimizePlacement(pods, nodes)
	assert.Equal(t, "node-2", placement["pod-1"])
}

func TestOptimizePlacement_SkipsFullNodes(t *testing.T) {
	s := AIScheduler{}
	pods := []Pod{
		{Name: "pod-1", CPUReq: 10, MemReq: 4},
	}
	nodes := []Node{
		{ID: "node-1", Region: "us-east", Capacity: 10, Free: 5, Health: 0.95},
		{ID: "node-2", Region: "us-east", Capacity: 10, Free: 12, Health: 0.8},
	}
	placement := s.OptimizePlacement(pods, nodes)
	assert.Equal(t, "node-2", placement["pod-1"])
}

func TestOptimizePlacement_MultiplePods(t *testing.T) {
	s := AIScheduler{}
	pods := []Pod{
		{Name: "pod-1", CPUReq: 2, MemReq: 4},
		{Name: "pod-2", CPUReq: 3, MemReq: 2},
		{Name: "pod-3", CPUReq: 1, MemReq: 1},
	}
	nodes := []Node{
		{ID: "node-1", Region: "us-east", Capacity: 10, Free: 8, Health: 0.9},
		{ID: "node-2", Region: "us-west", Capacity: 10, Free: 8, Health: 0.7},
	}
	placement := s.OptimizePlacement(pods, nodes)
	assert.Len(t, placement, 3)
	assert.Equal(t, "node-1", placement["pod-1"])
	assert.Equal(t, "node-1", placement["pod-2"])
	assert.Equal(t, "node-1", placement["pod-3"])
}

func TestOptimizePlacement_UnplaceablePod(t *testing.T) {
	s := AIScheduler{}
	pods := []Pod{
		{Name: "pod-1", CPUReq: 100, MemReq: 100},
	}
	nodes := []Node{
		{ID: "node-1", Region: "us-east", Capacity: 10, Free: 8, Health: 0.9},
	}
	placement := s.OptimizePlacement(pods, nodes)
	assert.Empty(t, placement)
}

func TestOptimizePlacement_ZeroHealth(t *testing.T) {
	s := AIScheduler{}
	pods := []Pod{
		{Name: "pod-1", CPUReq: 2, MemReq: 4},
	}
	nodes := []Node{
		{ID: "node-1", Region: "us-east", Capacity: 10, Free: 8, Health: 0.0},
	}
	placement := s.OptimizePlacement(pods, nodes)
	assert.Equal(t, "node-1", placement["pod-1"])
}

func TestOptimizePlacement_MixedPlaceability(t *testing.T) {
	s := AIScheduler{}
	pods := []Pod{
		{Name: "pod-small", CPUReq: 2, MemReq: 4},
		{Name: "pod-big", CPUReq: 50, MemReq: 50},
		{Name: "pod-medium", CPUReq: 5, MemReq: 5},
	}
	nodes := []Node{
		{ID: "node-1", Region: "us-east", Capacity: 10, Free: 10, Health: 0.9},
	}
	placement := s.OptimizePlacement(pods, nodes)
	assert.Equal(t, "node-1", placement["pod-small"])
	assert.Equal(t, "node-1", placement["pod-medium"])
	assert.Empty(t, placement["pod-big"])
}

func TestOptimizePlacement_NoDuplicates(t *testing.T) {
	s := AIScheduler{}
	pods := []Pod{{Name: "pod-1", CPUReq: 2, MemReq: 4}}
	nodes := []Node{
		{ID: "node-1", Region: "us-east", Capacity: 10, Free: 10, Health: 0.9},
		{ID: "node-2", Region: "us-west", Capacity: 10, Free: 10, Health: 0.9},
	}
	placement := s.OptimizePlacement(pods, nodes)
	assert.Len(t, placement, 1)
}
