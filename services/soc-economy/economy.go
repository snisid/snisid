package economy

import (
	"sort"
	"time"
)

type SOCAgent struct {
	ID         string
	Capability float64 // 0.0 to 1.0
	Score      float64
	LastActive time.Time
}

type IncidentEconomy struct {
	Agents []*SOCAgent
}

func (e *IncidentEconomy) AssignIncident(incidentID string, complexity float64) *SOCAgent {
	// Bidding logic: find the agent with the highest capability/score ratio
	sort.Slice(e.Agents, func(i, j int) bool {
		return (e.Agents[i].Capability * e.Agents[i].Score) > (e.Agents[j].Capability * e.Agents[j].Score)
	})

	if len(e.Agents) > 0 {
		winner := e.Agents[0]
		winner.LastActive = time.Now()
		return winner
	}
	return nil
}

func (e *IncidentEconomy) RewardAgent(agentID string, success bool) {
	for _, a := range e.Agents {
		if a.ID == agentID {
			if success {
				a.Score += 10.0
			} else {
				a.Score -= 5.0
			}
		}
	}
}
