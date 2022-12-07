package message

import (
	"infra/game/commons"
	"infra/game/state"
)

type TradeMessage interface {
	sealedTradeMessage()
}

type TradeRequest struct {
	counterPartyID commons.ID
	offer          TradeOffer
	demand         TradeDemand
}

type TradeResponse interface {
	TradeMessage
	sealedTradeResponse()
}

type TradeBargain struct {
	tradeID commons.TradeID
	offer   TradeOffer
	demand  TradeDemand
}

type TradeAccept struct {
	tradeID commons.TradeID
}

type TradeReject struct {
	tradeID commons.TradeID
}

type TradeOffer struct {
	itemtype commons.ItemType
	item     state.Item
	isEmpty  bool
}

type TradeDemand struct {
	itemType commons.ItemType
	minValue uint
}

type TradeCondition struct {
	offer  TradeOffer
	demand TradeDemand
}

func (t TradeRequest) sealedTradeMessage() {}

func (t TradeBargain) sealedTradeMessage()  {}
func (t TradeBargain) sealedTradeResponse() {}

func (t TradeAccept) sealedTradeMessage()  {}
func (t TradeAccept) sealedTradeResponse() {}
func (t TradeAccept) GetTradeID() commons.TradeID {
	return t.tradeID
}

func (t TradeReject) sealedTradeMessage()  {}
func (t TradeReject) sealedTradeResponse() {}
func (t TradeReject) GetTradeID() commons.TradeID {
	return t.tradeID
}
