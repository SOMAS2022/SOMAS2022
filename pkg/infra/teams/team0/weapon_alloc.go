package team0

// AllocateLoot
/**
* This default function allocates loot randomly.
 */
//func AllocateLoot(globalState state.State, weaponLoot []uint, shieldLoot []uint, hpPotionLoot []uint, stPotionLoot []uint) *state.State {
//	allocatedState := globalState
//
//	for agentID, agentState := range allocatedState.AgentState {
//		allocatedState = loot.AllocateHPPotion(allocatedState, hpPotionLoot, agentID, rand.Intn(len(hpPotionLoot)))
//		allocatedState = loot.AllocateSTPotion(allocatedState, stPotionLoot, agentID, rand.Intn(len(stPotionLoot)))
//		allocatedState.AgentState[agentID] = agentState
//	}
//
//	for agentID, agentState := range allocatedState.AgentState {
//		allocatedWeaponIdx := rand.Intn(len(weaponLoot))
//		allocatedShieldIdx := rand.Intn(len(shieldLoot))
//
//		// add W to global InventoryMap and this agent's inventory
//		wid := uuid.NewString()
//		weaponValue := weaponLoot[allocatedWeaponIdx]
//		allocatedState.InventoryMap.Weapons[wid] = weaponValue
//		allocatedWeapon := state.InventoryItem{ID: wid, Value: weaponValue}
//		agentState.AddWeapon(allocatedWeapon)
//
//		// add S to global InventoryMap and this agent's inventory
//		sid := uuid.NewString()
//		shieldValue := shieldLoot[allocatedShieldIdx]
//		allocatedState.InventoryMap.Shields[sid] = shieldValue
//		allocatedShield := state.InventoryItem{ID: sid, Value: shieldValue}
//		agentState.AddShield(allocatedShield)
//
//		allocatedState.AgentState[agentID] = agentState
//
//		// remove W and S from unallocated loot
//
//		weaponLoot, _ = commons.DeleteElFromSlice(weaponLoot, allocatedWeaponIdx)
//		shieldLoot, _ = commons.DeleteElFromSlice(shieldLoot, allocatedShieldIdx)
//	}
//
//	return &allocatedState
//}
