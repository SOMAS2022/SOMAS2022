package agent

import (
	"github.com/benbjohnson/immutable"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"math/rand"
)

type RandomAgent struct {
	bravery int
}

func (r RandomAgent) Default() decision.FightAction {
	//TODO implement me
	panic("implement me")
}

func NewRandomAgent() *RandomAgent {
	return &RandomAgent{bravery: 0}
}

func (RandomAgent) HandleFight(_ *state.View, _ BaseAgent, decisionC chan<- decision.FightAction, _ *immutable.Map[commons.ID, decision.FightAction]) {
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

func (r RandomAgent) HandleFightMessage(m message.TaggedMessage, view *state.View, agent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) *decision.FightAction {
	fight := rand.Intn(3)
	switch fight {
	case 0:
		return decision.CowerPtr()
	case 1:
		return decision.AttackPtr()
	default:
		return decision.DefendPtr()
	}
}
