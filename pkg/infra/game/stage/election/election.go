package election

import (
	"fmt"
	"infra/game/agent"
	"infra/game/decision"
	"infra/game/state"
	"infra/logging"
)

func HandleElection(state *state.State, agents map[string]agent.Agent) decision.Ballot {
	logging.Log.Info(fmt.Sprintf("Number of agents left: %d", len(agents)))
	decisions := make(map[string]decision.Ballot)
	channels := make(map[string]chan decision.Ballot)

	// Filter out dead agents
	aliveAgents := make(map[string]agent.Agent)
	for id, agent := range agents {
		if (*state).AgentState[id].Hp > 0 {
			logging.Log.Debug(fmt.Sprintf("[%s] Hp: %4d", id, (*state).AgentState[id].Hp))
			aliveAgents[id] = agent
		}
	}
	logging.Log.Debug(fmt.Sprintf("%4d agents are still alive", len(aliveAgents)))

	// Create channels to each agents
	for i, a := range aliveAgents {
		channels[i] = startAgentElectionHandlers(*state, a)
	}

	// Collect ballots
	for i, dChan := range channels {
		decisions[i] = <-dChan
		close(dChan)
	}

	// Generate result
	var result decision.Ballot

	return result
}

// Create channel to a specific agent
func startAgentElectionHandlers(state state.State, a agent.Agent) chan decision.Ballot {
	decisionChan := make(chan decision.Ballot)
	go a.Strategy.HandleElection(state, a.BaseAgent, decisionChan)
	return decisionChan
}
