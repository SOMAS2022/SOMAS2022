package election

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
	"infra/logging"
	"math/rand"
)

func ChooseLeaderByChooseOne(ballots map[commons.ID]decision.Ballot) (commons.ID, uint) {
	// Count number of votes collected for each candidate
	votes := make(map[commons.ID]uint)
	for _, ballot := range ballots {
		votes[ballot[0]]++
	}

	// Find the candidate(s) with max number of votes
	var maxNumVotes uint
	var winners []commons.ID
	for agentId, numVotes := range votes {
		if numVotes > maxNumVotes {
			maxNumVotes = numVotes
			winners = []commons.ID{agentId}
		} else if numVotes == maxNumVotes {
			winners = append(winners, agentId)
		}
	}

	// Randomly pick one if no valid votes
	if maxNumVotes == 0 {
		candidateIds := make([]commons.ID, len(ballots))
		i := uint(0)
		for id := range ballots {
			candidateIds[i] = id
			i++
		}
		logging.Log(logging.Debug, nil, "No Valid Votes")
		return candidateIds[len(ballots)], 0
	}

	// Randomly choose one if there are more than one winner
	var winner commons.ID
	if len(winners) > 1 {
		logging.Log(
			logging.Debug,
			logging.LogField{"winners": winners},
			"Multiple candidates with a winning number of votes",
		)
		randIdx := rand.Intn(len(winners))
		winner = winners[uint(randIdx)]
	} else {
		winner = winners[0]
	}

	return winner, 100 * maxNumVotes / uint(len(ballots))
}

func HandleElection(state *state.State, agents map[commons.ID]agent.Agent) (commons.ID, uint) {
	ballots := make(map[commons.ID]decision.Ballot)
	channels := make(map[commons.ID]chan decision.Ballot)

	// Make immutable view of current state
	view := state.ToView()

	// Create channels to each agents
	for id, agent := range agents {
		channels[id] = startAgentElectionHandlers(view, agent)
	}

	// Collect ballots
	for i, dChan := range channels {
		ballots[i] = <-dChan
		close(dChan)
	}

	// Generate result
	winner, percentage := ChooseLeaderByChooseOne(ballots)

	return winner, percentage
}

// Create channel to a specific agent
func startAgentElectionHandlers(view *state.View, a agent.Agent) chan decision.Ballot {
	decisionChan := make(chan decision.Ballot)
	go a.Strategy.HandleElection(view, a.BaseAgent, decisionChan)
	return decisionChan
}
