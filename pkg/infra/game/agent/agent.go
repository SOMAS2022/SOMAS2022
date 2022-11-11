package agent

import (
	"github.com/benbjohnson/immutable"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
)

type Strategy interface {
	HandleFight(state state.State, baseAgent BaseAgent, decisionC chan<- decision.FightAction, log *immutable.Map[string, decision.FightAction])
}

type Agent struct {
	BaseAgent BaseAgent
	Strategy  Strategy
}

type BaseAgent struct {
	Communication commons.Communication
	Id            string
}
