package election

import (
	"sync"

	"infra/game/agent"
	"infra/game/commons"
	"infra/game/state"
	"infra/game/state/proposal"
)

func HandleElection(state *state.State, agents map[commons.ID]agent.Agent, strategy proposal.VotingStrategy, numberOfPreferences uint) (
	commons.ID, proposal.Manifesto,
) {
	// Get manifestos from agents
	agentManifestos := make(map[commons.ID]proposal.Manifesto)

	for id, a := range agents {
		agentManifesto := *a.SubmitManifesto(state.AgentState[id])
		if agentManifesto.Runnning() {
			agentManifestos[id] = agentManifesto
		}
	}

	if len(agentManifestos) < 1 {
		return "", proposal.Manifesto{}
	}

	agentIDs := make([]commons.ID, len(agents))

	i := 0
	for k := range agents {
		agentIDs[i] = k
		i++
	}

	ballots := make([]proposal.Ballot, 0)

	ballotChan := make(chan proposal.Ballot)

	params := proposal.NewElectionParams(agentManifestos, strategy, numberOfPreferences)

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
	case proposal.VotingStrategy(proposal.SingleChoicePlurality):
		winningID := singleChoicePlurality(ballots)
		winningManifesto := agentManifestos[winningID]

		return winningID, winningManifesto

	case proposal.VotingStrategy(proposal.BordaCount):
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
func startAgentElectionHandlers(agentState state.AgentState, a agent.Agent, params *proposal.ElectionParams, dChan chan<- proposal.Ballot, wg *sync.WaitGroup) {
	go func(group *sync.WaitGroup) {
		dChan <- a.HandleElection(agentState, params)

		group.Done()
	}(wg)
}
