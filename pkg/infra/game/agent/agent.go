package agent

import (
	"github.com/benbjohnson/immutable"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
)

type Strategy interface {
	HandleFight(view *state.View, baseAgent BaseAgent, decisionC chan<- decision.FightAction, log *immutable.Map[ID, decision.FightAction])
}

type Agent struct {
	BaseAgent BaseAgent
	Strategy  Strategy
}

type BaseAgent struct {
	Communication commons.Communication
	Id            ID
}

type ID = string
