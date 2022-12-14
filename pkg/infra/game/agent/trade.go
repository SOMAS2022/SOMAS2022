package agent

import (
	"infra/game/message"
)

type Trade interface {
	// HandleTradeNegotiation given a map of trade negotiations, respond to one of them or start a new trade negotiation
	HandleTradeNegotiation(baseAgent BaseAgent, Info message.TradeInfo) message.TradeMessage
}
