package team6

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"infra/logging"
	"math/rand"

	"github.com/benbjohnson/immutable"
)

type Team6Agent struct {
	bravery int
}

func (r *Team6Agent) CreateManifesto(view *state.View, baseAgent agent.BaseAgent) *decision.Manifesto {
	manifesto := decision.NewManifesto(true, false, 10, 50)
	return manifesto
}

func (r *Team6Agent) HandleConfidencePoll(view *state.View, baseAgent agent.BaseAgent) decision.Intent {
	switch rand.Intn(3) {
	case 0:
		return decision.Abstain
	case 1:
		return decision.Negative
	default:
		return decision.Positive
	}
}

func (r *Team6Agent) HandleFightInformation(_ message.TaggedMessage, _ *state.View, agent agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) {
	agent.Log(logging.Trace, logging.LogField{"bravery": r.bravery, "hp": agent.ViewState().Hp}, "Cowering")
}

func (r *Team6Agent) HandleFightRequest(_ message.TaggedMessage, _ *state.View, _ *immutable.Map[commons.ID, decision.FightAction]) message.Payload {
	return nil
}

func (r *Team6Agent) CurrentAction() decision.FightAction {
	fight := rand.Intn(3)
	switch fight {
	case 0:
		return decision.Cower
	case 1:
		return decision.Attack
	default:
		return decision.Defend
	}
}

func (r *Team6Agent) HandleElectionBallot(view *state.View, _ agent.BaseAgent, _ *decision.ElectionParams) decision.Ballot {
	// Extract ID of alive agents
	agentState := view.AgentState()
	aliveAgentIds := make([]string, agentState.Len())
	i := 0
	itr := agentState.Iterator()
	for !itr.Done() {
		id, a, ok := itr.Next()
		if ok && a.Hp > 0 {
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

	return ballot
}

func NewRandomAgent() agent.Strategy {
	return &Team6Agent{bravery: 0}
}
