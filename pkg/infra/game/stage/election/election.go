package election

import (
	"fmt"
	"infra/game/agent"
	"infra/game/decision"
	"infra/game/state"
	"infra/logging"
	"math/rand"
)

func CountVotesByChooseOne(decisions map[uint]decision.Ballot, agents map[uint]agent.Agent) string {
	// Count votes for each candidate
	candidates := make(map[string]uint)
	var totalVotes uint
	for _, d := range decisions {
		candidates[d[0]]++
		totalVotes++
	}
	logging.Log.Debug(fmt.Sprintf("Collect %4d votes for %4d candidates in total from %4d agents", totalVotes, len(candidates), len(agents)))

	// Find out the candidate/s with max #votes
	var maxNumberOfVote uint
	var winners []string
	for candidate, votes := range candidates {
		logging.Log.Trace(fmt.Sprintf("%s: %4d", candidate, votes))
		if votes > maxNumberOfVote {
			maxNumberOfVote = votes
			winners = []string{candidate}
		} else if votes == maxNumberOfVote {
			winners = append(winners, candidate)
		}
	}

	// randomly return one agent if no valid vote
	if maxNumberOfVote == 0 {
		randIdx := rand.Intn(len(agents))
		return agents[uint(randIdx)].BaseAgent.Id
	}

	// Return winner with max #votes,
	// if more than one candidates with max #votes, randomly pick one
	var winner string
	if len(winners) > 1 {
		winner = winners[rand.Intn(len(winners))]
	} else {
		winner = winners[0]
	}

	return winner
}

func HandleElection(state *state.State, agents map[uint]agent.Agent) string {
	decisions := make(map[uint]decision.Ballot)
	channels := make(map[uint]chan decision.Ballot)

	// filter out dead agents
	aliveAgents := make(map[uint]agent.Agent)
	for _, agent := range state.AgentState {
		if agent.Hp > 0 {
			aliveAgents = append(aliveAgents, agent)
		}
	}

	// Create channels to each agents
	for i, a := range agents {
		channels[i] = startAgentElectionHandlers(*state, agents, a)
	}

	// Collect ballots
	for i, dChan := range channels {
		decisions[i] = <-dChan
		close(dChan)
	}

	// Count votes
	result := CountVotesByChooseOne(decisions, agents)

	return result
}

// Create channel for a specific agent
func startAgentElectionHandlers(state state.State, agents map[uint]agent.Agent, a agent.Agent) chan decision.Ballot {
	decisionChan := make(chan decision.Ballot)
	go a.Strategy.HandleElection(state, agents, a.BaseAgent, decisionChan)
	return decisionChan
}
