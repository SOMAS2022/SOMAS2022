package state

import (
	"infra/game/commons"
)

type InventoryMap struct {
	Shields map[commons.ID]uint
	Weapons map[commons.ID]uint
}

type InventoryItem struct {
	ID    commons.ItemID
	Value uint
}
