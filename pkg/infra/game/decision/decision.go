package decision

import "infra/game/state"

type Decision interface {
	decisionSealed()
}

type FightDecision struct {
	Cower  bool
	Attack uint
	Defend uint
}

func (d *FightDecision) ValidateDecision(s state.AgentState) {
	if d.Cower {
		return
	}
	if !(d.Attack <= s.TotalAttack() && d.Defend <= s.TotalDefense() && d.Attack+d.Defend <= s.AbilityPoints) {
		d.Cower = true
	}
}

func (FightDecision) decisionSealed() {}

type LootDecision struct{}

func (LootDecision) decisionSealed() {}

type HPPoolDecision struct{}

func (HPPoolDecision) decisionSealed() {}
