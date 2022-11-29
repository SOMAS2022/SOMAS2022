package agent

import (
	"infra/game/decision"
	"infra/game/state"
	"math/rand"
	// "fmt"
)

type FivAgent struct {
}

func (r FivAgent) HandleFight(state state.State, baseAgent BaseAgent) {
	confidenceI := 0.6
	sampleTime := int(float64(len(state.AgentState)) * confidenceI)
	// fmt.Println(sampleTime)
	approxAttack := 0
	for i := 0; i < sampleTime; i++ {
		approxAttack += rand.Intn(int(state.AgentState[baseAgent.Id].TotalAttack()))
	}
	fight := (uint(approxAttack) > state.MonsterHealth)

	var action decision.FightAction
	if fight {
		attackVal := rand.Intn(int(state.AgentState[baseAgent.Id].TotalAttack()))
		defendVal := rand.Intn(int(state.AgentState[baseAgent.Id].TotalDefense()))
		action = decision.Fight{Attack: uint(attackVal), Defend: uint(defendVal)}
	} else {
		action = decision.Cower{}
	}
	baseAgent.Communication.Sender <- decision.FightDecision{Action: action}
}
