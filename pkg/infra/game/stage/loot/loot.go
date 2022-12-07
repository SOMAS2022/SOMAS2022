package loot

import (
	"infra/game/decision"
	"infra/game/message"
	"infra/game/tally"
	"infra/logging"
	"sync"
	"time"

	"github.com/benbjohnson/immutable"

	"infra/game/agent"
	"infra/game/commons"
	"infra/game/state"
)

type agentStateUpdate struct {
	commons.ID
	state.AgentState
}

func UpdateItems(s state.State, agents map[commons.ID]agent.Agent) *state.State {
	updatedState := s
	var wg sync.WaitGroup
	updatedStates := make(chan agentStateUpdate)
	for id, a := range agents {
		wg.Add(1)
		id := id
		a := a
		agentState := s.AgentState[id]
		go func(id commons.ID, a agent.Agent, sender chan<- agentStateUpdate, wait *sync.WaitGroup) {
			weaponId := a.HandleUpdateWeapon(agentState)
			shieldId := a.HandleUpdateShield(agentState)
			agentState.ChangeWeaponInUse(weaponId)
			agentState.ChangeShieldInUse(shieldId)
			sender <- agentStateUpdate{
				ID:         id,
				AgentState: agentState,
			}
			wait.Done()
		}(id, a, updatedStates, &wg)
	}
	go func(group *sync.WaitGroup) {
		group.Wait()
		close(updatedStates)
	}(&wg)

	for update := range updatedStates {
		updatedState.AgentState[update.ID] = update.AgentState
	}

	return &updatedState
}

func AgentLootDecisions(
	state state.State,
	availableLoot state.LootPool,
	agents map[commons.ID]agent.Agent,
	channelsMap map[commons.ID]chan message.TaggedMessage,
) *tally.Tally[decision.LootAction] {
	proposalVotes := make(chan commons.ProposalID)
	proposalSubmission := make(chan message.Proposal[decision.LootAction])
	tallyClosure := make(chan struct{})

	propTally := tally.NewTally(proposalVotes, proposalSubmission, tallyClosure)
	go propTally.HandleMessages()
	closures := make(map[commons.ID]chan<- struct{})
	starts := make(map[commons.ID]chan<- message.StartLoot)
	for id, a := range agents {
		a := a
		closure := make(chan struct{})
		closures[id] = closure

		start := make(chan message.StartLoot)
		starts[id] = start

		agentState := state.AgentState[a.BaseAgent.ID()]
		if a.BaseAgent.ID() == state.CurrentLeader {
			go (&a).HandleLoot(agentState, proposalVotes, proposalSubmission, closure, start)
		} else {
			go (&a).HandleLoot(agentState, proposalVotes, nil, closure, start)
		}
	}

	startLootMessage := *message.NewStartLoot(availableLoot)
	for _, start := range starts {
		start <- startLootMessage
	}

	time.Sleep(100 * time.Millisecond)
	for id, c := range channelsMap {
		closures[id] <- struct{}{}
		go func(recv <-chan message.TaggedMessage) {
			for m := range recv {
				switch m.Message().(type) {
				case message.Request:
					// todo: respond with nil thing here as we're closing! Or do we need to?
					// maybe because we're closing there's no point...
				default:
				}
			}
		}(c)
	}

	for _, c := range channelsMap {
		close(c)
	}

	tallyClosure <- struct{}{}
	close(tallyClosure)
	return propTally
}

func HandleLootAllocation(globalState state.State, allocation *immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]], pool *state.LootPool) *state.State {
	weaponSet := itemListToSet(pool.Weapons())
	shieldSet := itemListToSet(pool.Shields())
	hpPotionSet := itemListToSet(pool.HpPotions())
	staminaPotionSet := itemListToSet(pool.StaminaPotions())

	allocationIterator := allocation.Iterator()

	for !allocationIterator.Done() {
		agentID, items, _ := allocationIterator.Next()
		itemIterator := items.Iterator()
		for !itemIterator.Done() {
			item, _, _ := itemIterator.Next()
			agentState := globalState.AgentState[agentID]

			if val, ok := weaponSet[item]; ok {
				globalState.InventoryMap.Weapons[item] = val
				agentState.AddWeapon(*state.NewItem(item, val))
			} else if val, ok := shieldSet[item]; ok {
				globalState.InventoryMap.Shields[item] = val
				agentState.AddWeapon(*state.NewItem(item, val))
			} else if val, ok := hpPotionSet[item]; ok {
				agentState.Hp += val
			} else if val, ok := staminaPotionSet[item]; ok {
				agentState.Stamina += val
			} else {
				logging.Log(logging.Warn, nil, "unknown item attempted to be allocated")
			}
			globalState.AgentState[agentID] = agentState
		}
	}
	return &globalState
}

func itemListToSet(
	list *commons.ImmutableList[state.Item],
) map[commons.ItemID]uint {
	iterator := list.Iterator()
	res := make(map[commons.ItemID]uint)
	for !iterator.Done() {
		next, _ := iterator.Next()
		res[next.Id()] = next.Value()
	}
	return res
}
