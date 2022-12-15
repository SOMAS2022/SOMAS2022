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
	"math/rand"
	"strings"

	"github.com/benbjohnson/immutable"
)

type Agent5 struct {
	//FightInformation
	// InternalState   internalState
	LootInformation commons5.Loot
	SocialNetwork   SocialNetwork
	LeaderCombo     uint
	bravery         int
	preHealth       uint
	prePopNum       uint
	exploreRate     float32
	qtable          *Qtable
	ttable          *TrustTable
}

// --------------- Election ---------------

func (t5 *Agent5) CreateManifesto(baseAgent agent.BaseAgent) *decision.Manifesto {
	returnType := decision.NewManifesto(false, true, t5.LeaderCombo+1, 50)
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
	if t5.ttable.EstimateLeadTrust(currentLeader) < 0 {
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
	var allCount float32 = 0.0
	var cowCount float32 = 0.0
	itr := proposals.Iterator()
	for !itr.Done() {
		_, fight, _ := itr.Next()
		allCount += 1
		if fight == decision.Cower {
			cowCount += 1
		}
	}
	percentCow := cowCount / allCount
	return percentCow <= 0.4
}

func (t5 *Agent5) HandleFightProposal(_ message.Proposal[decision.FightAction], _ agent.BaseAgent) decision.Intent {
	intent := rand.Intn(2)
	if intent == 0 {
		return decision.Positive
	} else {
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
}

func (t5 *Agent5) HandleLootRequest(m message.TaggedRequestMessage[message.LootRequest]) message.LootInform {
	return nil
}

func (t5 *Agent5) HandleLootProposal(_ message.Proposal[decision.LootAction], _ agent.BaseAgent) decision.Intent {
	switch rand.Intn(3) {
	case 0:
		return decision.Positive
	case 1:
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

func (t5 *Agent5) HandleTradeNegotiation(_ agent.BaseAgent, _ message.TradeInfo) message.TradeMessage {
	return message.TradeRequest{}
}

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
	hpLossLastRound := t5.preHealth - baseAgent.AgentState().Hp
	if hpLossLastRound <= 0 {
		t5.ttable.PositiveLeaderEvent(currentLeader, 1)
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
	t5.bravery += rand.Intn(10)
	t5.UpdateQ(a)
	t5.UpdateTrust(a)
	log <- logging.AgentLog{
		Name: a.Name(),
		ID:   a.ID(),
		Properties: map[string]float32{
			"bravery":       float32(t5.bravery),
			t5.qtable.Log(): 0,
		},
	}
}

func NewAgent5() agent.Strategy {
	return &Agent5{bravery: rand.Intn(5), exploreRate: float32(0.25), qtable: NewQTable(0.25, 0.75), ttable: NewTrustTable()}
}
