package team5

import (
	"infra/game/commons"
	"infra/game/message"
	"infra/game/state"
)

type ClashItem struct {
	ItemID          commons.ID
	Type            string
	Value           uint
	RequestedAgents []commons.ID
	winner          commons.ID
}

// handles incomming message, whether
func AllocMessageHandler(m *message.TaggedMessage, s *state.View) *message.TaggedMessage {
	clashSet := message2ClashLoot(m)
	//function to be decided depending on the messaging method;
	//-----------------------------//

	//-----------------------------//
	//AllocMessageHandle handles all allocation message and return allocation decisions.
	//-----------------------------//
	winnerSet := FindWinnerSet(s, clashSet)
	//-----------------------------//
	AllocationMessage := clashLoot2Message(winnerSet)
	return AllocationMessage
}

// may not be necessary depending on infra
func message2ClashLoot(m *message.TaggedMessage) []ClashItem {
	var ClashLoot []ClashItem
	//from message to ClashLoot
	ClashLoot = FindClashLoot()
	return ClashLoot
}
func clashLoot2Message(cl []ClashItem) *message.TaggedMessage {
	var m *message.TaggedMessage
	//convert struct to message
	return m
}

func FindClashLoot() []ClashItem {
	var ClashLoot []ClashItem
	//build clashloot struct from all information/messages
	return ClashLoot
}

func FindWinner(s *state.View, cl ClashItem) ClashItem {
	var Winner ClashItem
	//function fot resolving clash one by one
	return Winner
}

func FindWinnerSet(s *state.View, cl []ClashItem) []ClashItem {
	var Winner []ClashItem
	//function fot resolving all clashes at the same time
	return Winner
}

type TotalResource struct{}

func OptimalAlloc(totalAlives uint, totalResource TotalResource) {}
