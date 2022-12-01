package team0

import (
	"math/rand"

	"infra/game/commons"
	"infra/game/state"

	"github.com/google/uuid"
)

// AllocateLoot
/**
* This default function allocates loot randomly.
 */
func AllocateLoot(globalState state.State, weaponLoot []uint, shieldLoot []uint) state.State {
	allocatedState := globalState

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
