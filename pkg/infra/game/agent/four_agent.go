package agent

import (
	"infra/game/decision"
	"infra/game/state"
	"math/rand"
)

type FourAgent struct {
}

func randFloats(min, max float64, n int) []float64 {
    res := make([]float64, n)
    for i := range res {
        res[i] = min + rand.Float64() * (max - min)
    }
    return res
}

func (r FourAgent) HandleFight(state state.State, baseAgent BaseAgent) {
	
	//We fight almost 8 times out of 10! :)
	threshold := 0.8
	fight := randFloats(0,1)
	

	var action decision.FightAction
	if (fight > threshold) {
		attackVal := rand.Intn(int(state.AgentState[baseAgent.Id].TotalAttack()))
		defendVal := rand.Intn(int(state.AgentState[baseAgent.Id].TotalDefense()))
		action = decision.Fight{Attack: uint(attackVal), Defend: uint(defendVal)}
	} else {
		action = decision.Cower{}
	}
	baseAgent.Communication.Sender <- decision.FightDecision{Action: action}
}
