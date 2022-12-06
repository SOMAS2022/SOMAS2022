package allocation

import (
	"infra/game/commons"
	"infra/game/state"
)

func AllocMessageHandler() {}

type ClashItem struct {
	ItemID          commons.ID
	Type            string
	Value           uint
	RequestedAgents []commons.ID
	winner          commons.ID
}

// may not be neccessary depending on infra
func FindClashLoot(s *state.View) []ClashItem {
	var ClashLoot []ClashItem

	return ClashLoot
}

type TotalResource struct{}

func OptimalAlloc(totalAlives uint, totalResource TotalResource) {}
