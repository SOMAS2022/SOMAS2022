package state

import (
	"infra/game/commons"
	"infra/game/decision"

	"github.com/benbjohnson/immutable"
)

type AgentState struct {
	Hp          uint
	Stamina     uint
	Attack      uint
	Defense     uint
	WeaponInUse commons.ItemID
	ShieldInUse commons.ItemID
	Weapons     immutable.List[InventoryItem]
	Shields     immutable.List[InventoryItem]
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

func (s *AgentState) AddWeapon(weapon InventoryItem) {
	s.Weapons = Add2Inventory(s.Weapons, weapon)
}

func (s *AgentState) AddShield(shield InventoryItem) {
	s.Shields = Add2Inventory(s.Shields, shield)
}

func (s *AgentState) ChangeWeaponInUse(weaponIdx decision.ItemIdx) {
	if int(weaponIdx) < s.Weapons.Len() {
		s.WeaponInUse = s.Weapons.Get(int(weaponIdx)).ID
	}
}

func (s *AgentState) ChangeShieldInUse(shieldIdx decision.ItemIdx) {
	if int(shieldIdx) < s.Shields.Len() {
		s.ShieldInUse = s.Shields.Get(int(shieldIdx)).ID
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
