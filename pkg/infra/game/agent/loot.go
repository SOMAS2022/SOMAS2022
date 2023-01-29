package agent

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"

	"github.com/benbjohnson/immutable"
)

type Loot interface {
	HandleLootInformation(m message.TaggedInformMessage[message.LootInform], baseAgent BaseAgent)
	HandleLootRequest(m message.TaggedRequestMessage[message.LootRequest]) message.LootInform
	HandleLootProposal(r message.Proposal[decision.LootAction], baseAgent BaseAgent) decision.Intent
	HandleLootProposalRequest(proposal message.Proposal[decision.LootAction], baseAgent BaseAgent) bool
	LootAllocation(
		baseAgent BaseAgent,
		proposal message.Proposal[decision.LootAction],
		proposedAllocations map[commons.ID]map[commons.ItemID]struct{},
	) immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]]
	LootActionNoProposal(baseAgent BaseAgent) immutable.SortedMap[commons.ItemID, struct{}]
	LootAction(baseAgent BaseAgent, proposedLoot immutable.SortedMap[commons.ItemID, struct{}], acceptedProposal message.Proposal[decision.LootAction]) immutable.SortedMap[commons.ItemID, struct{}]
	PruneAgentList(agentMap map[commons.ID]Agent) map[commons.ID]Agent
	chooseItem(BaseAgent BaseAgent) (bool, bool, bool, bool)
}
