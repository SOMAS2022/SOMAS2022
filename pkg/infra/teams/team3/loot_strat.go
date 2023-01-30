package team3

import (
	"fmt"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/message/proposal"
	"math/rand"
	"sort"

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

// this is really poor from infra - i'd just ignore tbh
func (a *AgentThree) HandleLootInformation(m message.TaggedInformMessage[message.LootInform], baseAgent agent.BaseAgent) {
	// submit a proposal to the leader
	fmt.Println("Made it here")
	switch m.Message().(type) {
	case message.LootInform:
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

// forcibly call at start of loot phase to begin proceedings
func (a *AgentThree) RequestLootProposal(baseAgent agent.BaseAgent) { // put your logic here now, instead
	sendProposal := rand.Intn(100)
	if sendProposal > a.personality {
		return
	}
	// general and send a loot proposal at the start of every turn
	baseAgent.SendLootProposalToLeader(a.generateLootProposal())
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

func (a *AgentThree) ChooseItem(baseAgent agent.BaseAgent,
	items map[string]struct{}, weaponSet map[string]uint, shieldSet map[string]uint, hpPotionSet map[string]uint, staminaPotionSet map[string]uint) string {
	// function to calculate the agents choice of loot

	// get group average stats
	avHP, avST, avATT, avDEF := GetGroupAv(baseAgent)
	// normalise the group stats
	groupAvHP, groupAvST, groupAvATT, groupAvDEF := normalize4El(avHP, avST, avATT, avDEF)
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

	// create an array of the above, order them
	diffs := []float64{diffHP, diffST, diffATT, diffDEF}
	sortedDiffs := make([]float64, len(diffs))
	copy(sortedDiffs, diffs)
	sort.Slice(sortedDiffs, func(i, j int) bool {
		return sortedDiffs[i] > sortedDiffs[j]
	})
	var item string
	// return the item that is needed most (out of the items available)
	for _, val := range sortedDiffs {
		if val == 0 {
			// if val is zero, everyone the same so take arbitrary loot
			for id := range items {
				item = id
				break
			}
			return item
		} else if val == diffHP {
			//search of item in corresponding set
			item = searchForItem(hpPotionSet, items)
			// if excists, then return the item
			if len(item) > 0 {
				return item
			}
		} else if val == diffST {
			item = searchForItem(staminaPotionSet, items)
			if len(item) > 0 {
				return item
			}
		} else if val == diffATT {
			item = searchForItem(weaponSet, items)
			if len(item) > 0 {
				return item
			}
		} else if val == diffDEF {
			item = searchForItem(shieldSet, items)
			if len(item) > 0 {
				return item
			}
		}
	}
	if item == "" {
		// if got nothing, take arbitrary loot
		for id := range items {
			item = id
			break
		}
	}
	return item
}

func searchForItem(set map[string]uint, items map[string]struct{}) string {
	for item := range items {
		if _, ok := set[item]; ok {
			return item
		}
	}
	return ""
}

func GetGroupAv(baseAgent agent.BaseAgent) (float64, float64, float64, float64) {
	avHP := AverageArray(GetHealthAllAgents(baseAgent))
	avST := AverageArray(GetStaminaAllAgents(baseAgent))
	avATT := AverageArray(GetAttackAllAgents(baseAgent))
	avDEF := AverageArray(GetDefenceAllAgents(baseAgent))

	return avHP, avST, avATT, avDEF
}
