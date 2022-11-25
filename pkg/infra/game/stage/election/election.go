package election

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
	"sync"
)

func HandleElection(state *state.State, agents map[commons.ID]agent.Agent, strategy decision.VotingStrategy, numberOfPreferences uint) (commons.ID, uint) {
	ballots := make([]decision.Ballot, len(agents))

	// Make immutable view of current state
	view := state.ToView()
	candidateList := make([]commons.ID, len(agents))
	ballotChan := make(chan decision.Ballot)
	i := 0
	for id := range agents {
		candidateList[i] = id
		i++
	}

	params := decision.NewElectionParams(candidateList, strategy, numberOfPreferences)
	var wg sync.WaitGroup
	for id, a := range agents {
		wg.Add(1)
		startAgentElectionHandlers(state.AgentState[id], view, a, params, ballotChan, &wg)
	}

	go func(group *sync.WaitGroup) {
		group.Wait()
		close(ballotChan)
	}(&wg)

	for ballot := range ballotChan {
		ballots = append(ballots, ballot)
	}

	switch strategy {
	case decision.VotingStrategy(decision.SingleChoicePlurality):
		return singleChoicePlurality(ballots)
	default:
		return singleChoicePlurality(ballots)
	}
}

// Create channel to a specific agent
func startAgentElectionHandlers(agentState state.AgentState, view *state.View, a agent.Agent, params *decision.ElectionParams, dChan chan<- decision.Ballot, wg *sync.WaitGroup) {
	go func(group *sync.WaitGroup) {
		dChan <- a.HandleElection(agentState, view, a.BaseAgent, params)
		group.Done()
	}(wg)
}
