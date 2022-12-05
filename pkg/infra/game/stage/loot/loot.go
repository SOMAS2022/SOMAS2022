package loot

import (
	"fmt"
	"math/rand"
	"sync"

	"infra/game/agent"
	"infra/game/commons"
	"infra/game/state"
	"infra/logging"

	"github.com/google/uuid"
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

// AllocateHPPotion HPi and STi are the index of HP potion slice and ST potion slice that is allocated to the agent. Pass one at a time.
func AllocateHPPotion(globalState state.State, loot []uint, agentID commons.ID, HPi int) state.State {
	allocatedState := globalState
	hpPotionValue := loot[HPi]
	v := allocatedState
	a := allocatedState.AgentState[agentID]
	a.Hp = v.AgentState[agentID].Hp + hpPotionValue
	allocatedState.AgentState[agentID] = a
	return allocatedState
}

func AllocateSTPotion(globalState state.State, loot []uint, agentID commons.ID, STi int) state.State {
	allocatedState := globalState
	stPotionValue := loot[STi]
	v := allocatedState
	a := allocatedState.AgentState[agentID]
	a.Stamina = v.AgentState[agentID].Stamina + stPotionValue
	allocatedState.AgentState[agentID] = a
	return allocatedState
}

type PotionList struct {
	HPPotionList *commons.ImmutableList[uint]
	STPotionList *commons.ImmutableList[uint]
}

// AllocateLoot immutable list for communication with agent only.
// a slice is generated from state, action is done on the slice
// immutable list is generated upon temporary slice.
func AllocateLoot(globalState state.State, weaponLoot []uint, shieldLoot []uint, hpPotionLoot []uint, stPotionLoot []uint) *state.State {
	allocatedState := globalState

	idx := 0
	rand.Shuffle(len(hpPotionLoot), func(i, j int) { hpPotionLoot[i], hpPotionLoot[j] = hpPotionLoot[j], hpPotionLoot[i] })
	rand.Shuffle(len(stPotionLoot), func(i, j int) { stPotionLoot[i], stPotionLoot[j] = stPotionLoot[j], stPotionLoot[i] })

	for agentID, agentState := range allocatedState.AgentState {
		allocatedState = AllocateHPPotion(allocatedState, hpPotionLoot, agentID, idx)
		allocatedState = AllocateSTPotion(allocatedState, stPotionLoot, agentID, idx)
		allocatedState.AgentState[agentID] = agentState
		idx++
	}

	for agentID, agentState := range allocatedState.AgentState {
		allocatedWeaponIdx := rand.Intn(len(weaponLoot))
		allocatedShieldIdx := rand.Intn(len(shieldLoot))

		// add W to global InventoryMap and this agent's inventory
		wid := uuid.NewString()
		weaponValue := weaponLoot[allocatedWeaponIdx]
		allocatedState.InventoryMap.Weapons[wid] = weaponValue
		allocatedWeapon := state.InventoryItem{ID: wid, Value: weaponValue}
		agentState.AddWeapon(allocatedWeapon)

		// add S to global InventoryMap and this agent's inventory
		sid := uuid.NewString()
		shieldValue := shieldLoot[allocatedShieldIdx]
		allocatedState.InventoryMap.Shields[sid] = shieldValue
		allocatedShield := state.InventoryItem{ID: sid, Value: shieldValue}
		agentState.AddShield(allocatedShield)

		allocatedState.AgentState[agentID] = agentState

		// remove W and S from unallocated loot

		weaponLoot, _ = commons.DeleteElFromSlice(weaponLoot, allocatedWeaponIdx)
		shieldLoot, _ = commons.DeleteElFromSlice(shieldLoot, allocatedShieldIdx)
	}

	logging.Log(
		logging.Info,
		logging.LogField{
			"weapons": len(globalState.InventoryMap.Weapons),
			"shields": len(globalState.InventoryMap.Shields),
		},
		fmt.Sprintf("%6d Weapons, %6d Shields in InventoryMap", len(globalState.InventoryMap.Weapons), len(globalState.InventoryMap.Shields)),
	)

	return &allocatedState
}
