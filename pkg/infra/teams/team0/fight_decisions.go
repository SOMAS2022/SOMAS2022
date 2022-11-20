package team0

import (
	"infra/game/agent"
	"infra/game/decision"
)

func AllDefend(agents map[commons.ID]agent.Agent) map[commons.ID]decision.FightAction {
	decisionMap := make(map[commons.ID]decision.FightAction)

	for i, _ := range agents {
		decisionMap[i] = decision.Defend
	}

	return decisionMap
}
