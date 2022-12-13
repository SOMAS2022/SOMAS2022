package team5

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/message/proposal"
	"infra/game/state"
	"math/rand"

	"github.com/benbjohnson/immutable"
)

type FivAgent struct {
	qtable *Qtable
}

// --------------------------------- Election Stage ---------------------------------

func (fiv *FivAgent) CreateManifesto(_ agent.BaseAgent) *decision.Manifesto {
	manifesto := decision.NewManifesto(false, true, 10, 50)
	return manifesto
}

func (fiv *FivAgent) HandleConfidencePoll(_ agent.BaseAgent) decision.Intent {
	switch rand.Intn(3) {
	case 0:
		return decision.Abstain
	case 1:
		return decision.Negative
	default:
		return decision.Positive
	}
}

func (fiv *FivAgent) HandleElectionBallot(b agent.BaseAgent, _ *decision.ElectionParams) decision.Ballot {
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
	numCandidate := rand.Intn(numAliveAgents)
	for i := 0; i < numCandidate; i++ {
		randomIdx := rand.Intn(numAliveAgents)
		randomCandidate := aliveAgentIDs[uint(randomIdx)]
		ballot = append(ballot, randomCandidate)
	}

	return ballot
}

// --------------------------------- Fight Stage ---------------------------------

func (fiv *FivAgent) HandleFightInformation(_ message.TaggedInformMessage[message.FightInform], baseAgent agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) {
	// baseAgent.Log(logging.Trace, logging.LogField{"bravery": r.bravery, "hp": baseAgent.AgentState().Hp}, "Cowering")
	makesProposal := rand.Intn(100)

	if makesProposal > 80 {
		prop := fiv.FightResolution(baseAgent)
		_ = baseAgent.SendFightProposalToLeader(prop)
	}
}

func (fiv *FivAgent) HandleFightRequest(_ message.TaggedRequestMessage[message.FightRequest], _ *immutable.Map[commons.ID, decision.FightAction]) message.FightInform {
	return nil
}

func (fiv *FivAgent) FightResolution(_ agent.BaseAgent) commons.ImmutableList[proposal.Rule[decision.FightAction]] {
	rules := make([]proposal.Rule[decision.FightAction], 0)

	rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Attack,
		proposal.NewAndCondition(*proposal.NewComparativeCondition(proposal.Health, proposal.GreaterThan, 1000),
			*proposal.NewComparativeCondition(proposal.Stamina, proposal.GreaterThan, 1000)),
	))

	rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Defend,
		proposal.NewComparativeCondition(proposal.TotalDefence, proposal.GreaterThan, 1000),
	))

	rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Cower,
		proposal.NewComparativeCondition(proposal.Health, proposal.LessThan, 1),
	))

	rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Attack,
		proposal.NewComparativeCondition(proposal.Stamina, proposal.GreaterThan, 10),
	))

	return *commons.NewImmutableList(rules)
}

func (fiv *FivAgent) HandleFightProposal(_ message.Proposal[decision.FightAction], _ agent.BaseAgent) decision.Intent {
	intent := rand.Intn(2)
	if intent == 0 {
		return decision.Positive
	} else {
		return decision.Negative
	}
}

func (fiv *FivAgent) HandleFightProposalRequest(
	_ message.Proposal[decision.FightAction],
	_ agent.BaseAgent,
	_ *immutable.Map[commons.ID, decision.FightAction],
) bool {
	// leader fight function
	switch rand.Intn(2) {
	case 0:
		return true
	default:
		return false
	}
}

func (fiv *FivAgent) FightAction(baseAgent agent.BaseAgent) decision.FightAction {
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

// --------------------------------- Loot Stage ---------------------------------

func (fiv *FivAgent) HandleLootInformation(m message.TaggedInformMessage[message.LootInform], agent agent.BaseAgent) {
}

func (fiv *FivAgent) HandleLootRequest(m message.TaggedRequestMessage[message.LootRequest]) message.LootInform {
	//TODO implement me
	panic("implement me")
}

func (fiv *FivAgent) HandleLootProposal(_ message.Proposal[decision.LootAction], _ agent.BaseAgent) decision.Intent {
	switch rand.Intn(3) {
	case 0:
		return decision.Positive
	case 1:
		return decision.Negative
	default:
		return decision.Abstain
	}
}

func (fiv *FivAgent) HandleLootProposalRequest(_ message.Proposal[decision.LootAction], _ agent.BaseAgent) bool {
	// leader loot function
	switch rand.Intn(2) {
	case 0:
		return true
	default:
		return false
	}
}

func (fiv *FivAgent) LootAction() immutable.List[commons.ItemID] {
	return *immutable.NewList[commons.ItemID]()
}

func (fiv *FivAgent) LootAllocation(ba agent.BaseAgent) immutable.Map[commons.ID, immutable.List[commons.ItemID]] {
	lootAllocation := make(map[commons.ID][]commons.ItemID)
	view := ba.View()
	ids := commons.ImmutableMapKeys(view.AgentState())
	iterator := ba.Loot().Weapons().Iterator()
	allocateRandomly(iterator, ids, lootAllocation)
	iterator = ba.Loot().Shields().Iterator()
	allocateRandomly(iterator, ids, lootAllocation)
	iterator = ba.Loot().HpPotions().Iterator()
	allocateRandomly(iterator, ids, lootAllocation)
	iterator = ba.Loot().StaminaPotions().Iterator()
	allocateRandomly(iterator, ids, lootAllocation)
	mMapped := make(map[commons.ID]immutable.List[commons.ItemID])
	for id, itemIDS := range lootAllocation {
		mMapped[id] = commons.ListToImmutable(itemIDS)
	}
	return commons.MapToImmutable(mMapped)
}

func allocateRandomly(iterator commons.Iterator[state.Item], ids []commons.ID, lootAllocation map[commons.ID][]commons.ItemID) {
	for !iterator.Done() {
		next, _ := iterator.Next()
		toBeAllocated := ids[rand.Intn(len(ids))]
		if l, ok := lootAllocation[toBeAllocated]; ok {
			l = append(l, next.Id())
			lootAllocation[toBeAllocated] = l
		} else {
			l := make([]commons.ItemID, 0)
			l = append(l, next.Id())
			lootAllocation[toBeAllocated] = l
		}
	}
}

// --------------------------------- Update Stage ---------------------------------

func (fiv *FivAgent) UpdateInternalState(_ agent.BaseAgent, _ *commons.ImmutableList[decision.ImmutableFightResult], _ *immutable.Map[decision.Intent, uint]) {
}

func (fiv *FivAgent) HandleUpdateWeapon(_ agent.BaseAgent) decision.ItemIdx {
	// weapons := b.AgentState().weapons
	// return decision.ItemIdx(rand.Intn(weapons.Len() + 1))

	// 0th weapon has the greatest attack points
	return decision.ItemIdx(0)
}

func (fiv *FivAgent) HandleUpdateShield(_ agent.BaseAgent) decision.ItemIdx {
	// shields := b.AgentState().Shields
	// return decision.ItemIdx(rand.Intn(shields.Len() + 1))
	return decision.ItemIdx(0)
}

// --------------------------------- Hp pool Stage ---------------------------------

func (fiv *FivAgent) DonateToHpPool(baseAgent agent.BaseAgent) uint {
	// return uint(rand.Intn(int(baseAgent.AgentState().Hp)))
	return uint(int(0.012 * float32(baseAgent.AgentState().Hp)))
}

// --------------------------------- Misc ---------------------------------

func NewFivAgent() agent.Strategy {
	fiv := new(FivAgent)
	fiv.qtable = NewQTable(0.1, 0.6)
	return fiv
}
