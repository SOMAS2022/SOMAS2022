package team0

import (
	"infra/game/agent"
	"infra/game/decision"
)

func AllDefend(agents map[string]agent.Agent) map[string]decision.FightAction {
	decisionMap := make(map[string]decision.FightAction)

	for i, _ := range agents {
		decisionMap[i] = decision.Defend
	}

	return decisionMap
}
