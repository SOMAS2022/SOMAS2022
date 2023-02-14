package loot

import (
	"infra/game/decision"
	"infra/game/message"
	"infra/game/tally"

	// "math"
	"sync"
	"time"

	// "github.com/benbjohnson/immutable"

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
			go (&a).HandleLoot(agentState, proposalVotes, proposalSubmission, closure, start) // <----- threaded for voting
		} else {
			go (&a).HandleLoot(agentState, proposalVotes, nil, closure, start)
		}
	}

	startLootMessage := *message.NewStartLoot(availableLoot)
	for _, start := range starts {
		start <- startLootMessage
	}

	time.Sleep(25 * time.Millisecond)
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

func HandleLootAllocation(globalState state.State, allocation map[commons.ID]map[commons.ItemID]struct{}, pool *state.LootPool, agentMap map[commons.ID]agent.Agent) *state.State {
	weaponSet := itemListToSet(pool.Weapons())
	shieldSet := itemListToSet(pool.Shields())
	hpPotionSet := itemListToSet(pool.HpPotions())
	staminaPotionSet := itemListToSet(pool.StaminaPotions())

	// each agent can only take 1 item
	// calc diff of user between their normalized average and health/stamina/attack/defense, get highest diff
	// and use it as a boolean param for item selection

	// averageHP, averageST, averageATT, averageDEF := getAverageStats(globalState)

	for agentID, items := range allocation {
		agentState := globalState.AgentState[agentID]
		a := agentMap[agentID]

		// if items is of length 1, then take allocation
		if len(items) == 1 {
			// assign the only piece of loot they are eligable for
			for item := range items {
				assignChosenItem(item, weaponSet, shieldSet, hpPotionSet, staminaPotionSet, &agentState)
			}
		} else {
			// choose the most needed item from the list of allocated items
			item := a.ChooseItem(*a.BaseAgent, items, weaponSet, shieldSet, hpPotionSet, staminaPotionSet)

			// asign the most needed item to the agent
			assignChosenItem(item, weaponSet, shieldSet, hpPotionSet, staminaPotionSet, &agentState)
		}

		globalState.AgentState[agentID] = agentState

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

// assign loot function
func assignChosenItem(item string, weaponSet map[string]uint, shieldSet map[string]uint, hpPotionSet map[string]uint, staminaPotionSet map[string]uint, agentState *state.AgentState) {
	// check the item id and assign to the agent
	// then delete the item from the weaponSet

	if val, ok := weaponSet[item]; ok {
		// globalState.InventoryMap.Weapons[item] = val
		agentState.AddWeapon(*state.NewItem(item, val))
		delete(weaponSet, item)
		// delete(globalState.InventoryMap.Weapons, item)
	} else if val, ok := shieldSet[item]; ok {
		// globalState.InventoryMap.Shields[item] = val
		agentState.AddShield(*state.NewItem(item, val))
		delete(shieldSet, item)
		// delete(globalState.InventoryMap.Shields, item)
	} else if val, ok := hpPotionSet[item]; ok {
		agentState.Hp += val
		delete(hpPotionSet, item)
	} else if val, ok := staminaPotionSet[item]; ok {
		agentState.Stamina += val
		delete(staminaPotionSet, item)
	}
}

//TODO
// func getEligibleItems(agent agent.Agent) {
// 	preference := agent.GetLootPreferenceOrder()
// 	for _, pref := range preference {
// 		// TODO: return available items for preference
// 		// BRANCH: if item length > 0 then return ELSE: keep looping
// 	}
// 	// FALLBACK: return empty array
// }

// func getAverageStats(globalState state.State) (float64, float64, float64, float64) {
// 	var averageHP float64 = 0
// 	var averageST float64 = 0
// 	var averageATT float64 = 0
// 	var averageDEF float64 = 0

// 	agentLen := float64(len(globalState.AgentState))
// 	for _, state := range globalState.AgentState {
// 		averageHP += float64(state.Hp)
// 		averageST += float64(state.Stamina)
// 		averageATT += float64(state.BonusAttack())
// 		averageDEF += float64(state.BonusDefense())
// 	}

// 	averageHP /= agentLen
// 	averageST /= agentLen
// 	averageATT /= agentLen
// 	averageDEF /= agentLen
// 	// fmt.Println(averageHP, averageST, averageATT, averageDEF)
// 	meanAverageHP, meanAverageST, meanAverageATT, meanAverageDEF := normalize4El(averageHP, averageST, averageATT, averageDEF)
// 	// fmt.Println(meanAverageHP, meanAverageST, meanAverageATT, meanAverageDEF)
// 	return meanAverageHP, meanAverageST, meanAverageATT, meanAverageDEF
// }

// func chooseItem(agent state.AgentState, averageHP float64, averageST float64, averageATT float64, averageDEF float64) (bool, bool, bool, bool) {

// 	HP := float64(agent.Hp)
// 	ST := float64(agent.Stamina)
// 	ATT := float64(agent.BonusAttack())
// 	DEF := float64(agent.BonusDefense())
// 	// fmt.Println(HP, ST, ATT, DEF)
// 	meanHP, meanST, meanATT, meanDEF := normalize4El(HP, ST, ATT, DEF)
// 	// fmt.Println(meanHP, meanST, meanATT, meanDEF)
// 	diffHP := averageHP - meanHP
// 	diffST := averageST - meanST
// 	diffATT := averageATT - meanATT
// 	diffDEF := averageDEF - meanDEF
// 	// fmt.Println(diffHP, diffST, diffATT, diffDEF)
// 	// get largest diff = var most in need
// 	if diffHP > diffST && diffHP > diffATT && diffHP > diffDEF {
// 		return true, false, false, false // HP highest diff
// 	} else if diffST > diffATT && diffST > diffDEF {
// 		return false, true, false, false // ST highest diff
// 	} else if diffATT > diffDEF {
// 		return false, false, true, false // ATT highest diff
// 	} else {
// 		return false, false, false, true // DEF highest diff
// 	}

// }

// // works bc normalization changes the data distribution, so small sheild/weapon difference values are significant enough now
// func normalize4El(x, y, z, w float64) (float64, float64, float64, float64) {
// 	maxVal := minMax4(true, [...]float64{x, y, z, w})
// 	minVal := minMax4(false, [...]float64{x, y, z, w})
// 	return (x - minVal) / (maxVal - minVal), (y - minVal) / (maxVal - minVal), (z - minVal) / (maxVal - minVal), (w - minVal) / (maxVal - minVal)
// }

// func minMax4(isMax bool, nums [4]float64) float64 {
// 	ans := nums[0]
// 	for _, num := range nums[1:] {
// 		if isMax {
// 			ans = math.Max(num, ans)
// 		} else {
// 			ans = math.Min(num, ans)
// 		}
// 	}
// 	return ans
// }

// didn't work as mean scaling just recenters the distribution -> sheild/weapon values were too small compared to the rest
// func meanScale4El(el1 float64, el2 float64, el3 float64, el4 float64) (float64, float64, float64, float64) {
// 	var mean float64 = (el1 + el2 + el3 + el4) / 4.0
// 	// fmt.Println(el1, el2, el3, el4)
// 	el1 /= mean
// 	el2 /= mean
// 	el3 /= mean
// 	el4 /= mean

// 	return el1, el2, el3, el4
// }
