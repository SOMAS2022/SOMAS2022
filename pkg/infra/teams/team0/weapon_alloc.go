package team0

import (
	"math/rand"

	"infra/game/commons"
	"infra/game/state"
)

// AllocateLoot
/**
* This default function allocates loot randomly.
 */
func AllocateLoot(globalState state.State, weaponLoot []uint, shieldLoot []uint) state.State {
	allocatedState := globalState

	for _, agentState := range allocatedState.AgentState {
		allocatedWeapon := rand.Intn(len(weaponLoot))
		allocatedShield := rand.Intn(len(shieldLoot))

		agentState.BonusAttack = weaponLoot[allocatedWeapon]
		agentState.BonusDefense = shieldLoot[allocatedShield]
		weaponLoot, _ = commons.DeleteElFromSlice(weaponLoot, allocatedWeapon)
		shieldLoot, _ = commons.DeleteElFromSlice(shieldLoot, allocatedShield)
	}

	return allocatedState
}
