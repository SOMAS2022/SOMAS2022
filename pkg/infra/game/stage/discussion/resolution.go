package discussion

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message/proposal"
	"infra/game/state"
	"infra/game/tally"
	"math/rand"

	"golang.org/x/exp/maps"

	"github.com/benbjohnson/immutable"
)

func ResolveFightDiscussion(gs state.State, agentMap map[commons.ID]agent.Agent, currentLeader agent.Agent, manifesto decision.Manifesto, tally *tally.Tally[decision.FightAction]) decision.FightResult {
	fightActions := make(map[commons.ID]decision.FightAction)
	// todo: cleanup the nil check that acts to check if the leader died in combat
	var prop commons.ImmutableList[proposal.Rule[decision.FightAction]]
	if manifesto.FightImposition() && currentLeader.Strategy != nil {
		prop = currentLeader.Strategy.FightResolution(*currentLeader.BaseAgent)
	} else {
		// get proposal with most votes
		prop = tally.GetMax().Rules()
	}

	predicate := proposal.ToSinglePredicate(prop)
	if predicate == nil {
		for id, a := range agentMap {
			fightActions[id] = a.FightActionNoProposal(*a.BaseAgent)
		}
	} else {
		for id, a := range agentMap {
			expectedFightAction := predicate(gs, a.AgentState())
			if gs.Defection {
				fightActions[id] = a.FightAction(*a.BaseAgent, expectedFightAction)
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

	return decision.FightResult{
		Choices:         fightActions,
		AttackingAgents: nil,
		ShieldingAgents: nil,
		CoweringAgents:  nil,
		AttackSum:       0,
		ShieldSum:       0,
	}
}

func ResolveLootDiscussion(
	gs state.State,
	agentMap map[commons.ID]agent.Agent,
	pool *state.LootPool,
	leader agent.Agent,
	manifesto decision.Manifesto,
	tally *tally.Tally[decision.LootAction],
) immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]] {
	if manifesto.LootImposition() && leader.Strategy != nil {
		// todo: eliminate double allocations here
		return leader.Strategy.LootAllocation(*leader.BaseAgent)
	} else {
		predicate := proposal.ToMultiPredicate(tally.GetMax().Rules())
		if predicate == nil {
			// either leader died or no proposal was made
			return handleNilLootAllocation(agentMap)
		}
		getsWeapon, getsShield, getsHealthPotion, getsStaminaPotion := demandList(gs, agentMap, predicate)
		m := make(map[commons.ID]map[commons.ItemID]struct{})
		buildAllocation(pool.Weapons(), getsWeapon, m)
		buildAllocation(pool.Shields(), getsShield, m)
		buildAllocation(pool.HpPotions(), getsHealthPotion, m)
		buildAllocation(pool.StaminaPotions(), getsStaminaPotion, m)

		mMapped := make(map[commons.ID]immutable.SortedMap[commons.ItemID, struct{}])

		anyDetractors := false
		wantedItems := make(map[commons.ItemID]map[commons.ID]struct{})

		for id, itemIDS := range m {
			mMapped[id] = commons.MapToSortedImmutable[commons.ItemID, struct{}](itemIDS)
			addWantedLootToItemAllocMap(mMapped[id], wantedItems, id)
			agentLoot := agentMap[id].Strategy.LootAction(*agentMap[id].BaseAgent, mMapped[id])
			if !commons.ImmutableSetEquality(mMapped[id], agentLoot) {
				anyDetractors = true
				defector := gs.AgentState[id].Defector
				defector.SetLoot(true)
			}
		}

		if !anyDetractors {
			return commons.MapToImmutable(mMapped)
		} else {
			return convertAllocationMapToImmutable(formAllocationFromConflicts(wantedItems))
		}
	}
}

func demandList(
	gs state.State,
	agentMap map[commons.ID]agent.Agent,
	predicate func(state.State, state.AgentState) map[decision.LootAction]struct{},
) ([]commons.ID, []commons.ID, []commons.ID, []commons.ID) {
	getsWeapon := make([]commons.ID, 0)
	getsShield := make([]commons.ID, 0)
	getsHealthPotion := make([]commons.ID, 0)
	getsStaminaPotion := make([]commons.ID, 0)
	for id := range agentMap {
		actions := predicate(gs, gs.AgentState[id])
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

func handleNilLootAllocation(agentMap map[commons.ID]agent.Agent) immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]] {
	wantedItems := make(map[commons.ItemID]map[commons.ID]struct{})
	for id, a := range agentMap {
		wantedLoot := a.Strategy.LootActionNoProposal(*agentMap[id].BaseAgent)
		addWantedLootToItemAllocMap(wantedLoot, wantedItems, id)
	}
	allocations := formAllocationFromConflicts(wantedItems)

	return convertAllocationMapToImmutable(allocations)
}

func convertAllocationMapToImmutable(allocations map[commons.ID]map[commons.ItemID]struct{}) immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]] {
	mMapped := make(map[commons.ID]immutable.SortedMap[commons.ItemID, struct{}])
	for id, itemIDS := range allocations {
		mMapped[id] = commons.MapToSortedImmutable(itemIDS)
	}
	return commons.MapToImmutable(mMapped)
}

func formAllocationFromConflicts(wantedItems map[commons.ItemID]map[commons.ID]struct{}) map[commons.ID]map[commons.ItemID]struct{} {
	allocations := make(map[commons.ID]map[commons.ItemID]struct{})
	for item, agentSet := range wantedItems {
		agents := maps.Keys(agentSet)
		if len(agents) > 0 {
			rand.Shuffle(len(agents), func(i, j int) { agents[i], agents[j] = agents[j], agents[i] })
			if m, ok := allocations[agents[0]]; ok {
				m[item] = struct{}{}
			} else {
				m := make(map[commons.ItemID]struct{})
				m[item] = struct{}{}
				allocations[agents[0]] = m
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

func buildAllocation(pool *commons.ImmutableList[state.Item], proposedLooters []commons.ID, allocation map[commons.ID]map[commons.ItemID]struct{}) {
	idx := 0
	iterator := pool.Iterator()
	for !iterator.Done() {
		if idx >= len(proposedLooters) {
			break
		}
		next, _ := iterator.Next()
		if m, ok := allocation[proposedLooters[idx]]; ok {
			m[next.Id()] = struct{}{}
		} else {
			m := make(map[commons.ItemID]struct{})
			m[next.Id()] = struct{}{}
			allocation[proposedLooters[idx]] = m
		}
	}
}
