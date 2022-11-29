package team0

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"math/rand"

	"github.com/benbjohnson/immutable"
)

type ProbabilisticAgent struct {
	fightDecisionCDF []float32
	fightAction      decision.FightAction
}

func (r *ProbabilisticAgent) Default() decision.FightAction {
	//TODO implement me
	panic("implement me")
}

// Demonstrate creating a strategy with input parameters
func CreateAggressiveAgent() agent.Strategy {
	return NewProbabilisticAgent(0.1, 0.8, 0.1)
}

func CreateDefensiveAgent() agent.Strategy {
	return NewProbabilisticAgent(0.1, 0.1, 0.8)
}

func CreateCoweringAgent() agent.Strategy {
	return NewProbabilisticAgent(0.8, 0.1, 0.1)
}

/**
 * Create agent with given probability of cowering, attacking, defending
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

func (r *ProbabilisticAgent) CreateManifesto(view *state.View, baseAgent agent.BaseAgent) *decision.Manifesto {
	manifesto := decision.NewManifesto(true, false, 10, 50)
	return manifesto
}

func (r *ProbabilisticAgent) HandleConfidencePoll(view *state.View, baseAgent agent.BaseAgent) decision.Intent {
	switch rand.Intn(3) {
	case 0:
		return decision.Abstain
	case 1:
		return decision.Negative
	default:
		return decision.Positive
	}
}

func (r *ProbabilisticAgent) HandleFightInformation(_ message.TaggedMessage, _ *state.View, agent agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) {
	dice := rand.Float32()

	fight := 0
	for dice > r.fightDecisionCDF[fight] {
		fight++
	}
	switch fight {
	case 0:
		r.fightAction = decision.Cower
	case 1:
		r.fightAction = decision.Attack
	default:
		r.fightAction = decision.Defend
	}
}

func (r *ProbabilisticAgent) HandleFightRequest(_ message.TaggedMessage, _ *state.View, _ *immutable.Map[commons.ID, decision.FightAction]) message.Payload {
	return nil
}

func (r *ProbabilisticAgent) CurrentAction() decision.FightAction {
	return r.fightAction
}

func (r *ProbabilisticAgent) HandleElectionBallot(view *state.View, _ agent.BaseAgent, _ *decision.ElectionParams) decision.Ballot {
	// Extract ID of alive agents
	agentState := view.AgentState()
	aliveAgentIds := make([]string, agentState.Len())
	i := 0
	itr := agentState.Iterator()
	for !itr.Done() {
		id, a, ok := itr.Next()
		if ok && a.Hp > 0 {
			aliveAgentIds[i] = id
			i++
		}
	}

	// Randomly fill the ballot
	var ballot decision.Ballot
	numAliveAgents := len(aliveAgentIds)
	numCandidate := 2
	for i := 0; i < numCandidate; i++ {
		randomIdx := rand.Intn(numAliveAgents)
		randomCandidate := aliveAgentIds[uint(randomIdx)]
		ballot = append(ballot, randomCandidate)
	}

	return ballot
}
