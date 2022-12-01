package loot

import (
	"fmt"
	"math/rand"

	"infra/game/commons"
	"infra/game/state"
	"infra/logging"

	"github.com/google/uuid"
)

func AllocateLoot(globalState state.State, weaponLoot []uint, shieldLoot []uint) state.State {
	allocatedState := globalState

	for agentID, agentState := range allocatedState.AgentState {
		allocatedWeapon := rand.Intn(len(weaponLoot))
		allocatedShield := rand.Intn(len(shieldLoot))

		// add W to global InentoryMap and this agent's inventory
		wid := uuid.New().String()
		allocatedState.InventoryMap.Weapons[wid] = weaponLoot[allocatedWeapon]
		agentState.AddWeapon(wid)

		// add S to global InentoryMap and this agent's inventory
		sid := uuid.New().String()
		allocatedState.InventoryMap.Shields[sid] = shieldLoot[allocatedShield]
		agentState.AddShield(sid)

		allocatedState.AgentState[agentID] = agentState

		// remove W and S from unallocated loot

		weaponLoot, _ = commons.DeleteElFromSlice(weaponLoot, allocatedWeapon)
		shieldLoot, _ = commons.DeleteElFromSlice(shieldLoot, allocatedShield)
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
