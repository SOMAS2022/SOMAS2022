package election

import (
	"sync"

	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
)

func HandleElection(state *state.State, agents map[commons.ID]agent.Agent, strategy decision.VotingStrategy, numberOfPreferences uint) (
	commons.ID, decision.Manifesto,
) {
	// Get manifestos from agents
	agentManifestos := make(map[commons.ID]decision.Manifesto)

	for id, a := range agents {
		agentManifestos[id] = *a.SubmitManifesto(state.AgentState[id])
	}

	agentIDs := make([]commons.ID, len(agents))

	i := 0
	for k := range agents {
		agentIDs[i] = k
		i++
	}

	ballots := make([]decision.Ballot, 0)

	ballotChan := make(chan decision.Ballot)

	params := decision.NewElectionParams(agentManifestos, strategy, numberOfPreferences)

	var wg sync.WaitGroup
	for id, a := range agents {
		wg.Add(1)
		startAgentElectionHandlers(state.AgentState[id], a, params, ballotChan, &wg)
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
		winningID := singleChoicePlurality(ballots)
		winningManifesto := agentManifestos[winningID]

		return winningID, winningManifesto

	case decision.VotingStrategy(decision.BordaCount):
		winningID := BordaCount(ballots, agentIDs)
		winningManifesto := agentManifestos[winningID]

		return winningID, winningManifesto
	default:
		winningID := singleChoicePlurality(ballots)
		winningManifesto := agentManifestos[winningID]

		return winningID, winningManifesto
	}
}

// Create channel to a specific agent.
func startAgentElectionHandlers(agentState state.AgentState, a agent.Agent, params *decision.ElectionParams, dChan chan<- decision.Ballot, wg *sync.WaitGroup) {
	go func(group *sync.WaitGroup) {
		dChan <- a.HandleElection(agentState, params)

		group.Done()
	}(wg)
}
