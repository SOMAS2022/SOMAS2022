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
	// makesProposal := rand.Intn(100)

	// if makesProposal > 80 {
	// 	rules := make([]proposal.Rule[decision.FightAction], 0)

	// 	rules = append(rules, *proposal.NewRule(decision.Attack,
	// 		proposal.NewAndCondition(*proposal.NewComparativeCondition(proposal.Health, proposal.GreaterThan, 1000),
	// 			*proposal.NewComparativeCondition(proposal.Stamina, proposal.GreaterThan, 1000)),
	// 	))

	// 	rules = append(rules, *proposal.NewRule(decision.Defend,
	// 		proposal.NewComparativeCondition(proposal.TotalDefence, proposal.GreaterThan, 1000),
	// 	))

	// 	rules = append(rules, *proposal.NewRule(decision.Cower,
	// 		proposal.NewComparativeCondition(proposal.Health, proposal.LessThan, 1),
	// 	))

	// 	rules = append(rules, *proposal.NewRule(decision.Attack,
	// 		proposal.NewComparativeCondition(proposal.Stamina, proposal.GreaterThan, 10),
	// 	))

	// 	prop := *commons.NewImmutableList(rules)
	// 	_ = baseAgent.SendFightProposalToLeader(prop)
	// }
}

func (a *Team6Agent) HandleFightRequest(m message.TaggedRequestMessage[message.FightRequest], log *immutable.Map[commons.ID, decision.FightAction]) message.FightInform {
	return nil
}

func (a *Team6Agent) FightResolution(agent agent.BaseAgent, prop commons.ImmutableList[proposal.Rule[decision.FightAction]], proposedActions immutable.Map[commons.ID, decision.FightAction]) immutable.Map[commons.ID, decision.FightAction] {
	return proposedActions
}

func (a *Team6Agent) HandleFightProposal(proposal message.Proposal[decision.FightAction], baseAgent agent.BaseAgent) decision.Intent {
	intent := rand.Intn(2)
	if intent == 0 {
		return decision.Positive
	} else {
		return decision.Negative
	}
}

// HandleFightProposalRequest only called as leader
func (a *Team6Agent) HandleFightProposalRequest(proposal message.Proposal[decision.FightAction], baseAgent agent.BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) bool {
	// Do we want to forward proposal to other agents to get opinion?

	// TODO: Replace one of these with agent's own proposal
	return proposalSimilarity(proposal.Rules(), proposal.Rules()) > 0.6
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
