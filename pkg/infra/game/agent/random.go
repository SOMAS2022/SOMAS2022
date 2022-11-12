package agent

import (
	"github.com/benbjohnson/immutable"
	"infra/game/decision"
	"infra/game/state"
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

func (RandomAgent) HandleElection(view *state.View, _ BaseAgent, decisionC chan<- decision.Ballot) {
	// Extract ID of alive agents
	agentState := *(*view).AgentState
	aliveAgentIds := make([]string, agentState.Len())
	i := 0
	itr := agentState.Iterator()
	for !itr.Done() {
		id, agent, ok := itr.Next()
		if ok && agent.Hp > 0 {
			aliveAgentIds[i] = id
			i++
		}
	}

	// Randomly fill the ballot
	var ballot decision.Ballot
	numAliveAgents := len(aliveAgentIds)
	numCandidate := 2
	for i := 0; i < numCandidate; i++ {
		randomIdx := rand.Intn(numAliveAgents)
		randomCandidate := aliveAgentIds[uint(randomIdx)]
		ballot = append(ballot, randomCandidate)
	}

	// Send ballot to receiver
	decisionC <- ballot
}
