/*
	// assuming a decision was already made before, the following is run inside a fight round loop to drink potions and change our decision if necessary

	agentState := view.AgentState()[a.BaseAgent.Id]; // our agent's current state
	// WARNING : not sure how go works exactly, this should be a reference so that we always point to the current value in the following statements


	// Potions
	if agentState.Hp < 25 && hasHPpotion {
		// drink health potion
	}

	=====================TASKS=======================
	!!!FUNCTIONS TO IMPLEMENT!!!

========================

	Tim:

	// assume the following function that estimates if the amount of other agents is enough to kill the monsters (so we can cower if we're low)
	func enoughOtherAgentsFighting(...)
	input args : global state map
	output : bool

	for example, remember the amount of agents that fought last round and how much damage was dealt, to get a mean damage/agent
	then estimate this round's damage based on the number of fighting agents WITH HEALTH NOT LOW (otherwise assume they'll be cowering to survive instead)


	// cower if low health and not critically needed
	if agentState.Hp < 15 {
		if enoughOtherAgentsFighting(...) {
			return decision.Cower
		}
	}

========================

========================

	Giorgos:

	1. func History() {}
	2. func FreuencyOfDecisions() {}

=========================

=========================

	Vasos:

	// drink potions if attacking or defending
	if decision != decision.Cower {
		if hasAttackPotion && decision == decision.Attack {
			// drink it
		}
		if hasShieldPotion {
			// drink it
		}
	}

==========================
*/
package agent

import (
	"github.com/benbjohnson/immutable"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"infra/logging"
	"math/rand"
)

// Agent2 type : private attributes of agent
type Agent2 struct {
	HistoryMap []map[uint]map[commons.ID]decision.FightAction
}

// NewAgent2 : Constructor of Agent2
func NewAgent2() Strategy {
	return &Agent2{}
}

func History(baseAgent BaseAgent, currentLevel uint) map[uint]map[commons.ID]decision.FightAction {
	draftMap := make(map[commons.ID]decision.FightAction)
	draftLevelMap := make(map[uint]map[commons.ID]decision.FightAction)
	draftMap[baseAgent.id] = decision.Undecided
	draftLevelMap[currentLevel] = draftMap
	return draftLevelMap
}

// HandleFightInformation TODO: Implement me!
// Description: Used to extract agent information
// Return:		nil
func (a Agent2) HandleFightInformation(m message.TaggedMessage, view *state.View, agent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) {
	agent.Log(logging.Trace, logging.LogField{}, "Something")
	currentHistory := History(agent, view.CurrentLevel())
	a.HistoryMap = append(a.HistoryMap, currentHistory)
}

// HandleFightRequest TODO: Implement me!
// Description: Used for comms to request p2p message probably? Not Sure!
// Return		Message Payload
func (a Agent2) HandleFightRequest(m message.TaggedMessage, view *state.View, log *immutable.Map[commons.ID, decision.FightAction]) message.Payload {
	return nil
}

// CurrentAction TODO: Implement me!
// Description: Logic of Fighting Action Decision Making.
// Return:		Cower, Defend or Attack decision.
func (a Agent2) CurrentAction() decision.FightAction {

	fight := rand.Intn(10)
	switch {
	case fight == 0:
		return decision.Cower
	case (fight <= 4) && (fight > 0):
		return decision.Defend
	default:
		return decision.Attack
	}
}

// CreateManifesto TODO: Implement me!
// Description: Used to give Manifesto Information if elected Leader.
// Return:		The Manifesto with FightImposition, LootImposition, term length and overthrow threshold.
func (a Agent2) CreateManifesto(view *state.View, baseAgent BaseAgent) *decision.Manifesto {
	draftManifesto := decision.NewManifesto(true, false, 5, 40)
	return draftManifesto
}

// HandleConfidencePoll TODO: Implement me!
// Description: Used for voting on confidence for Leader.
// Return:		Positive, Negative, or Abstain decision.
func (a Agent2) HandleConfidencePoll(view *state.View, baseAgent BaseAgent) decision.Intent {
	return decision.Positive
}

// HandleElectionBallot TODO: Implement me!
// Description: Used to elect a Leader.
// Return:  	A single Commons.ID for choose-one voting or an array of commons.ID of top leader choices for ranked-voting.
func (a Agent2) HandleElectionBallot(view *state.View, baseAgent BaseAgent, params *decision.ElectionParams) decision.Ballot {
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

// Default TODO: Add Default logic to be taken by agent
// Description: Default action of agent
func (a Agent2) Default() decision.FightAction {
	return decision.Undecided
}
