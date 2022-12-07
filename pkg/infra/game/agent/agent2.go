package agent

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	_ "infra/game/state"
	"math"
	"math/rand"

	"github.com/benbjohnson/immutable"
	"github.com/google/uuid"
)

// Agent2 type : private attributes of agent
type Agent2 struct {
	// Parameters for decision-making
	// Can be randomly generated per agent, tests need to be conducted to see
	// which values are the most efficient
	personalTendency    float64 // in [0,1] Tendency to fight, defend or cower
	replacementTendency float64 // in [0,1] Tendency to replace cowering agents on the battlefield
	estimationTendency  float64 // Tendency to go fight if it thinks not enough agents are fighting still

	decisionMap   []immutable.Map[commons.ID, decision.FightAction]
	baseAgentMap  []BaseAgent
	viewMap       []state.View
	agentStateMap []immutable.Map[commons.ID, state.HiddenAgentState]
	leaderMap     []commons.ID

	// Updates per Level
	baseAgentPerLevel   []BaseAgent
	fightResultPerLevel []commons.ImmutableList[decision.ImmutableFightResult]
	voteResultPerLevel  []immutable.Map[decision.Intent, uint]
}

// NewAgent2 : Constructor of Agent2 FIXME:!!!INITIALISE AGENT WITH VALUES!!!
func NewAgent2() Strategy {
	return &Agent2{}
}

func (a *Agent2) updateBaseAgentPerLevel(agent BaseAgent) {
	a.baseAgentPerLevel = append(a.baseAgentPerLevel, agent)
}

func (a *Agent2) updateFightResultPerLevel(result commons.ImmutableList[decision.ImmutableFightResult]) {
	a.fightResultPerLevel = append(a.fightResultPerLevel, result)
}
func (a *Agent2) updateVoteResultPerLevel(result immutable.Map[decision.Intent, uint]) {
	a.voteResultPerLevel = append(a.voteResultPerLevel, result)
}

// Description: The function is used to extract the agents' decisions made in the previous rounds
// Returns: If multi == true => array of multiple decisions maps (previous rounds), else => array with a single element (decision map of the last round)
func (a *Agent2) getDecisionHelper(multi bool) []immutable.Map[commons.ID, decision.FightAction] {
	if multi == true {
		return a.decisionMap
	} else {
		singleArray := a.decisionMap[len(a.decisionMap)-1:]
		return singleArray
	}
}

// Description: The function returns our agent's (for a specific instance of our agent type) BaseAgent structures from the previous rounds
// Returns: If multi ==  true => array of the previous BaseAgent structures, else => array with a single element (BaseAgent struct of the last round)
func (a *Agent2) getBaseHelper(multi bool) []BaseAgent {
	if multi == true {
		return a.baseAgentMap
	} else {
		singleArray := a.baseAgentMap[len(a.baseAgentMap)-1:]
		return singleArray
	}
}

// Description: The function is used to extract the general game info (a.k.a View) of the previous rounds
// Returns: If multi => array of the previous View structures, else => array with a single element (View struct of the last round)
func (a *Agent2) getViewHelper(multi bool) []state.View {
	if multi == true {
		return a.viewMap
	} else {
		singleArray := a.viewMap[len(a.viewMap)-1:]
		return singleArray
	}
}

// Description: The function is used to extract the agents' HiddenAgentState struct of the previous rounds
// Returns: if multi => array of the previous HiddenAgentState structures, else => array with a single element (HiddenAgentState struct of the last round)
func (a *Agent2) getAgentStateHelper(multi bool) []immutable.Map[commons.ID, state.HiddenAgentState] {
	if multi == true {
		return a.agentStateMap
	} else {
		singleArray := a.agentStateMap[len(a.agentStateMap)-1:]
		return singleArray
	}
}

// Description: The function is used to extract the previous leaders' ids
// Returns: if multi => an array with all the previous leaders' ids, else => array with a single element (id of the last round's leader)
func (a *Agent2) getLeaderHelper(multi bool) []commons.ID {
	if multi == true {
		return a.leaderMap
	} else {
		singleArray := a.leaderMap[len(a.leaderMap)-1:]
		return singleArray
	}
}

// Description: The function is used to return the current Agent's health (Hp)
// Returns: uint
func (a *Agent2) getCurrentHp() uint {
	return a.getBaseHelper(false)[0].latestState.Hp
}

// Description: The function is used to return the current Agent's stamina
// Returns: uint
func (a *Agent2) getCurrentStamina() uint {
	return a.getBaseHelper(false)[0].latestState.Stamina
}

// Description: The function is used to return the current Agent's defence points
// Returns: uint
func (a *Agent2) getCurrentDefense() uint {
	return a.getBaseHelper(false)[0].latestState.Defense
}

// Description: The function is used to return the current Agent's attack points
// Returns: uint
func (a *Agent2) getCurrentAttack() uint {
	return a.getBaseHelper(false)[0].latestState.Attack
}

// Description: The function is used to return the current Agent's bonus defense points
// Returns: uint
func (a *Agent2) getBonusDefense() uint {
	AllAgentStates := a.getAgentStateHelper(false)[0]
	AgentID := a.getBaseHelper(false)[0].id
	AgentState, _ := AllAgentStates.Get(AgentID)
	return AgentState.BonusDefense
}

// Description: The function is used to return the current Agent's bonus attack points
// Returns: uint
func (a *Agent2) getBonusAttack() uint {
	AllAgentStates := a.getAgentStateHelper(false)[0]
	AgentID := a.getBaseHelper(false)[0].id
	AgentState, _ := AllAgentStates.Get(AgentID)
	return AgentState.BonusAttack
}

func (a *Agent2) updateDecisionHelper(log immutable.Map[commons.ID, decision.FightAction]) {
	a.decisionMap = append(a.decisionMap, log)
}
func (a *Agent2) updateBaseHelper(baseAgent BaseAgent) {
	a.baseAgentMap = append(a.baseAgentMap, baseAgent)
}
func (a *Agent2) updateViewHelper(view state.View) {
	a.viewMap = append(a.viewMap, view)
}
func (a *Agent2) updateAgentStateHelper(agentState immutable.Map[commons.ID, state.HiddenAgentState]) {
	a.agentStateMap = append(a.agentStateMap, agentState)
}

func (a *Agent2) updateLeaderHelper(leader commons.ID) {
	a.leaderMap = append(a.leaderMap, leader)
}

// HandleFightInformation TODO: Implement me!
// Description: Used to extract agent information
// Return:		nil
func (a *Agent2) HandleFightInformation(m message.TaggedInformMessage[message.FightInform], baseAgent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) {
	// Update Logs
	a.updateDecisionHelper(*log)

	//v, ok := view.AgentState().Get("asdf")
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

// HandleFightProposal: TODO: Implement me!
// Description: Throught that function our agent votes on a broadcastes proposal
func (a *Agent2) HandleFightProposal(proposal message.FightProposalMessage, baseAgent BaseAgent) decision.Intent {
	intent := rand.Intn(2)
	if intent == 0 {
		return decision.Positive
	} else {
		return decision.Negative
	}
}

// FightResolution: TODO: Implement me!
// Description: Through that function our agent provides a proposal
func (a *Agent2) FightResolution(baseAgent BaseAgent) message.MapProposal[decision.FightAction] {

	// Info update per round
	view := baseAgent.View()
	a.updateBaseHelper(baseAgent)
	a.updateViewHelper(view)
	a.updateAgentStateHelper(baseAgent.view.AgentState())
	a.updateLeaderHelper(view.CurrentLeader())

	actions := make(map[commons.ID]decision.FightAction)

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

// HandleFightProposalRequest TODO: Implement me!
// Description: Only called as a leader: True for broadcasting the proposal / False for declining the proposal
// Return:		Bool: True/False
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

// UpdateInternalState TODO: Implement me!
// Description: the function is called at the end of each level (provides a list of type FightResult / can be thought as raw & processed overall game info)
func (a *Agent2) UpdateInternalState(baseAgent BaseAgent, fightResult *commons.ImmutableList[decision.ImmutableFightResult], voteResult *immutable.Map[decision.Intent, uint]) {
	a.updateBaseAgentPerLevel(baseAgent)
	a.updateFightResultPerLevel(*fightResult)
	a.updateVoteResultPerLevel(*voteResult)
}

// DonateToHpPool TODO: Implement me!
// Description: The function returns the amount of Hp that our agent is willing to donate to the HpPool
func (a *Agent2) DonateToHpPool(baseAgent BaseAgent) uint {
	return uint(rand.Intn(int(baseAgent.AgentState().Hp)))
}
