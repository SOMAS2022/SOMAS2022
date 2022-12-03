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

func (s *AgentState) BonusAttack(state State) uint {
	if val, ok := state.InventoryMap.Weapons[s.WeaponInUse]; ok {
		return val
	}
	return 0
}

func (s *AgentState) BonusDefense(state State) uint {
	if val, ok := state.InventoryMap.Shields[s.ShieldInUse]; ok {
		return val
	}
	return 0
}

func (s *AgentState) TotalAttack(state State) uint {
	return s.Attack + s.BonusAttack(state)
}

func (s *AgentState) TotalDefense(state State) uint {
	return s.Defense + s.BonusDefense(state)
}

func (s *AgentState) AddWeapon(weapon InventoryItem) {
	s.Weapons = addToInventory(s.Weapons, weapon)
}

func (s *AgentState) AddShield(shield InventoryItem) {
	s.Shields = addToInventory(s.Shields, shield)
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

type PotionSlice struct {
	HPpotion []uint //index is potion's id and element cotent is the potion's value
	STpotion []uint
}

type State struct {
	CurrentLevel    uint
	HpPool          uint
	MonsterHealth   uint
	MonsterAttack   uint
	AgentState      map[commons.ID]AgentState
	InventoryMap    InventoryMap
	PotionSlice     PotionSlice
	CurrentLeader   commons.ID
	LeaderManifesto decision.Manifesto
}
