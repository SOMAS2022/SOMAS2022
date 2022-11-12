package election

import (
	"infra/game/agent"
	"infra/game/decision"
	"infra/game/state"
	"math/rand"
)

func ChooseLeaderByChooseOne(ballots map[string]decision.Ballot) string {
	// Count number of votes collected for each candidate
	votes := make(map[string]uint)
	for _, ballot := range ballots {
		votes[ballot[0]]++
	}

	// Find the candidate(s) with max number of votes
	var maxNumVotes uint
	var winners []string
	for agentId, numVotes := range votes {
		if numVotes > maxNumVotes {
			maxNumVotes = numVotes
			winners = []string{agentId}
		} else if numVotes == maxNumVotes {
			winners = append(winners, agentId)
		}
	}

	// Randomly pick one if no valid votes
	if maxNumVotes == 0 {
		candidateIds := make([]string, len(ballots))
		i := uint(0)
		for id := range ballots {
			candidateIds[i] = id
			i++
		}
		return candidateIds[len(ballots)]
	}

	// Randomly choose one if there are more than one winner
	var winner string
	if len(winners) > 1 {
		randIdx := rand.Intn(len(winners))
		winner = winners[uint(randIdx)]
	} else {
		winner = winners[0]
	}

	return winner
}

func HandleElection(state *state.State, agents map[string]agent.Agent) (uint, string) {
	ballots := make(map[string]decision.Ballot)
	channels := make(map[string]chan decision.Ballot)

	// Make immutable view of current state
	view := state.ToView()

	// Filter out dead agents
	aliveAgents := make(map[string]agent.Agent)
	for id, agent := range agents {
		if (*state).AgentState[id].Hp > 0 {
			aliveAgents[id] = agent
		}
	}

	// Create channels to each agents
	for i, a := range aliveAgents {
		channels[i] = startAgentElectionHandlers(view, a)
	}

	// Collect ballots
	for i, dChan := range channels {
		ballots[i] = <-dChan
		close(dChan)
	}

	// Generate result
	result := ChooseLeaderByChooseOne(ballots)

	return uint(len(aliveAgents)), result
}

// Create channel to a specific agent
func startAgentElectionHandlers(view *state.View, a agent.Agent) chan decision.Ballot {
	decisionChan := make(chan decision.Ballot)
	go a.Strategy.HandleElection(view, a.BaseAgent, decisionChan)
	return decisionChan
}
