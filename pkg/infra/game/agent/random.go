package agent

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"math/rand"

	"github.com/benbjohnson/immutable"
)

type RandomAgent struct {
	bravery int
}

func NewRandomAgent() *RandomAgent {
	return &RandomAgent{bravery: 0}
}

func (r RandomAgent) GenerateActionDecision() decision.FightAction {
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

func (r RandomAgent) ProcessStartOfRound(view state.View, log immutable.Map[commons.ID, decision.FightAction]) {

}

func (r RandomAgent) ProcessFightDecisionRequestMessage(
	FightDecisionRequestMessage message.Message) message.FightDecisionMessage {
	return message.FightDecisionMessage{FightDecision: decision.Undecided}
}

func (r RandomAgent) ProcessFightDecisionMessage(message.FightDecisionMessage) {

}
