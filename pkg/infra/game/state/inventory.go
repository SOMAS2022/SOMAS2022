package state

import (
	"infra/game/commons"
)

type InventoryMap struct {
	Shields map[commons.ID]uint
	Weapons map[commons.ID]uint
}
