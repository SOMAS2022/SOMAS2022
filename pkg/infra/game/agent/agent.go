package agent

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
)

type Strategy interface {
	HandleFight(state state.State, baseAgent BaseAgent, decision chan<- decision.FightAction)
}

type Agent struct {
	BaseAgent BaseAgent
	Strategy  Strategy
}

type BaseAgent struct {
	Communication commons.Communication
	Id            uint
}
