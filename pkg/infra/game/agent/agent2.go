package agent

import (
	"github.com/benbjohnson/immutable"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/message/proposal"
	"infra/game/state"
	_ "infra/game/state"
	"infra/logging"
	"math"
	"math/rand"
	"reflect"
	"sort"
)

// Agent2 type : private attributes of agent
type Agent2 struct {
	// Parameters for decision-making
	// Can be randomly generated per agent, tests need to be conducted to see
	// which values are the most efficient
	personalTendency    float64 // in [0,1] Tendency to fight, defend or cower
	replacementTendency float64 // in [0,1] Tendency to replace cowering agents on the battlefield
	estimationTendency  float64 // in [0,1] Tendency to go fight if it thinks not enough agents are fighting still

	decisionMap   []immutable.Map[commons.ID, decision.FightAction]
	baseAgentMap  []BaseAgent
	viewMap       []state.View
	agentStateMap []immutable.Map[commons.ID, state.HiddenAgentState]
	leaderMap     []commons.ID

	// Updates per Level
	baseAgentPerLevel   []BaseAgent
	fightResultPerLevel []commons.ImmutableList[decision.ImmutableFightResult]
	voteResultPerLevel  []immutable.Map[decision.Intent, uint]
	governmentTimeline  []leaderInfo
	haveElections       bool
	// CurrentLevelAverages

	avgHp      uint
	avgDefend  uint
	avgAttack  uint
	avgStamina uint

	// To decide how to vote in no-confidence vote at the end of each level, use a social capital framework with weighted factors and a binary activation function to decide yes/no
	// These are:
	avgSurvivalCurrTerm       float64 // average % of agents alive at the end of a level during current leadership term (+ve relationship, high weighting)
	avgSurvivalPastTerms      float64 // average % of agents alive at the end of a level from past leadership terms of that agent (+ve)
	avgSurvival               float64 // average % of agents alive at the end of a level from all past leadership terms (for comparison - namely normalize by this amount)
	avgBroadcastRateCurrTerm  float64 // % of the proposals we submitted that were actually accepted/broadcast by the leader, in current term (+ve, high weighting)
	avgBroadcastRatePastTerms float64 // % of the proposals we submitted that were actually accepted/broadcast by the leader, from past terms of that leader (+ve)
	avgBroadcastRate          float64 // % of the proposals we submitted that were actually accepted/broadcast, from all past leadership terms (again, normalize by this)
	leadershipXp              float64 // fraction of levels up to now that they were leader (+ve)
	noConfRate                float64 // fraction of their terms they were voted out prematurely (-ve))
	avgLeadershipXp           float64 // avg fraction of levels up to now that any one agent is leader
	avgNoConfRate             float64 // avg fraction of an agent's leaderships terms that he is voted out
	// These variables are marked with -- below
	// For these, firstly we need a history data helper function that returns an array of the form:
	// leader_timeline_array [{id, manifesto, duration, leader_stats}, {id, manifesto, duration, leader_stats}, ...]
	// The object of type leader_stats will contain the following items, corresponding JUST to that elapsed leadership term:
	avgTermSurvival      float64 // (calculate for each level of their leadership and average)
	avgTermBroadcastRate float64 // (calculate for each round/level? of their leadership and average)
	noConf               bool    // whether they were voted out of that term
	// This array is best created in the election function that is only called at the end of one leadership term / start of another
	// It's best to have private attributes that accrue raw data and then reset - some every new term, some every new level
	// These are used by the confidence function at the end of every level to actually yield the no conf vote, and by the election function at the end of a term to calculate stats and append to leader_timeline_array, and to vote
	// Namely, the ones we reset after every level:
	numAgentsBeginLevel uint //(actually, do we only have list of agent IDs?)
	numAgentsEndLevel   uint
	proposalsTotal      uint // how many proposals we put forward that level (necessarily equal to rounds?)
	proposalsBroadcast  uint // how many of these were broadcast
	// And variables we re-calculate every level, but reset every election (no need for arrays for the raw data from which we calculate them):
	survivalRates  []float64 // array of % of agents alive at the end of each level (this array is appended to at the end of every level, and resets every election, so that each elem corresponds to a level in a leadership term)
	broadcastRates []float64 // % of the proposals we submitted during the level that were actually accepted/broadcast (ditto)
	// And variables we re-calculate every level but never reset:
	// Then the ones we reset every election:
	termBeginLevel uint // level at which the leadership term began (can read from viewMap every time election func is called)
	termEndLevel   uint // level at which the leadership term ended (again, from viewMap, in election func)
	// And variables we calculate every election (using all previous vars), to add to leader_timeline_array (esp leader_stats):
	termDuration uint //  number of levels that leadership term lasted before elapsed or deposed (term_end_level - term_begin_level)

}

// NewAgent2 : Constructor of Agent2 FIXME:!!!INITIALISE AGENT WITH (MEANINGFUL) VALUES!!!
func NewAgent2() Strategy {
	// testing : random between [0.25,0.75]
	personalTendency := rand.Float64()*0.25 + 0.5
	replacementTendency := rand.Float64()*0.25 + 0.5
	estimationTendency := rand.Float64()*0.25 + 0.5
	return &Agent2{
		personalTendency:    personalTendency,
		replacementTendency: replacementTendency,
		estimationTendency:  estimationTendency,
		haveElections:       false,
	}
}

/* ---- HELPER FUNCTIONS ---- */
type leaderInfo struct {
	id         commons.ID
	manifesto  decision.Manifesto
	duration   uint
	overthrown bool
}

func (a *Agent2) newGovernmentTimeline(agent BaseAgent, haveElections bool) {
	view := agent.View()
	if haveElections == true {
		newLeaderInfo := leaderInfo{
			id:         view.CurrentLeader(),
			manifesto:  view.LeaderManifesto(),
			duration:   1,
			overthrown: false,
		}
		a.governmentTimeline = append(a.governmentTimeline, newLeaderInfo)
		a.haveElections = false
	} else {
		currentLeaderInfo := a.governmentTimeline[len(a.governmentTimeline)-1]
		a.governmentTimeline[len(a.governmentTimeline)-1] =
			leaderInfo{
				id:         currentLeaderInfo.id,
				manifesto:  currentLeaderInfo.manifesto,
				duration:   currentLeaderInfo.duration + uint(1),
				overthrown: false,
			}
	}

}

//var past_terms_of_curr_leader := make([]term_struct, 0)
//for leadership_term in leader_term_timeline_array {
//	if leadership_term[id] == curr_leader["id"] {
//		past_terms_of_curr_leader = append(past_terms_of_curr_leader, leadership_term) // will have redundant id key but whatever
//	}
//}
/* ---- UPDATES FUNCTIONS ---- */

func (a *Agent2) updateSocialCapital(agent BaseAgent) {
	view := agent.View()
	// Pseudocode for how stats are calculated from raw data, elsewhere in the program:
	a.survivalRates = append(a.survivalRates, float64(a.numAgentsEndLevel)/float64(a.numAgentsBeginLevel))
	a.broadcastRates = append(a.broadcastRates, float64(a.proposalsBroadcast)/float64(a.proposalsTotal))
	a.avgSurvivalCurrTerm = avg(a.survivalRates)
	a.avgBroadcastRateCurrTerm = avg(a.broadcastRates)
	a.avgSurvival = (a.avgSurvival*float64(view.CurrentLevel()-1) + float64(a.numAgentsEndLevel)/float64(a.numAgentsBeginLevel)) / float64(view.CurrentLevel())
	a.avgBroadcastRate = (a.avgBroadcastRate*float64(view.CurrentLevel()-1) + float64(a.proposalsBroadcast)/float64(a.proposalsTotal)) / float64(view.CurrentLevel())
}

func avg(array []float64) float64 {
	sum := 0.0
	for _, item := range array {
		sum += item
	}
	return sum / float64(len(array))
}

func updateAverages(agent BaseAgent) (uint, uint, uint, uint) {
	sumHp, sumDefend, sumAttack, sumStamina := uint(0), uint(0), uint(0), uint(0)
	view := agent.View()
	agentState := view.AgentState()
	numOfAgents := uint(0)
	for _, id := range commons.ImmutableMapKeys(agentState) {
		state, ok := agentState.Get(id)
		if ok && state.Hp > 0 {
			sumHp = sumHp + uint(state.Hp)
			sumDefend = sumDefend + state.BonusDefense
			sumAttack = sumAttack + state.BonusAttack
			sumStamina = sumStamina + uint(state.Stamina)
			numOfAgents += 1
		}
	}
	avgHp := sumHp / numOfAgents
	avgDefend := sumDefend / numOfAgents
	avgAttack := sumAttack / numOfAgents
	avgStamina := sumStamina / numOfAgents

	return avgHp, avgDefend, avgAttack, avgStamina
}

func minDefend(agent BaseAgent) proposal.Value {
	avgDefend := 0.0
	view := agent.View()
	agentState := view.AgentState()
	for _, id := range commons.ImmutableMapKeys(agentState) {
		state, ok := agentState.Get(id)
		if ok {
			avgDefend = avgDefend + float64(state.BonusDefense)
		}
	}
	return uint(0.50 * avgDefend)
}

func minAttack(agent BaseAgent) proposal.Value {
	state := agent.AgentState()
	attack := float64(state.BonusAttack())
	return uint(0.20 * attack)
}

func baseHealth(agent BaseAgent) proposal.Value {
	state := agent.AgentState()
	hp := float64(state.Hp)
	return uint(hp)
}

func minStamina(agent BaseAgent) proposal.Value {
	state := agent.AgentState()
	stamina := float64(state.Stamina)
	return uint(0.20 * stamina)
}

// Returns minimum Health that a healthy agent should have.
func minHealth(agent BaseAgent) proposal.Value {
	state := agent.AgentState()
	hp := float64(state.Hp)
	return uint(0.50 * hp)
}

// Returns Manifesto Effectiveness based on History
func weightedManifestoEffectiveness(agent BaseAgent, weight float64) float64 {
	return weight * 0.0
}

// Returns Number of Terms Agent2 served
// and was overthrown
func weightedOverthrowPercentage(agent BaseAgent, weight float64) float64 {
	return weight * 0.0
}

// Returns Adjusted Expertise with new mean and std provided.
func adjustedExpertise(agent BaseAgent, from float64, to float64) float64 {
	return from + (to-from)*expertise(agent)
}

func wasOverthrown(agent BaseAgent) bool {
	return false
}

func lastFightDecisionPower(agent BaseAgent, bias float64) float64 {
	if true {
		return bias
	} else {
		return 0.0
	}
}

func lastLootDecisionPower(agent BaseAgent, bias float64) float64 {
	if true {
		return bias
	} else {
		return 0.0
	}
}

func leaderElectedBefore(agent BaseAgent, bias float64) float64 {
	if true {
		return bias
	} else {
		return 0.0
	}
}

// SOT : [0,1] : map :(Overthrow_i_ranked+Term_i_ranked)_ranked
func SOT(agent BaseAgent) float64 {
	return 0.0
}
// returning prospect leader's score
func prospectLeaderScore(agent BaseAgent, par1 float64, par2 float64, par3 float64) float64 {
	return par1 + par2 + par3 + SOT(agent)
}

// weightedFracTermsDeposed : NumOfTimesDeposed / NumOfTimesElected
func weightedFracTermsDeposed(agent BaseAgent, id commons.ID) float64 {
	return 0.0
}

// weightedAvgSurRateUnderLeader : Sum_terms(Sum_levels(survival_rate))/(Sum_terms(Sum_levels(1)))
func weightedAvgSurRateUnderLeader(agent BaseAgent, id commons.ID) float64 {
	return 0.0
}

// Experience of agent [0,1]
func expertise(agent BaseAgent) float64 {
	return weightedOverthrowPercentage(agent, -0.4) + weightedManifestoEffectiveness(agent, 1.)
}

func similarityGeneralScore(agent BaseAgent) float64 {
	return 0.0
}

func similarityTargetedScore(agent BaseAgent) float64 {
	return 0.0
}

func lootAllocationScore(agent BaseAgent, leader bool) float64 {
	if leader {
		return similarityGeneralScore(agent)
	} else {
		return similarityGeneralScore(agent) + similarityTargetedScore(agent)
	}
}

func dynamicDonation(agent BaseAgent) uint {
	return 0
}

func getAliveAgents(agent BaseAgent) uint {
	view := agent.View()
	agentState := view.AgentState()
	i := 0
	itr := agentState.Iterator()
	for !itr.Done() {
		_, a, ok := itr.Next()
		if ok && a.Hp > 0 {
			i++
		}
	}
	return uint(i)
}

func getAgentStateSize(agent BaseAgent) uint {
	view := agent.View()
	agentState := view.AgentState()
	return uint(agentState.Len())
}

func elasticity(agent BaseAgent, w1 uint, w2 uint) uint {
	view := agent.View()
	agentState := agent.AgentState()
	return w1*(getAliveAgents(agent)/getAgentStateSize(agent)) + view.HpPool()/view.CurrentLevel() + w2*agentState.Hp
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

// Logistic function
func logistic(x float64, k float64, x0 float64) float64 {
	return 1 / (1 + math.Exp(-k*(x-x0)))
}

/* ---- STRATEGY ---- */

// HandleUpdateWeapon return the index of the weapon you want to use in AgentState.Weapons
func (a *Agent2) HandleUpdateWeapon(agent BaseAgent) decision.ItemIdx {
	weaponsInventory := agent.AgentState().Weapons
	if weaponsInventory.Len() != 0 {
		desiredWeaponIndex := decision.ItemIdx(weaponsInventory.Len() - 1)
		for i := 0; i < weaponsInventory.Len(); i++ {
			if weaponsInventory.Get(i).Value() < agent.AgentState().Stamina {
				desiredWeaponIndex = decision.ItemIdx(i)
				break
			}
		}
		for i := 0; i < weaponsInventory.Len(); i++ {
			if weaponsInventory.Get(i).Value() < agent.AgentState().Stamina {
				if weaponsInventory.Get(int(desiredWeaponIndex)).Value() < weaponsInventory.Get(i).Value() {
					desiredWeaponIndex = decision.ItemIdx(i)
				} else {
					continue
				}
			}
		}
		return desiredWeaponIndex
	} else {
		return decision.ItemIdx(0)
	}
}

// HandleUpdateShield return the index of the shield you want to use in AgentState.Shields
func (a *Agent2) HandleUpdateShield(agent BaseAgent) decision.ItemIdx {

	shieldInventory := agent.AgentState().Shields
	if shieldInventory.Len() != 0 {
		desiredShieldIndex := decision.ItemIdx(shieldInventory.Len() - 1)

		for i := 0; i < shieldInventory.Len(); i++ {
			if shieldInventory.Get(i).Value() < agent.AgentState().Stamina {
				desiredShieldIndex = decision.ItemIdx(i)
				break
			}
		}
		for i := 0; i < shieldInventory.Len(); i++ {
			if shieldInventory.Get(i).Value() < agent.AgentState().Stamina {
				if shieldInventory.Get(int(desiredShieldIndex)).Value() < shieldInventory.Get(i).Value() {
					desiredShieldIndex = decision.ItemIdx(i)
				} else {
					continue
				}
			}
		}
		return desiredShieldIndex
	} else {
		return decision.ItemIdx(0)
	}
}

// UpdateInternalState TODO: Implement me!
// Description: the function is called at the end of each level (provides a list of type FightResult / can be thought as raw & processed overall game info)
func (a *Agent2) UpdateInternalState(baseAgent BaseAgent, fightResult *commons.ImmutableList[decision.ImmutableFightResult], voteResult *immutable.Map[decision.Intent, uint], logChan chan<- logging.AgentLog) {
	a.updateBaseAgentPerLevel(baseAgent)
	a.updateFightResultPerLevel(*fightResult)
	a.updateVoteResultPerLevel(*voteResult)
	a.avgHp, a.avgDefend, a.avgAttack, a.avgStamina = updateAverages(baseAgent)
	a.updateSocialCapital(baseAgent)
	a.newGovernmentTimeline(baseAgent, a.haveElections)
}

/* ---- ELECTION ---- */

// CreateManifesto
// Description: Used to give Manifesto Information if elected Leader.
// Return:		The Manifesto with FightImposition, LootImposition, term length and overthrow threshold.
func (a *Agent2) CreateManifesto(agent BaseAgent) *decision.Manifesto {

	fightThreshold := 2.5
	lootThreshold := 2.5
	fightDecisionPower := false // default value

	if !wasOverthrown(agent) {
		if (adjustedExpertise(agent, 0, 5) + lastFightDecisionPower(agent, 2.5)) > fightThreshold {
			fightDecisionPower = true
		}
	} else {
		if adjustedExpertise(agent, 0, 5) > fightThreshold {
			fightDecisionPower = true
		}
	}

	lootDecisionPower := false

	if !wasOverthrown(agent) {
		if (adjustedExpertise(agent, 0, 5) + lastLootDecisionPower(agent, 2.5)) > lootThreshold {
			lootDecisionPower = true
		}
	} else {
		if adjustedExpertise(agent, 0, 5) > lootThreshold {
			lootDecisionPower = true
		}
	}

	termLength := uint(adjustedExpertise(agent, 0, 4) + 1)

	overthrowPercentage := uint(51)
	if wasOverthrown(agent) {
		overthrowPercentage = uint(float64(overthrowPercentage) + adjustedExpertise(agent, -10, 10))
	}

	manifesto := decision.NewManifesto(fightDecisionPower, lootDecisionPower, termLength, overthrowPercentage)
	return manifesto
}

// HandleConfidencePoll
// Description: Used for voting on confidence for Leader.
// Return:		Positive, Negative, or Abstain decision.
func (a *Agent2) HandleConfidencePoll(baseAgent BaseAgent) decision.Intent {
	w0, w1, w2, w3, w4, w5 := 1.0, 1.0, 1.0, 1.0, 1.0, 1.0
	avgSurvivalCurrTermNorm := (a.avgSurvivalCurrTerm - a.avgSurvival) / a.avgSurvival
	avgSurvivalPastTermsNorm := (a.avgSurvivalPastTerms - a.avgSurvival) / a.avgSurvival
	avgBroadcastRateCurrTermNorm := (a.avgBroadcastRateCurrTerm - a.avgBroadcastRate) / a.avgBroadcastRate
	avgBroadcastRatePastTermNorm := (a.avgBroadcastRatePastTerms - a.avgBroadcastRate) / a.avgBroadcastRate
	leadershipXpNorm := (a.leadershipXp - a.avgLeadershipXp) / a.avgLeadershipXp
	noConfRateNorm := (a.noConfRate - a.avgNoConfRate) / a.avgNoConfRate
	sum := w0*avgSurvivalCurrTermNorm + w1*avgSurvivalPastTermsNorm + w2*avgBroadcastRateCurrTermNorm + w3*avgBroadcastRatePastTermNorm + w4*leadershipXpNorm + w5*noConfRateNorm

	if sum >= 0 {
		return decision.Positive
	} else {
		return decision.Negative
	}
}

// HandleElectionBallot
// Description: Used to elect a Leader.
// Return:  	A single Commons.ID for choose-one voting or an array of commons.ID of top leader choices for ranked-voting.
func (a *Agent2) HandleElectionBallot(baseAgent BaseAgent, params *decision.ElectionParams) decision.Ballot {
	// Extract ID of alive agents
	view := baseAgent.View()
	agentState := view.AgentState()
	// Updating Leader Parameters
	a.termEndLevel = view.CurrentLevel()                    //level_temp is another priv attribute initialized to 0/1?
	a.termDuration = view.CurrentLevel() - a.termBeginLevel // term_begin_level was last updated at the beginning of the term that is now ending
	a.termBeginLevel = view.CurrentLevel()                  // now that it has been used, can reset to track new leadership
	a.avgTermSurvival = a.avgSurvivalCurrTerm               //(assert len(a.survival_rates)==term_duration )
	a.avgTermBroadcastRate = a.avgBroadcastRateCurrTerm     //(assert len(a.broadcast_rates)==term_duration )
	a.haveElections = true
	lastLeaderInfo := a.governmentTimeline[len(a.governmentTimeline)-1]
	if lastLeaderInfo.duration < lastLeaderInfo.manifesto.TermLength() {
		lastLeaderInfo.overthrown = true
	}
	a.governmentTimeline[len(a.governmentTimeline)-1] = lastLeaderInfo
	aliveAgentIds := make([]commons.ID, agentState.Len())
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
	agentScores := make(map[commons.ID]float64, numAliveAgents)
	for i := 0; i < numAliveAgents; i++ {
		par1 := leaderElectedBefore(baseAgent, weightedFracTermsDeposed(baseAgent, aliveAgentIds[i])+weightedAvgSurRateUnderLeader(baseAgent, aliveAgentIds[i]))
		par2 := lastFightDecisionPower(baseAgent, 5)
		par3 := lastLootDecisionPower(baseAgent, 5)
		agentScores[aliveAgentIds[i]] = prospectLeaderScore(baseAgent, par1, par2, par3)

	}
	sort.SliceStable(aliveAgentIds, func(i, j int) bool {
		return agentScores[aliveAgentIds[i]] < agentScores[aliveAgentIds[j]]
	})
	for i := uint(0); i < params.NumberOfPreferences(); i++ {
		ballot = append(ballot, aliveAgentIds[i])
	}
	return ballot
}

/* ---- FIGHT ---- */

// HandleFightInformation
// Description: Called every time a fight information message is received (I believe it could be from a leader for providing a proposal or another agent for providing fight info (e.g proposal directly to them?)
// Return:		nil
func (a *Agent2) HandleFightInformation(m message.TaggedInformMessage[message.FightInform], baseAgent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) {
	// baseAgent.Log(logging.Trace, logging.LogField{"bravery": r.bravery, "hp": baseAgent.AgentState().Hp}, "Cowering")
	rules := make([]proposal.Rule[decision.FightAction], 0)

	rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Cower,
		proposal.NewComparativeCondition(proposal.Health, proposal.LessThan, minHealth(baseAgent)),
	))

	rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Cower,
		proposal.NewAndCondition(*proposal.NewComparativeCondition(proposal.Health, proposal.GreaterThan, minHealth(baseAgent)),
			*proposal.NewComparativeCondition(proposal.Stamina, proposal.LessThan, minStamina(baseAgent))),
	))

	rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Attack,
		proposal.NewAndCondition(*proposal.NewComparativeCondition(proposal.Health, proposal.GreaterThan, baseHealth(baseAgent)),
			*proposal.NewComparativeCondition(proposal.TotalAttack, proposal.GreaterThan, minAttack(baseAgent))),
	))

	rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Defend,
		proposal.NewAndCondition(*proposal.NewComparativeCondition(proposal.Health, proposal.GreaterThan, baseHealth(baseAgent)),
			*proposal.NewComparativeCondition(proposal.TotalDefence, proposal.GreaterThan, minDefend(baseAgent))),
	))

	prop := *commons.NewImmutableList(rules)
	_ = baseAgent.SendFightProposalToLeader(prop)
}

// HandleFightRequest
// Description: Called every time a fight request message is received
// Return		Message Payload
func (a *Agent2) HandleFightRequest(m message.TaggedRequestMessage[message.FightRequest], log *immutable.Map[commons.ID, decision.FightAction]) message.FightInform {
	return nil
}

// TODO: Implement me!
func (a *Agent2) FightResolution(agent BaseAgent, prop commons.ImmutableList[proposal.Rule[decision.FightAction]],
	proposedActions immutable.Map[commons.ID, decision.FightAction],
) immutable.Map[commons.ID, decision.FightAction] {
	view := agent.View()
	builder := immutable.NewMapBuilder[commons.ID, decision.FightAction](nil)
	for _, id := range commons.ImmutableMapKeys(view.AgentState()) {
		var fightAction decision.FightAction
		switch rand.Intn(3) {
		case 0:
			fightAction = decision.Attack
		case 1:
			fightAction = decision.Defend
		default:
			fightAction = decision.Cower
		}
		builder.Set(id, fightAction)
	}
	return *builder.Map()
}

// HandleFightProposal: TODO: Implement me!
// Description: Throught that function our agent votes on a broadcastes proposal
func (a *Agent2) HandleFightProposal(proposal message.Proposal[decision.FightAction], baseAgent BaseAgent) decision.Intent {
	intent := rand.Intn(2)
	if intent == 0 {
		return decision.Positive
	} else {
		return decision.Negative
	}
}

// HandleFightProposalRequest TODO: Implement me!
// Description: Only called as a leader: True for broadcasting the proposal / False for declining the proposal
// Return:		Bool: True/False
func (a *Agent2) HandleFightProposalRequest(prop message.Proposal[decision.FightAction], baseAgent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) bool {
	propRules := prop.Rules()
	itr := propRules.Iterator()
	for !itr.Done() {
		a, ok := itr.Next()
		if ok {
			if reflect.TypeOf(a.Condition()) == reflect.TypeOf(proposal.NewComparativeCondition) {
				a.Condition()
			}
		}
	}

	switch rand.Intn(2) {
	case 0:
		return true
	default:
		return false
	}
}

func (a *Agent2) FightActionNoProposal(agent BaseAgent) decision.FightAction {
	// If not enough Stamina, no choice
	attack := agent.AgentState().Attack
	defense := agent.AgentState().Defense
	stamina := agent.AgentState().Stamina
	bonusAttack := a.getBonusAttack()
	bonusDefense := a.getBonusDefense()
	if stamina < attack+bonusAttack && stamina < defense+bonusDefense {
		return decision.Cower
	}

	currentDecision := a.initialDecision(agent)

	if currentDecision == decision.Cower {
		currentDecision = a.replaceDecision(agent, 10) // Second argument is the number of previous rounds to consider
	}
	/* Removed because damage per round is not relevant
	if currentDecision == decision.Cower {
		currentDecision = a.estimateDecision(baseAgent)
	}
	*/
	return currentDecision
}

// Description : Compare defense and attack potential, output a decision
// Return:		Cower, Defend or Attack decision.
func (a *Agent2) initialDecision(baseAgent BaseAgent) decision.FightAction {

	// method to retrieve state ?

	attack := baseAgent.AgentState().Attack
	defense := baseAgent.AgentState().Defense
	health := baseAgent.AgentState().Hp
	stamina := baseAgent.AgentState().Stamina
	bonusAttack := a.getBonusAttack()
	bonusDefense := a.getBonusDefense()

	// Bravery is a function of health and stamina
	bravery := 0.5*logistic(float64(health), 0.01, 500) + 0.5*logistic(float64(stamina), 0.005, 1000)

	// If current bravery is higher than parameter tendency, do something
	if bravery >= a.personalTendency {
		if stamina >= attack+bonusAttack {
			return decision.Attack
		} else if stamina >= defense+bonusDefense {
			return decision.Defend
		}
	}

	return decision.Cower // Else cower
}

func countDecisionInMaps(value decision.FightAction, decisionMaps []immutable.Map[commons.ID, decision.FightAction]) int {

	N := 0
	for _, m := range decisionMaps {
		itr := m.Iterator()
		for !itr.Done() {
			_, v, _ := itr.Next()
			if v == value {
				N += 1
			}
		}
	}
	return N
}

// Description : Compare current number of cowering agents to previous numbers
// and possibly replace them
// Return:		Cower, Defend or Attack decision.
func (a *Agent2) replaceDecision(baseAgent BaseAgent, N int) decision.FightAction {
	decisionHistory := a.getDecisionHelper(true)

	meanNLastRounds := countDecisionInMaps(decision.Attack, decisionHistory[len(decisionHistory)-N:]) / N
	lastFighting := countDecisionInMaps(decision.Attack, decisionHistory[len(decisionHistory)-1:])

	if lastFighting < meanNLastRounds {

		attack := baseAgent.AgentState().Attack
		defense := baseAgent.AgentState().Defense
		health := baseAgent.AgentState().Hp
		stamina := baseAgent.AgentState().Stamina
		bonusAttack := a.getBonusAttack()
		bonusDefense := a.getBonusDefense()

		// Bravery is a function of health and stamina
		bravery := 0.5*logistic(float64(health), 0.01, 500) + 0.5*logistic(float64(stamina), 0.005, 1000)

		if bravery >= a.replacementTendency {
			if stamina >= attack+bonusAttack {
				return decision.Attack
			} else if stamina >= defense+bonusDefense {
				return decision.Defend
			}
		}
	}
	return decision.Cower
}

// Description : Estimate current damage and possibly change decision
// Return:		Cower, Defend or Attack decision.
/*
func (a *Agent2) estimateDecision(baseAgent BaseAgent, N int) decision.FightAction {

	meanNLastTotalAttack := 0
	meanNLastTotalDefense := 0
	lastN := a.fightResultPerLevel[len(a.fightResultPerLevel)-N:]
	itr := lastN.Iterator()
	for !itr.Done() {
		_, result, _ := itr.Next()
		meanNLastTotalAttack += result.attackSum / N
		meanNLastTotalDefense += result.shieldSum / N
	}

	lastTotalAttack := a.fightResultPerLevel[len(a.fightResultPerLevel)-1].attackSum
	lastTotalDefense := a.fightResultPerLevel[len(a.fightResultPerLevel)-1].shieldSum

	diffAttack := lastTotalAttack - meanNLastTotalAttack
	diffDefense := lastTotalDefense - meanNLastTotalDefense

	// If there is less attack or defense
	if diffAttack < 0 || diffDefense < 0 {
		attack := baseAgent.AgentState().Attack
		defense := baseAgent.AgentState().Defense
		health := baseAgent.AgentState().Hp
		stamina := baseAgent.AgentState().Stamina
		bonusAttack := a.getBonusAttack()
		bonusDefense := a.getBonusDefense()

		// Bravery is a function of health and stamina
		bravery := 0.5*logistic(float64(health), 0.01, 500) + 0.5*logistic(float64(stamina), 0.005, 1000)

		if bravery >= a.estimationTendency {
			if stamina >= attack+bonusAttack {
				return decision.Attack
			} else if stamina >= defense+bonusDefense {
				return decision.Defend
			}
		}
	}
	return decision.Cower
}
*/

// FightAction
// Description: Logic of Fighting Action Decision-Making.
// Return:		Cower, Defend or Attack decision.
func (a *Agent2) FightAction(baseAgent BaseAgent, proposedAction decision.FightAction, acceptedProposal message.Proposal[decision.FightAction]) decision.FightAction {
	// If not enough Stamina, no choice
	attack := baseAgent.AgentState().Attack
	defense := baseAgent.AgentState().Defense
	stamina := baseAgent.AgentState().Stamina
	bonusAttack := a.getBonusAttack()
	bonusDefense := a.getBonusDefense()
	if stamina < attack+bonusAttack && stamina < defense+bonusDefense {
		return decision.Cower
	}

	currentDecision := a.initialDecision(baseAgent)

	if currentDecision == decision.Cower {
		currentDecision = a.replaceDecision(baseAgent, 10) // Second argument is the number of previous rounds to consider
	}
	/* Removed because damage per round is not relevant
	if currentDecision == decision.Cower {
		currentDecision = a.estimateDecision(baseAgent)
	}
	*/

	return currentDecision
}

/* ---- LOOT ---- */

func (a *Agent2) HandleLootInformation(m message.TaggedInformMessage[message.LootInform], agent BaseAgent) {
	rules := make([]proposal.Rule[decision.LootAction], 0)

	rules = append(rules, *proposal.NewRule[decision.LootAction](decision.HealthPotion,
		proposal.NewComparativeCondition(proposal.Health, proposal.LessThan, minHealth(agent)),
	))

	rules = append(rules, *proposal.NewRule[decision.LootAction](decision.StaminaPotion,
		proposal.NewAndCondition(*proposal.NewComparativeCondition(proposal.Health, proposal.GreaterThan, minHealth(agent)),
			*proposal.NewComparativeCondition(proposal.Stamina, proposal.LessThan, minStamina(agent))),
	))

	rules = append(rules, *proposal.NewRule[decision.LootAction](decision.Weapon,
		proposal.NewAndCondition(*proposal.NewComparativeCondition(proposal.Health, proposal.GreaterThan, baseHealth(agent)),
			*proposal.NewComparativeCondition(proposal.TotalAttack, proposal.LessThan, minAttack(agent))),
	))

	rules = append(rules, *proposal.NewRule[decision.LootAction](decision.Shield,
		proposal.NewAndCondition(*proposal.NewComparativeCondition(proposal.Health, proposal.GreaterThan, baseHealth(agent)),
			*proposal.NewComparativeCondition(proposal.TotalDefence, proposal.LessThan, minDefend(agent))),
	))

	prop := *commons.NewImmutableList(rules)
	_ = agent.SendLootProposalToLeader(prop)
}

// TODO: Implement me!
func (a *Agent2) HandleLootRequest(m message.TaggedRequestMessage[message.LootRequest]) message.LootInform {
	return nil
}

// TODO: Implement me!
func (a *Agent2) HandleLootProposal(r message.Proposal[decision.LootAction], agent BaseAgent) decision.Intent {
	switch rand.Intn(3) {
	case 0:
		return decision.Positive
	case 1:
		return decision.Negative
	default:
		return decision.Abstain
	}
}

// TODO: Implement me!
func (a *Agent2) HandleLootProposalRequest(proposal message.Proposal[decision.LootAction], agent BaseAgent) bool {
	switch rand.Intn(2) {
	case 0:
		return true
	default:
		return false
	}
}

func (a *Agent2) LootAllocation(
	baseAgent BaseAgent,
	proposal message.Proposal[decision.LootAction],
	proposedAllocation immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]],
) immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]] {

	lootAllocation := make(map[commons.ID][]commons.ItemID)
	view := baseAgent.View()
	agentState := view.AgentState()
	healthIDs := make(map[commons.ID]uint, agentState.Len())
	staminaIDs := make(map[commons.ID]uint, agentState.Len())
	weapons := make(map[commons.ID]uint, agentState.Len())
	shields := make(map[commons.ID]uint, agentState.Len())
	itr := agentState.Iterator()
	for !itr.Done() {
		id, a, ok := itr.Next()
		if ok && a.Hp > 0 {
			healthIDs[id] = uint(a.Hp)
			staminaIDs[id] = uint(a.Stamina)
			weapons[id] = a.BonusAttack
			shields[id] = a.BonusDefense
		}
	}

	ids := commons.ImmutableMapKeys(view.AgentState())
	iterator := baseAgent.Loot().Weapons().Iterator()
	allocateEgaliterian(iterator, ids, weapons, lootAllocation)
	iterator = baseAgent.Loot().Shields().Iterator()
	allocateEgaliterian(iterator, ids, shields, lootAllocation)
	iterator = baseAgent.Loot().HpPotions().Iterator()
	allocateEgaliterian(iterator, ids, healthIDs, lootAllocation)
	iterator = baseAgent.Loot().StaminaPotions().Iterator()
	allocateEgaliterian(iterator, ids, staminaIDs, lootAllocation)
	mMapped := make(map[commons.ID]immutable.SortedMap[commons.ItemID, struct{}])
	for id, itemIDS := range lootAllocation {
		mMapped[id] = commons.ListToImmutableSortedSet(itemIDS)
	}
	return commons.MapToImmutable(mMapped)
}

func allocateEgaliterian(iterator commons.Iterator[state.Item], ids []commons.ID, values map[commons.ID]uint, lootAllocation map[commons.ID][]commons.ItemID) {
	sort.SliceStable(ids, func(i, j int) bool {
		return values[ids[i]] < values[ids[j]]
	})
	i := 0
	for !iterator.Done() {
		next, _ := iterator.Next()
		toBeAllocated := ids[i]
		if l, ok := lootAllocation[toBeAllocated]; ok {
			l = append(l, next.Id())
			lootAllocation[toBeAllocated] = l
		} else {
			l := make([]commons.ItemID, 0)
			l = append(l, next.Id())
			lootAllocation[toBeAllocated] = l
		}
	}
}

func (a *Agent2) LootAction(
	_ BaseAgent,
	proposedLoot immutable.SortedMap[commons.ItemID, struct{}],
	_ message.Proposal[decision.LootAction],
) immutable.SortedMap[commons.ItemID, struct{}] {
	return proposedLoot
}

func (a *Agent2) LootActionNoProposal(baseAgent BaseAgent) immutable.SortedMap[commons.ItemID, struct{}] {
	loot := baseAgent.Loot()
	weapons := loot.Weapons().Iterator()
	shields := loot.Shields().Iterator()
	hpPotions := loot.HpPotions().Iterator()
	staminaPotions := loot.StaminaPotions().Iterator()

	builder := immutable.NewSortedMapBuilder[commons.ItemID, struct{}](nil)

	for !weapons.Done() {
		weapon, _ := weapons.Next()
		if rand.Int()%2 == 0 {
			builder.Set(weapon.Id(), struct{}{})
		}
	}

	for !shields.Done() {
		shield, _ := shields.Next()
		if rand.Int()%2 == 0 {
			builder.Set(shield.Id(), struct{}{})
		}
	}

	for !hpPotions.Done() {
		pot, _ := hpPotions.Next()
		if rand.Int()%2 == 0 {
			builder.Set(pot.Id(), struct{}{})
		}
	}

	for !staminaPotions.Done() {
		pot, _ := staminaPotions.Next()
		if rand.Int()%2 == 0 {
			builder.Set(pot.Id(), struct{}{})
		}
	}

	return *builder.Map()
}

/* ---- HPPOOL ---- */

// DonateToHpPool
// Description: The function returns the amount of Hp that our agent is willing to donate to the HpPool
func (a *Agent2) DonateToHpPool(agent BaseAgent) uint {
	agentState := agent.AgentState()
	if agentState.Hp < minHealth(agent) {
		return 0
	} else {
		return 10 + dynamicDonation(agent)
	}
}

/* ---- TRADE ----- */
// TODO: Implement me!
func (a *Agent2) HandleTradeNegotiation(_ BaseAgent, _ message.TradeInfo) message.TradeMessage {
	return message.TradeRequest{}
}
