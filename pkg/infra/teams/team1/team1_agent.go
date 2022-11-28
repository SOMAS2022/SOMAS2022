/*******************************************************
* Copyright (C) 2022 Team 1 @ SOMAS2022
*
* This file is part of SOMAS2022.
*
* This file or its contents can not be copied and/or used
* without the express permission of Team 1, SOMAS2022
*******************************************************/

package team1

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"infra/teams/team1/utils"
	"math/rand"
	"time"

	"github.com/benbjohnson/immutable"
)

type ProbabilisticAgent struct {
	fightDecisionCDF []float32

	// current fightDecision
	currentDecision decision.FightAction

	// Metadata
	battleUtility float32

	socialCapital    map[string][4]float64
	lastLevelUpdated uint // The level at which the socialCapital was last updated

	// Four metrics for each agent's perception of other agents. Three metrics are borrowed from Ostrom-Ahn
	// social capital model: (The ordering below is the same as the ordering in the array)
	// 1. Institutions: Influenced by institutional actions such as voting. Also influenced by any elected
	//positions the agent currently holds.
	// 2. Networks: Messages received about an agent from other agents.
	// The impact which a message from another agent has on the network metric depends on the overall perception of
	// the agent sending the message.
	// 3. Trustworthiness: Measures how likely the agent is to cooperate based on historic decisions.
	// Finally, our addition to the model is
	// 4. Favours: A metric of how many favours an agent has done for another agent, compared to the amount of
	//favours the other agent has done for them. This metric is supposed to capture social expectations, such
	//as if I went out with the trash last time, then I feel that next time someone else should do it. The value
	// of the  favours metric will over time impact the trustworthiness metric of another agent. In many ways,
	// favours is a short time horizon equivalent of trustworthiness.

	// Forgiveness: In order to enable even the bitterest of enemies to start cooperating again in the long term,
	// the metrics will decay with time. However, negative metrics will decay far faster than positive metrics.
	// This is similar to how you can have a massive fight and hate someone one day, and then be friends again the
	// next. On the other hand, you can go years without seeing a friend and still have a very positive impression
	// of them.

	// Messages we should have: Denounce and Praise
}

func (r ProbabilisticAgent) Default() decision.FightAction {
	//TODO implement me
	panic("implement me")
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

func (r ProbabilisticAgent) CurrentAction() decision.FightAction {
	return r.currentDecision
}

func (r ProbabilisticAgent) HandleFightRequest(m message.TaggedMessage, view *state.View, log *immutable.Map[commons.ID, decision.FightAction]) message.Payload {
	return nil
}

func (r ProbabilisticAgent) HandleFightInformation(m message.TaggedMessage, view *state.View, agent agent.BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) {
	r.UpdateMetadata(agent)

	r.updateSocialCapital(m, view, agent, log)

	// Calculate utility value of each action
	utilCower := r.utilityValue(decision.Cower, view, agent)
	utilAttack := r.utilityValue(decision.Attack, view, agent)
	utilDefend := r.utilityValue(decision.Defend, view, agent)

	// Apply softmax to get probabilities
	softArray := utils.Softmax([3]float64{utilCower, utilAttack, utilDefend})

	// Make number representation incremental
	probArray := utils.MakeIncremental(softArray)

	/*it := view.AgentState().Iterator()
	nextId, _, _ := it.Next()
	if agent.Id == nextId {
		fmt.Println(utilCower)
		fmt.Println([3]float64{utilCower, utilAttack, utilDefend})
		fmt.Println(softArray)
		fmt.Println(probArray)
	}*/

	// Initialise a random seed
	rand.Seed(time.Now().UnixNano())

	// Do action with probability based on utility value
	switch random := rand.Float64(); {
	case 0.0 < random && random < probArray[0]:
		r.currentDecision = decision.Cower
	case probArray[0] < random && random < probArray[1]:
		r.currentDecision = decision.Attack
	case probArray[1] < random && random < probArray[2]:
		r.currentDecision = decision.Defend
	default:
		r.currentDecision = decision.Attack
	}

	return
}

func (r ProbabilisticAgent) UpdateMetadata(self agent.BaseAgent) {
	r.battleUtility = utils.AgentBattleUtility(self.ViewState())
}

// Calculate utility value of different decisions
func (r ProbabilisticAgent) utilityValue(action decision.FightAction, _ *state.View, agent agent.BaseAgent) float64 {
	// Utility of each action is dependent on relationship with others. If agent hates all other agents, then
	// will only act in its own interest.

	agentStats := agent.ViewState()

	switch action {
	case decision.Cower:
		// Goes down with health, and down with stamina and down with high social capital of others
		return 0.005 * float64(1000-int(agentStats.Hp))
	case decision.Attack:
		// Goes up with health, and up with stamina and up with high social capital of others
		return 0.005 * float64(int(agentStats.Hp))
	case decision.Defend:
		return 0.005 * float64(int(agentStats.Hp))
	default:
		return 1
	}
}

// Called any time a message is received, initialises or updates the socialCapital map
func (r ProbabilisticAgent) updateSocialCapital(_ message.TaggedMessage, view *state.View, agent agent.BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) {

	// Ensure that socialCapital map is initialised
	agentState := view.AgentState()
	agentStateLength := agentState.Len()
	if len(r.socialCapital) == 0 && agentStateLength > 1 {
		// Create empty map
		r.socialCapital = map[string][4]float64{}

		// Populate map with every currently living agent, and calculate socialCapital based on log
		itr := agentState.Iterator()
		for !itr.Done() {
			key, _, _ := itr.Next()

			action, exists := log.Get(key)
			if exists { // If agent exists in log, calculate socialCapital
				r.socialCapital[key] = utils.BoundArray(utils.ActionSentiment(action))
			} else { // Else initialize socialCapital to 0
				r.socialCapital[key] = [4]float64{0.0, 0.0, 0.0, 0.0}
			}
		}

		// Delete the agents own id from the socialCapital array
		delete(r.socialCapital, agent.Name())

		// Set the lastLevelUpdated variable
		r.lastLevelUpdated = view.CurrentLevel()
	} else if r.lastLevelUpdated < view.CurrentLevel() { // socialCapital variable already exists
		for key := range r.socialCapital {

			// Remove any agents that have died from socialCapital map (Might be unnecessary as it adds a lot of computation)
			_, exists := agentState.Get(key)
			if !exists {
				delete(r.socialCapital, key)
			}

			// Decay socialCapital values
			r.socialCapital[key] = utils.DecayArray(r.socialCapital[key])

			// TODO: Update of socialCaptial should be dependent on the agents own action (especially for favours)
			// Update socialCapital based on log
			action, exists := log.Get(key)
			if exists {
				r.socialCapital[key] = utils.AddArrays(r.socialCapital[key], utils.BoundArray(utils.ActionSentiment(action)))
			}
		}

		// Set lastLevelUpdated to current level
		r.lastLevelUpdated = view.CurrentLevel()
	}

	// Ensure all socialCapital values are between -1 and 1
	for key := range r.socialCapital {
		r.socialCapital[key] = utils.BoundArray(r.socialCapital[key])
	}
}

func (r ProbabilisticAgent) CreateManifesto(view *state.View, baseAgent agent.BaseAgent) *decision.Manifesto {
	manifesto := decision.NewManifesto(true, false, 10, 50)
	return manifesto
}

func (r ProbabilisticAgent) HandleConfidencePoll(view *state.View, baseAgent agent.BaseAgent) decision.Intent {
	switch rand.Intn(3) {
	case 0:
		return decision.Abstain
	case 1:
		return decision.Negative
	default:
		return decision.Positive
	}
}

func (r ProbabilisticAgent) HandleElectionBallot(view *state.View, _ agent.BaseAgent, _ *decision.ElectionParams) decision.Ballot {
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
