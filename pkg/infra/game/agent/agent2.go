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
	HistoryMap []map[uint][]immutable.Map[commons.ID, decision.FightAction]
}

// NewAgent2 : Constructor of Agent2
func NewAgent2() Strategy {
	return &Agent2{}
}

func (a *Agent2) History(log immutable.Map[commons.ID, decision.FightAction], currentLevel uint) map[uint][]immutable.Map[commons.ID, decision.FightAction] {
	draftMap := make([]immutable.Map[commons.ID, decision.FightAction], 0)
	draftLevelMap := make(map[uint][]immutable.Map[commons.ID, decision.FightAction])
	draftMap = append(draftMap, log)
	draftLevelMap[currentLevel] = draftMap
	return draftLevelMap
}

// HandleFightInformation TODO: Implement me!
// Description: Used to extract agent information
// Return:		nil
func (a *Agent2) HandleFightInformation(m message.TaggedMessage, view *state.View, agent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) {
	agent.Log(logging.Trace, logging.LogField{}, "Something")
	currentHistory := a.History(*log, view.CurrentLevel())
	a.HistoryMap = append(a.HistoryMap, currentHistory)
}

// HandleFightRequest TODO: Implement me!
// Description: Used for comms to request p2p message probably? Not Sure!
// Return		Message Payload
func (a *Agent2) HandleFightRequest(m message.TaggedMessage, view *state.View, log *immutable.Map[commons.ID, decision.FightAction]) message.Payload {
	return nil
}

// CurrentAction TODO: Implement me!
// Description: Logic of Fighting Action Decision-Making.
// Return:		Cower, Defend or Attack decision.
func (a *Agent2) CurrentAction() decision.FightAction {
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
func (a *Agent2) CreateManifesto(view *state.View, baseAgent BaseAgent) *decision.Manifesto {
	draftManifesto := decision.NewManifesto(false, true, 4, 51)
	return draftManifesto
}

// HandleConfidencePoll TODO: Implement me!
// Description: Used for voting on confidence for Leader.
// Return:		Positive, Negative, or Abstain decision.
func (a *Agent2) HandleConfidencePoll(view *state.View, baseAgent BaseAgent) decision.Intent {
	return decision.Positive
}

// HandleElectionBallot TODO: Implement me!
// Description: Used to elect a Leader.
// Return:  	A single Commons.ID for choose-one voting or an array of commons.ID of top leader choices for ranked-voting.
func (a *Agent2) HandleElectionBallot(view *state.View, baseAgent BaseAgent, params *decision.ElectionParams) decision.Ballot {
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
func (a *Agent2) Default() decision.FightAction {
	/*
	   - Agent attacks by default (50% of the time)
	   - Agent defends 40% of the time
	   - Agent cowers 10% of time
	*/
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
