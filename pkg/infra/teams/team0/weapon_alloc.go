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

	for _, agentState := range allocatedState.AgentState {
		allocatedWeapon := rand.Intn(len(weaponLoot))
		allocatedShield := rand.Intn(len(shieldLoot))

		wid := uuid.New().String()
		globalState.InventoryMap.Weapons[wid] = weaponLoot[allocatedWeapon]
		agentState.Weapons = append(agentState.Weapons, wid)

		sid := uuid.New().String()
		globalState.InventoryMap.Shields[sid] = shieldLoot[allocatedShield]
		agentState.Shields = append(agentState.Shields, sid)

		weaponLoot, _ = commons.DeleteElFromSlice(weaponLoot, allocatedWeapon)
		shieldLoot, _ = commons.DeleteElFromSlice(shieldLoot, allocatedShield)
	}

	return allocatedState
}
