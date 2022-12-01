package team3

import (
	"math/rand"

	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/tally"
	"infra/logging"

	"github.com/benbjohnson/immutable"
	"github.com/google/uuid"
)

type AgentThree struct {
	bravery int
}

func (a *AgentThree) FightResolution(baseAgent agent.BaseAgent) tally.Proposal[decision.FightAction] {

	baseAgent.Log(logging.Trace, logging.LogField{"bravery": a.bravery, "hp": baseAgent.AgentState().Hp}, "Cowering")

	actions := make(map[commons.ID]decision.FightAction)
	view := baseAgent.View()
	agentState := view.AgentState()
	itr := agentState.Iterator()
	for !itr.Done() {
		id, _, ok := itr.Next()
		if !ok {
			break
		}
		rNum := rand.Intn(3)
		switch rNum {
		case 0:
			actions[id] = decision.Attack
		case 1:
			actions[id] = decision.Defend
		default:
			actions[id] = decision.Cower
		}
	}
	newUUID, _ := uuid.NewUUID()
	prop := tally.NewProposal(newUUID.String(), commons.MapToImmutable(actions))
	return *prop
}

func (a *AgentThree) CreateManifesto(_ agent.BaseAgent) *decision.Manifesto {
	manifesto := decision.NewManifesto(true, false, 10, 50)
	return manifesto
}

func (a *AgentThree) HandleConfidencePoll(_ agent.BaseAgent) decision.Intent {
	switch rand.Intn(3) {
	case 0:
		return decision.Abstain
	case 1:
		return decision.Negative
	default:
		return decision.Positive
	}
}

func (a *AgentThree) HandleFightInformation(_ message.TaggedMessage, baseAgent agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) {
	baseAgent.Log(logging.Trace, logging.LogField{"bravery": a.bravery, "hp": baseAgent.AgentState().Hp}, "Cowering")
	makesProposal := rand.Intn(100)

	if makesProposal > 80 {
		prop := a.FightResolution(baseAgent)
		view := baseAgent.View()
		_ = baseAgent.SendBlockingMessage(view.CurrentLeader(), *message.NewMessage(message.Proposal, *message.NewProposalPayload(prop.Proposal())))
	}
}

func (a *AgentThree) HandleFightRequest(_ message.TaggedMessage, _ *immutable.Map[commons.ID, decision.FightAction]) message.Payload {
	return nil
}

func (a *AgentThree) CurrentAction() decision.FightAction {
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

func (a *AgentThree) HandleElectionBallot(b agent.BaseAgent, _ *decision.ElectionParams) decision.Ballot {
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

func (a *AgentThree) HandleFightProposal(_ *message.FightProposalMessage, _ agent.BaseAgent) decision.Intent {
	intent := rand.Intn(2)
	if intent == 0 {
		return decision.Positive
	} else {
		return decision.Negative
	}
}

func (a *AgentThree) HandleFightProposalRequest(_ *message.FightProposalMessage, _ agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) bool {
	switch rand.Intn(2) {
	case 0:
		return true
	default:
		return false
	}
}

func NewAgentThree() agent.Strategy {
	return &AgentThree{bravery: 0}
}
