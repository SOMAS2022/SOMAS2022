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
func AllocateHPPotion(globalState state.State, agentID commons.ID, HPi int) state.State {
	allocatedState := globalState
	hpPotionValue := allocatedState.PotionSlice.HPpotion[HPi]
	v := allocatedState
	a := allocatedState.AgentState[agentID]
	a.Hp = v.AgentState[agentID].Hp + hpPotionValue
	allocatedState.AgentState[agentID] = a
	allocatedState.PotionSlice.HPpotion, _ = commons.DeleteElFromSlice(allocatedState.PotionSlice.HPpotion, HPi)
	return allocatedState
}

func AllocateSTPotion(globalState state.State, agentID commons.ID, STi int) state.State {
	allocatedState := globalState
	stPotionValue := allocatedState.PotionSlice.STpotion[STi]
	v := allocatedState
	a := allocatedState.AgentState[agentID]
	a.Stamina = v.AgentState[agentID].Hp + stPotionValue
	allocatedState.AgentState[agentID] = a
	allocatedState.PotionSlice.STpotion, _ = commons.DeleteElFromSlice(allocatedState.PotionSlice.STpotion, STi)
	return allocatedState
}

//Use simple append function to append to the potion slice when generating new loot potions.

func AllocateLoot(globalState state.State, weaponLoot []uint, shieldLoot []uint, HPpotionloot []uint, STpotionloot []uint) state.State {
	allocatedState := globalState

	//allocate potion
	allocatedState.PotionSlice.HPpotion = make([]uint, len(HPpotionloot))
	allocatedState.PotionSlice.STpotion = make([]uint, len(STpotionloot))

	idx := 0

	for agentID, agentState := range allocatedState.AgentState {
		allocatedState = AllocateHPPotion(allocatedState, agentID, rand.Intn(len(HPpotionloot)-idx))
		allocatedState = AllocateSTPotion(allocatedState, agentID, rand.Intn(len(STpotionloot)-idx))
		allocatedState.AgentState[agentID] = agentState
		idx++
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
