package team3

import (
	"infra/game/agent"
	"infra/game/message"
)

func (a *AgentThree) HandleTradeNegotiation(_ agent.BaseAgent, _ message.TradeInfo) message.TradeMessage {
	return message.TradeRequest{}
}
