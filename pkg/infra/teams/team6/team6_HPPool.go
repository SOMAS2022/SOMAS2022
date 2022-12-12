package team6

import (
	"infra/game/agent"
	"math/rand"
)

func (a *Team6Agent) DonateToHpPool(baseAgent agent.BaseAgent) uint {
	return uint(rand.Intn(int(baseAgent.AgentState().Hp)))
}
