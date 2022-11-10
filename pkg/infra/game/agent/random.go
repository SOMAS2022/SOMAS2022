package agent

import (
	"infra/game/decision"
	"infra/game/state"
	"math/rand"
)

type RandomAgent struct {
}

func (RandomAgent) HandleFight(state state.State, agent Agent, prevDecisions map[uint]decision.FightDecision) {
	fight := rand.Intn(2) == 0
	var action decision.FightDecision
	if fight {
		attackVal := rand.Intn(int(agent.State.TotalAttack()))
		defendVal := rand.Intn(int(agent.State.TotalDefense()))
		action = decision.FightDecision{Cower: false, Attack: uint(attackVal), Defend: uint(defendVal)}

	} else {
		action = decision.FightDecision{Cower: true}
	}
	// baseAgent.Communication.Sender <- decision.FightDecision{Action: action}
	agent.FightDecisionChannel <- action
}
