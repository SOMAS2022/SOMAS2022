package agent

import (
	"infra/server/commons"
	"infra/server/state"
)

type Logic interface {
	HandleFight(state state.State, baseAgent BaseAgent)
}

type BaseAgent struct {
	Communication commons.Communication
	Id            uint
	Logic         Logic
}
