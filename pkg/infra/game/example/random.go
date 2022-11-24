package example

import (
	"infra/game/agent"
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

func (r RandomAgent) HandleFightInformation(_ message.TaggedMessage, _ *state.View, agent agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) {
	agent.Log(logging.Trace, logging.LogField{"bravery": r.bravery}, "Cowering")
}

func (r RandomAgent) HandleFightRequest(_ message.TaggedMessage, _ *state.View, _ *immutable.Map[commons.ID, decision.FightAction]) message.Payload {
	return nil
}

func (r RandomAgent) CurrentAction() decision.FightAction {
	fight := rand.Intn(3)
	switch fight {
	case 0:
		return decision.Cower
	case 1:
		return decision.Attack
	default:
		return decision.Defend
	}
}

func NewRandomAgent() *RandomAgent {
	return &RandomAgent{bravery: 0}
}
