package election

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
)

func HandleElection(state *state.State, agents map[commons.ID]agent.Agent, strategy decision.VotingStrategy, numberOfPreferences uint) (commons.ID, uint) {
	ballots := make(map[commons.ID]decision.Ballot)
	channels := make(map[commons.ID]chan decision.Ballot)

	// Make immutable view of current state
	view := state.ToView()
	candidateList := make([]commons.ID, len(agents))
	decisionChan := make(chan decision.Ballot)
	i := 0
	for id := range agents {
		candidateList[i] = id
		i++
	}

	params := decision.ElectionParams{candidateList: candidateList, strategy: strategy, numberOfPreferences: numberOfPreferences}
	// Create channels to each agents
	for id, agent := range agents {
		channels[id] = startAgentElectionHandlers(state.AgentState[id], view, agent, params)
	}

	// Collect ballots
	for i, dChan := range channels {
		ballots[i] = <-dChan
		close(dChan)
	}

	switch strategy {
	case decision.VotingStrategy(decision.SingleChoicePlurality):
		return singleChoicePlurality(ballots)
	default:
		return singleChoicePlurality(ballots)
	}
}

// Create channel to a specific agent
func startAgentElectionHandlers(agentState state.AgentState, view *state.View, a agent.Agent, params decision.ElectionParams) chan decision.Ballot {
	go a.HandleElection(agentState, view, a.BaseAgent, params)
	return decisionChan
}
