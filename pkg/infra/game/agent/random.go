package agent

import (
	"github.com/benbjohnson/immutable"
	"infra/game/decision"
	"infra/game/state"
	"math/rand"
)

type RandomAgent struct {
}

func (RandomAgent) HandleFight(_ state.State, _ BaseAgent, decisionC chan<- decision.FightAction, log *immutable.Map[uint, decision.FightAction]) {
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
