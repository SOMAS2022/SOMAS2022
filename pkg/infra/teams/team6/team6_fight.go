package team6

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/message/proposal"
	"math/rand"

	"github.com/benbjohnson/immutable"
)

func (a *Team6Agent) HandleFightInformation(m message.TaggedInformMessage[message.FightInform], baseAgent agent.BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) {
	switch m.Message().(type) {
	case *message.StartFight:
		a.generateFightProposal()
		sendsProposal := rand.Intn(100)
		if sendsProposal > 90 {
			baseAgent.SendFightProposalToLeader(a.fightProposal)
		}
	default:
		return
	}
}

func (a *Team6Agent) HandleFightRequest(m message.TaggedRequestMessage[message.FightRequest], log *immutable.Map[commons.ID, decision.FightAction]) message.FightInform {
	return nil
}

func (a *Team6Agent) FightResolution(agent agent.BaseAgent, prop commons.ImmutableList[proposal.Rule[decision.FightAction]], proposedActions immutable.Map[commons.ID, decision.FightAction]) immutable.Map[commons.ID, decision.FightAction] {
	return proposedActions
}

func (a *Team6Agent) HandleFightProposal(proposal message.Proposal[decision.FightAction], baseAgent agent.BaseAgent) decision.Intent {
	a.currentProposalsReceived++

	similarity := proposalSimilarity(a.fightProposal, proposal.Rules())

	//Update similarity SC value
	init := SafeMapReadOrDefault(a.similarity, proposal.ProposerID(), 50)
	diff := int(10 * (similarity - 0.5))
	if diff < 0 {
		a.similarity[proposal.ProposerID()] = commons.SaturatingSub(init, uint(-diff))
	} else {
		a.similarity[proposal.ProposerID()] = SCSaturatingAdd(init, uint(diff), 100)
	}

	if similarity >= 0.8 {
		return decision.Positive
	} else {
		return decision.Negative
	}
}

// HandleFightProposalRequest only called as leader
func (a *Team6Agent) HandleFightProposalRequest(proposal message.Proposal[decision.FightAction], baseAgent agent.BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) bool {
	similarity := proposalSimilarity(a.fightProposal, proposal.Rules())

	//Update similarity SC value
	init := SafeMapReadOrDefault(a.similarity, proposal.ProposerID(), 50)
	diff := int(10 * (similarity - 0.5))
	if diff < 0 {
		a.similarity[proposal.ProposerID()] = commons.SaturatingSub(init, uint(-diff))
	} else {
		a.similarity[proposal.ProposerID()] = SCSaturatingAdd(init, uint(diff), 100)
	}

	return similarity > 0.6
}

func (a *Team6Agent) FightActionNoProposal(baseAgent agent.BaseAgent) decision.FightAction {
	state := baseAgent.AgentState()
	if state.Hp > uint((a.HPThreshold*float32(startingHP-1))) && state.Attack > uint((a.ATThreshold*float32(startingAT-1))) && state.Stamina > uint((a.STThreshold*float32(startingST-1))) {
		return decision.Attack
	} else if state.Hp > uint((a.HPThreshold*float32(startingHP-1))) && state.Stamina > uint((a.STThreshold*float32(startingST-1))) {
		return decision.Defend
	} else {
		return decision.Cower
	}
}

func (a *Team6Agent) FightAction(baseAgent agent.BaseAgent, proposedAction decision.FightAction, acceptedProposal message.Proposal[decision.FightAction]) decision.FightAction {
	/*
		If our agent thinks they should cower, let's let them trust their intuition over the group's.
	*/
	a.lastFightRound = 0
	if a.FightActionNoProposal(baseAgent) == decision.Cower && proposedAction != decision.Cower {
		return decision.Cower
	} else {
		return a.FightActionNoProposal(baseAgent)
	}
}
