package inventory

import (
	"bytes"
	"infra/game/commons"
)

type Item struct {
	Itemtype string
	ItemID   commons.ID
	Itemstat uint
}

// ShieldOn is the shield that is held on hand and visible to others from ShieldAll, the complite shield inventory
type ShieldInventory struct {
	ShieldOn  uint
	ShieldAll []Item
}

type WeaponInventory struct {
	WeaponOn  uint
	WeaponAll []Item
}

type InventoryState struct {
	Shield ShieldInventory
	Weapon WeaponInventory
}

func changeItemOn(item Item, iventorystate InventoryState) (err string) {

	switch item.Itemtype {
	case "Shield":
		shieldindex := uint(bytes.IndexAny(inventorystate.Shield.ShieldAll, item))
		if shieldindex >= 0 {
			inventorystate.Shield.ShieldOn = shieldindex
		}

	case "Weapon":
		if i == 0 {
			inventorystate.Weapon.ShieldOn == 0
		}
		if (i < len(inventorystate.Weapon.WeaponAll)) && (i > 0) {
			inventorystate.Weapon.WeaponOn == i
		}
	default:
		inventorystate = inventorystate
	}
}

func ItemTransfer(itemA2B Item, agentA InventoryState, agentB InventoryState) {

}
