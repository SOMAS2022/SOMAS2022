package agent

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	_ "infra/game/state"
	"infra/logging"
	"math"
	"math/rand"

	"github.com/benbjohnson/immutable"
	"github.com/google/uuid"
)

// Agent2 type : private attributes of agent
type Agent2 struct {
	HistoryMap []map[uint][]immutable.Map[commons.ID, decision.FightAction]

	// Parameters for decision making
	// Can be randomly generated per agent, tests need to be conducted to see
	// which values are the most efficient
	personalTendency    float64 // in [0,1] Tendency to fight, defend or cower
	replacementTendency float64 // in [0,1] Tendency to replace cowering agents on the battlefield
	estimationTendency  float64 // Tendency to go fight if it thinks not enough agents are fighting still
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
func (a *Agent2) HandleFightInformation(m message.TaggedInformMessage[message.FightInform], baseAgent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) {
	baseAgent.Log(logging.Trace, logging.LogField{}, "Something")
	currentHistory := a.History(*log, baseAgent.view.CurrentLevel())
	a.HistoryMap = append(a.HistoryMap, currentHistory)
}

// HandleFightRequest TODO: Implement me!
// Description: Used for comms to request p2p message probably? Not Sure!
// Return		Message Payload
func (a *Agent2) HandleFightRequest(m message.TaggedRequestMessage[message.FightRequest], log *immutable.Map[commons.ID, decision.FightAction]) message.FightInform {
	return nil
}

// Logistic function
func logistic(x float64, k float64, x0 float64) float64 {
	return 1 / (1 + math.Exp(-k*(x-x0)))
}

/*
==============

Work In Progress Tim

// Description : Compare defense and attack potential, output a decision
// Return:		Cower, Defend or Attack decision.
func (a *Agent2) initialDecision() decision.FightAction {

	// method to retrieve state ?

	// Bravery is a function of health
	bravery := logistic(Health, 0.1, 50)

	// If current bravery is higher than parameter tendency, do something
	if bravery >= a.personalTendency {
		if Stamina >= Attack+AttackBonus {
			return decision.Attack
		} else if Stamina >= Defense+DefenseBonus {
			return decision.Defend
		}
	}

	return decision.Cower // Else cower
}

// Description : Compare current number of cowering agents to previous numbers
// and possibly replace them
// Return:		Cower, Defend or Attack decision.
func (a *Agent2) replaceDecision() decision.FightAction {
	mean10LastRounds := ...
	currentFighting := ...

	bravery := logistic(Health, 0.1, 50)

	if bravery >= a.replacementTendency {
		if Stamina >= Attack+AttackBonus {
			return decision.Attack
		} else if Stamina >= Defense+DefenseBonus {
			return decision.Defend
		}
	}

	return decision.Cower
}

// Description : Estimate current damage and possibly change decision
// Return:		Cower, Defend or Attack decision.
func (a *Agent2) estimateDecision() decision.FightAction {
	lastTotalAttack := ...
	lastAgents := ...
	lastTotalDefense := ...

	currentAgents := ...
	estimatedTotalAttack := currentAgents * lastTotalAttack / lastAgents
	estimatedTotalDefense := currentAgents * lastTotalDefense / lastAgents

	diffAttack = estimatedTotalAtack - lastTotalAttack
	diffDefense = estimatedTotalDefense - lastTotalDefense

	// Cower if there are more attackers and defenders
	if(diffAttack > 0 && diffDefense > 0){
		return decision.Cower
	}

	bravery := logistic(Health, 0.1, 50)

	if bravery >= a.estimationTendency {
		if Stamina >= Attack+AttackBonus {
			return decision.Attack
		} else if Stamina >= Defense+DefenseBonus {
			return decision.Defend
		}
	}
}

// CurrentAction TODO: Implement me!
// Description: Logic of Fighting Action Decision-Making.
// Return:		Cower, Defend or Attack decision.
func (a *Agent2) CurrentAction() decision.FightAction {

	// If not enough Stamina, no choice
	if Stamina < Attack+AttackBonus && Stamina < Defend+DefendBonus {
		return decision.Cower
	}

	currentDecision := a.initialDecision()

	if currentDecision == decision.Cower {
		currentDecision = a.replaceDecision()
	}
	if currentDecision == decision.Cower {
		currentDecision = a.estimateDecision()
	}

	return currentDecision
}

=============
*/

// DEFAULT FUNCTION (TO DELETE)
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
func (a *Agent2) CreateManifesto(baseAgent BaseAgent) *decision.Manifesto {
	draftManifesto := decision.NewManifesto(false, true, 4, 51)
	return draftManifesto
}

// HandleConfidencePoll TODO: Implement me!
// Description: Used for voting on confidence for Leader.
// Return:		Positive, Negative, or Abstain decision.
func (a *Agent2) HandleConfidencePoll(baseAgent BaseAgent) decision.Intent {
	switch rand.Intn(3) {
	case 0:
		return decision.Abstain
	case 1:
		return decision.Negative
	default:
		return decision.Positive
	}
}

// HandleElectionBallot TODO: Implement me!
// Description: Used to elect a Leader.
// Return:  	A single Commons.ID for choose-one voting or an array of commons.ID of top leader choices for ranked-voting.
func (a *Agent2) HandleElectionBallot(baseAgent BaseAgent, params *decision.ElectionParams) decision.Ballot {
	// Extract ID of alive agents
	view := baseAgent.View()
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

func (a *Agent2) HandleFightProposal(proposal message.FightProposalMessage, baseAgent BaseAgent) decision.Intent {
	intent := rand.Intn(2)
	if intent == 0 {
		return decision.Positive
	} else {
		return decision.Negative
	}
}
func (a *Agent2) FightResolution(baseAgent BaseAgent) message.MapProposal[decision.FightAction] {
	actions := make(map[commons.ID]decision.FightAction)
	view := baseAgent.View()
	agentState := view.AgentState()
	itr := agentState.Iterator()
	for !itr.Done() {
		id, _, ok := itr.Next()
		if !ok {
			break
		}

		switch rand.Intn(3) {
		case 0:
			actions[id] = decision.Attack
		case 1:
			actions[id] = decision.Defend
		default:
			actions[id] = decision.Cower
		}
	}
	prop := message.NewProposal(uuid.NewString(), commons.MapToImmutable(actions))
	return *prop
}

// HandleFightProposalRequest only called as leader
func (a *Agent2) HandleFightProposalRequest(proposal message.FightProposalMessage, baseAgent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) bool {
	switch rand.Intn(2) {
	case 0:
		return true
	default:
		return false
	}
}

// HandleUpdateWeapon return the index of the weapon you want to use in AgentState.Weapons
func (a *Agent2) HandleUpdateWeapon(baseAgent BaseAgent) decision.ItemIdx {
	// weapons := b.AgentState().Weapons
	// return decision.ItemIdx(rand.Intn(weapons.Len() + 1))

	// 0th weapon has the greatest attack points
	return decision.ItemIdx(0)
}

// HandleUpdateShield return the index of the shield you want to use in AgentState.Shields
func (a *Agent2) HandleUpdateShield(baseAgent BaseAgent) decision.ItemIdx {
	// shields := b.AgentState().Shields
	// return decision.ItemIdx(rand.Intn(shields.Len() + 1))
	return decision.ItemIdx(0)
}
func (a *Agent2) UpdateInternalState(baseAgent BaseAgent, fightResult *commons.ImmutableList[decision.ImmutableFightResult], voteResult *immutable.Map[decision.Intent, uint]) {
}
func (a *Agent2) DonateToHpPool(baseAgent BaseAgent) uint {
	return uint(rand.Intn(int(baseAgent.AgentState().Hp)))
}