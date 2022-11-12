package agent

import (
	"github.com/benbjohnson/immutable"
	"infra/game/decision"
	"infra/game/state"
	"infra/logging"
	"math/rand"
)

type RandomAgent struct {
}

func (RandomAgent) HandleFight(_ *state.View, _ BaseAgent, decisionC chan<- decision.FightAction, log *immutable.Map[string, decision.FightAction]) {
	fight := rand.Intn(3)
	switch fight {
	case 0:
		decisionC <- decision.Cower
	case 1:
		decisionC <- decision.Attack
	case 2:
		decisionC <- decision.Defend
	}
}

func (RandomAgent) HandleElection(gameState state.State, _ BaseAgent, decisionC chan<- decision.Ballot) {
	var ballot decision.Ballot

	// Extract ID of alive agents
	aliveAgentIds := make([]string, len(gameState.AgentState))
	i := 0
	for id, agent := range gameState.AgentState {
		if agent.Hp > 0 {
			aliveAgentIds[i] = id
			i++
		}
	}

	numAliveAgents := len(aliveAgentIds)
	numCandidate := 2

	logging.Log.Debug(numAliveAgents)

	for i := 0; i < numCandidate; i++ {
		randomIdx := rand.Intn(numAliveAgents)
		randomCandidate := aliveAgentIds[uint(randomIdx)]
		ballot = append(ballot, randomCandidate)
	}

	// Send ballot to receiver
	decisionC <- ballot
}
