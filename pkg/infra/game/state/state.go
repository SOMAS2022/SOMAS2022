package state

import (
	"infra/game/commons"
	"infra/game/decision"

	"github.com/benbjohnson/immutable"
)

type Defector struct {
	fight bool
	loot  bool
}

func (d *Defector) SetFight(fight bool) {
	d.fight = fight
}

func (d *Defector) SetLoot(loot bool) {
	d.loot = loot
}

func NewDefector() *Defector {
	return &Defector{}
}

func (d *Defector) IsDefector() bool {
	return d.fight || d.loot
}

type AgentState struct {
	Hp          uint
	Stamina     uint
	Attack      uint
	Defense     uint
	WeaponInUse commons.ItemID
	ShieldInUse commons.ItemID
	Weapons     immutable.List[Item]
	Shields     immutable.List[Item]
	Defector    Defector
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

func (s *AgentState) AddWeapon(weapon Item) {
	s.Weapons = addToInventory(s.Weapons, weapon)
}

func (s *AgentState) AddShield(shield Item) {
	s.Shields = addToInventory(s.Shields, shield)
}

func (s *AgentState) ChangeWeaponInUse(weaponIdx decision.ItemIdx) {
	if int(weaponIdx) < s.Weapons.Len() {
		s.WeaponInUse = s.Weapons.Get(int(weaponIdx)).Id()
	}
}

func (s *AgentState) ChangeShieldInUse(shieldIdx decision.ItemIdx) {
	if int(shieldIdx) < s.Shields.Len() {
		s.ShieldInUse = s.Shields.Get(int(shieldIdx)).Id()
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
	Defection       bool
}
