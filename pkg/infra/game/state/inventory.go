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

// return a sorted immutable.List with 0th InventoryItem has greatest value
func updateInventory(items immutable.List[InventoryItem], item InventoryItem) (newItems immutable.List[InventoryItem]) {
	// convert immutable.List to slice
	itemList := []InventoryItem{item}
	itr := items.Iterator()
	for !itr.Done() {
		_, value := itr.Next()
		itemList = append(itemList, value)
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
