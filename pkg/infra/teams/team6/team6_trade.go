package team6

import (
	"infra/game/agent"
	"infra/game/message"
)

func (a *Team6Agent) HandleTradeNegotiation(_ agent.BaseAgent, _ message.TradeInfo) message.TradeMessage {
	return message.TradeRequest{}
}
