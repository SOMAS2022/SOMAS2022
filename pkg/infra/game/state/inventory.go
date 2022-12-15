package state

import (
	"sort"

	"infra/game/commons"

	"github.com/benbjohnson/immutable"
)

type InventoryMap struct {
	Shields map[commons.ItemID]uint
	Weapons map[commons.ItemID]uint
}

// Add an InventoryItem to an immutable list of InventoryItem.
// return a sorted immutable.List with 0th InventoryItem has the greatest value.
func addToInventory(items immutable.List[Item], item Item) immutable.List[Item] {
	// convert immutable.List to slice
	itemList := []Item{item}
	itr := items.Iterator()
	for !itr.Done() {
		_, inventoryItem := itr.Next()
		itemList = append(itemList, inventoryItem)
	}

	// sort slice
	sort.SliceStable(itemList, func(i, j int) bool {
		return itemList[i].value > itemList[j].value
	})

	// convert sorted slice to immutable
	b := immutable.NewListBuilder[Item]()
	for _, w := range itemList {
		b.Append(w)
	}

	return *b.List()
}

// Remove an InventoryItem from an immutable list of InventoryItem.
// return a sorted immutable.List with 0th InventoryItem has the greatest value.
func removeFromInventory(items immutable.List[Item], itemID commons.ItemID) immutable.List[Item] {
	b := immutable.NewListBuilder[Item]()

	// filter inventoryItem with ID == itemID
	itr := items.Iterator()
	for !itr.Done() {
		_, item := itr.Next()
		if item.id != itemID {
			b.Append(item)
		}
	}

	return *b.List()
}
