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
	return nil
}

func (a *Team6Agent) HandleLootProposal(proposal message.Proposal[decision.LootAction], _ agent.BaseAgent) decision.Intent {
	a.currentProposalsReceived++

	similarity := proposalSimilarity(a.lootProposal, proposal.Rules())

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

func (a *Team6Agent) HandleLootProposalRequest(proposal message.Proposal[decision.LootAction], _ agent.BaseAgent) bool {
	similarity := proposalSimilarity(a.lootProposal, proposal.Rules())

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

func (a *Team6Agent) LootAllocation(baseAgent agent.BaseAgent, proposal message.Proposal[decision.LootAction], proposedAllocations immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]]) immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]] {
	return proposedAllocations
}
