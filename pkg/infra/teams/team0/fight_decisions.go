package team0

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
)

func AllDefend(agents map[commons.ID]agent.Agent) map[commons.ID]decision.FightAction {
	decisionMap := make(map[commons.ID]decision.FightAction)

	for i := range agents {
		decisionMap[i] = decision.Defend
	}

	return decisionMap
}
