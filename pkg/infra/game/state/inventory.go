package state

import (
	"infra/game/commons"
	"sort"

	"github.com/benbjohnson/immutable"
)

type InventoryMap struct {
	Shields map[commons.ID]uint
	Weapons map[commons.ID]uint
}

type InventoryItem struct {
	ID    commons.ItemID
	Value uint
}

// Add an InventoryItem to an immutable list of InventoryItem.
// return a sorted immutable.List with 0th InventoryItem has greatest value.
func Add2Inventory(items immutable.List[InventoryItem], item InventoryItem) (newItems immutable.List[InventoryItem]) {
	// convert immutable.List to slice
	itemList := []InventoryItem{item}
	itr := items.Iterator()
	for !itr.Done() {
		_, inventoryItem := itr.Next()
		itemList = append(itemList, inventoryItem)
	}

	// sort slice
	sort.SliceStable(itemList, func(i, j int) bool {
		return itemList[i].Value > itemList[j].Value
	})

	// convert sorted slice to immutable
	b := immutable.NewListBuilder[InventoryItem]()
	for _, w := range itemList {
		b.Append(w)
	}

	return *b.List()
}

// Remove an InventoryItem from an immutable list of InventoryItem.
// return a sorted immutable.List with 0th InventoryItem has greatest value.
func RemoveFromInventory(items immutable.List[InventoryItem], itemID commons.ItemID) (newItems immutable.List[InventoryItem]) {
	b := immutable.NewListBuilder[InventoryItem]()

	// filter inventoryItem with ID == itemID
	itr := items.Iterator()
	for !itr.Done() {
		_, item := itr.Next()
		if item.ID != itemID {
			b.Append(item)
		}
	}

	return *b.List()
}
