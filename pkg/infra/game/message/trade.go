package message

import (
	"infra/game/commons"

	"github.com/benbjohnson/immutable"
)

type TradeItemType uint

const (
	Shield TradeItemType = iota
	Weapon
)

type TradeItem struct {
	itemtype TradeItemType
	id       commons.ItemID
}

type TradeDemand struct {
	itemType TradeItemType
	minValue uint
}

type TradeCondition struct {
	offer  TradeItem
	demand TradeDemand
}

type TradeNegotiation struct {
	initialised bool
	// terminate negotiation when this number decrement to 0
	roundNum   uint
	conditions immutable.Map[commons.ID, TradeCondition]
}

func NewTradeItem(itemType TradeItemType, id commons.ItemID) TradeItem {
	return TradeItem{itemtype: itemType, id: id}
}

func NewTradeDemand(itemType TradeItemType, minValue uint) TradeDemand {
	return TradeDemand{itemType: itemType, minValue: minValue}
}

func NewTradeNegotiation(agentID commons.ID, offerType TradeItemType, itemID commons.ItemID, demandType TradeItemType, minValue uint) TradeNegotiation {
	conditions := immutable.NewMap[commons.ID, TradeCondition](nil).Set(agentID, TradeCondition{
		offer:  NewTradeItem(offerType, itemID),
		demand: NewTradeDemand(demandType, minValue),
	})
	return TradeNegotiation{
		initialised: true,
		roundNum:    3,
		conditions:  *conditions,
	}
}

func (negotiation TradeNegotiation) sealedMessage() {}

func (negotiation TradeNegotiation) GetRoundNum() uint {
	return negotiation.roundNum
}

func (negotiation TradeNegotiation) GetOffer(agentID commons.ID) (TradeItem, bool) {
	condition, ok := negotiation.conditions.Get(agentID)
	if ok {
		return condition.offer, true
	}
	return TradeItem{}, false
}

func (negotiation TradeNegotiation) GetCounterParty(agentID commons.ID) (commons.ID, bool) {
	itr := negotiation.conditions.Iterator()
	for !itr.Done() {
		id, _, _ := itr.Next()
		if id != agentID {
			return id, true
		}
	}
	return "", false
}

func (negotiation TradeNegotiation) GetDemand(agentID commons.ID) (demand TradeDemand, ok bool) {
	condition, ok := negotiation.conditions.Get(agentID)
	if ok {
		return condition.demand, true
	}
	return TradeDemand{}, false
}

func (negotiation *TradeNegotiation) UpdateRoundNum() {
	negotiation.roundNum = negotiation.roundNum - 1
}

func (negotiation *TradeNegotiation) UpdateOffer(agentID commons.ID, offerType TradeItemType, itemID commons.ItemID) {
	condition, ok := negotiation.conditions.Get(agentID)
	if ok {
		negotiation.conditions.Set(agentID, TradeCondition{
			offer:  NewTradeItem(offerType, itemID),
			demand: condition.demand,
		})
	}
}

func (negotiation *TradeNegotiation) UpdateDemand(agentID commons.ID, demandType TradeItemType, minValue uint) {
	condition, ok := negotiation.conditions.Get(agentID)
	if ok {
		negotiation.conditions.Set(agentID, TradeCondition{
			offer:  condition.offer,
			demand: NewTradeDemand(demandType, minValue),
		})
	}
}