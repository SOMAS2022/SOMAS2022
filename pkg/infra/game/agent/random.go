package agent

import (
	"infra/game/decision"
	"infra/game/state"
	"math/rand"
)

type RandomAgent struct {
}

func (RandomAgent) HandleFight(state state.State, baseAgent BaseAgent, decisionC chan<- decision.FightAction) {
	fight := rand.Intn(3)
	switch fight {
	case 0:
		decisionC <- decision.Cower
	case 1:
		decisionC <- decision.Attack
	case 2:
		decisionC <- decision.Defend
	}
}
