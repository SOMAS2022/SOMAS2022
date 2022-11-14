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
	var fightDecisionToTake *decision.FightAction

	if m.Sender == "server" {
		r.bravery += rand.Intn(5)
	}

	currentState, _ := view.AgentState().Get(agent.Id)
	r.bravery += rand.Intn(10)

	lastFightDecision, _ := log.Get(agent.Id)
	if lastFightDecision == decision.Attack {
		r.bravery += rand.Intn(5)
	}

	if float32(r.bravery) > float32(rand.Intn(5)/10)*float32(view.MonsterHealth()) {
		fightDecisionToTake = new(decision.FightAction)
		*fightDecisionToTake = decision.Attack
	}
	if float32(currentState.Hp) < float32(rand.Intn(5)/10)*float32(view.MonsterHealth()) {
		fightDecisionToTake = new(decision.FightAction)
		*fightDecisionToTake = decision.Defend
	}

	if rand.Intn(100) < 15 {
		fightDecisionToTake = new(decision.FightAction)
		*fightDecisionToTake = decision.Cower
	}

	return fightDecisionToTake
}
