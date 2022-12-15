package decision

import (
	"infra/game/commons"

	"github.com/benbjohnson/immutable"
)

type FightResult struct {
	Choices         map[commons.ID]FightAction
	AttackingAgents []commons.ID
	ShieldingAgents []commons.ID
	CoweringAgents  []commons.ID
	AttackSum       uint
	ShieldSum       uint
}

type ImmutableFightResult struct {
	choices         immutable.Map[commons.ID, FightAction]
	attackingAgents immutable.List[commons.ID]
	shieldingAgents immutable.List[commons.ID]
	coweringAgents  immutable.List[commons.ID]
	attackSum       uint
	shieldSum       uint
	round           uint
}

func (ifr *ImmutableFightResult) Choices() immutable.Map[commons.ID, FightAction] {
	return ifr.choices
}

func (ifr *ImmutableFightResult) AttackingAgents() immutable.List[commons.ID] {
	return ifr.attackingAgents
}

func (ifr *ImmutableFightResult) ShieldingAgents() immutable.List[commons.ID] {
	return ifr.shieldingAgents
}

func (ifr *ImmutableFightResult) CoweringAgents() immutable.List[commons.ID] {
	return ifr.coweringAgents
}

func (ifr *ImmutableFightResult) AttackSum() uint {
	return ifr.attackSum
}

func (ifr *ImmutableFightResult) ShieldSum() uint {
	return ifr.shieldSum
}

func (ifr *ImmutableFightResult) Round() uint {
	return ifr.round
}

func NewImmutableFightResult(fightResult FightResult, round uint) *ImmutableFightResult {
	return &ImmutableFightResult{
		choices:         commons.MapToImmutable(fightResult.Choices),
		attackingAgents: commons.ListToImmutableList(fightResult.AttackingAgents),
		shieldingAgents: commons.ListToImmutableList(fightResult.ShieldingAgents),
		coweringAgents:  commons.ListToImmutableList(fightResult.CoweringAgents),
		attackSum:       fightResult.AttackSum,
		shieldSum:       fightResult.ShieldSum,
		round:           round,
	}
}
