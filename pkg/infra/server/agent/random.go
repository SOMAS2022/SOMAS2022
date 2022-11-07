package agent

import (
	"infra/server/decision"
	"infra/server/state"
)

type RandomAgent struct {
}

func (r RandomAgent) HandleFight(_ state.State, baseAgent BaseAgent) {
	//todo: make random action based on state instead of always cower
	baseAgent.Communication.Sender <- decision.FightDecision{Choice: decision.Cower{}}
}
