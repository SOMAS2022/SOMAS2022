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

func (s *AgentState) HasItem(itemType commons.ItemType, itemID commons.ItemID) bool {
	var inventory immutable.List[Item]
	if itemType == commons.Weapon {
		inventory = s.Weapons
	} else if itemType == commons.Shield {
		inventory = s.Shields
	}
	itr := inventory.Iterator()
	for !itr.Done() {
		_, item := itr.Next()
		if item.Id() == itemID {
			return true
		}
	}
	return false
}

func (s *AgentState) BonusAttack() uint {
	iterator := s.Weapons.Iterator()
	for !iterator.Done() {
		_, value := iterator.Next()
		if value.Id() == s.WeaponInUse {
			return value.Value()
		}
	}
	return 0
}

func (s *AgentState) BonusDefense() uint {
	iterator := s.Shields.Iterator()
	for !iterator.Done() {
		_, value := iterator.Next()
		if value.Id() == s.ShieldInUse {
			return value.Value()
		}
	}
	return 0
}

func (s *AgentState) TotalAttack() uint {
	return s.Attack + s.BonusAttack()
}

func (s *AgentState) TotalDefense() uint {
	return s.Defense + s.BonusDefense()
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
