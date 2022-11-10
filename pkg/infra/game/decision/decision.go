package decision

import "infra/game/state"

type LootDecision struct{}

func (LootDecision) decisionSealed() {}

type HPPoolDecision struct{}

func (HPPoolDecision) decisionSealed() {}

type FightAction interface {
	HandleAction(state.AgentState) FightAction
}

type Cower struct{}

func (c Cower) HandleAction(state.AgentState) FightAction {
	return c
}

type Fight struct {
	Attack uint
	Defend uint
}

func (f Fight) HandleAction(s state.AgentState) FightAction {
	if f.Attack <= s.TotalAttack() && f.Defend <= s.TotalDefense() {
		return f
	} else {
		return Cower{}
	}
}
