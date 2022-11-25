package election

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
)

func HandleElection(state *state.State, agents map[commons.ID]agent.Agent, strategy decision.VotingStrategy, qtyPreferences uint) (commons.ID, uint) {
	ballots := make(map[commons.ID]decision.Ballot)
	channels := make(map[commons.ID]chan decision.Ballot)

	// Make immutable view of current state
	view := state.ToView()
	candidateList := make([]commons.ID, len(agents))
	i := 0
	for id := range agents {
		candidateList[i] = id
		i++
	}

	// Create channels to each agents
	for id, agent := range agents {
		channels[id] = startAgentElectionHandlers(view, agent, candidateList, strategy, qtyPreferences)
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
func startAgentElectionHandlers(view *state.View, a agent.Agent, candidateList []commons.ID, strategy decision.VotingStrategy, qtyPreferences uint) chan decision.Ballot {
	decisionChan := make(chan decision.Ballot)
	go a.Strategy.HandleElection(view, a.BaseAgent, decisionChan, candidateList, strategy, qtyPreferences)
	return decisionChan
}
