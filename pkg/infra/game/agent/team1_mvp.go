package agent

import (
	"github.com/benbjohnson/immutable"
	"infra/game/decision"
	"infra/game/state"
	"math/rand"
)

type Team1MVP struct {
}

func (Team1MVP) HandleFight(_ state.State, _ BaseAgent, decisionC chan<- decision.FightAction, log *immutable.Map[uint, decision.FightAction]) {
	/*
		Either attack or defend
	*/
	fight := rand.Intn(2)
	switch fight {
	case 0:
		decisionC <- decision.Attack
	case 1:
		decisionC <- decision.Defend
	}
}
