package agent

import (
	"infra/game/decision"
	"infra/game/state"
)

type Strategy interface {
	HandleFight(state state.State, agent Agent, prevDecisions map[uint]decision.FightDecision)
}

type Agent struct {
	Strategy             Strategy
	StateChannel         chan state.State
	FightDecisionChannel chan decision.FightDecision
	State                state.AgentState
}

// type BaseAgent struct {
// 	Communication commons.Communication
// 	Id            uint
// }
