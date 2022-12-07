package message

import (
	"infra/game/commons"
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
	id       commons.ItemID
	itemtype commons.ItemType
	value    uint
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

func (t TradeReject) sealedTradeMessage()  {}
func (t TradeReject) sealedTradeResponse() {}
