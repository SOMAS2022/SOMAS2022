package state

import (
	"infra/game/commons"
	"infra/game/decision"
)

type AgentState struct {
	Hp          uint
	Stamina     uint
	Attack      uint
	Defense     uint
	WeaponInUse commons.ItemID
	ShieldInUse commons.ItemID
	Weapons     []commons.ItemID
	Shields     []commons.ItemID
}

func (a AgentState) BonusAttack(state State) uint {
	if val, ok := state.InventoryMap.Weapons[a.WeaponInUse]; ok {
		return val
	}
	return 0
}

func (a AgentState) BonusDefense(state State) uint {
	if val, ok := state.InventoryMap.Shields[a.ShieldInUse]; ok {
		return val
	}
	return 0
}

func (a AgentState) TotalAttack(state State) uint {
	return a.Attack + a.BonusAttack(state)
}

func (a AgentState) TotalDefense(state State) uint {
	return a.Defense + a.BonusDefense(state)
}

func (s *AgentState) AddWeapon(weaponID commons.ItemID) {
	s.Weapons = append(s.Weapons, weaponID)
}

func (s *AgentState) AddShield(shieldID commons.ItemID) {
	s.Shields = append(s.Shields, shieldID)
}

func (s *AgentState) ChangeWeaponInUse(weaponIdx decision.ItemIdx) {
	if int(weaponIdx) < len(s.Weapons) {
		s.WeaponInUse = s.Weapons[weaponIdx]
	}
}

func (s *AgentState) ChangeShieldInUse(shieldIdx decision.ItemIdx) {
	if int(shieldIdx) < len(s.Shields) {
		s.ShieldInUse = s.Shields[shieldIdx]
	}
}

type State struct {
	CurrentLevel    uint
	HpPool          uint
	MonsterHealth   uint
	MonsterAttack   uint
	AgentState      map[commons.ID]AgentState
	InventoryMap    InventoryMap
	CurrentLeader   commons.ID
	LeaderManifesto decision.Manifesto
}
