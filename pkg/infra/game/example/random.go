package example

import (
	"math/rand"

	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"infra/game/tally"
	"infra/logging"

	"github.com/benbjohnson/immutable"
	"github.com/google/uuid"
)

type RandomAgent struct {
	bravery int
}

func (r *RandomAgent) FightResolution(baseAgent agent.BaseAgent) tally.Proposal[decision.FightAction] {
	actions := make(map[commons.ID]decision.FightAction)
	view := baseAgent.View()
	agentState := view.AgentState()
	itr := agentState.Iterator()
	for !itr.Done() {
		id, _, ok := itr.Next()
		if !ok {
			break
		}

		switch rand.Intn(3) {
		case 0:
			actions[id] = decision.Attack
		case 1:
			actions[id] = decision.Defend
		default:
			actions[id] = decision.Cower
		}
	}
	newUUID, _ := uuid.NewUUID()
	prop := tally.NewProposal[decision.FightAction](newUUID.String(), commons.MapToImmutable(actions))
	return *prop
}

func (r *RandomAgent) CreateManifesto(_ agent.BaseAgent) *decision.Manifesto {
	manifesto := decision.NewManifesto(false, false, 10, 50)
	return manifesto
}

func (r *RandomAgent) HandleConfidencePoll(_ agent.BaseAgent) decision.Intent {
	switch rand.Intn(3) {
	case 0:
		return decision.Abstain
	case 1:
		return decision.Negative
	default:
		return decision.Positive
	}
}

func (r *RandomAgent) HandleFightInformation(_ message.TaggedMessage, baseAgent agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) {
	baseAgent.Log(logging.Trace, logging.LogField{"bravery": r.bravery, "hp": baseAgent.AgentState().Hp}, "Cowering")
	makesProposal := rand.Intn(100)

	if makesProposal > 80 {
		prop := r.FightResolution(baseAgent)
		view := baseAgent.View()
		_ = baseAgent.SendBlockingMessage(view.CurrentLeader(), *message.NewMessage(message.Proposal, *message.NewProposalPayload(prop.Proposal())))
	}
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
	aliveAgentIDs := make([]string, agentState.Len())
	i := 0
	itr := agentState.Iterator()
	for !itr.Done() {
		id, a, ok := itr.Next()
		if ok && a.Hp > 0 {
			aliveAgentIDs[i] = id
			i++
		}
	}

	// Randomly fill the ballot
	var ballot decision.Ballot
	numAliveAgents := len(aliveAgentIDs)
	numCandidate := 2
	for i := 0; i < numCandidate; i++ {
		randomIdx := rand.Intn(numAliveAgents)
		randomCandidate := aliveAgentIDs[uint(randomIdx)]
		ballot = append(ballot, randomCandidate)
	}

	return ballot
}

func (r *RandomAgent) HandleFightProposal(_ *message.FightProposalMessage, _ agent.BaseAgent) decision.Intent {
	intent := rand.Intn(2)
	if intent == 0 {
		return decision.Positive
	} else {
		return decision.Negative
	}
}

func (r *RandomAgent) HandleFightProposalRequest(_ *message.FightProposalMessage, _ agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) bool {
	switch rand.Intn(2) {
	case 0:
		return true
	default:
		return false
	}
}

func (r *RandomAgent) HandleUpdateWeapon(view *state.View, b agent.BaseAgent) decision.ItemIdx {
	// weapons := b.AgentState().Weapons
	// return decision.ItemIdx(rand.Intn(weapons.Len() + 1))

	// 0th weapon has greatest attack points
	return decision.ItemIdx(0)
}

func (r *RandomAgent) HandleUpdateShield(view *state.View, b agent.BaseAgent) decision.ItemIdx {
	// shields := b.AgentState().Shields
	// return decision.ItemIdx(rand.Intn(shields.Len() + 1))
	return decision.ItemIdx(0)
}

func NewRandomAgent() agent.Strategy {
	return &RandomAgent{bravery: 0}
}
