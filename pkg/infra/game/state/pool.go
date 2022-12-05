package state

import (
	"infra/game/commons"
)

type Item struct {
	id    commons.ItemID
	value uint
}

func (i Item) Id() commons.ItemID {
	return i.id
}

func (i Item) Value() uint {
	return i.value
}

func NewItem(id commons.ItemID, value uint) *Item {
	return &Item{id: id, value: value}
}

type LootPool struct {
	weapons        *commons.ImmutableList[Item]
	shields        *commons.ImmutableList[Item]
	hpPotions      *commons.ImmutableList[Item]
	staminaPotions *commons.ImmutableList[Item]
}

func (l LootPool) Weapons() *commons.ImmutableList[Item] {
	return l.weapons
}

func (l LootPool) Shields() *commons.ImmutableList[Item] {
	return l.shields
}

func (l LootPool) HpPotions() *commons.ImmutableList[Item] {
	return l.hpPotions
}

func (l LootPool) StaminaPotions() *commons.ImmutableList[Item] {
	return l.staminaPotions
}

func NewLootPool(weapons *commons.ImmutableList[Item], shields *commons.ImmutableList[Item], hpPotions *commons.ImmutableList[Item], staminaPotions *commons.ImmutableList[Item]) *LootPool {
	return &LootPool{weapons: weapons, shields: shields, hpPotions: hpPotions, staminaPotions: staminaPotions}
}
