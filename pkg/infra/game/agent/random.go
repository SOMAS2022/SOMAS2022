package agent

import (
	"infra/game/decision"
	"infra/game/state"
	"math/rand"
)

type RandomAgent struct {
}

func (RandomAgent) HandleFight(state state.State, baseAgent BaseAgent, decisionC chan<- decision.FightAction) {
	fight := rand.Intn(2) == 0
	if fight {
		attackVal := rand.Intn(int(state.AgentState[baseAgent.Id].TotalAttack()))
		defendVal := rand.Intn(int(state.AgentState[baseAgent.Id].TotalDefense()))
		decisionC <- decision.Fight{Attack: uint(attackVal), Defend: uint(defendVal)}
	} else {
		decisionC <- decision.Cower{}
	}
}
