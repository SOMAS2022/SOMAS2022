package message

import (
	"infra/game/commons"
	"infra/game/state"
)

type TradeMessage interface {
	sealedTradeMessage()
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
	IsEmpty  bool
}

type TradeDemand struct {
	ItemType commons.ItemType
	MinValue uint
}

type TradeCondition struct {
	Offer  TradeOffer
	Demand TradeDemand
}

func (t TradeRequest) sealedTradeMessage() {}

func (t TradeBargain) sealedTradeMessage()  {}
func (t TradeBargain) sealedTradeResponse() {}

func (t TradeAccept) sealedTradeMessage()  {}
func (t TradeAccept) sealedTradeResponse() {}

func (t TradeReject) sealedTradeMessage()  {}
func (t TradeReject) sealedTradeResponse() {}
