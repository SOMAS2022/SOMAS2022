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
	bravery      int
	utilityScore map[commons.ID]int
}

// Update Utility
func (a *AgentThree) UpdateUtility(baseAgent agent.BaseAgent) {

	view := baseAgent.View()
	agentState := view.AgentState()
	itr := agentState.Iterator()
	for !itr.Done() {
		id, _, ok := itr.Next()
		if !ok {
			break
		}

		// We are already cool, don't need the utility score for ourselves
		if id != baseAgent.ID() {
			a.utilityScore[id] = rand.Intn(10)
		}
	}
}

// Create proposal for fight decisions
func (a *AgentThree) FightResolution(baseAgent agent.BaseAgent) tally.Proposal[decision.FightAction] {

	actions := make(map[commons.ID]decision.FightAction)
	view := baseAgent.View()
	agentState := view.AgentState()
	itr := agentState.Iterator()
	for !itr.Done() {
		id, _, ok := itr.Next()
		if !ok {
			break
		}

		// Check for our agent and assign what we want to do
		if id == baseAgent.ID() {
			actions[id] = a.CurrentAction()
			baseAgent.Log(logging.Trace, logging.LogField{"bravery": a.bravery, "hp": baseAgent.AgentState().Hp, "choice": a.CurrentAction()}, "Intent")
		} else {
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
	}
	newUUID, _ := uuid.NewUUID()
	prop := tally.NewProposal(newUUID.String(), commons.MapToImmutable(actions))
	return *prop
}

// Manifesto
func (a *AgentThree) CreateManifesto(_ agent.BaseAgent) *decision.Manifesto {
	manifesto := decision.NewManifesto(false, false, 10, 50)
	return manifesto
}

// Handle No Confidence vote
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

// Send proposal to leader
func (a *AgentThree) HandleFightInformation(_ message.TaggedMessage, baseAgent agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) {
	// baseAgent.Log(logging.Trace, logging.LogField{"bravery": a.bravery, "hp": baseAgent.AgentState().Hp}, "Cowering")
	makesProposal := rand.Intn(100)

	if makesProposal > 80 {
		prop := a.FightResolution(baseAgent)
		view := baseAgent.View()
		_ = baseAgent.SendBlockingMessage(view.CurrentLeader(), *message.NewMessage(message.Proposal, *message.NewProposalPayload(prop.Proposal())))
	}
}

// Calculate our agents action
func (a *AgentThree) CurrentAction() decision.FightAction {

	if a.bravery > 2 {
		fight := rand.Intn(2)
		switch fight {
		case 0:
			return decision.Attack
		default:
			return decision.Defend
		}
	} else {
		return decision.Cower
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

// Vote on proposal
func (a *AgentThree) HandleFightProposal(m *message.FightProposalMessage, baseAgent agent.BaseAgent) decision.Intent {

	prop := m.Proposal()
	act, _ := prop.Get(baseAgent.ID())

	// Selfish, only agree if our decision is ok
	if act == a.CurrentAction() {
		return decision.Positive
	} else {
		return decision.Negative
	}
}

// Leader function to grant the floor
func (a *AgentThree) HandleFightProposalRequest(_ *message.FightProposalMessage, _ agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) bool {
	switch rand.Intn(2) {
	case 0:
		return true
	default:
		return false
	}
}

func (a *AgentThree) HandleFightRequest(_ message.TaggedMessage, _ *immutable.Map[commons.ID, decision.FightAction]) message.Payload {
	return nil
}

func NewAgentThree() agent.Strategy {
	return &AgentThree{
		bravery:      rand.Intn(10),
		utilityScore: make(map[string]int),
	}
}
