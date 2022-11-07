package agent

import (
	"infra/server/commons"
	"infra/server/state"
)

type Strategy interface {
	HandleFight(state state.State, baseAgent BaseAgent)
}

type Agent struct {
	BaseAgent BaseAgent
	Strategy  Strategy
}

type BaseAgent struct {
	Communication commons.Communication
	Id            uint
}
