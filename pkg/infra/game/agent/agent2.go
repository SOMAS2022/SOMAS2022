package agent

import (
	"github.com/benbjohnson/immutable"
	"infra/game/decision"
	"infra/game/state"
	"math/rand"
    "infra/game/message"
    "infra/game/commons"
)

type Agent2 struct {
    bravery int
}

func (a Agent2) Default() decision.FightAction {
    //TODO implement me
    panic("implement me")
}

func NewAgent2() *Agent2 {
    return &Agent2{bravery: 0}
}

func (a Agent2) HandleFightMessage(m message.TaggedMessage, view *state.View, agent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) decision.FightAction {
/*
Fighting is similar to Random Agent for MVP.
- Agent attacks by default (50% of the time)
- Agent defends 40% of the time
- Agent cowers 10% of time
*/
	fight := rand.Intn(10)
	switch {
	case fight == 0:
		return decision.Cower
	case (fight <= 4) && (fight > 0) :
		return decision.Defend
    default:
        return decision.Attack
	}
}
