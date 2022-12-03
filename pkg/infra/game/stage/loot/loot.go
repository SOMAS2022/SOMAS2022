package loot

import (
	"fmt"
	"math/rand"

	"infra/game/agent"
	"infra/game/commons"
	"infra/game/state"
	"infra/logging"

	"github.com/google/uuid"
)

func UpdateItems(state state.State, agents map[commons.ID]agent.Agent) *state.State {
	updatedState := state
	for _, a := range agents {
		agentState := updatedState.AgentState[a.BaseAgent.ID()]
		weaponId := a.HandleUpdateWeapon(agentState)
		shieldId := a.HandleUpdateShield(agentState)
		agentState.ChangeWeaponInUse(weaponId)
		agentState.ChangeShieldInUse(shieldId)
		updatedState.AgentState[a.BaseAgent.ID()] = agentState
	}
	return &updatedState
}

// AllocateHPPotion HPi and STi are the index of HP potion slice and ST potion slice that is allocated to the agent. Pass one at a time.
func AllocateHPPotion(globalState state.State, agentID commons.ID, HPi int) (state.State, *commons.ImmutableList[uint]) {
	allocatedState := globalState
	hpPotionValue := allocatedState.PotionSlice.HPPotion[HPi]
	v := allocatedState
	a := allocatedState.AgentState[agentID]
	a.Hp = v.AgentState[agentID].Hp + hpPotionValue
	allocatedState.AgentState[agentID] = a
	allocatedState.PotionSlice.HPPotion[HPi] = 0
	HPPotionList := commons.NewImmutableList(allocatedState.PotionSlice.HPPotion)
	return allocatedState, HPPotionList
}

func AllocateSTPotion(globalState state.State, agentID commons.ID, STi int) (state.State, *commons.ImmutableList[uint]) {
	allocatedState := globalState
	stPotionValue := allocatedState.PotionSlice.STPotion[STi]
	v := allocatedState
	a := allocatedState.AgentState[agentID]
	a.Stamina = v.AgentState[agentID].Stamina + stPotionValue
	allocatedState.AgentState[agentID] = a
	allocatedState.PotionSlice.STPotion[STi] = 0
	STPotionList := commons.NewImmutableList(allocatedState.PotionSlice.STPotion)
	return allocatedState, STPotionList
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
	var PotionList PotionList

	// allocate potion
	allocatedState.PotionSlice.HPPotion = make([]uint, len(hpPotionLoot))
	allocatedState.PotionSlice.STPotion = make([]uint, len(stPotionLoot))

	idx := 0

	for agentID, agentState := range allocatedState.AgentState {
		allocatedState, PotionList.HPPotionList = AllocateHPPotion(allocatedState, agentID, rand.Intn(len(hpPotionLoot)-idx))
		allocatedState, PotionList.STPotionList = AllocateSTPotion(allocatedState, agentID, rand.Intn(len(stPotionLoot)-idx))
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
