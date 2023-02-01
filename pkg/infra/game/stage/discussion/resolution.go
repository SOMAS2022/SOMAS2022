package discussion

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/message/proposal"
	"infra/game/state"
	"infra/game/tally"
	"math/rand"

	"golang.org/x/exp/maps"

	"github.com/benbjohnson/immutable"
)

func ResolveFightDiscussion(gs state.State, agentMap map[commons.ID]agent.Agent, currentLeader agent.Agent, manifesto decision.Manifesto, tally *tally.Tally[decision.FightAction]) decision.FightResult {
	fightActions := make(map[commons.ID]decision.FightAction)
	prop := tally.GetMax()
	rules := prop.Rules()

	predicate := proposal.ToSinglePredicate(&rules)
	if predicate == nil {
		for id, a := range agentMap {
			fightActions[id] = a.FightActionNoProposal(*a.BaseAgent)
		}
	} else {
		for id, a := range agentMap {
			expectedFightAction := predicate(a.AgentState())
			if gs.Defection {
				fightActions[id] = a.FightAction(*a.BaseAgent, expectedFightAction, prop)
				if expectedFightAction != fightActions[id] {
					agentState := gs.AgentState[id]
					agentState.Defector.SetFight(true)
					gs.AgentState[id] = agentState
				}
			} else {
				fightActions[id] = expectedFightAction
			}
		}
	}

	if manifesto.FightDecisionPower() && currentLeader.Strategy != nil {
		resolution := currentLeader.Strategy.FightResolution(*currentLeader.BaseAgent, rules, commons.MapToImmutable(fightActions))
		handleDefectionFight(gs, agentMap, resolution, fightActions, prop)
	}

	return decision.FightResult{
		Choices:         fightActions,
		AttackingAgents: nil,
		ShieldingAgents: nil,
		CoweringAgents:  nil,
		AttackSum:       0,
		ShieldSum:       0,
	}
}

func handleDefectionFight(gs state.State, agentMap map[commons.ID]agent.Agent, resolution immutable.Map[commons.ID, decision.FightAction], fightActions map[commons.ID]decision.FightAction, prop message.Proposal[decision.FightAction]) {
	for id, a := range agentMap {
		value, ok := resolution.Get(id)
		if ok {
			actualAction := a.FightAction(*a.BaseAgent, value, prop)
			if actualAction != value {
				agentState := gs.AgentState[id]
				agentState.Defector.SetFight(true)
				gs.AgentState[id] = agentState
			}
			fightActions[id] = actualAction
		} else {
			fightActions[id] = a.FightActionNoProposal(*a.BaseAgent)
		}
	}
}

func ResolveLootDiscussion(
	gs state.State,
	agentMap map[commons.ID]agent.Agent,
	pool *state.LootPool,
	leader agent.Agent,
	manifesto decision.Manifesto,
	tally *tally.Tally[decision.LootAction],
) map[commons.ID]map[commons.ItemID]struct{} {
	prop := tally.GetMax()
	allocation := getAllocation(gs, agentMap, pool, prop)
	// if manifesto.LootDecisionPower() && leader.Strategy != nil {
	// 	leaderAllocation := leader.Strategy.LootAllocation(*leader.BaseAgent, prop, allocation)
	// 	iterator := leaderAllocation.Iterator()
	// 	actualAllocation := make(map[commons.ID]immutable.SortedMap[commons.ItemID, struct{}])

	// 	handleDefectionLoot(gs, agentMap, iterator, actualAllocation, prop)

	// 	defectedAllocation := commons.MapToImmutable(actualAllocation)
	// 	iterator = defectedAllocation.Iterator()
	// 	wantedItems := make(map[commons.ItemID]map[commons.ID]struct{}) // item mapped to (map of people who want it)
	// 	for !iterator.Done() {
	// 		agentID, items, _ := iterator.Next()
	// 		itemIterator := items.Iterator()
	// 		for !itemIterator.Done() {
	// 			item, _, _ := itemIterator.Next()
	// 			if m, ok := wantedItems[item]; ok {
	// 				m[agentID] = struct{}{}
	// 			} else {
	// 				m := make(map[commons.ID]struct{})
	// 				m[agentID] = struct{}{}
	// 				wantedItems[item] = m
	// 			}
	// 		}
	// 	}

	// 	return (formAllocationFromConflicts(wantedItems))
	// } else {
	return allocation
	// }
}

func getAllocation(gs state.State, agentMap map[commons.ID]agent.Agent, pool *state.LootPool, prop message.Proposal[decision.LootAction]) map[commons.ID]map[commons.ItemID]struct{} {
	predicate := proposal.ToMultiPredicate(prop.Rules())
	if predicate == nil {
		// either leader died or no proposal was made
		return handleNilLootAllocation(agentMap)
	}
	getsWeapon, getsShield, getsHealthPotion, getsStaminaPotion := demandList(gs, agentMap, predicate)
	agentToElibilityList := make(map[commons.ID]map[commons.ItemID]struct{})
	buildEligibility(pool.Weapons(), getsWeapon, agentToElibilityList)
	buildEligibility(pool.Shields(), getsShield, agentToElibilityList)
	buildEligibility(pool.HpPotions(), getsHealthPotion, agentToElibilityList)
	buildEligibility(pool.StaminaPotions(), getsStaminaPotion, agentToElibilityList)

	return agentToElibilityList

}

// 	wantedItems := make(map[commons.ItemID]map[commons.ID]struct{})

// 	for id, itemIDS := range m {
// 		alloc := commons.MapToSortedImmutable(itemIDS)
// 		if gs.Defection {
// 			agentLoot := agentMap[id].Strategy.LootAction(*agentMap[id].BaseAgent, alloc, prop)
// 			addWantedLootToItemAllocMap(agentLoot, wantedItems, id)
// 			if !commons.ImmutableSetEquality(alloc, agentLoot) {
// 				defector := gs.AgentState[id].Defector
// 				defector.SetLoot(true)
// 			}
// 		} else {
// 			addWantedLootToItemAllocMap(alloc, wantedItems, id)
// 		}
// 	}
// 	return formAllocationFromConflicts(wantedItems)
// }

// func handleDefectionLoot(
// 	gs state.State,
// 	agentMap map[commons.ID]agent.Agent,
// 	iterator *immutable.MapIterator[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]],
// 	actualAllocation map[commons.ID]immutable.SortedMap[commons.ItemID, struct{}],
// 	prop message.Proposal[decision.LootAction],
// ) {
// 	for !iterator.Done() {
// 		agentID, allocation, _ := iterator.Next()
// 		a := agentMap[agentID]
// 		newAllocation := a.LootAction(*a.BaseAgent, allocation, prop)
// 		if !commons.ImmutableSetEquality(newAllocation, allocation) {
// 			defector := gs.AgentState[agentID].Defector
// 			defector.SetLoot(true)
// 		}
// 		actualAllocation[agentID] = newAllocation
// 	}
// }

func demandList(
	gs state.State,
	agentMap map[commons.ID]agent.Agent,
	predicate func(state.AgentState) map[decision.LootAction]struct{},
) ([]commons.ID, []commons.ID, []commons.ID, []commons.ID) {
	getsWeapon := make([]commons.ID, 0)
	getsShield := make([]commons.ID, 0)
	getsHealthPotion := make([]commons.ID, 0)
	getsStaminaPotion := make([]commons.ID, 0)
	for id := range agentMap {
		actions := predicate(gs.AgentState[id])
		if _, ok := actions[decision.Weapon]; ok {
			getsWeapon = append(getsWeapon, id)
		}
		if _, ok := actions[decision.Shield]; ok {
			getsShield = append(getsShield, id)
		}
		if _, ok := actions[decision.HealthPotion]; ok {
			getsHealthPotion = append(getsHealthPotion, id)
		}
		if _, ok := actions[decision.StaminaPotion]; ok {
			getsStaminaPotion = append(getsStaminaPotion, id)
		}
	}
	return getsWeapon, getsShield, getsHealthPotion, getsStaminaPotion
}

func handleNilLootAllocation(agentMap map[commons.ID]agent.Agent) map[commons.ID]map[commons.ItemID]struct{} {
	wantedItems := make(map[commons.ItemID]map[commons.ID]struct{})
	for id, a := range agentMap {
		wantedLoot := a.Strategy.LootActionNoProposal(*agentMap[id].BaseAgent)
		addWantedLootToItemAllocMap(wantedLoot, wantedItems, id)
	}
	allocations := formAllocationFromConflicts(wantedItems)

	return (allocations)
}

// func convertAllocationMapToImmutable(allocations map[commons.ID]map[commons.ItemID]struct{}) immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]] {
// 	mMapped := make(map[commons.ID]immutable.SortedMap[commons.ItemID, struct{}])
// 	for id, itemIDS := range allocations {
// 		mMapped[id] = commons.MapToSortedImmutable(itemIDS)
// 	}
// 	return commons.MapToImmutable(mMapped)
// }

func formAllocationFromConflicts(wantedItems map[commons.ItemID]map[commons.ID]struct{}) map[commons.ID]map[commons.ItemID]struct{} {
	allocations := make(map[commons.ID]map[commons.ItemID]struct{})
	for item, agentSet := range wantedItems {
		agents := maps.Keys(agentSet)
		if len(agents) > 0 {
			// for a number of iterations
			for i := 0; i < len(agents); i++ {
				randID := rand.Intn(len(agents))
				// choose a random agent, if they are not in the allocation, then assign them the item
				if _, ok := allocations[agents[randID]]; !ok {
					m := make(map[commons.ItemID]struct{})
					m[item] = struct{}{}
					allocations[agents[randID]] = m
					break
				} else {
					// otherwise, choose again
					continue
				}
			}
		}
	}
	return allocations
}

func addWantedLootToItemAllocMap(wantedLoot immutable.SortedMap[commons.ItemID, struct{}], wantedItems map[commons.ItemID]map[commons.ID]struct{}, id commons.ID) {
	iterator := wantedLoot.Iterator()
	for !iterator.Done() {
		value, _, _ := iterator.Next()
		if s, ok := wantedItems[value]; ok {
			s[id] = struct{}{}
		} else {
			s := make(map[commons.ID]struct{})
			s[id] = struct{}{}
			wantedItems[value] = s
		}
	}
}

// this is now useless - just place the whole pool inside each looter's key???
func buildEligibility(pool *commons.ImmutableList[state.Item], proposedLooters []commons.ID, allocation map[commons.ID]map[commons.ItemID]struct{}) {
	if len(proposedLooters) == 0 {
		return
	}
	iterator := pool.Iterator()
	// iterate over all items in pool
	for !iterator.Done() {
		// get next item in pool
		next, _ := iterator.Next()
		for _, looter := range proposedLooters {
			// if agent is in allocation already, add item to their key
			if m, ok := allocation[looter]; ok {
				m[next.Id()] = struct{}{}
			} else {
				// else, add agent to allocation and assign item to their key
				m := make(map[commons.ItemID]struct{})
				m[next.Id()] = struct{}{}
				allocation[looter] = m
			}
		}
	}
}
