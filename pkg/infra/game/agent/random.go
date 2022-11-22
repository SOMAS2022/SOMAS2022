package agent

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"infra/logging"
	"math/rand"

	"github.com/benbjohnson/immutable"
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

func (r RandomAgent) HandleFightMessage(m message.TaggedMessage, view *state.View, agent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) decision.FightAction {
	fight := rand.Intn(3)
	switch fight {
	case 0:
		agent.log(logging.Trace, logging.LogField{"bravery": r.bravery}, "Cowering")
		return decision.Cower
	case 1:
		return decision.Attack
	default:
		return decision.Defend
	}
}
