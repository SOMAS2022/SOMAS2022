package internal

import (
	"infra/game/commons"
	"infra/game/state"
)

type Inventory struct {
	weapons map[commons.ID][]state.Item
	shields map[commons.ID][]state.Item
}

func (i *Inventory) Weapons() map[commons.ID][]state.Item {
	return i.weapons
}

func (i *Inventory) Shields() map[commons.ID][]state.Item {
	return i.shields
}

func NewInventory(weapons map[commons.ID][]state.Item, shields map[commons.ID][]state.Item) *Inventory {
	return &Inventory{weapons: weapons, shields: shields}
}
