package agent

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/message/proposal"

	"github.com/benbjohnson/immutable"
)

type Strategy interface {
	Fight
	Election
	Loot
	HPPool
	// HandleUpdateWeapon return the index of the weapon you want to use in AgentState.weapons
	HandleUpdateWeapon(baseAgent BaseAgent) decision.ItemIdx
	// HandleUpdateShield return the index of the shield you want to use in AgentState.Shields
	HandleUpdateShield(baseAgent BaseAgent) decision.ItemIdx

	UpdateInternalState(baseAgent BaseAgent, fightResult *commons.ImmutableList[decision.ImmutableFightResult], voteResult *immutable.Map[decision.Intent, uint])
}

type Election interface {
	CreateManifesto(baseAgent BaseAgent) *decision.Manifesto
	HandleConfidencePoll(baseAgent BaseAgent) decision.Intent
	HandleElectionBallot(baseAgent BaseAgent, params *decision.ElectionParams) decision.Ballot
}

type Fight interface {
	HandleFightInformation(m message.TaggedInformMessage[message.FightInform], baseAgent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction])
	HandleFightRequest(m message.TaggedRequestMessage[message.FightRequest], log *immutable.Map[commons.ID, decision.FightAction]) message.FightInform
	FightResolution(agent BaseAgent, prop commons.ImmutableList[proposal.Rule[decision.FightAction]]) immutable.Map[commons.ID, decision.FightAction]
	HandleFightProposal(proposal message.Proposal[decision.FightAction], baseAgent BaseAgent) decision.Intent
	// HandleFightProposalRequest only called as leader
	HandleFightProposalRequest(proposal message.Proposal[decision.FightAction], baseAgent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) bool
	FightActionNoProposal(baseAgent BaseAgent) decision.FightAction
	FightAction(baseAgent BaseAgent, proposedAction decision.FightAction) decision.FightAction
}

type Loot interface {
	HandleLootInformation(m message.TaggedInformMessage[message.LootInform], agent BaseAgent)
	HandleLootRequest(m message.TaggedRequestMessage[message.LootRequest]) message.LootInform
	HandleLootProposal(r message.Proposal[decision.LootAction], agent BaseAgent) decision.Intent
	HandleLootProposalRequest(proposal message.Proposal[decision.LootAction], agent BaseAgent) bool
	LootAllocation(agent BaseAgent) immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]]
	LootActionNoProposal(baseAgent BaseAgent) immutable.SortedMap[commons.ItemID, struct{}]
	LootAction(baseAgent BaseAgent, proposedLoot immutable.SortedMap[commons.ItemID, struct{}]) immutable.SortedMap[commons.ItemID, struct{}]
}

type HPPool interface {
	DonateToHpPool(baseAgent BaseAgent) uint
}
