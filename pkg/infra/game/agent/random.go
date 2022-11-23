package agent

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"infra/game/strategy"
	"math/rand"

	"github.com/benbjohnson/immutable"
)

type RandomAgent struct {
	bravery int
}

func (r RandomAgent) ProcessStartOfRound(ba *message.BaseAgent, view *state.View, log *immutable.Map[commons.ID, decision.FightAction]) {

}

func (r RandomAgent) ProcessFightDecisionRequestMessage(ba *message.BaseAgent, FightDecisionRequestMessage message.Message, view *state.View, log *immutable.Map[commons.ID, decision.FightAction]) strategy.FightDecisionMessage {
	return strategy.FightDecisionMessage{FightDecision: decision.Undecided}
}

func (r RandomAgent) ProcessFightDecisionMessage(ba *message.BaseAgent, m strategy.FightDecisionMessage, view *state.View, log *immutable.Map[commons.ID, decision.FightAction]) {

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

//func (r RandomAgent) ProcessStartOfRound(ba *message.BaseAgent, view *state.View, log *immutable.Map[commons.ID, decision.FightAction]) {
//
//}
//
//func (r RandomAgent) ProcessFightDecisionRequestMessage(
//	ba *message.BaseAgent,
//	FightDecisionRequestMessage message.Message) strategy.FightDecisionMessage {
//	return strategy.FightDecisionMessage{FightDecision: decision.Undecided}
//}
//
//func (r RandomAgent) ProcessFightDecisionMessage(ba *message.BaseAgent, strategy.FightDecisionMessage) {
//
//}
