package team5

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/message/proposal"
	"infra/game/state"
	"infra/logging"
	"infra/teams/team5/commons5"
	"math"
	"math/rand"
	"strings"

	"github.com/benbjohnson/immutable"
)

type Agent5 struct {
	lootInformation   commons5.Loot
	socialNetwork     SocialNetwork
	t5Manifesto       T5Manifesto
	myAgentState      commons5.MyAgentState
	myFightResolution immutable.Map[commons.ID, decision.FightAction]
	preHealth         uint
	prePopNum         uint
	exploreRate       float32
	qtable            *Qtable
	ttable            *TrustTable
	round             int
}

// --------------- Election ---------------

func (t5 *Agent5) CreateManifesto(baseAgent agent.BaseAgent) *decision.Manifesto {
	returnType := decision.NewManifesto(false, true, uint(t5.t5Manifesto.LeaderCombo+1), 50)
	return returnType
}

func (t5 *Agent5) HandleElectionBallot(baseAgent agent.BaseAgent, _ *decision.ElectionParams) decision.Ballot {
	var ballot decision.Ballot
	var trustee commons.ID
	fstCheck := true

	myview := baseAgent.View()
	globalStates := myview.AgentState()
	itr := globalStates.Iterator()
	for !itr.Done() {
		id, _, _ := itr.Next()
		if fstCheck || t5.ttable.EstimateLeadTrust(id) > t5.ttable.EstimateLeadTrust(trustee) {
			trustee = id
		}
	}
	ballot = append(ballot, trustee)

	return ballot
}

func (t5 *Agent5) HandleConfidencePoll(baseAgent agent.BaseAgent) decision.Intent {
	myview := baseAgent.View()
	currentLeader := myview.CurrentLeader()
	if t5.ttable.EstimateLeadTrust(currentLeader) < 5 {
		return decision.Negative
	}
	return decision.Positive
}

// --------------- Fight ---------------

func (t5 *Agent5) HandleUpdateWeapon(_ agent.BaseAgent) decision.ItemIdx {
	// 0th weapon has the greatest attack points
	return decision.ItemIdx(0)
}

func (t5 *Agent5) HandleUpdateShield(_ agent.BaseAgent) decision.ItemIdx {
	// 0th weapon has the greatest shield points
	return decision.ItemIdx(0)
}

func (t5 *Agent5) HandleFightRequest(_ message.TaggedRequestMessage[message.FightRequest], _ *immutable.Map[commons.ID, decision.FightAction]) message.FightInform {
	return nil
}

func (t5 *Agent5) HandleFightInformation(_ message.TaggedInformMessage[message.FightInform], baseAgent agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) {
	// Purpose by finding highest Q state associated with an action
	myview := baseAgent.View()
	globalStates := myview.AgentState()
	var globalATMax float32
	var globalSHMax float32
	for _, id := range commons.ImmutableMapKeys(globalStates) {
		agState, _ := globalStates.Get(id)
		if agState.Attack+agState.BonusAttack > uint(globalATMax) {
			globalATMax = float32(agState.Attack + agState.BonusAttack)
		}
		if agState.Defense+agState.BonusDefense > uint(globalSHMax) {
			globalATMax = float32(agState.Defense + agState.BonusDefense)
		}
	}

	rules := make([]proposal.Rule[decision.FightAction], 0)

	cowerState := t5.qtable.GetMaxQAction("Cower")
	if cowerState == "NoSaPairAvailable" {
		rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Cower,
			proposal.NewComparativeCondition(proposal.Health, proposal.LessThan, 10),
		))
	} else {
		cowerStateSplit := strings.Split(cowerState, "-")
		rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Cower,
			proposal.NewAndCondition(
				proposal.NewAndCondition(t5.findProposalHealth(cowerStateSplit[0]), t5.findProposalStamina(cowerStateSplit[1])),
				proposal.NewAndCondition(t5.findProposalAT(cowerStateSplit[1], globalATMax), t5.findProposalSH(cowerStateSplit[2], globalSHMax)))))
	}

	attckState := t5.qtable.GetMaxQAction("Attck")
	if attckState == "NoSaPairAvailable" {
		rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Attack,
			proposal.NewComparativeCondition(proposal.Stamina, proposal.GreaterThan, 500),
		))
	} else {
		attckStateSplit := strings.Split(attckState, "-")
		rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Attack,
			proposal.NewAndCondition(
				proposal.NewAndCondition(t5.findProposalHealth(attckStateSplit[0]), t5.findProposalStamina(attckStateSplit[1])),
				proposal.NewAndCondition(t5.findProposalAT(attckStateSplit[1], globalATMax), t5.findProposalSH(attckStateSplit[2], globalSHMax)))))
	}

	defndState := t5.qtable.GetMaxQAction("Defnd")
	if defndState == "NoSaPairAvailable" {
		rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Defend,
			proposal.NewComparativeCondition(proposal.Stamina, proposal.GreaterThan, 500),
		))
	} else {
		defndStateSplit := strings.Split(defndState, "-")
		rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Defend,
			proposal.NewAndCondition(
				proposal.NewAndCondition(t5.findProposalHealth(defndStateSplit[0]), t5.findProposalStamina(defndStateSplit[1])),
				proposal.NewAndCondition(t5.findProposalAT(defndStateSplit[1], globalATMax), t5.findProposalSH(defndStateSplit[2], globalSHMax)))))
	}

	prop := *commons.NewImmutableList(rules)
	_ = baseAgent.SendFightProposalToLeader(prop)
}

func (t5 *Agent5) HandleFightProposalRequest(
	// Leader fight function
	_ message.Proposal[decision.FightAction],
	_ agent.BaseAgent,
	proposals *immutable.Map[commons.ID, decision.FightAction],
) bool {
	allCount := 0
	cowCount := 0
	itr := proposals.Iterator()
	for !itr.Done() {
		_, fight, _ := itr.Next()
		allCount += 1
		if fight == decision.Cower {
			cowCount += 1
		}
	}
	percentCow := float64(cowCount) / float64(allCount)
	return percentCow <= 0.4
}

func (t5 *Agent5) HandleFightProposal(proposal message.Proposal[decision.FightAction], baseAgent agent.BaseAgent) decision.Intent {
	iter := proposal.Rules().Iterator()
	//var hamming_distance float32 = 0
	hamming_distance_placeholder := 0.0
	biwise_xor_sum_placeholder := 0.0
	entry_counter := 0.0
	id := proposal.ProposerID()
	if !iter.Done() {
		//using xor bitwise operator to find the hamming distance of their proposal and ours
		//hamming_distance += t5.myFightResolution.Get(decision.FightAction)^iter.Next()
		entry_counter += 1
		biwise_xor_sum_placeholder += rand.Float64()
	}
	hamming_distance_placeholder = biwise_xor_sum_placeholder / entry_counter

	if hamming_distance_placeholder <= 0.5 {
		t5.socialNetwork.UpdatePersonality(id, 0.1, 0)
		return decision.Positive
	} else if hamming_distance_placeholder <= 0.7 {
		t5.socialNetwork.UpdatePersonality(id, 0.05, 0)
		return decision.Abstain
	} else {
		t5.socialNetwork.UpdatePersonality(id, -0.01, 0)
		return decision.Negative
	}
}

// resolve fight

func (t5 *Agent5) FightActionNoProposal(baseAgent agent.BaseAgent) decision.FightAction {
	if t5.qtable.saTaken.state != "" {
		t5.UpdateQ(baseAgent)
	}
	t5.preHealth = baseAgent.AgentState().Hp
	myview := baseAgent.View()
	globalStates := myview.AgentState()
	t5.prePopNum = uint(globalStates.Len())
	qstate := t5.CurrentQState(baseAgent)
	if rand.Float32() < t5.exploreRate || len(t5.qtable.table) == 0 {
		return t5.Explore(qstate)
	}
	return t5.Exploit(qstate)
}

func (t5 *Agent5) FightAction(
	baseAgent agent.BaseAgent,
	proposedAction decision.FightAction,
	acceptedProposal message.Proposal[decision.FightAction],
) decision.FightAction {
	return t5.FightActionNoProposal(baseAgent)
}

func (t5 *Agent5) FightResolution(baseAgent agent.BaseAgent, prop commons.ImmutableList[proposal.Rule[decision.FightAction]], proposedActions immutable.Map[commons.ID, decision.FightAction]) immutable.Map[commons.ID, decision.FightAction] {
	view := baseAgent.View()
	builder := immutable.NewMapBuilder[commons.ID, decision.FightAction](nil)
	for _, id := range commons.ImmutableMapKeys(view.AgentState()) {
		var fightAction decision.FightAction
		agentAttack_placeholder := rand.Float64()
		agentShield_placeholder := rand.Float64()
		agentHP_placeholder := rand.Float64()      // *250 to convert hidden state to estimation
		agentStamina_placeholder := rand.Float64() //*500 to convert hiddent state to estimation
		fight_max := math.Max(agentAttack_placeholder, agentShield_placeholder)
		HPST_min := math.Max(agentHP_placeholder, agentStamina_placeholder)
		if fight_max > HPST_min {
			if agentAttack_placeholder > agentShield_placeholder {
				fightAction = decision.Attack
			} else {
				fightAction = decision.Defend
			}
		} else {
			fightAction = decision.Cower
		}

		builder.Set(id, fightAction)
	}
	//update agent's own fight resolution stored
	t5.myFightResolution = *builder.Map()
	return *builder.Map()
}

// --------------- Loot ---------------

// Loot allocation

func (t5 *Agent5) LootAllocation(baseAgent agent.BaseAgent, proposal message.Proposal[decision.LootAction], proposedAllocations immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]]) immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]] {
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

//Peseudo code for implementing Allocate according to agent state and aget personalities

//with the same type of loot:
//comparator is the agents' state field, which is to be sorted by it's value, and the loot is sorted by it's value as well;
//cases are set up to treat different groups of agent with different personlities,
//with all the agents that's qualified for this case, then within this case will do as following:
//the loot value is to be compared to the comparator value to iterate the allocation
//for example:
//case 1 qualifies Good-Lawful||Good-StrategyNeutral||GoodwillNeutral-Lawfull agents
//in this case, only those agents qualified but with low health or stamina<=certain value excluded,
//then agents are allocated loot with the logic:
// agent in the rest with the lowest shield points gets weapon in the rest with the highest attack points
//then case by case...

// func allocateAccordingly(iterator commons.Iterator[state.Item], comparator commons.ItemType, ids []commons.ID, lootAllocation map[commons.ID][]commons.ItemID) {
// 	for iterator != Done(){
// 		next, _ := iterator.Next()
// 		toBeAllocated := ids[sort(iterator.value)]
// 		switch agent.personality{
// 			case Good-Lawful||Good-StrategyNeutral||GoodwillNeutral-Lawfull:
// 				for id in range(view.agentState.HP.High){
// 					sort(comparator[id].value
// 					for id in range(sorted comparator value){
// 						if view.state[id].stamina >= 5*toBeAllocated
// 						lootAllocation[toBeAllocated] := ids
// 					}
// 				}
// 		}
// 	}
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

// resolve loot discussion

func (t5 *Agent5) LootActionNoProposal(baseAgent agent.BaseAgent) immutable.SortedMap[commons.ItemID, struct{}] {
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

func (t5 *Agent5) LootAction(
	_ agent.BaseAgent,
	proposedLoot immutable.SortedMap[commons.ItemID, struct{}],
	_ message.Proposal[decision.LootAction],
) immutable.SortedMap[commons.ItemID, struct{}] {
	return proposedLoot
}

func (t5 *Agent5) HandleLootInformation(m message.TaggedInformMessage[message.LootInform], agent agent.BaseAgent) {
	mystate := agent.AgentState()
	wants := make([]proposal.Rule[decision.LootAction], 0)

	wants = append(wants, *proposal.NewRule[decision.LootAction](decision.HealthPotion, proposal.NewComparativeCondition(proposal.Health, proposal.LessThan, 750)))
	wants = append(wants, *proposal.NewRule[decision.LootAction](decision.StaminaPotion, proposal.NewComparativeCondition(proposal.Stamina, proposal.LessThan, 1500)))
	wants = append(wants, *proposal.NewRule[decision.LootAction](decision.Weapon, proposal.NewComparativeCondition(proposal.TotalAttack, proposal.LessThan, mystate.TotalAttack())))
	wants = append(wants, *proposal.NewRule[decision.LootAction](decision.Shield, proposal.NewComparativeCondition(proposal.TotalDefence, proposal.LessThan, mystate.TotalDefense())))

	prop := *commons.NewImmutableList(wants)
	_ = agent.SendLootProposalToLeader(prop)
}

func (t5 *Agent5) HandleLootRequest(m message.TaggedRequestMessage[message.LootRequest]) message.LootInform {
	return nil
}

func (t5 *Agent5) HandleLootProposal(m message.Proposal[decision.LootAction], _ agent.BaseAgent) decision.Intent {
	id := m.ProposalID()

	switch t5.socialNetwork.AgentProfile[id].Strategy {
	case Lawful:
		return decision.Positive
	case Chaotic:
		return decision.Negative
	default:
		return decision.Abstain
	}
}

func (t5 *Agent5) HandleLootProposalRequest(_ message.Proposal[decision.LootAction], _ agent.BaseAgent) bool {
	switch rand.Intn(2) {
	case 0:
		return true
	default:
		return false
	}
}

// --------------- Trade ---------------

//func (t5 *Agent5) HandleTradeNegotiation(_ agent.BaseAgent, _ message.TradeInfo) message.TradeMessage {
//	return message.TradeRequest{}
//}

// --------------- Hp Pool ---------------

func (t5 *Agent5) DonateToHpPool(baseAgent agent.BaseAgent) uint {
	return 0
}

// --------------- Internal State ---------------

func (t5 *Agent5) UpdateTrust(baseAgent agent.BaseAgent) {
	t5.ttable.Decay()

	myview := baseAgent.View()
	globalStates := myview.AgentState()
	currentLeader := myview.CurrentLeader()

	// lead trust loss due to population death
	deathLastRound := t5.prePopNum - uint(globalStates.Len())
	if deathLastRound > 0 {
		t5.ttable.NegativeLeaderEvent(currentLeader, float32(deathLastRound))
	}
	// lead trust gain due to no health loss (being protected)
	hpLossLastRound := float32(t5.preHealth) - float32(baseAgent.AgentState().Hp)
	if hpLossLastRound > 0 {
		t5.ttable.NegativeLeaderEvent(currentLeader, -hpLossLastRound)
	}
	if hpLossLastRound <= 0 {
		t5.ttable.PositiveLeaderEvent(currentLeader, -hpLossLastRound)
	}

	for _, id := range commons.ImmutableMapKeys(globalStates) {
		agState, _ := globalStates.Get(id)
		// individual trust loss due to low health (not performing well)
		if uint(agState.Hp) == state.LowHealth {
			t5.ttable.NegativeIndivlEvent(id, 1)
		}
		// individual trust gain due to high health (performing well)
		if uint(agState.Hp) == state.HighHealth {
			t5.ttable.PositiveIndivlEvent(id, 1)
		}
	}
}

func (t5 *Agent5) UpdateInternalState(a agent.BaseAgent, _ *commons.ImmutableList[decision.ImmutableFightResult], _ *immutable.Map[decision.Intent, uint], log chan<- logging.AgentLog) {
	view := a.View()
	mas := a.AgentState()
	as := view.AgentState()
	iter := as.Iterator()
	leaderID := view.CurrentLeader()

	if !t5.myAgentState.Initilised {
		t5.myAgentState.InitMyAgentState()
	} else {
		t5.myAgentState = commons5.MyAgentState{
			MyAttackPoint: mas.Attack,
			MyShieldPoint: mas.Defense,
			MyHP:          mas.Hp,
			MyStamina:     mas.Stamina,
		}
	}

	//initialise social network after all agent has been spawned
	if !t5.socialNetwork.Initilised {
		t5.socialNetwork.InitSocialNetwork(a)
	}
	//if defect, goodwill += 0.2*(0.5-leader's goodwill), strategy += 0.2*(0.5-leader's strategy)
	//it not defect, goodwill +0.03
	for !iter.Done() {
		id, as, _ := iter.Next()

		//update trusts based on leader's personalities:
		leaderPersonalityCategorised := t5.socialNetwork.AgentProfile[leaderID].Goodwill + Goodwill(t5.socialNetwork.AgentProfile[leaderID].Strategy)
		if leaderPersonalityCategorised >= 4 {
			if as.Defector.IsDefector() {
				t5.socialNetwork.UpdatePersonality(id, -0.1, -0.1)
			} else {
				t5.socialNetwork.UpdatePersonality(id, 0.02, 0.02)
			}
		} else if leaderPersonalityCategorised >= 2 {
			if as.Defector.IsDefector() {
				t5.socialNetwork.UpdatePersonality(id, -0.05, -0.05)
			} else {
				t5.socialNetwork.UpdatePersonality(id, 0.02, 0.02)
			}
		} else if leaderPersonalityCategorised == 0 {
			if as.Defector.IsDefector() {
				t5.socialNetwork.UpdatePersonality(id, 0.1, 0.1)
			} else {
				t5.socialNetwork.UpdatePersonality(id, -0.02, -0.02)
			}
		}

		//update trusts based on trust scores(alternatice to the method above):
		//if defect, goodwill += 0.2*(0.5-leader's goodwill), strategy += 0.2*(0.5-leader's strategy)
		//it not defect, goodwill +0.03
		//==========code===========
		// if as.Defector.IsDefector() {
		// 	t5.socialNetwork.UpdatePersonality(id, 0.2*(0.5-t5.socialNetwork.AgentProfile[leaderID].Trusts.GoodwillScore), 0.2*(0.5-t5.socialNetwork.AgentProfile[leaderID].Trusts.StrategyScore))
		// } else {
		// 	t5.socialNetwork.UpdatePersonality(id, 0, 0.03)
		// }
	}

	t5.socialNetwork.Log(a.ID(), view.CurrentLevel())

	t5.t5Manifesto.updateLeaderCombo(a)

	t5.UpdateQ(a)
	t5.UpdateTrust(a)
	log <- logging.AgentLog{
		Name: a.Name(),
		ID:   a.ID(),
		Properties: map[string]float32{
			t5.qtable.Log(): 0,
			"trustToLeader": t5.ttable.EstimateLeadTrust(view.CurrentLeader()),
		},
	}
}

func NewAgent5() agent.Strategy {
	return &Agent5{
		socialNetwork: SocialNetwork{
			AgentProfile: make(map[commons.ID]AgentProfile),
			LawfullMin:   0.8,
			ChaoticMax:   0.2,
			GoodMin:      0.8,
			EvilMax:      0.2,
		},
		exploreRate: float32(0.25),
		qtable:      NewQTable(0.25, 0.75),
		ttable:      NewTrustTable(),
		round:       int(-1),
	}
}
