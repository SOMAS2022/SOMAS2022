package team0

import (
	"infra/game/commons"
	"infra/game/stage/loot"
	"infra/game/state"
	"math/rand"

	"github.com/google/uuid"
)

// AllocateLoot
/**
* This default function allocates loot randomly.
 */
func AllocateLoot(globalState state.State, weaponLoot []uint, shieldLoot []uint, HPpotionloot []uint, STpotionloot []uint) state.State {
	allocatedState := globalState

	//allocate potion
	allocatedState.PotionSlice.HPpotion = nil
	allocatedState.PotionSlice.STpotion = nil

	allocatedState.PotionSlice.HPpotion = make([]uint, len(HPpotionloot))
	allocatedState.PotionSlice.HPpotion = make([]uint, len(STpotionloot))

	for agentID, agentState := range allocatedState.AgentState {
		allocatedState = loot.AllocateHPPotion(allocatedState, agentID, rand.Intn(len(HPpotionloot)))
		allocatedState = loot.AllocateSTPotion(allocatedState, agentID, rand.Intn(len(STpotionloot)))
		allocatedState.AgentState[agentID] = agentState
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

	return allocatedState
}
