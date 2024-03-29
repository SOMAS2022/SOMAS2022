package team4

import (
	"math/rand"

	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/message/proposal"
	"infra/game/state"
	"infra/logging"

	"github.com/benbjohnson/immutable"
)

// Struct for AgentFour
type AgentFour struct {
	HP                int
	ST                int
	AT                int
	SH                int
	C                 int
	bravery           int
	hasVotedThisRound bool
	utilityScore      map[commons.ID]int
}

// bravery and utility score defined
func NewAgentFour() agent.Strategy {
	return &AgentFour{
		bravery:      rand.Intn(10),
		utilityScore: make(map[string]int),
	}
}

// *********************************** STRATEGY INTERFACE FUNCTIONS ***********************************

// we always pick our best shield
func (a *AgentFour) HandleUpdateShield(baseAgent agent.BaseAgent) decision.ItemIdx {
	// first shield has highest points, since inventory is sorted
	return decision.ItemIdx(0)
}

// we always pick the best weapon
func (a *AgentFour) HandleUpdateWeapon(baseAgent agent.BaseAgent) decision.ItemIdx {
	// first weapon has highest points, since inventory is sorted
	return decision.ItemIdx(0)
}

// Define and update the attributes for agent four
func (a *AgentFour) UpdateInternalState(baseAgent agent.BaseAgent, _ *commons.ImmutableList[decision.ImmutableFightResult], _ *immutable.Map[decision.Intent, uint], log chan<- logging.AgentLog) {
	a.UpdateUtility(baseAgent)
	a.HP = int(baseAgent.AgentState().Hp)
	a.ST = int(baseAgent.AgentState().Stamina)
	a.AT = int(baseAgent.AgentState().Attack)
	//a.SH = int(baseAgent.AgentState().Shields)
	a.C = 0
}

// *********************************** ELECTION INTERFACE FUNCTIONS ***********************************

func (a *AgentFour) CreateManifesto(_ agent.BaseAgent) *decision.Manifesto {
	manifesto := decision.NewManifesto(false, false, 10, 5)
	return manifesto
}

func (a *AgentFour) HandleConfidencePoll(baseAgent agent.BaseAgent) decision.Intent {
	switch rand.Intn(3) {
	case 0:
		return decision.Abstain
	case 1:
		return decision.Negative
	default:
		return decision.Positive
	}
}

func (a *AgentFour) HandleElectionBallot(baseAgent agent.BaseAgent, _ *decision.ElectionParams) decision.Ballot {
	view := baseAgent.View()
	agentState := view.AgentState()
	aliveAgentIDs := make([]string, agentState.Len())
	i := 0
	itr := agentState.Iterator()
	for !itr.Done() {
		id, a, ok := itr.Next()
		if ok && a.Hp > 0 {
			aliveAgentIDs[i] = id
			i++
		}
	}

	var ballot decision.Ballot
	numAliveAgents := len(aliveAgentIDs)
	numCandidate := rand.Intn(numAliveAgents)
	for i := 0; i < numCandidate; i++ {
		randomIdx := rand.Intn(numAliveAgents)
		randomCandidate := aliveAgentIDs[uint(randomIdx)]
		ballot = append(ballot, randomCandidate)
	}

	return ballot
}

// *********************************** FIGHT INTERFACE FUNCTIONS ***********************************

func (a *AgentFour) HandleFightInformation(_ message.TaggedInformMessage[message.FightInform], baseAgent agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) {
	makesProposal := rand.Intn(100)

	if makesProposal > 80 {
		rules := make([]proposal.Rule[decision.FightAction], 0)

		rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Attack,
			proposal.NewAndCondition(*proposal.NewComparativeCondition(proposal.Health, proposal.GreaterThan, 1000),
				*proposal.NewComparativeCondition(proposal.Stamina, proposal.GreaterThan, 1000)),
		))

		rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Defend,
			proposal.NewComparativeCondition(proposal.TotalDefence, proposal.GreaterThan, 1000),
		))

		rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Cower,
			proposal.NewComparativeCondition(proposal.Health, proposal.LessThan, 1),
		))

		rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Attack,
			proposal.NewComparativeCondition(proposal.Stamina, proposal.GreaterThan, 10),
		))

		prop := *commons.NewImmutableList(rules)
		_ = baseAgent.SendFightProposalToLeader(prop)
	}
}

func (a *AgentFour) HandleFightRequest(_ message.TaggedRequestMessage[message.FightRequest], _ *immutable.Map[commons.ID, decision.FightAction]) message.FightInform {
	return nil
}

func (a *AgentFour) FightResolution(baseAgent agent.BaseAgent, prop commons.ImmutableList[proposal.Rule[decision.FightAction]], proposedActions immutable.Map[string, decision.FightAction]) immutable.Map[string, decision.FightAction] { // Attack-Defend-Cower Strategy
	view := baseAgent.View()
	Agentstate := baseAgent.AgentState()
	builder := immutable.NewMapBuilder[commons.ID, decision.FightAction](nil)
	TotalAttack := Agentstate.Attack + Agentstate.BonusAttack()
	TotalDefense := Agentstate.Defense + Agentstate.BonusDefense()

	//var fightRes = decision.FightResult.CoweringAgents
	damage := int(view.MonsterAttack()) / rand.Intn(20)

	for _, id := range commons.ImmutableMapKeys(view.AgentState()) {
		var fightAction decision.FightAction

		if a.HP > (damage + 1) {
			if float64(TotalAttack) >= float64(TotalDefense)*0.8 {
				fightAction = decision.Attack
			} else {
				fightAction = decision.Defend
			}
		} else {
			fightAction = decision.Cower
		}
		builder.Set(id, fightAction)
	}
	return *builder.Map()
}

func (a *AgentFour) HandleFightProposal(m message.Proposal[decision.FightAction], baseAgent agent.BaseAgent) decision.Intent {
	intent := rand.Intn(2)
	if intent == 0 {
		return decision.Positive
	} else {
		return decision.Negative
	}
}

func (a *AgentFour) HandleFightProposalRequest(_ message.Proposal[decision.FightAction], _ agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) bool {
	switch rand.Intn(2) {
	case 0:
		return true
	default:
		return false
	}
}

func (a *AgentFour) FightActionNoProposal(baseAgent agent.BaseAgent) decision.FightAction {
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

func (a *AgentFour) FightAction(baseAgent agent.BaseAgent, proposedAction decision.FightAction, acceptedProposal message.Proposal[decision.FightAction]) decision.FightAction {
	return a.FightActionNoProposal(baseAgent)
}

// *********************************** LOOT INTERFACE FUNCTIONS ***********************************

func (a *AgentFour) HandleLootInformation(m message.TaggedInformMessage[message.LootInform], agent agent.BaseAgent) {
	//
}

func (a *AgentFour) HandleLootRequest(m message.TaggedRequestMessage[message.LootRequest]) message.LootInform {
	panic("implement me")
}

func (a *AgentFour) HandleLootProposal(_ message.Proposal[decision.LootAction], _ agent.BaseAgent) decision.Intent {
	switch rand.Intn(3) {
	case 0:
		return decision.Positive
	case 1:
		return decision.Negative
	default:
		return decision.Abstain
	}
}

func (a *AgentFour) HandleLootProposalRequest(_ message.Proposal[decision.LootAction], _ agent.BaseAgent) bool {
	switch rand.Intn(2) {
	case 0:
		return true
	default:
		return false
	}
}

func (a *AgentFour) LootAllocation(baseAgent agent.BaseAgent, proposal message.Proposal[decision.LootAction], proposedAllocations immutable.Map[string, immutable.SortedMap[string, struct{}]]) immutable.Map[string, immutable.SortedMap[string, struct{}]] {
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

func (a *AgentFour) LootActionNoProposal(baseAgent agent.BaseAgent) immutable.SortedMap[commons.ItemID, struct{}] {
	return *immutable.NewSortedMap[commons.ItemID, struct{}](nil)
}

func (a *AgentFour) LootAction(baseAgent agent.BaseAgent, proposedLoot immutable.SortedMap[commons.ItemID, struct{}], acceptedProposal message.Proposal[decision.LootAction]) immutable.SortedMap[commons.ItemID, struct{}] {
	return proposedLoot
}

// func (a *AgentFour) LootManifesto(baseAgent agent.BaseAgent) {
// 	Agentstate := baseAgent.AgentState()
// 	TotalAttack := Agentstate.Attack + Agentstate.BonusAttack()
// 	TotalDefense := Agentstate.Defense + Agentstate.BonusDefense()

// 	ratio_agents_HPLow := get_HP_levels.sliceOfAgentsWithLowHealth / logging.LevelStats.NumberOfAgents
// 	ratio_agents_HPNormal := get_HP_levels.sliceOfAgentsWithMidHealth / logging.LevelStats.NumberOfAgents
// 	ratio_agents_HPHigh := get_HP_levels.sliceOfAgentsWithHighHealth / logging.LevelStats.NumberOfAgents

// 	ratio_agents_STLow := get_ST_levels.sliceOfAgentsWithLowST / logging.LevelStats.NumberOfAgents
// 	ratio_agents_STNormal := get_ST_levels.sliceOfAgentsWithMidST / logging.LevelStats.NumberOfAgents
// 	ratio_agents_STHigh := get_ST_levels.sliceOfAgentsWithHighST / logging.LevelStats.NumberOfAgents

// 	thresh_attack := decision.NewImmutableFightResult().AttackSum() / logging.LevelStats.NumberOfAgents
// 	thresh_defend := decision.NewImmutableFightResult().ShieldSum() / logging.LevelStats.NumberOfAgents

// 	threshold_fight_HP := ratio_agents_HPLow*(250) + ratio_agents_HPNormal*(500) + ratio_agents_HPHigh*(750)
// 	threshold_fight_ST := ratio_agents_STLow*(500) + ratio_agents_STNormal*(1000) + ratio_agents_STHigh*(1500)

// 	switch {
// 	case a.HP > threshold_fight_HP && a.ST > threshold_fight_ST:
// 		if TotalAttack < thresh_attack {
// 			//Can get sword (any level)
// 		}

// 		if TotalDefense < thresh_defend {
// 			//Can get shield (any level)
// 		}

// 	case a.HP < threshold_fight_HP:
// 		//Can get HP Potion

// 	case a.ST < threshold_fight_ST:
// 		//Can get ST Potion
// 	}
// }

// *********************************** HPPOOL INTERFACE FUNCTIONS ***********************************

// HP pool donation
func (a *AgentFour) DonateToHpPool(baseAgent agent.BaseAgent) uint {
	C_thresh_HP := 1
	donation := float64(0)
	initialHP := 1000
	first := 0.8
	second := 0.1

	if float64(a.HP) > first*float64(initialHP) && a.C < C_thresh_HP {
		donation = float64(a.HP) * second
		a.C += 1
	} else {
		donation = 0
	}
	return uint(donation)
}

// *********************************** TRADE INTERFACE FUNCTIONS ***********************************

func (a *AgentFour) HandleTradeNegotiation(theAgent agent.BaseAgent, m message.TradeInfo) message.TradeMessage {
	return message.TradeRequest{}
}

// *********************************** OTHER FUNCTIONS ***********************************
//

// PrefightComms
// func (a *AgentFour) PrefightComms(state state.State, agent agent.BaseAgent) {
// 	//actions array stores all agent actions from the beginning to the end of the game.
// 	actions := make([]commons.ID, 0)              //[[agent_id, [action_round1, action_round2,.....]]]
// 	agents_contributions := make([]commons.ID, 0) //[C1, C2, ……]
// 	num_agents_fight := 0
// 	num_agents_defend := 0
// 	num_agents_cower := 0
// 	i := 0

// 	NumberOfAgents := baseAgent.AgentState()

// 	for i < NumberOfAgents {
// 		//asks what is the action of agent[i]
// 		switch {
// 		case actions[i] == fight:
// 			num_agents_fight += 1
// 			actions[agent_id] == decision.Attack
// 			if a.HP == state.LowHealth || a.ST == state.LowStamina {
// 				agents_contributions[agent_ID] += 2
// 			}

// 		case action[i] == defend:
// 			num_agents_defend += 1
// 			actions[agent_id] == decision.Defend
// 			agents_contributions[agent_ID] += 1
// 			if a.HP == state.LowHealth || a.ST == state.LowStamina {
// 				agents_contributions[agent_ID] += 2
// 			}

// 		case action[i] == cower:
// 			num_agents_cower += 1
// 			actions[agent_id] = decision.Cower
// 			if a.HP == state.HighHealth || a.ST == state.HighStamina {
// 				agents_contributions[agent_ID] -= 1
// 			}
// 		}
// 	}
// }

func (a *AgentFour) HPLevels(agent agent.BaseAgent) [][]string {
	view := agent.View()
	agentState := view.AgentState()
	agentStateIterator := agentState.Iterator()
	sliceOfAgentsWithLowHealth := make([]commons.ID, 0)
	sliceOfAgentsWithMidHealth := make([]commons.ID, 0)
	sliceOfAgentsWithHighHealth := make([]commons.ID, 0)
	for !agentStateIterator.Done() {
		key, value, _ := agentStateIterator.Next()
		switch {
		case value.Hp == state.HealthRange(state.LowHealth):
			sliceOfAgentsWithLowHealth = append(sliceOfAgentsWithLowHealth, key)
		case value.Hp == state.HealthRange(state.MidHealth):
			sliceOfAgentsWithMidHealth = append(sliceOfAgentsWithMidHealth, key)
		case value.Hp == state.HealthRange(state.HighHealth):
			sliceOfAgentsWithHighHealth = append(sliceOfAgentsWithHighHealth, key)
		}
	}
	var l = [][]string{sliceOfAgentsWithLowHealth, sliceOfAgentsWithMidHealth, sliceOfAgentsWithHighHealth}
	return l
}

func (a *AgentFour) STLevels(agent agent.BaseAgent) [][]string {
	view := agent.View()
	agentState := view.AgentState()
	agentStateIterator := agentState.Iterator()
	sliceOfAgentsWithLowST := make([]commons.ID, 0)
	sliceOfAgentsWithMidST := make([]commons.ID, 0)
	sliceOfAgentsWithHighST := make([]commons.ID, 0)
	for !agentStateIterator.Done() {
		key, value, _ := agentStateIterator.Next()
		switch {
		case value.Hp == state.HealthRange(state.LowStamina):
			sliceOfAgentsWithLowST = append(sliceOfAgentsWithLowST, key)
		case value.Hp == state.HealthRange(state.MidStamina):
			sliceOfAgentsWithMidST = append(sliceOfAgentsWithMidST, key)
		case value.Hp == state.HealthRange(state.HighStamina):
			sliceOfAgentsWithHighST = append(sliceOfAgentsWithHighST, key)
		}
	}
	var l = [][]string{sliceOfAgentsWithLowST, sliceOfAgentsWithMidST, sliceOfAgentsWithHighST}
	return l
}

// Attack-Defend-Cower Strategy
func (a *AgentFour) AttackDefendCower(state state.AgentState, baseAgent agent.BaseAgent, fightResult int) *decision.FightAction {
	view := baseAgent.View()
	Agentstate := baseAgent.AgentState()
	TotalAttack := Agentstate.Attack + Agentstate.BonusAttack()
	TotalDefense := Agentstate.Defense + Agentstate.BonusDefense()
	var action decision.FightAction
	damage := view.MonsterAttack() / uint(fightResult)

	if uint(a.HP) > (damage + 1) {
		if float64(TotalAttack) >= float64(TotalDefense)*0.8 {
			action = decision.Attack
		} else {
			action = decision.Defend
		}
	} else {
		action = decision.Cower
	}
	return &action
}

// FightManifesto
func (a *AgentFour) FightManifesto(baseAgent agent.BaseAgent, prop commons.ImmutableList[proposal.Rule[decision.FightAction]]) immutable.Map[commons.ID, decision.FightAction] {
	Agentstate := baseAgent.AgentState()
	TotalAttack := Agentstate.Attack + Agentstate.BonusAttack()
	TotalDefense := Agentstate.Defense + Agentstate.BonusDefense()
	view := baseAgent.View()
	builder := immutable.NewMapBuilder[commons.ID, decision.FightAction](nil)
	var manifesto_decision decision.FightAction
	rand_prob := 0.5

	//var get_HP_levels = a.HPLevels(baseAgent)
	var HP_levels_list = a.HPLevels(baseAgent)
	var ST_levels_list = a.STLevels(baseAgent)
	thresh_fight := 0.3

	var agent_map = view.AgentState()

	ratio_agents_HPLow := len(HP_levels_list[0]) / agent_map.Len()
	ratio_agents_HPNormal := len(HP_levels_list[1]) / agent_map.Len()
	ratio_agents_HPHigh := len(HP_levels_list[2]) / agent_map.Len()

	ratio_agents_STLow := len(ST_levels_list[0]) / agent_map.Len()
	ratio_agents_STNormal := len(ST_levels_list[1]) / agent_map.Len()
	ratio_agents_STHigh := len(ST_levels_list[2]) / agent_map.Len()

	thresh_attack := rand.Intn(50) / agent_map.Len()
	thresh_defend := rand.Intn(50) / agent_map.Len()
	threshold_fight_HP := ratio_agents_HPLow*(250) + ratio_agents_HPNormal*(500) + ratio_agents_HPHigh*(750)
	threshold_fight_ST := ratio_agents_STLow*(500) + ratio_agents_STNormal*(1000) + ratio_agents_STHigh*(1500)

	var fightRes = (rand.Intn(20)) + 1

	//var FightMethod = a.AttackDefendCower(Agentstate, baseAgent, &fightActions)
	var FightMethod = a.AttackDefendCower(Agentstate, baseAgent, fightRes)

	for _, id := range commons.ImmutableMapKeys(view.AgentState()) {
		if a.HP > threshold_fight_HP && a.ST > threshold_fight_ST {
			if TotalAttack > uint(thresh_attack) && TotalDefense > uint(thresh_defend) {
				switch {
				case rand_prob >= 0.4:
					manifesto_decision = decision.Defend
				case rand_prob <= 0.6:
					manifesto_decision = decision.Attack
				}
			}
			if TotalAttack > uint(thresh_attack) {
				manifesto_decision = decision.Attack
			}
			if TotalDefense > uint(thresh_defend) {
				manifesto_decision = decision.Defend
			}
		} else {
			manifesto_decision = decision.Cower
		}

		if *FightMethod == decision.Cower && manifesto_decision == decision.Attack {
			if rand_prob < thresh_fight {
				threshold_fight_HP = a.HP + 10
				a.C -= 1
			} else {
				a.C += 1
			}
		}
		builder.Set(id, *FightMethod)
	}
	return *builder.Map()
}

func (a *AgentFour) VoteFightManifesto(baseAgent agent.BaseAgent) {
	Agentstate := baseAgent.AgentState()
	view := baseAgent.View()
	var HP_levels_list = a.HPLevels(baseAgent)
	var ST_levels_list = a.STLevels(baseAgent)
	var agent_map = view.AgentState()

	ratio_agents_HPLow := len(HP_levels_list[0]) / agent_map.Len()
	ratio_agents_HPNormal := len(HP_levels_list[1]) / agent_map.Len()
	ratio_agents_HPHigh := len(HP_levels_list[2]) / agent_map.Len()
	ratio_agents_STLow := len(ST_levels_list[0]) / agent_map.Len()
	ratio_agents_STNormal := len(ST_levels_list[1]) / agent_map.Len()
	ratio_agents_STHigh := len(ST_levels_list[2]) / agent_map.Len()

	threshold_fight_HP := ratio_agents_HPLow*(250) + ratio_agents_HPNormal*(500) + ratio_agents_HPHigh*(750)
	threshold_fight_ST := ratio_agents_STLow*(500) + ratio_agents_STNormal*(1000) + ratio_agents_STHigh*(1500)

	thresh_attack := rand.Intn(20) / agent_map.Len()
	thresh_defend := rand.Intn(20) / agent_map.Len()

	TotalAttack := Agentstate.Attack + Agentstate.BonusAttack()
	TotalDefense := Agentstate.Defense + Agentstate.BonusDefense()
	v_tol := 0.6
	if float64(a.HP) >= v_tol*float64(threshold_fight_HP) && float64(a.HP) <= (1+v_tol)*float64(threshold_fight_HP) && float64(a.ST) >= v_tol*float64(threshold_fight_ST) && float64(a.ST) <= (1+v_tol)*float64(threshold_fight_ST) && float64(TotalAttack) >= v_tol*float64(thresh_attack) && float64(TotalAttack) <= (1+v_tol)*float64(thresh_attack) && float64(TotalDefense) >= v_tol*float64(thresh_defend) && float64(TotalDefense) <= (1+v_tol)*float64(thresh_defend) {
		a.hasVotedThisRound = true //vote YES
	} else {
		a.hasVotedThisRound = false //vote NO
	}
}

func (a *AgentFour) VoteLootManifesto(baseAgent agent.BaseAgent) {
	Agentstate := baseAgent.AgentState()
	view := baseAgent.View()
	var HP_levels_list = a.HPLevels(baseAgent)
	var ST_levels_list = a.STLevels(baseAgent)
	var agent_map = view.AgentState()

	ratio_agents_HPLow := len(HP_levels_list[0]) / agent_map.Len()
	ratio_agents_HPNormal := len(HP_levels_list[1]) / agent_map.Len()
	ratio_agents_HPHigh := len(HP_levels_list[2]) / agent_map.Len()
	ratio_agents_STLow := len(ST_levels_list[0]) / agent_map.Len()
	ratio_agents_STNormal := len(ST_levels_list[1]) / agent_map.Len()
	ratio_agents_STHigh := len(ST_levels_list[2]) / agent_map.Len()

	threshold_fight_HP := ratio_agents_HPLow*(250) + ratio_agents_HPNormal*(500) + ratio_agents_HPHigh*(750)
	threshold_fight_ST := ratio_agents_STLow*(500) + ratio_agents_STNormal*(1000) + ratio_agents_STHigh*(1500)

	thresh_attack := rand.Intn(20) / agent_map.Len()
	thresh_defend := rand.Intn(20) / agent_map.Len()

	TotalAttack := Agentstate.Attack + Agentstate.BonusAttack()
	TotalDefense := Agentstate.Defense + Agentstate.BonusDefense()
	v_tol := 0.6

	if float64(a.HP) >= v_tol*float64(threshold_fight_HP) && float64(a.HP) <= (1+v_tol)*float64(threshold_fight_HP) && float64(a.ST) >= v_tol*float64(threshold_fight_ST) && float64(a.ST) <= (1+v_tol)*float64(threshold_fight_ST) && float64(TotalAttack) >= v_tol*float64(thresh_attack) && float64(TotalAttack) <= (1+v_tol)*float64(thresh_attack) && float64(TotalDefense) >= v_tol*float64(thresh_defend) && float64(TotalDefense) <= (1+v_tol)*float64(thresh_defend) {
		a.hasVotedThisRound = true //vote YES
	} else {
		a.hasVotedThisRound = false //vote NO
	}
}

// func (a *AgentFour) DisobeyFightManifesto(baseAgent agent.BaseAgent) {
// 	Agentstate := baseAgent.AgentState()
// 	var fightRes = rand.Intn(20)
// 	var FightMethod = a.AttackDefendCower(Agentstate, baseAgent, fightRes)
// 	var manifesto_decision = a.FightActionNoProposal(baseAgent)
// 	//Receive fight manifesto_decision from the server
// 	definite_decision := manifesto_decision //Default
// 	if *FightMethod != manifesto_decision {
// 		if manifesto_decision == decision.Attack || decision.Defend {
// 			if float64(a.HP) < 0.1 || float64(a.ST) < 0.1 {
// 				definite_decision = *FightMethod
// 				a.C = a.C - 2
// 			}
// 		}
// 	}
// }

//For each fight round
// itr := fightResult.Iterator()
// 	for !itr.Done() { // For each fight round
// 	}

//Alternative FightManifesto method
// func FightManifesto(agents map[commons.ID]agent.Agent) map[commons.ID]decision.FightAction {
// 	decisionMap := make(map[commons.ID]decision.FightAction)

// 	for i := range agents {
// 		decisionMap[i] = decision.Defend
// 	}

// 	return decisionMap
// }

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

func (a *AgentFour) UpdateUtility(baseAgent agent.BaseAgent) {
}
