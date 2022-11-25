package example

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

type RandomAgent struct {
	bravery int
}

func (r RandomAgent) HandleFightInformation(_ message.TaggedMessage, _ *state.View, agent agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) {
	agent.Log(logging.Trace, logging.LogField{"bravery": r.bravery, "hp": agent.ViewState().Hp}, "Cowering")
}

func (r RandomAgent) HandleFightRequest(_ message.TaggedMessage, _ *state.View, _ *immutable.Map[commons.ID, decision.FightAction]) message.Payload {
	return nil
}

func (r RandomAgent) CurrentAction() decision.FightAction {
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

func (r RandomAgent) HandleElection(view *state.View, baseAgent agent.BaseAgent, decisionC chan<- decision.Ballot, candidateList []commons.ID, strategy decision.VotingStrategy, qtyPreferences uint) {
	// Extract ID of alive agents
	agentState := view.AgentState()
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

func NewRandomAgent() *RandomAgent {
	return &RandomAgent{bravery: 0}
}
