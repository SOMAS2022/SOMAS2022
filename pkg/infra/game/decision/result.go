package decision

import (
	"infra/game/commons"
)

type FightResult struct {
	Choices         map[commons.ID]FightAction
	AttackingAgents []commons.ID
	ShieldingAgents []commons.ID
	CoweringAgents  []commons.ID
	AttackSum       uint
	ShieldSum       uint
}
