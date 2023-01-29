package team3

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/message/proposal"
	"math/rand"

	"github.com/benbjohnson/immutable"
)

func (a *AgentThree) LootActionNoProposal(baseAgent agent.BaseAgent) immutable.SortedMap[commons.ItemID, struct{}] {
	loot := baseAgent.Loot()
	weapons := loot.Weapons().Iterator()
	shields := loot.Shields().Iterator()
	hpPotions := loot.HpPotions().Iterator()
	staminaPotions := loot.StaminaPotions().Iterator()

	builder := immutable.NewSortedMapBuilder[commons.ItemID, struct{}](nil)

	for !weapons.Done() {
		weapon, _ := weapons.Next()
		if rand.Int()%2 == 0 {
			builder.Set(weapon.Id(), struct{}{})
		}
	}

	for !shields.Done() {
		shield, _ := shields.Next()
		if rand.Int()%2 == 0 {
			builder.Set(shield.Id(), struct{}{})
		}
	}

	for !hpPotions.Done() {
		pot, _ := hpPotions.Next()
		if rand.Int()%2 == 0 {
			builder.Set(pot.Id(), struct{}{})
		}
	}

	for !staminaPotions.Done() {
		pot, _ := staminaPotions.Next()
		if rand.Int()%2 == 0 {
			builder.Set(pot.Id(), struct{}{})
		}
	}

	return *builder.Map()
}

func (a *AgentThree) LootAction(
	baseAgent agent.BaseAgent,
	proposedLoot immutable.SortedMap[commons.ItemID, struct{}],
	acceptedProposal message.Proposal[decision.LootAction],
) immutable.SortedMap[commons.ItemID, struct{}] {
	return proposedLoot
}

func (a *AgentThree) HandleLootInformation(m message.TaggedInformMessage[message.LootInform], baseAgent agent.BaseAgent) {
	// submit a proposal to the leader
	switch m.Message().(type) {
	case *message.StartLoot:
		// Send Proposal?
		sendProposal := rand.Intn(100)
		if sendProposal < a.personality {
			// general and send a loot proposal
			baseAgent.SendLootProposalToLeader(a.generateLootProposal())
		}
	default:
		return
	}
}

func (a *AgentThree) HandleLootProposal(_ message.Proposal[decision.LootAction], _ agent.BaseAgent) decision.Intent {
	// vote on the loot proposal

	// do i vote?
	toVote := rand.Intn(100)
	if toVote < a.personality {
		// Enter logic for evaluating a loot proposal here
		switch rand.Intn(2) {
		case 0:
			return decision.Positive
		default:
			return decision.Negative
		}
	} else {
		// abstain
		return decision.Abstain
	}
}

func (a *AgentThree) generateLootProposal() commons.ImmutableList[proposal.Rule[decision.LootAction]] {
	rules := make([]proposal.Rule[decision.LootAction], 0)

	rules = append(rules, *proposal.NewRule(decision.HealthPotion,
		proposal.NewComparativeCondition(proposal.Health, proposal.LessThan, uint(0.5*float64(GetStartingHP())))))

	rules = append(rules, *proposal.NewRule(decision.StaminaPotion,
		proposal.NewComparativeCondition(proposal.Stamina, proposal.LessThan, uint(0.5*float64(GetStartingStamina())))))

	rules = append(rules, *proposal.NewRule(decision.Weapon,
		proposal.NewComparativeCondition(proposal.TotalAttack, proposal.LessThan, uint(0.5*float64(GetStartingHP())))))

	rules = append(rules, *proposal.NewRule(decision.Shield,
		proposal.NewComparativeCondition(proposal.TotalDefence, proposal.LessThan, uint(0.5*float64(GetStartingHP())))))

	return *commons.NewImmutableList(rules)
}

func (a *AgentThree) chooseItem(baseAgent agent.BaseAgent) (bool, bool, bool, bool) {
	// function to calculate the agents choice of loot

	// calculate the average stats of the group
	avHP := AverageArray(GetHealthAllAgents(baseAgent))
	avST := AverageArray(GetStaminaAllAgents(baseAgent))
	avATT := AverageArray(GetAttackAllAgents(baseAgent))
	avDEF := AverageArray(GetDefenceAllAgents(baseAgent))

	// normalise the group stats
	groupAvHP, groupAvST, groupAvATT, groupAvDEF := normalize4El(avHP, avST, avATT, avDEF)

	// calculate the average stats
	// get agent
	agentState := baseAgent.AgentState()

	HP := float64(agentState.Hp)
	ST := float64(agentState.Stamina)
	ATT := float64(agentState.BonusAttack())
	DEF := float64(agentState.BonusDefense())

	// normalise the agent stats
	meanHP, meanST, meanATT, meanDEF := normalize4El(HP, ST, ATT, DEF)

	// cal differences
	diffHP := groupAvHP - meanHP
	diffST := groupAvST - meanST
	diffATT := groupAvATT - meanATT
	diffDEF := groupAvDEF - meanDEF

	// end loot logic here
	if diffHP > diffST && diffHP > diffATT && diffHP > diffDEF {
		return true, false, false, false // HP highest diff
	} else if diffST > diffATT && diffST > diffDEF {
		return false, true, false, false // ST highest diff
	} else if diffATT > diffDEF {
		return false, false, true, false // ATT highest diff
	} else {
		return false, false, false, true // DEF highest diff
	}
}
