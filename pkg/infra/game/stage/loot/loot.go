package loot

import (
	"fmt"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/state"
	"infra/logging"
	"math/rand"

	"github.com/benbjohnson/immutable"
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
func AllocateHPPotion(globalState state.State, agentID commons.ID, HPi int) (state.State, immutable.Map[int, uint]) {
	allocatedState := globalState
	hpPotionValue := allocatedState.PotionSlice.HPpotion[HPi]
	v := allocatedState
	a := allocatedState.AgentState[agentID]
	a.Hp = v.AgentState[agentID].Hp + hpPotionValue
	allocatedState.AgentState[agentID] = a
	allocatedState.PotionSlice.HPpotion[HPi] = 0
	hpPotionMap := commons.Slice2Map(allocatedState.PotionSlice.HPpotion)
	HPPotionMap := commons.MapToImmutable(hpPotionMap)
	return allocatedState, HPPotionMap
}

func AllocateSTPotion(globalState state.State, agentID commons.ID, STi int) (state.State, immutable.Map[int, uint]) {
	allocatedState := globalState
	stPotionValue := allocatedState.PotionSlice.STpotion[STi]
	v := allocatedState
	a := allocatedState.AgentState[agentID]
	a.Stamina = v.AgentState[agentID].Hp + stPotionValue
	allocatedState.AgentState[agentID] = a
	allocatedState.PotionSlice.STpotion[STi] = 0
	stPotionMap := commons.Slice2Map(allocatedState.PotionSlice.STpotion)
	STPotionMap := commons.MapToImmutable(stPotionMap)
	return allocatedState, STPotionMap
}

//Use simple append function to append to the potion slice when generating new loot potions.

func AllocateLoot(globalState state.State, weaponLoot []uint, shieldLoot []uint, HPpotionloot []uint, STpotionloot []uint) state.State {
	allocatedState := globalState

	//allocate potion
	allocatedState.PotionSlice.HPpotion = make([]uint, len(HPpotionloot))
	allocatedState.PotionSlice.STpotion = make([]uint, len(STpotionloot))

	idx := 0

	for agentID, agentState := range allocatedState.AgentState {
		allocatedState, _ = AllocateHPPotion(allocatedState, agentID, rand.Intn(len(HPpotionloot)-idx))
		allocatedState, _ = AllocateSTPotion(allocatedState, agentID, rand.Intn(len(STpotionloot)-idx))
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

// HPi and STi are the index of HP potion slice and ST potion slice that is allocate to the agent. Pass one at a time.
func AllocatePotion(globalState state.State, Potion state.PotionSlice, agent state.AgentState, HPi int, STi int) state.State {
	agent.Hp += globalState.PotionSlice.HPpotion[HPi]
	agent.Stamina += globalState.PotionSlice.STpotion[STi]
	globalState.PotionSlice.HPpotion, _ = commons.DeleteElFromSlice(globalState.PotionSlice.HPpotion, HPi)
	globalState.PotionSlice.STpotion, _ = commons.DeleteElFromSlice(globalState.PotionSlice.STpotion, STi)
	return globalState
}

//Use simple append function to append to the postion slice when generating new loot potions.
