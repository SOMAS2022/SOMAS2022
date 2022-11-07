package agent

import (
	"infra/game/decision"
	"infra/game/state"
	"math/rand"
)

type RandomAgent struct {
}

func (r RandomAgent) HandleFight(state state.State, baseAgent BaseAgent) {
	fight := rand.Intn(1)%2 == 0
	var action decision.FightAction
	if fight {
		action = decision.Cower{}
	} else {
		attackVal := rand.Intn(int(state.AgentState[baseAgent.Id].TotalAttack()))
		defendVal := rand.Intn(int(state.AgentState[baseAgent.Id].TotalDefense()))
		action = decision.Fight{Attack: uint(attackVal), Defend: uint(defendVal)}
	}
	baseAgent.Communication.Sender <- decision.FightDecision{Action: action}
}
