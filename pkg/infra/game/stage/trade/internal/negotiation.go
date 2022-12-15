package internal

import (
	"infra/game/commons"
	"infra/game/message"
)

type Info struct {
	negotiations map[commons.TradeID]message.TradeNegotiation
	Inventory
}

func (n *Info) Negotiations() map[commons.TradeID]message.TradeNegotiation {
	return n.negotiations
}

func NewInfo(negotiations map[commons.TradeID]message.TradeNegotiation, inventory Inventory) *Info {
	return &Info{negotiations: negotiations, Inventory: inventory}
}
