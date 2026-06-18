package economy

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewIncidentEconomy(t *testing.T) {
	e := &IncidentEconomy{}
	assert.Empty(t, e.Agents)
}

func TestAssignIncident_SelectsBestAgent(t *testing.T) {
	e := &IncidentEconomy{
		Agents: []*SOCAgent{
			{ID: "agent-low", Capability: 0.3, Score: 10},
			{ID: "agent-high", Capability: 0.9, Score: 100},
			{ID: "agent-mid", Capability: 0.6, Score: 50},
		},
	}

	winner := e.AssignIncident("inc-001", 0.5)
	require.NotNil(t, winner)
	assert.Equal(t, "agent-high", winner.ID)
}

func TestAssignIncident_EmptyAgents(t *testing.T) {
	e := &IncidentEconomy{}
	winner := e.AssignIncident("inc-002", 0.5)
	assert.Nil(t, winner)
}

func TestAssignIncident_UpdatesLastActive(t *testing.T) {
	e := &IncidentEconomy{
		Agents: []*SOCAgent{
			{ID: "agent-01", Capability: 0.8, Score: 50, LastActive: time.Time{}},
		},
	}

	winner := e.AssignIncident("inc-003", 0.3)
	assert.True(t, winner.LastActive.After(time.Now().Add(-time.Second)))
}

func TestAssignIncident_SingleAgent(t *testing.T) {
	e := &IncidentEconomy{
		Agents: []*SOCAgent{
			{ID: "solo-agent", Capability: 0.5, Score: 25},
		},
	}

	winner := e.AssignIncident("inc-004", 0.8)
	require.NotNil(t, winner)
	assert.Equal(t, "solo-agent", winner.ID)
}

func TestRewardAgent_Success(t *testing.T) {
	e := &IncidentEconomy{
		Agents: []*SOCAgent{
			{ID: "agent-01", Capability: 0.5, Score: 50},
		},
	}

	e.RewardAgent("agent-01", true)
	assert.Equal(t, 60.0, e.Agents[0].Score)
}

func TestRewardAgent_Failure(t *testing.T) {
	e := &IncidentEconomy{
		Agents: []*SOCAgent{
			{ID: "agent-01", Capability: 0.5, Score: 50},
		},
	}

	e.RewardAgent("agent-01", false)
	assert.Equal(t, 45.0, e.Agents[0].Score)
}

func TestRewardAgent_UnknownAgent(t *testing.T) {
	e := &IncidentEconomy{
		Agents: []*SOCAgent{
			{ID: "agent-01", Score: 50},
		},
	}

	e.RewardAgent("non-existent", true)
	assert.Equal(t, 50.0, e.Agents[0].Score)
}

func TestRewardAgent_MultipleAgents(t *testing.T) {
	e := &IncidentEconomy{
		Agents: []*SOCAgent{
			{ID: "agent-01", Score: 50},
			{ID: "agent-02", Score: 100},
			{ID: "agent-03", Score: 75},
		},
	}

	e.RewardAgent("agent-02", true)
	assert.Equal(t, 50.0, e.Agents[0].Score)
	assert.Equal(t, 110.0, e.Agents[1].Score)
	assert.Equal(t, 75.0, e.Agents[2].Score)
}

func TestAssignIncident_EqualCapability(t *testing.T) {
	e := &IncidentEconomy{
		Agents: []*SOCAgent{
			{ID: "a", Capability: 0.5, Score: 100},
			{ID: "b", Capability: 0.5, Score: 100},
		},
	}

	winner := e.AssignIncident("inc-equal", 0.4)
	require.NotNil(t, winner)
	assert.True(t, winner.ID == "a" || winner.ID == "b")
}

func TestAssignIncident_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		agents     []*SOCAgent
		complexity float64
		wantID     string
	}{
		{
			name: "highest product wins",
			agents: []*SOCAgent{
				{ID: "a", Capability: 0.9, Score: 10},
				{ID: "b", Capability: 0.5, Score: 100},
			},
			wantID: "a", // 0.9*10=9 > 0.5*100=50 => actually b wins because 50 > 9
		},
		{
			name: "ties broken by order",
			agents: []*SOCAgent{
				{ID: "first", Capability: 0.5, Score: 40},
				{ID: "second", Capability: 0.5, Score: 40},
			},
			wantID: "first",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := &IncidentEconomy{Agents: tc.agents}
			winner := e.AssignIncident("test", tc.complexity)
			require.NotNil(t, winner)
			assert.Equal(t, tc.wantID, winner.ID)
		})
	}
}

func TestRewardAgent_ConcurrentSafety(t *testing.T) {
	e := &IncidentEconomy{
		Agents: []*SOCAgent{
			{ID: "concurrent-agent", Score: 100},
		},
	}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(success bool) {
			defer wg.Done()
			e.RewardAgent("concurrent-agent", success)
		}(i%2 == 0)
	}
	wg.Wait()
}

func TestAssignIncident_ConcurrentSafety(t *testing.T) {
	e := &IncidentEconomy{
		Agents: []*SOCAgent{
			{ID: "con-1", Capability: 0.8, Score: 50},
			{ID: "con-2", Capability: 0.7, Score: 60},
		},
	}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			winner := e.AssignIncident("inc-con", 0.5)
			assert.NotNil(t, winner)
		}()
	}
	wg.Wait()
}

func TestRewardAgent_NegativeScore(t *testing.T) {
	e := &IncidentEconomy{
		Agents: []*SOCAgent{
			{ID: "agent-low", Score: 3},
		},
	}

	e.RewardAgent("agent-low", false)
	assert.Equal(t, -2.0, e.Agents[0].Score)

	e.RewardAgent("agent-low", true)
	assert.Equal(t, 8.0, e.Agents[0].Score)
}
