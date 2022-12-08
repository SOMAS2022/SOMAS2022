package message

import (
	"infra/game/commons"
	"infra/game/state"

	"github.com/benbjohnson/immutable"
)

type TradeMessage interface {
	sealedTradeMessage()
}

type TradeAbstain struct {
	TradeMessage
}

type TradeRequest struct {
	CounterPartyID commons.ID
	Offer          TradeOffer
	Demand         TradeDemand
}

type TradeResponse interface {
	TradeMessage
	sealedTradeResponse()
}

type TradeBargain struct {
	TradeID commons.TradeID
	Offer   TradeOffer
	Demand  TradeDemand
}

type TradeAccept struct {
	TradeID commons.TradeID
}

type TradeReject struct {
	TradeID commons.TradeID
}

type TradeOffer struct {
	ItemType commons.ItemType
	Item     state.Item
	IsValid  bool
}

type TradeDemand struct {
	ItemType commons.ItemType
	MinValue uint
}

type TradeCondition struct {
	Offer  TradeOffer
	Demand TradeDemand
}

type TradeInfo struct {
	Negotiations map[commons.TradeID]TradeNegotiation
	Weapons      immutable.List[state.Item]
	Shields      immutable.List[state.Item]
}

func (t TradeAbstain) sealedTradeMessage() {}

func (t TradeRequest) sealedTradeMessage() {}

func (t TradeBargain) sealedTradeMessage()  {}
func (t TradeBargain) sealedTradeResponse() {}

func (t TradeAccept) sealedTradeMessage()  {}
func (t TradeAccept) sealedTradeResponse() {}

func (t TradeReject) sealedTradeMessage()  {}
func (t TradeReject) sealedTradeResponse() {}

func NewTradeOffer(itemType commons.ItemType, idx uint, weapon immutable.List[state.Item], shield immutable.List[state.Item]) (offer TradeOffer, ok bool) {
	var inventory immutable.List[state.Item]
	if itemType == commons.Weapon {
		inventory = weapon
	} else if itemType == commons.Shield {
		inventory = shield
	}
	if idx > uint(inventory.Len()) {
		return TradeOffer{}, false
	}
	item := inventory.Get(int(idx))
	return TradeOffer{ItemType: itemType, Item: item, IsValid: true}, true
}

func NewTradeDemand(itemType commons.ItemType, minValue uint) TradeDemand {
	return TradeDemand{ItemType: itemType, MinValue: minValue}
}
