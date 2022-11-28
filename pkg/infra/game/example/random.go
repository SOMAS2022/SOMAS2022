package example

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/logging"
	"math/rand"

	"github.com/benbjohnson/immutable"
)

type RandomAgent struct {
	bravery int
}

func (r *RandomAgent) CreateManifesto(baseAgent agent.BaseAgent) *decision.Manifesto {
	manifesto := decision.NewManifesto(true, false, 10, 50)
	return manifesto
}

func (r *RandomAgent) HandleConfidencePoll(baseAgent agent.BaseAgent) decision.Intent {
	switch rand.Intn(3) {
	case 0:
		return decision.Abstain
	case 1:
		return decision.Negative
	default:
		return decision.Positive
	}
}

func (r *RandomAgent) HandleFightInformation(_ message.TaggedMessage, agent agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) {
	agent.Log(logging.Trace, logging.LogField{"bravery": r.bravery, "hp": agent.AgentState().Hp}, "Cowering")
}

func (r *RandomAgent) HandleFightRequest(_ message.TaggedMessage, _ *immutable.Map[commons.ID, decision.FightAction]) message.Payload {
	return nil
}

func (r *RandomAgent) CurrentAction() decision.FightAction {
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

func (r *RandomAgent) HandleElectionBallot(b agent.BaseAgent, _ *decision.ElectionParams) decision.Ballot {
	// Extract ID of alive agents
	view := b.View()
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

func (r *RandomAgent) HandleFightProposal(proposal *message.FightProposalMessage, baseAgent agent.BaseAgent) decision.Intent {
	intent := rand.Intn(2)
	if intent == 0 {
		return decision.Positive
	} else {
		return decision.Negative
	}
}

func NewRandomAgent() agent.Strategy {
	return &RandomAgent{bravery: 0}
}
