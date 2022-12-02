package loot

import (
	"fmt"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/state"
	"infra/logging"
	"math/rand"

	"github.com/google/uuid"
)

func UpdateItems(state state.State, agents map[commons.ID]agent.Agent) state.State {
	updatedState := state
	for _, agent := range agents {
		agentState := updatedState.AgentState[agent.BaseAgent.ID()]
		weaponId := agent.HandleUpdateWeapon(agentState, state.ToView())
		shieldId := agent.HandleUpdateShield(agentState, state.ToView())
		agentState.ChangeWeaponInUse(weaponId)
		agentState.ChangeShieldInUse(shieldId)
		updatedState.AgentState[agent.BaseAgent.ID()] = agentState
	}
	return updatedState
}

// HPi and STi are the index of HP potion slice and ST potion slice that is allocate to the agent. Pass one at a time.
func AllocatePotion(globalState state.State, agentID uint, HPi int, STi int) state.State {
	allocatedState := globalState
	agent.Hp += allocatedState.PotionSlice.HPpotion[HPi]
	agent.Stamina += allocatedState.PotionSlice.STpotion[STi]
	allocatedState.PotionSlice.HPpotion, _ = commons.DeleteElFromSlice(allocatedState.PotionSlice.HPpotion, HPi)
	allocatedState.PotionSlice.STpotion, _ = commons.DeleteElFromSlice(allocatedState.PotionSlice.STpotion, STi)
	return allocatedState
}

//Use simple append function to append to the potion slice when generating new loot potions.

func AllocateLoot(globalState state.State, weaponLoot []uint, shieldLoot []uint) state.State {
	allocatedState := globalState

	//allocate potion
	allocatedState.PotionSlice.HPpotion = nil
	allocatedState.PotionSlice.STpotion = nil

	for agentID, agentState := range allocatedState.AgentState {
		agentState.H = AllocatePotion(allocatedState, agent state.AgentState, HPi int, STi int) state.State
	}

	

	for agentID, agentState := range allocatedState.AgentState {
		allocatedWeaponIdx := rand.Intn(len(weaponLoot))
		allocatedShieldIdx := rand.Intn(len(shieldLoot))

		// add W to global InventoryMap and this agent's inventory
		wid := uuid.New().String()
		weaponValue := weaponLoot[allocatedWeaponIdx]
		allocatedState.InventoryMap.Weapons[wid] = weaponValue
		allocatedWeapon := state.InventoryItem{ID: wid, Value: weaponValue}
		agentState.AddWeapon(allocatedWeapon)

		// add S to global InventoryMap and this agent's inventory
		sid := uuid.New().String()
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

	return allocatedState
}