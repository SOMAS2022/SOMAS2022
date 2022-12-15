package team6

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"math/rand"

	"github.com/benbjohnson/immutable"
)

func (a *Team6Agent) LootActionNoProposal(baseAgent agent.BaseAgent) immutable.SortedMap[commons.ItemID, struct{}] {
	return *immutable.NewSortedMap[commons.ItemID, struct{}](nil)
}

func (a *Team6Agent) LootAction(
	baseAgent agent.BaseAgent,
	proposedLoot immutable.SortedMap[commons.ItemID, struct{}],
	acceptedProposal message.Proposal[decision.LootAction],
) immutable.SortedMap[commons.ItemID, struct{}] {
	return proposedLoot
}

func (a *Team6Agent) HandleLootInformation(m message.TaggedInformMessage[message.LootInform], baseAgent agent.BaseAgent) {
	switch m.Message().(type) {
	case *message.StartLoot:
		a.generateLootProposal()
		sendsProposal := rand.Intn(100)
		if sendsProposal > 90 {
			baseAgent.SendLootProposalToLeader(a.lootProposal)
		}
	default:
		return
	}
}

func (a *Team6Agent) HandleLootRequest(m message.TaggedRequestMessage[message.LootRequest]) message.LootInform {
	//TODO implement me
	panic("implement me")
}

func (a *Team6Agent) HandleLootProposal(proposal message.Proposal[decision.LootAction], _ agent.BaseAgent) decision.Intent {
	similarity := proposalSimilarity(a.lootProposal, proposal.Rules())
	if similarity >= 0.8 {
		return decision.Positive
	} else {
		return decision.Negative
	}
}

func (a *Team6Agent) HandleLootProposalRequest(proposal message.Proposal[decision.LootAction], _ agent.BaseAgent) bool {
	// TODO: Replace one of these with agent's own proposal
	return proposalSimilarity(proposal.Rules(), proposal.Rules()) > 0.6
}

func (a *Team6Agent) LootAllocation(baseAgent agent.BaseAgent, proposal message.Proposal[decision.LootAction], proposedAllocations immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]]) immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]] {
	return proposedAllocations
}
