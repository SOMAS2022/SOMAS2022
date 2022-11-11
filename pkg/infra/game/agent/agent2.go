//Base template for Agent Group 2
package agent

import (
	"infra/game/decision"
	"infra/game/state"
    "math/rand"
)

type Agent2 struct {

}

func (Agent2) HandleFight(state state.State, baseAgent BaseAgent, decisionC chan<- decision.FightAction) {
/*
Fighting is similar to Random Agent for MVP.
- Agent attacks by default (50% of the time)
- Agent defends 40% of the time
- Agent cowers 10% of time

*/
	fight := rand.Intn(10)
	switch {
	case fight == 0:
		decisionC <- decision.Cower
	case (fight <= 4) && (fight > 0) :
		decisionC <- decision.Defend
    default:
        decisionC <- decision.Attack
	}
}
