package team2

import (
	"infra/game/agent"
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
	baseAgentMap  []agent.BaseAgent
	viewMap       []state.View
	agentStateMap []immutable.Map[commons.ID, state.HiddenAgentState]
	leaderMap     []commons.ID

	// Updates per Level
	baseAgentPerLevel   []agent.BaseAgent
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

// NewAgent2 : Constructor of Agent2
func NewAgent2() agent.Strategy {
	// testing : random between [0.25,0.75]
	personalTendency := rand.Float64()*0.25 + 0.5
	replacementTendency := rand.Float64()*0.25 + 0.5
	estimationTendency := rand.Float64()*0.25 + 0.5
	return &Agent2{
		personalTendency:          personalTendency,
		replacementTendency:       replacementTendency,
		estimationTendency:        estimationTendency,
		decisionMap:               make([]immutable.Map[commons.ID, decision.FightAction], 100),
		baseAgentMap:              make([]agent.BaseAgent, 100),
		viewMap:                   make([]state.View, 100),
		agentStateMap:             make([]immutable.Map[commons.ID, state.HiddenAgentState], 100),
		leaderMap:                 make([]commons.ID, 100),
		baseAgentPerLevel:         make([]agent.BaseAgent, 100),
		fightResultPerLevel:       make([]commons.ImmutableList[decision.ImmutableFightResult], 100),
		voteResultPerLevel:        make([]immutable.Map[decision.Intent, uint], 100),
		governmentTimeline:        make([]leaderInfo, 100),
		haveElections:             false,
		avgHp:                     0,
		avgDefend:                 0,
		avgAttack:                 0,
		avgStamina:                0,
		avgSurvivalCurrTerm:       0.0,
		avgSurvivalPastTerms:      0.0,
		avgSurvival:               0.0,
		avgBroadcastRateCurrTerm:  0.0,
		avgBroadcastRatePastTerms: 0.0,
		avgBroadcastRate:          0.0,
		leadershipXp:              0.0,
		noConfRate:                0.0,
		avgLeadershipXp:           0.0,
		avgNoConfRate:             0.0,
		avgTermSurvival:           0.0,
		avgTermBroadcastRate:      0.0,
		noConf:                    false,
		numAgentsBeginLevel:       0,
		numAgentsEndLevel:         0,
		proposalsTotal:            0,
		proposalsBroadcast:        0,
		survivalRates:             make([]float64, 100),
		broadcastRates:            make([]float64, 100),
		termBeginLevel:            0,
		termEndLevel:              0,
		termDuration:              0,
	}
}

/* ---- HELPER FUNCTIONS ---- */
type leaderInfo struct {
	id         commons.ID
	manifesto  decision.Manifesto
	duration   uint
	overthrown bool
}

func (a *Agent2) newGovernmentTimeline(agent agent.BaseAgent, haveElections bool) {
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

/* ---- UPDATES FUNCTIONS ---- */

func (a *Agent2) updateSocialCapital(fightResult commons.ImmutableList[decision.ImmutableFightResult], agent agent.BaseAgent) {
	view := agent.View()
	agents := view.AgentState()
	// Pseudocode for how stats are calculated from raw data, elsewhere in the program:
	numOfAgents := make([]uint, agents.Len())
	itr := fightResult.Iterator()
	for !itr.Done() {
		a, ok := itr.Next()
		if ok {
			attacking := a.AttackingAgents()
			defending := a.ShieldingAgents()
			cowering := a.CoweringAgents()
			numOfAgents = append(numOfAgents, uint(attacking.Len()+defending.Len()+cowering.Len()))
		}
	}
	a.numAgentsBeginLevel = numOfAgents[0]
	a.numAgentsEndLevel = numOfAgents[len(numOfAgents)-1]
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

func updateAverages(agent agent.BaseAgent) (uint, uint, uint, uint) {
	sumHp, sumDefend, sumAttack, sumStamina := uint(0), uint(0), uint(0), uint(0)
	view := agent.View()
	agentState := view.AgentState()
	numOfAgents := uint(0)
	for _, id := range commons.ImmutableMapKeys(agentState) {
		stat, ok := agentState.Get(id)
		if ok && stat.Hp > 0 {
			sumHp = sumHp + uint(stat.Hp)
			sumDefend = sumDefend + stat.BonusDefense
			sumAttack = sumAttack + stat.BonusAttack
			sumStamina = sumStamina + uint(stat.Stamina)
			numOfAgents += 1
		}
	}
	avgHp := sumHp / numOfAgents
	avgDefend := sumDefend / numOfAgents
	avgAttack := sumAttack / numOfAgents
	avgStamina := sumStamina / numOfAgents

	return avgHp, avgDefend, avgAttack, avgStamina
}

func minDefend(agent agent.BaseAgent) proposal.Value {
	avgDefend := 0.0
	view := agent.View()
	agentState := view.AgentState()
	for _, id := range commons.ImmutableMapKeys(agentState) {
		stat, ok := agentState.Get(id)
		if ok {
			avgDefend = avgDefend + float64(stat.BonusDefense)
		}
	}
	return uint(0.50 * avgDefend)
}

func minAttack(agent agent.BaseAgent) proposal.Value {
	avgAttack := 0.0
	view := agent.View()
	agentState := view.AgentState()
	for _, id := range commons.ImmutableMapKeys(agentState) {
		stat, ok := agentState.Get(id)
		if ok {
			avgAttack = avgAttack + float64(stat.BonusAttack)
		}
	}
	return uint(0.50 * avgAttack)
}

func baseHealth(agent agent.BaseAgent) proposal.Value {
	avgHealth := 0.0
	view := agent.View()
	agentState := view.AgentState()
	for _, id := range commons.ImmutableMapKeys(agentState) {
		stat, ok := agentState.Get(id)
		if ok {
			avgHealth = avgHealth + float64(stat.Hp)
		}
	}
	return uint(avgHealth)
}

func minStamina(agent agent.BaseAgent) proposal.Value {
	avgStamina := 0.0
	view := agent.View()
	agentState := view.AgentState()
	for _, id := range commons.ImmutableMapKeys(agentState) {
		stat, ok := agentState.Get(id)
		if ok {
			avgStamina = avgStamina + float64(stat.Stamina)
		}
	}
	return uint(0.50 * avgStamina)
}

// Returns minimum Health that a healthy agent should have.
func minHealth(agent agent.BaseAgent) proposal.Value {
	avgHealth := 0.0
	view := agent.View()
	agentState := view.AgentState()
	for _, id := range commons.ImmutableMapKeys(agentState) {
		stat, ok := agentState.Get(id)
		if ok {
			avgHealth = avgHealth + float64(stat.Hp)
		}
	}
	return uint(0.50 * avgHealth)
}

// Returns Manifesto Effectiveness based on History
func (a *Agent2) weightedManifestoEffectiveness(agent agent.BaseAgent, weight float64) float64 {
	view := agent.View()
	manifesto := view.LeaderManifesto()
	term := float64(manifesto.TermLength())
	return 0.5 * (weightedOverthrowPercentage(agent, weight) + term/a.avgTermSurvival)
}

// Returns Number of Terms Agent2 served
// and was overthrown
func weightedOverthrowPercentage(agent agent.BaseAgent, weight float64) float64 {
	view := agent.View()
	manifesto := view.LeaderManifesto()
	return weight * float64(manifesto.OverthrowThreshold())
}

// Returns Adjusted Expertise with new mean and std provided.
func (a *Agent2) adjustedExpertise(agent agent.BaseAgent, from float64, to float64) float64 {
	return from + (to-from)*a.expertise(agent)
}

func (a *Agent2) wasOverthrown(candidateID commons.ID) bool {
	for elected := len(a.governmentTimeline) - 1; elected >= 0; elected-- {
		if (a.governmentTimeline[elected].id == candidateID) && (a.governmentTimeline[elected].manifesto.TermLength() > a.governmentTimeline[elected].duration) {
			return true
		}
	}
	return false
}

func (a *Agent2) lastFightDecisionPower(candidateID commons.ID, bias float64) float64 {
	for elected := len(a.governmentTimeline) - 1; elected >= 0; elected-- {
		if (a.governmentTimeline[elected].id == candidateID) && (a.governmentTimeline[elected].manifesto.FightDecisionPower() == true) {
			return bias
		}
	}
	return 0.0
}

func (a *Agent2) lastLootDecisionPower(candidateID commons.ID, bias float64) float64 {
	for elected := len(a.governmentTimeline) - 1; elected >= 0; elected-- {
		if (a.governmentTimeline[elected].id == candidateID) && (a.governmentTimeline[elected].manifesto.LootDecisionPower() == true) {
			return bias
		}
	}
	return 0.0
}

func (a *Agent2) leaderElectedBefore(candidateID commons.ID, bias float64) float64 {
	for elected := len(a.governmentTimeline) - 1; elected >= 0; elected-- {
		if a.governmentTimeline[elected].id == candidateID {
			return bias
		}
	}
	return 0.0
}

// SOT : [0,1] : map :(Overthrow_i_ranked+Term_i_ranked)_ranked
func (a *Agent2) SOT(candidateID commons.ID, manifesto decision.Manifesto, w1 float64) float64 {
	return w1*(1-float64(manifesto.OverthrowThreshold())/20) + float64(manifesto.TermLength())
}

// returning prospect leader's score
func prospectLeaderScore(par1 float64, par2 float64, par3 float64, sot float64) float64 {
	return par1 + par2 + par3 + sot
}

// weightedFracTermsDeposed : NumOfTimesDeposed / NumOfTimesElected
func (a *Agent2) weightedFracTermsDeposed(w1 float64, candidateID commons.ID) float64 {
	timesElected := 0
	overthrown := 0
	for elected := len(a.governmentTimeline) - 1; elected >= 0; elected-- {
		if a.governmentTimeline[elected].id == candidateID {
			timesElected++
			if a.governmentTimeline[elected].duration < a.governmentTimeline[elected].manifesto.TermLength() {
				overthrown++
			}
		}
	}
	return w1 * float64(overthrown) / float64(timesElected)
}

// weightedAvgSurRateUnderLeader : Sum_terms(Sum_levels(survival_rate))/(Sum_terms(Sum_levels(1)))
func (a *Agent2) weightedAvgSurRateUnderLeader(w1 float64) float64 {
	return w1 * a.avgSurvivalCurrTerm
}

// Experience of agent [0,1]
func (a *Agent2) expertise(agent agent.BaseAgent) float64 {
	return a.weightedFracTermsDeposed(-1.0, agent.ID()) + a.weightedManifestoEffectiveness(agent, 1.0)
}

func similarityGeneralScore(agent agent.BaseAgent) float64 {
	return 0.0
}

func similarityTargetedScore(agent agent.BaseAgent) float64 {
	return 0.0
}

func lootAllocationScore(agent agent.BaseAgent, leader bool) float64 {
	if leader {
		return similarityGeneralScore(agent)
	} else {
		return similarityGeneralScore(agent) + similarityTargetedScore(agent)
	}
}

func dynamicDonation(agent agent.BaseAgent) uint {
	return 0
}

func getAliveAgents(agent agent.BaseAgent) uint {
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

func getAgentStateSize(agent agent.BaseAgent) uint {
	view := agent.View()
	agentState := view.AgentState()
	return uint(agentState.Len())
}

func elasticity(agent agent.BaseAgent, w1 uint, w2 uint) uint {
	view := agent.View()
	agentState := agent.AgentState()
	return w1*(getAliveAgents(agent)/getAgentStateSize(agent)) + view.HpPool()/view.CurrentLevel() + w2*agentState.Hp
}

func (a *Agent2) updateBaseAgentPerLevel(agent agent.BaseAgent) {
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
		if len(a.decisionMap)-1 > 0 {
			singleArray := a.decisionMap[len(a.decisionMap)-1:]
			return singleArray
		} else {
			return a.decisionMap
		}
	}
}

// Description: The function returns our agent's (for a specific instance of our agent type) BaseAgent structures from the previous rounds
// Returns: If multi ==  true => array of the previous BaseAgent structures, else => array with a single element (BaseAgent struct of the last round)
func (a *Agent2) getBaseHelper(multi bool) []agent.BaseAgent {
	if multi == true {
		return a.baseAgentMap
	} else {
		if len(a.baseAgentMap)-1 > 0 {
			singleArray := a.baseAgentMap[len(a.baseAgentMap)-1:]
			return singleArray
		} else {
			return a.baseAgentMap
		}
	}
}

// Description: The function is used to extract the agents' HiddenAgentState struct of the previous rounds
// Returns: if multi => array of the previous HiddenAgentState structures, else => array with a single element (HiddenAgentState struct of the last round)
func (a *Agent2) getAgentStateHelper(multi bool) []immutable.Map[commons.ID, state.HiddenAgentState] {
	if multi == true {
		return a.agentStateMap
	} else {
		if len(a.agentStateMap)-1 > 0 {
			singleArray := a.agentStateMap[len(a.agentStateMap)-1:]
			return singleArray
		} else {
			return a.agentStateMap
		}
	}
}

// Description: The function is used to return the current Agent's bonus defense points
// Returns: uint
func getBonusDefense(agent agent.BaseAgent) uint {
	agentState := agent.AgentState()
	return agentState.BonusAttack()
}

// Description: The function is used to return the current Agent's bonus attack points
// Returns: uint
func getBonusAttack(agent agent.BaseAgent) uint {
	agentState := agent.AgentState()
	return agentState.BonusAttack()
}

func (a *Agent2) updateDecisionHelper(log immutable.Map[commons.ID, decision.FightAction]) {
	a.decisionMap = append(a.decisionMap, log)
}
func (a *Agent2) updateBaseHelper(baseAgent agent.BaseAgent) {
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

// -------------

func (a *Agent2) FightResolution(agent agent.BaseAgent, prop commons.ImmutableList[proposal.Rule[decision.FightAction]],
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

func (r *Agent2) LootActionNoProposal(baseAgent agent.BaseAgent) immutable.SortedMap[commons.ItemID, struct{}] {
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

func (r *Agent2) LootAction(
	_ agent.BaseAgent,
	proposedLoot immutable.SortedMap[commons.ItemID, struct{}],
	_ message.Proposal[decision.LootAction],
) immutable.SortedMap[commons.ItemID, struct{}] {
	return proposedLoot
}

func (r *Agent2) FightActionNoProposal(_ agent.BaseAgent) decision.FightAction {
	fight := rand.Intn(3)
	switch fight {
	case 0:
		return decision.Cower
	case 1:
		return decision.Attack
	default:
		return decision.Defend
	}
}

func (r *Agent2) FightAction(
	baseAgent agent.BaseAgent,
	_ decision.FightAction,
	_ message.Proposal[decision.FightAction],
) decision.FightAction {
	return r.FightActionNoProposal(baseAgent)
}

func (r *Agent2) HandleLootInformation(m message.TaggedInformMessage[message.LootInform], _ agent.BaseAgent) {
}

func (r *Agent2) HandleLootRequest(m message.TaggedRequestMessage[message.LootRequest]) message.LootInform {
	return nil
}

func (r *Agent2) HandleLootProposal(_ message.Proposal[decision.LootAction], _ agent.BaseAgent) decision.Intent {
	switch rand.Intn(3) {
	case 0:
		return decision.Positive
	case 1:
		return decision.Negative
	default:
		return decision.Abstain
	}
}

func (r *Agent2) HandleLootProposalRequest(_ message.Proposal[decision.LootAction], _ agent.BaseAgent) bool {
	switch rand.Intn(2) {
	case 0:
		return true
	default:
		return false
	}
}

func (r *Agent2) LootAllocation(
	baseAgent agent.BaseAgent,
	proposal message.Proposal[decision.LootAction],
	proposedAllocation immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]],
) immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]] {
	lootAllocation := make(map[commons.ID][]commons.ItemID)
	view := baseAgent.View()
	ids := commons.ImmutableMapKeys(view.AgentState())
	iterator := baseAgent.Loot().Weapons().Iterator()
	allocateRandomly(iterator, ids, lootAllocation)
	iterator = baseAgent.Loot().Shields().Iterator()
	allocateRandomly(iterator, ids, lootAllocation)
	iterator = baseAgent.Loot().HpPotions().Iterator()
	allocateRandomly(iterator, ids, lootAllocation)
	iterator = baseAgent.Loot().StaminaPotions().Iterator()
	allocateRandomly(iterator, ids, lootAllocation)
	mMapped := make(map[commons.ID]immutable.SortedMap[commons.ItemID, struct{}])
	for id, itemIDS := range lootAllocation {
		mMapped[id] = commons.ListToImmutableSortedSet(itemIDS)
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

func (r *Agent2) DonateToHpPool(baseAgent agent.BaseAgent) uint {
	return uint(rand.Intn(int(baseAgent.AgentState().Hp)))
}

// UpdateInternalState
// Description: the function is called at the end of each level (provides a list of type FightResult / can be thought as raw & processed overall game info)
func (a *Agent2) UpdateInternalState(baseAgent agent.BaseAgent, fightResult *commons.ImmutableList[decision.ImmutableFightResult], voteResult *immutable.Map[decision.Intent, uint], logChan chan<- logging.AgentLog) {
	a.updateBaseAgentPerLevel(baseAgent)
	a.updateFightResultPerLevel(*fightResult)
	a.updateVoteResultPerLevel(*voteResult)
	a.avgHp, a.avgDefend, a.avgAttack, a.avgStamina = updateAverages(baseAgent)
	a.updateSocialCapital(*fightResult, baseAgent)
	a.newGovernmentTimeline(baseAgent, a.haveElections)
}

// CreateManifesto
// Description: Used to give Manifesto Information if elected Leader.
// Return:		The Manifesto with FightImposition, LootImposition, term length and overthrow threshold.
func (a *Agent2) CreateManifesto(agent agent.BaseAgent) *decision.Manifesto {
	fightThreshold := 2.5
	lootThreshold := 2.5
	fightDecisionPower := false // default value

	if !a.wasOverthrown(agent.ID()) {
		if (a.adjustedExpertise(agent, 0, 5) + a.lastFightDecisionPower(agent.ID(), 2.5)) > fightThreshold {
			fightDecisionPower = true
		}
	} else {
		if a.adjustedExpertise(agent, 0, 5) > fightThreshold {
			fightDecisionPower = true
		}
	}

	lootDecisionPower := false

	if !a.wasOverthrown(agent.ID()) {
		if (a.adjustedExpertise(agent, 0, 5) + a.lastLootDecisionPower(agent.ID(), 2.5)) > lootThreshold {
			lootDecisionPower = true
		}
	} else {
		if a.adjustedExpertise(agent, 0, 5) > lootThreshold {
			lootDecisionPower = true
		}
	}

	termLength := uint(a.adjustedExpertise(agent, 0, 4) + 1)

	overthrowPercentage := uint(51)
	if a.wasOverthrown(agent.ID()) {
		overthrowPercentage = uint(float64(overthrowPercentage) + a.adjustedExpertise(agent, -10, 10))
	}

	manifesto := decision.NewManifesto(fightDecisionPower, lootDecisionPower, termLength, overthrowPercentage)
	return manifesto
}

// HandleConfidencePoll
// Description: Used for voting on confidence for Leader.
// Return:		Positive, Negative, or Abstain decision.
func (a *Agent2) HandleConfidencePoll(baseAgent agent.BaseAgent) decision.Intent {
	w0, w1, w2, w3, w4, w5 := 2.0, 1.0, 2.0, 1.0, 1.0, -1.0
	avgSurvivalCurrTermNorm := (a.avgSurvivalCurrTerm - a.avgSurvival) / a.avgSurvival
	avgSurvivalPastTermsNorm := (a.avgSurvivalPastTerms - a.avgSurvival) / a.avgSurvival
	avgBroadcastRateCurrTermNorm := (a.avgBroadcastRateCurrTerm - a.avgBroadcastRate) / a.avgBroadcastRate
	avgBroadcastRatePastTermNorm := (a.avgBroadcastRatePastTerms - a.avgBroadcastRate) / a.avgBroadcastRate
	leadershipXpNorm := (a.leadershipXp - a.avgLeadershipXp) / a.avgLeadershipXp

	noConfRateNorm := (a.noConfRate - a.avgNoConfRate) / a.avgNoConfRate
	trustWorthyness := w0*avgSurvivalCurrTermNorm + w1*avgSurvivalPastTermsNorm
	networks := w2*avgBroadcastRateCurrTermNorm + w3*avgBroadcastRatePastTermNorm
	institutions := w4*leadershipXpNorm + w5*noConfRateNorm

	sum := trustWorthyness + networks + institutions

	if sum >= 0 {
		return decision.Positive
	} else {
		return decision.Negative
	}
}

// HandleFightInformation
// Description: Called every time a fight information message is received (I believe it could be from a leader for providing a proposal or another agent for providing fight info (e.g proposal directly to them?)
// Return:		nil
func (a *Agent2) HandleFightInformation(m message.TaggedInformMessage[message.FightInform], baseAgent agent.BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) {
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

// HandleElectionBallot
// Description: Used to elect a Leader.
// Return:  	A single Commons.ID for choose-one voting or an array of commons.ID of top leader choices for ranked-voting.
func (a *Agent2) HandleElectionBallot(baseAgent agent.BaseAgent, params *decision.ElectionParams) decision.Ballot {
	// Extract ID of alive agents
	view := baseAgent.View()
	//agentState := view.AgentState()
	cans := params.CandidateList()
	// Updating Leader Parameters
	a.termEndLevel = view.CurrentLevel()                    //level_temp is another priv attribute initialized to 0/1?
	a.termDuration = view.CurrentLevel() - a.termBeginLevel // term_begin_level was last updated at the beginning of the term that is now ending
	a.termBeginLevel = view.CurrentLevel()                  // now that it has been used, can reset to track new leadership
	a.avgTermSurvival = a.avgSurvivalCurrTerm               //(assert len(a.survival_rates)==term_duration )
	a.avgTermBroadcastRate = a.avgBroadcastRateCurrTerm     //(assert len(a.broadcast_rates)==term_duration )
	a.haveElections = true
	if len(a.governmentTimeline) > 0 {
		lastLeaderInfo := a.governmentTimeline[len(a.governmentTimeline)-1]
		if lastLeaderInfo.duration < lastLeaderInfo.manifesto.TermLength() {
			lastLeaderInfo.overthrown = true
		}
		a.governmentTimeline[len(a.governmentTimeline)-1] = lastLeaderInfo
	}

	candidates := make([]commons.ID, cans.Len())
	mainfestos := make([]decision.Manifesto, cans.Len())
	i := 0
	itr := cans.Iterator()
	for !itr.Done() {
		can, manifesto, ok := itr.Next()
		if ok {
			candidates[i] = can
			mainfestos[i] = manifesto
			i++
		}
	}

	// Randomly fill the ballot
	var ballot decision.Ballot
	numAliveAgents := len(candidates)
	agentScores := make(map[commons.ID]float64, numAliveAgents)
	for i := 0; i < numAliveAgents; i++ {
		par1 := a.leaderElectedBefore(candidates[i], a.weightedFracTermsDeposed(-1, candidates[i])+a.weightedAvgSurRateUnderLeader(1.0))
		par2 := a.lastFightDecisionPower(candidates[i], 5)
		par3 := a.lastLootDecisionPower(candidates[i], 5)
		sot := a.SOT(candidates[i], mainfestos[i], 0.10)
		agentScores[candidates[i]] = prospectLeaderScore(par1, par2, par3, sot)
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		return agentScores[candidates[i]] < agentScores[candidates[j]]
	})
	for i := uint(0); i < params.NumberOfPreferences(); i++ {
		ballot = append(ballot, candidates[i])
	}
	return ballot
}

// HandleFightProposal
// Description: Through that function our agent votes on a broadcasts proposal
func (a *Agent2) HandleFightProposal(proposal message.Proposal[decision.FightAction], baseAgent agent.BaseAgent) decision.Intent {
	intent := rand.Intn(2)
	if intent == 0 {
		return decision.Positive
	} else {
		return decision.Negative
	}
}

// HandleFightProposalRequest
// Description: Only called as a leader: True for broadcasting the proposal / False for declining the proposal
// Return:		Bool: True/False
func (a *Agent2) HandleFightProposalRequest(prop message.Proposal[decision.FightAction], baseAgent agent.BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) bool {
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

func (r *Agent2) HandleUpdateWeapon(_ agent.BaseAgent) decision.ItemIdx {
	// 0th weapon has the greatest attack points
	return decision.ItemIdx(0)
}

func (r *Agent2) HandleUpdateShield(_ agent.BaseAgent) decision.ItemIdx {
	// 0th weapon has the greatest shield points
	return decision.ItemIdx(0)
}

func (r *Agent2) HandleTradeNegotiation(_ agent.BaseAgent, _ message.TradeInfo) message.TradeMessage {
	return message.TradeRequest{}
}
