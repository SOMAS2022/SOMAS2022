package agent

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/message/proposal"
	"infra/game/state"
	_ "infra/game/state"
	"math"
	"math/rand"

	"github.com/benbjohnson/immutable"
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
}

// NewAgent2 : Constructor of Agent2 FIXME:!!!INITIALISE AGENT WITH (MEANINGFUL) VALUES!!!
func NewAgent2() Strategy {
	// testing : random between [0.25,0.75]
	personalTendency := rand.Float64()*0.25 + 0.5
	replacementTendency := rand.Float64()*0.25 + 0.5
	estimationTendency := rand.Float64()*0.25 + 0.5
	return &Agent2{personalTendency: personalTendency, replacementTendency: replacementTendency, estimationTendency: estimationTendency}
}

/* ---- HELPER FUNCTIONS ----*/
// Returns minimum Health that a healthy agent should have.
// Returns minimum Health that a healthy agent should have.
// Returns minimum Health that a healthy agent should have.
// Returns minimum Health that a healthy agent should have.
func minDefend() proposal.Value {
	return 1000
}

func minAttack() proposal.Value {
	return 1000
}

func baseHealth() proposal.Value {
	return 1000
}

func minStamina() proposal.Value {
	return 1000
}

// Returns minimum Health that a healthy agent should have.
func minHealth() proposal.Value {
	return 1000
}

// Returns Manifesto Effectiveness based on History
func manifestoEffectPercentage(agent BaseAgent) float64 {
	return 0.0
}

// Returns Number of Full Term Agent2 served
// without being overthrown
func fullTermPercentage(agent BaseAgent) float64 {
	return 0.0
}

// Returns Number of Terms Agent2 served
// and was overthrown
func overthrowPercentage(agent BaseAgent) float64 {
	return 0.0
}

// Returns T|F if agent was overthrown or not
func wasOverthrown(agent BaseAgent) bool {
	return false
}

// Returns T|F if previous agent manifesto imposed loot
func prevLootImpAbility(agent BaseAgent) bool {
	return false
}

// Returns T|F if previous agent manifesto imposed fight
func prevFightImpAbility(agent BaseAgent) bool {
	return false
}

// First time agent is leader
func firstTime(agent BaseAgent) bool {
	return true
}

// Experience of agent [0,1]
func experience(agent BaseAgent) float64 {
	overthrowWeight := 0.4
	fullTermWeight := 0.7
	manifestoEffectWeight := 0.7
	return -overthrowWeight*overthrowPercentage(agent) + fullTermWeight*fullTermPercentage(agent) + manifestoEffectWeight*manifestoEffectPercentage(agent)
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
func (a *Agent2) UpdateInternalState(baseAgent BaseAgent, fightResult *commons.ImmutableList[decision.ImmutableFightResult], voteResult *immutable.Map[decision.Intent, uint]) {
	a.updateBaseAgentPerLevel(baseAgent)
	a.updateFightResultPerLevel(*fightResult)
	a.updateVoteResultPerLevel(*voteResult)
}

/* ---- ELECTION ---- */

// CreateManifesto FIXME: Check me!
// Description: Used to give Manifesto Information if elected Leader.
// Return:		The Manifesto with FightImposition, LootImposition, term length and overthrow threshold.
func (a *Agent2) CreateManifesto(agent BaseAgent) *decision.Manifesto {
	/*
		CreateManifesto:
		- Term_Length =  Int(factor * experience) + Bias(Default term e.g. 1 term)
		- Overthrow(%) = first_time ? default -> 51% : experience*(mapping factor)* (-10,10)
		- FightImpositionDecision = experience + (prev_fight_imp_ability * !was_overthrown ? bias : 0 ) > threshold ? True : False
		- LootImposition = experience + (prev_loot_imp_ability * !was_overthrown ? bias : 0 ) > threshold ? True : False
	*/
	f1 := 4 * experience(agent)           // [0,4]
	f2 := 20.00*experience(agent) - 10.00 // [-10,10]
	defaultTerm := uint(1)                // 1 Term
	defaultOverthrow := uint(51)          // 51% Agents to overthrow
	overthrowThreshold := uint(0)
	termLength := uint(f1) + defaultTerm // [1,5]
	fightImposeThreshold := 0.0
	lootImposeThreshold := 0.0
	fightImposition := false
	lootImposition := false

	if firstTime(agent) == true {
		overthrowThreshold = defaultOverthrow
	} else {
		overthrowThreshold = uint(experience(agent) * f2)
	}
	if (prevFightImpAbility(agent) && !wasOverthrown(agent)) == true {
		fightImposeHistoryCheck := 0.5
		if (experience(agent) + fightImposeHistoryCheck) > fightImposeThreshold {
			fightImposition = true
		} else {
			fightImposition = false
		}
	} else {
		fightImposeHistoryCheck := 0.0
		if (experience(agent) + fightImposeHistoryCheck) > fightImposeThreshold {
			fightImposition = true
		} else {
			fightImposition = false
		}
	}
	if (prevLootImpAbility(agent) && !wasOverthrown(agent)) == true {
		lootImposeHistoryCheck := 0.5
		if (experience(agent) + lootImposeHistoryCheck) > lootImposeThreshold {
			lootImposition = true
		} else {
			lootImposition = false
		}
	} else {
		lootImposeHistoryCheck := 0.0
		if (experience(agent) + lootImposeHistoryCheck) > lootImposeThreshold {
			lootImposition = true
		} else {
			lootImposition = false
		}
	}
	Manifesto := decision.NewManifesto(fightImposition, lootImposition, termLength, overthrowThreshold)
	return Manifesto
}

// HandleConfidencePoll TODO: Implement me!
// Description: Used for voting on confidence for Leader.
// Return:		Positive, Negative, or Abstain decision.
func (a *Agent2) HandleConfidencePoll(baseAgent BaseAgent) decision.Intent {
	// To decide how to vote in no-confidence vote at the end of each level, use a social capital framework with weighted factors and a binary activation function to decide yes/no
	// These are:
	// - avg_survival_curr_term: average % of agents alive at the end of a level during current leadership term (+ve relationship, high weighting)
	// - avg_survival_past_terms: average % of agents alive at the end of a level from past leadership terms of that agent (+ve)
	// - avg_survival: average % of agents alive at the end of a level from all past leadership terms (for comparison - namely normalize by this amount)
	// - avg_broadcast_rate_curr_term: % of the proposals we submitted that were actually accepted/broadcast by the leader, in current term (+ve, high weighting)
	// - avg_broadcast_rate_past_terms: % of the proposals we submitted that were actually accepted/broadcast by the leader, from past terms of that leader (+ve)
	// - avg_broadcast_rate: % of the proposals we submitted that were actually accepted/broadcast, from all past leadership terms (again, normalize by this)
	// - leadership_xp: fraction of levels up to now that they were leader (+ve)
	// - no_conf_rate: fraction of their terms they were voted out prematurely (-ve))
	// - avg_leadership_xp: avg fraction of levels up to now that any one agent is leader
	// - avg_no_conf_rate: avg fraction of an agent's leaderships terms that he is voted out
	// - (fight imposition?) (-ve)
	// - (loot?)
	// These variables are marked with -- below
	// For these, firstly we need a history data helper function that returns an array of the form:
	// leader_timeline_array [{id, manifesto, duration, leader_stats}, {id, manifesto, duration, leader_stats}, ...]
	// The object of type leader_stats will contain the following items, corresponding JUST to that elapsed leadership term:
	// - avg_term_survival (calculate for each level of their leadership and average)
	// - avg_term_broadcast_rate (calculate for each round/level? of their leadership and average)
	// - bool no_conf: whether they were voted out of that term
	// - (fight/loot impositions?)
	// This array is best created in the election function that is only called at the end of one leadership term / start of another
	// It's best to have private attributes that accrue raw data and then reset - some every new term, some every new level
	// These are used by the confidence function at the end of every level to actually yield the no conf vote, and by the election function at the end of a term to calculate stats and append to leader_timeline_array, and to vote
	// Namely, the ones we reset after every level:
	// - num_agents_begin_level (actually, do we only have list of agent IDs?)
	// - num_agents_end_level
	// - proposals_total: how many proposals we put forward that level (necessarily equal to rounds?)
	// - proposals_broadcast: how many of these were broadcast
	// - loot/trade info?
	// And variables we re-calculate every level, but reset every election (no need for arrays for the raw data from which we calculate them):
	// - survival_rates: array of % of agents alive at the end of each level (this array is appended to at the end of every level, and resets every election, so that each elem corresponds to a level in a leadership term)
	// - broadcast_rates: % of the proposals we submitted during the level that were actually accepted/broadcast (ditto)
	// -- avg_survival_curr_term: avg of survival_rates (updated at end of every level - used as main measure of confidence in the current leader)
	// -- avg_broadcast_rate_curr_term: avg of broadcast_rates
	// And variables we re-calculate every level but never reset:
	// -- avg_survival: average % of agents alive at the end of a level from all levels
	// -- avg_broadcast_rate: % of the proposals we submitted that were actually accepted/broadcast, from all levels
	// Then the ones we reset every election:
	// - term_begin_level: level at which the leadership term began (can read from viewMap every time election func is called)
	// - term_end_level: level at which the leadership term ended (again, from viewMap, in election func)
	// - no_conf (bool): whether or not the term ended bc of no-confidence vote (where access this from?)
	// And variables we calculate every election (using all previous vars), to add to leader_timeline_array (esp leader_stats):
	// - term_duration: number of levels that leadership term lasted before elapsed or deposed (term_end_level - term_begin_level)
	// - avg_term_survival: survival_rates averaged over term_duration (which is actually the length of the survival_rates array - assertion?)
	// - avg_term_broadcast_rate: broadcast_rates averaged over term_duration (also length of broadcast_rates array)
	// - bool no_conf (no calculation needed)
	// These statistics are the 'condensed', useful form of the raw data
	// This leaves us, from all the vars involved in the confidence vote, with:
	// -- avg_survival_past_terms: average % of agents alive at the end of a level from past leadership terms of that agent
	// -- avg_broadcast_rate_past_terms: % of the proposals we submitted that were actually accepted/broadcast by the leader, from past terms of that leader
	// -- leadership_xp: fraction of levels up to now that they were leader (+ve)
	// -- no_conf_rate: fraction of their terms they were voted out prematurely (-ve))
	// -- avg_leadership_xp: avg fraction of levels up to now that any one agent is leader
	// -- avg_no_conf_rate: avg fraction of an agent's leaderships terms that he is voted out
	// These will not be calculated as private attributes for every single agent, but rather by looping through leader_timeline_array on an ad-hoc basis

	// Pseudocode for how stats are calculated from raw data, elsewhere in the program:
	// In a function that runs at the end of a every level:
	// - a.survival_rates = append(a.survival_rates, num_agents_end_level/num_agents_begin_level)
	// - a.broadcast_rates = append(a.broadcast_rates, proposals_broadcast/proposals_total)
	// - a.avg_survival_curr_term = avg(survival_rates)
	// - a.avg_broadcast_rate_curr_term = avg(broadcast_rates)
	// - a.avg_survival = (avg_survival*prev_level + num_agents_end_level/num_agents_begin_level)/level // where can get level?
	// - a.avg_broadcast_rate = (avg_broadcast_rate*prev_level + proposals_broadcast/proposals_total)/level
	// In election function:
	// - a.term_end_level = level //level_temp is another priv attribute initialized to 0/1?
	// - a.term_duration = level - term_begin_level // term_begin_level was last updated at the beginning of the term that is now ending
	// - a.term_begin_level = level // now that it has been used, can reset to track new leadership
	// - a.avg_term_survival = a.avg_survival_curr_term //(assert len(a.survival_rates)==term_duration )
	// - a.avg_term_broadcast_rate = a.avg_broadcast_rate_curr_term //(assert len(a.broadcast_rates)==term_duration )
	// - a.no_conf = // how get result of confidence poll? Do we need to calculate using if (manifesto_term!=term_duration) ?
	// Then need to construct leader_timeline_array...

	//var past_terms_of_curr_leader := make([]term_struct, 0)
	//for leadership_term in leader_term_timeline_array {
	//	if leadership_term[id] == curr_leader["id"] {
	//		past_terms_of_curr_leader = append(past_terms_of_curr_leader, leadership_term) // will have redundant id key but whatever
	//	}
	//}

	// avg_survival_curr_term_norm := (avg_survival_curr_term-avg_survival)/avg_survival
	// avg_survival_past_terms_norm := (avg_survival_past_terms-avg_survival)/avg_survival
	// avg_broadcast_rate_curr_term_norm := (avg_broadcast_rate_curr_term-avg_broadcast_rate)/avg_broadcast_rate
	// avg_broadcast_rate_past_terms_norm := (avg_broadcast_rate_past_terms-avg_broadcast_rate)/avg_broadcast_rate
	// leadership_xp_norm := (leadership_xp-avg_leadership_xp)/avg_leadership_xp
	// no_conf_rate_norm := (no_conf_rate-avg_no_conf_rate)/avg_no_conf_rate
	// Hm maybe tweak vars so that all count as +ve contrib to confidence
	// sum = w0*avg_survival_curr_term_norm + w1*avg_survival_past_terms_norm + w2*avg_broadcast_rate_curr_term_norm + w3*avg_broadcast_rate_past_terms_norm + w4*leadership_xp_norm + w5*no_conf_rate_norm

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

/* ---- FIGHT ---- */

// HandleFightInformation TODO: Implement me!
// Description: Called every time a fight information message is received (I believe it could be from a leader for providing a proposal or another agent for providing fight info (e.g proposal directly to them?)
// Return:		nil
func (a *Agent2) HandleFightInformation(m message.TaggedInformMessage[message.FightInform], baseAgent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) {
	a.updateDecisionHelper(*log)

	makesProposal := rand.Intn(100)

	if makesProposal > 80 {
		prop := a.FightResolution(baseAgent)
		_ = baseAgent.SendFightProposalToLeader(prop)
	}
}

// HandleFightRequest TODO: Implement me!
// Description: Called every time a fight request message is received
// Return		Message Payload
func (a *Agent2) HandleFightRequest(m message.TaggedRequestMessage[message.FightRequest], log *immutable.Map[commons.ID, decision.FightAction]) message.FightInform {
	return nil
}

// FightResolution: TODO: Implement me!
// Description: Through that function our agent provides a proposal
func (a *Agent2) FightResolution(agent BaseAgent) commons.ImmutableList[proposal.Rule[decision.FightAction]] {
	rules := make([]proposal.Rule[decision.FightAction], 0)

	rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Cower,
		proposal.NewComparativeCondition(proposal.Health, proposal.LessThan, minHealth()),
	))

	rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Cower,
		proposal.NewAndCondition(*proposal.NewComparativeCondition(proposal.Health, proposal.GreaterThan, minHealth()),
			*proposal.NewComparativeCondition(proposal.Stamina, proposal.LessThan, minStamina())),
	))

	rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Attack,
		proposal.NewAndCondition(*proposal.NewComparativeCondition(proposal.Health, proposal.GreaterThan, baseHealth()),
			*proposal.NewComparativeCondition(proposal.TotalAttack, proposal.GreaterThan, minAttack())),
	))

	rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Defend,
		proposal.NewAndCondition(*proposal.NewComparativeCondition(proposal.Health, proposal.GreaterThan, baseHealth()),
			*proposal.NewComparativeCondition(proposal.TotalDefence, proposal.GreaterThan, minDefend())),
	))

	return *commons.NewImmutableList(rules)
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
func (a *Agent2) HandleFightProposalRequest(proposal message.Proposal[decision.FightAction], baseAgent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) bool {
	switch rand.Intn(2) {
	case 0:
		return true
	default:
		return false
	}
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

// FightAction TODO: Implement me!
// Description: Logic of Fighting Action Decision-Making.
// Return:		Cower, Defend or Attack decision.
func (a *Agent2) FightAction(baseAgent BaseAgent) decision.FightAction {
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
	/*
		makesProposal := rand.Intn(100)

		if makesProposal > 80 {
			prop := a.LootAllocation(agent)
			_ = agent.SendLootProposalToLeader(prop)
		}
	*/
}

func (a *Agent2) HandleLootRequest(m message.TaggedRequestMessage[message.LootRequest]) message.LootInform {
	//TODO implement me
	panic("implement me")
}

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

func (a *Agent2) HandleLootProposalRequest(proposal message.Proposal[decision.LootAction], agent BaseAgent) bool {
	switch rand.Intn(2) {
	case 0:
		return true
	default:
		return false
	}
}

func (a *Agent2) LootAllocation(agent BaseAgent) immutable.Map[commons.ID, immutable.List[commons.ItemID]] {
	lootAllocation := make(map[commons.ID][]commons.ItemID)
	view := agent.View()
	ids := commons.ImmutableMapKeys(view.AgentState())
	iterator := agent.Loot().Weapons().Iterator()
	allocateRandomly(iterator, ids, lootAllocation)
	iterator = agent.Loot().Shields().Iterator()
	allocateRandomly(iterator, ids, lootAllocation)
	iterator = agent.Loot().HpPotions().Iterator()
	allocateRandomly(iterator, ids, lootAllocation)
	iterator = agent.Loot().StaminaPotions().Iterator()
	allocateRandomly(iterator, ids, lootAllocation)
	mMapped := make(map[commons.ID]immutable.List[commons.ItemID])
	for id, itemIDS := range lootAllocation {
		mMapped[id] = commons.ListToImmutable(itemIDS)
	}
	return commons.MapToImmutable(mMapped)
}

func allocateRandomly(iterator commons.Iterator[state.Item], ids []commons.ID, lootAllocation map[commons.ID][]commons.ItemID) {
	for !iterator.Done() {
		next, _ := iterator.Next()
		toBeAllocated := ids[rand.Intn(len(ids))]
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

func (a *Agent2) LootAction() immutable.List[commons.ItemID] {
	return *immutable.NewList[commons.ItemID]()
}

/* ---- HPPOOL ---- */

// DonateToHpPool TODO: Implement me!
// Description: The function returns the amount of Hp that our agent is willing to donate to the HpPool
func (a *Agent2) DonateToHpPool(baseAgent BaseAgent) uint {
	return uint(rand.Intn(int(baseAgent.AgentState().Hp)))
}
