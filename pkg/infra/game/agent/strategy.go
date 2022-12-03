package agent

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"

	"github.com/benbjohnson/immutable"
)

type Strategy interface {
	HandleFightInformation(m message.TaggedInformMessage[message.FightInform], baseAgent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction])
	HandleFightRequest(m message.TaggedRequestMessage[message.FightRequest], log *immutable.Map[commons.ID, decision.FightAction]) message.FightInform
	CurrentAction() decision.FightAction
	CreateManifesto(baseAgent BaseAgent) *decision.Manifesto
	HandleConfidencePoll(baseAgent BaseAgent) decision.Intent
	HandleElectionBallot(baseAgent BaseAgent, params *decision.ElectionParams) decision.Ballot
	HandleFightProposal(proposal message.FightProposalMessage, baseAgent BaseAgent) decision.Intent
	FightResolution(agent BaseAgent) message.MapProposal[decision.FightAction]
	// HandleFightProposalRequest only called as leader
	HandleFightProposalRequest(proposal message.FightProposalMessage, baseAgent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) bool
	// HandleUpdateWeapon return the index of the weapon you want to use in AgentState.Weapons
	HandleUpdateWeapon(baseAgent BaseAgent) decision.ItemIdx
	// HandleUpdateShield return the index of the shield you want to use in AgentState.Shields
	HandleUpdateShield(baseAgent BaseAgent) decision.ItemIdx
	UpdateInternalState(baseAgent BaseAgent, fightResult *commons.ImmutableList[decision.ImmutableFightResult], voteResult *immutable.Map[decision.Intent, uint])
	DonateToHpPool(baseAgent BaseAgent) uint
	// Only implementable in team experiments
	HandleCustomInformation(m message.TaggedInformMessage[message.CustomInform], baseAgent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction])
}
