package team0

import (
	"math/rand"

	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"

	"github.com/benbjohnson/immutable"
)

type ProbabilisticAgent struct {
	fightDecisionCDF []float32
}

func (r ProbabilisticAgent) Default() decision.FightAction {
	// TODO implement me
	panic("implement me")
}

/**
 * Create agent with given probability of cowering, attacking, defending.
 */
func NewProbabilisticAgent(pCower float32, pAttack float32, pDefend float32) *ProbabilisticAgent {
	// Ref: https://stackoverflow.com/questions/50507513/golang-choice-number-from-slice-array-with-given-probability
	pdf := []float32{pCower, pAttack, pDefend}
	// get cdf
	cdf := []float32{0.0, 0.0, 0.0}
	cdf[0] = pdf[0]
	for i := 1; i < 3; i++ {
		cdf[i] = cdf[i-1] + pdf[i]
	}
	return &ProbabilisticAgent{fightDecisionCDF: cdf}
}

func (r ProbabilisticAgent) HandleFightMessage(m message.TaggedMessage, view *state.View, agent agent.BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) decision.FightAction {
	dice := rand.Float32()

	fight := 0
	for dice > r.fightDecisionCDF[fight] {
		fight++
	}
	switch fight {
	case 0:
		return decision.Cower
	case 1:
		return decision.Attack
	default:
		return decision.Defend
	}
}
