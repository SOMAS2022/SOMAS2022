package team5

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/message/proposal"
	"infra/game/state"
	"infra/logging"
	"math/rand"
	"strings"

	"github.com/benbjohnson/immutable"
)

type FivAgent struct {
	bravery     int
	preHealth   uint
	prePopNum   uint
	exploreRate float32
	qtable      *Qtable
	ttable      *TrustTable
}

// --------------- Election ---------------

func (fiv *FivAgent) CreateManifesto(_ agent.BaseAgent) *decision.Manifesto {
	manifesto := decision.NewManifesto(false, false, 4, 10)
	return manifesto
}

func (fiv *FivAgent) HandleElectionBallot(baseAgent agent.BaseAgent, _ *decision.ElectionParams) decision.Ballot {
	var ballot decision.Ballot
	var trustee commons.ID
	fstCheck := true

	myview := baseAgent.View()
	globalStates := myview.AgentState()
	itr := globalStates.Iterator()
	for !itr.Done() {
		id, _, _ := itr.Next()
		if fstCheck || fiv.ttable.EstimateLeadTrust(id) > fiv.ttable.EstimateLeadTrust(trustee) {
			trustee = id
		}
	}
	ballot = append(ballot, trustee)

	return ballot
}

func (fiv *FivAgent) HandleConfidencePoll(baseAgent agent.BaseAgent) decision.Intent {
	myview := baseAgent.View()
	currentLeader := myview.CurrentLeader()
	if fiv.ttable.EstimateLeadTrust(currentLeader) < 5 {
		return decision.Negative
	}
	return decision.Positive
}

// --------------- Fight ---------------

func (fiv *FivAgent) HandleUpdateWeapon(_ agent.BaseAgent) decision.ItemIdx {
	// 0th weapon has the greatest attack points
	return decision.ItemIdx(0)
}

func (fiv *FivAgent) HandleUpdateShield(_ agent.BaseAgent) decision.ItemIdx {
	// 0th weapon has the greatest shield points
	return decision.ItemIdx(0)
}

func (fiv *FivAgent) HandleFightRequest(_ message.TaggedRequestMessage[message.FightRequest], _ *immutable.Map[commons.ID, decision.FightAction]) message.FightInform {
	return nil
}

func (fiv *FivAgent) HandleFightInformation(_ message.TaggedInformMessage[message.FightInform], baseAgent agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) {
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

	cowerState := fiv.qtable.GetMaxQAction("Cower")
	if cowerState == "NoSaPairAvailable" {
		rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Cower,
			proposal.NewComparativeCondition(proposal.Health, proposal.LessThan, 10),
		))
	} else {
		cowerStateSplit := strings.Split(cowerState, "-")
		rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Cower,
			proposal.NewAndCondition(
				proposal.NewAndCondition(fiv.findProposalHealth(cowerStateSplit[0]), fiv.findProposalStamina(cowerStateSplit[1])),
				proposal.NewAndCondition(fiv.findProposalAT(cowerStateSplit[1], globalATMax), fiv.findProposalSH(cowerStateSplit[2], globalSHMax)))))
	}

	attckState := fiv.qtable.GetMaxQAction("Attck")
	if attckState == "NoSaPairAvailable" {
		rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Attack,
			proposal.NewComparativeCondition(proposal.Stamina, proposal.GreaterThan, 500),
		))
	} else {
		attckStateSplit := strings.Split(attckState, "-")
		rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Attack,
			proposal.NewAndCondition(
				proposal.NewAndCondition(fiv.findProposalHealth(attckStateSplit[0]), fiv.findProposalStamina(attckStateSplit[1])),
				proposal.NewAndCondition(fiv.findProposalAT(attckStateSplit[1], globalATMax), fiv.findProposalSH(attckStateSplit[2], globalSHMax)))))
	}

	defndState := fiv.qtable.GetMaxQAction("Defnd")
	if defndState == "NoSaPairAvailable" {
		rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Defend,
			proposal.NewComparativeCondition(proposal.Stamina, proposal.GreaterThan, 500),
		))
	} else {
		defndStateSplit := strings.Split(defndState, "-")
		rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Defend,
			proposal.NewAndCondition(
				proposal.NewAndCondition(fiv.findProposalHealth(defndStateSplit[0]), fiv.findProposalStamina(defndStateSplit[1])),
				proposal.NewAndCondition(fiv.findProposalAT(defndStateSplit[1], globalATMax), fiv.findProposalSH(defndStateSplit[2], globalSHMax)))))
	}

	prop := *commons.NewImmutableList(rules)
	_ = baseAgent.SendFightProposalToLeader(prop)
}

func (fiv *FivAgent) HandleFightProposalRequest(
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

func (fiv *FivAgent) HandleFightProposal(_ message.Proposal[decision.FightAction], _ agent.BaseAgent) decision.Intent {
	intent := rand.Intn(2)
	if intent == 0 {
		return decision.Positive
	} else {
		return decision.Negative
	}
}

// resolve fight

func (fiv *FivAgent) FightActionNoProposal(baseAgent agent.BaseAgent) decision.FightAction {
	if fiv.qtable.saTaken.state != "" {
		fiv.UpdateQ(baseAgent)
	}
	fiv.preHealth = baseAgent.AgentState().Hp
	myview := baseAgent.View()
	globalStates := myview.AgentState()
	fiv.prePopNum = uint(globalStates.Len())
	qstate := fiv.CurrentQState(baseAgent)
	if rand.Float32() < fiv.exploreRate || len(fiv.qtable.table) == 0 {
		return fiv.Explore(qstate)
	}
	return fiv.Exploit(qstate)
}

func (fiv *FivAgent) FightAction(
	baseAgent agent.BaseAgent,
	proposedAction decision.FightAction,
	acceptedProposal message.Proposal[decision.FightAction],
) decision.FightAction {
	return fiv.FightActionNoProposal(baseAgent)
}

func (fiv *FivAgent) FightResolution(baseAgent agent.BaseAgent, prop commons.ImmutableList[proposal.Rule[decision.FightAction]], proposedActions immutable.Map[commons.ID, decision.FightAction]) immutable.Map[commons.ID, decision.FightAction] {
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

func (fiv *FivAgent) LootAllocation(baseAgent agent.BaseAgent, proposal message.Proposal[decision.LootAction], proposedAllocations immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]]) immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]] {
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

func (fiv *FivAgent) LootActionNoProposal(baseAgent agent.BaseAgent) immutable.SortedMap[commons.ItemID, struct{}] {
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

func (fiv *FivAgent) LootAction(
	_ agent.BaseAgent,
	proposedLoot immutable.SortedMap[commons.ItemID, struct{}],
	_ message.Proposal[decision.LootAction],
) immutable.SortedMap[commons.ItemID, struct{}] {
	return proposedLoot
}

func (fiv *FivAgent) HandleLootInformation(m message.TaggedInformMessage[message.LootInform], agent agent.BaseAgent) {
	mystate := agent.AgentState()
	wants := make([]proposal.Rule[decision.LootAction], 0)

	wants = append(wants, *proposal.NewRule[decision.LootAction](decision.HealthPotion, proposal.NewComparativeCondition(proposal.Health, proposal.LessThan, 750)))
	wants = append(wants, *proposal.NewRule[decision.LootAction](decision.StaminaPotion, proposal.NewComparativeCondition(proposal.Stamina, proposal.LessThan, 1500)))
	wants = append(wants, *proposal.NewRule[decision.LootAction](decision.Weapon, proposal.NewComparativeCondition(proposal.TotalAttack, proposal.LessThan, mystate.TotalAttack())))
	wants = append(wants, *proposal.NewRule[decision.LootAction](decision.Shield, proposal.NewComparativeCondition(proposal.TotalDefence, proposal.LessThan, mystate.TotalDefense())))

	prop := *commons.NewImmutableList(wants)
	_ = agent.SendLootProposalToLeader(prop)
}

func (fiv *FivAgent) HandleLootRequest(m message.TaggedRequestMessage[message.LootRequest]) message.LootInform {
	return nil
}

func (fiv *FivAgent) HandleLootProposal(_ message.Proposal[decision.LootAction], _ agent.BaseAgent) decision.Intent {
	switch rand.Intn(3) {
	case 0:
		return decision.Positive
	case 1:
		return decision.Negative
	default:
		return decision.Abstain
	}
}

func (fiv *FivAgent) HandleLootProposalRequest(_ message.Proposal[decision.LootAction], _ agent.BaseAgent) bool {
	switch rand.Intn(2) {
	case 0:
		return true
	default:
		return false
	}
}

// --------------- Trade ---------------

func (fiv *FivAgent) HandleTradeNegotiation(_ agent.BaseAgent, _ message.TradeInfo) message.TradeMessage {
	return message.TradeRequest{}
}

// --------------- Hp Pool ---------------

func (fiv *FivAgent) DonateToHpPool(baseAgent agent.BaseAgent) uint {
	return 0
}

// --------------- Internal State ---------------

func (fiv *FivAgent) UpdateTrust(baseAgent agent.BaseAgent) {
	fiv.ttable.Decay()

	myview := baseAgent.View()
	globalStates := myview.AgentState()
	currentLeader := myview.CurrentLeader()

	// lead trust loss due to population death
	deathLastRound := fiv.prePopNum - uint(globalStates.Len())
	if deathLastRound > 0 {
		fiv.ttable.NegativeLeaderEvent(currentLeader, float32(deathLastRound))
	}
	// lead trust gain due to no health loss (being protected)
	hpLossLastRound := float32(fiv.preHealth) - float32(baseAgent.AgentState().Hp)
	if hpLossLastRound > 0 {
		fiv.ttable.NegativeLeaderEvent(currentLeader, -hpLossLastRound)
	}
	if hpLossLastRound <= 0 {
		fiv.ttable.PositiveLeaderEvent(currentLeader, -hpLossLastRound)
	}

	for _, id := range commons.ImmutableMapKeys(globalStates) {
		agState, _ := globalStates.Get(id)
		// individual trust loss due to low health (not performing well)
		if uint(agState.Hp) == state.LowHealth {
			fiv.ttable.NegativeIndivlEvent(id, 1)
		}
		// individual trust gain due to high health (performing well)
		if uint(agState.Hp) == state.HighHealth {
			fiv.ttable.PositiveIndivlEvent(id, 1)
		}
	}
}

func (fiv *FivAgent) UpdateInternalState(a agent.BaseAgent, _ *commons.ImmutableList[decision.ImmutableFightResult], _ *immutable.Map[decision.Intent, uint], log chan<- logging.AgentLog) {
	fiv.bravery += rand.Intn(10)
	fiv.UpdateQ(a)
	fiv.UpdateTrust(a)
	myview := a.View()
	log <- logging.AgentLog{
		Name: a.Name(),
		ID:   a.ID(),
		Properties: map[string]float32{
			"bravery":        float32(fiv.bravery),
			fiv.qtable.Log(): 0,
			"trustToLeader":  fiv.ttable.EstimateLeadTrust(myview.CurrentLeader()),
		},
	}
}

func NewFivAgent() agent.Strategy {
	return &FivAgent{bravery: rand.Intn(5), exploreRate: float32(0.25), qtable: NewQTable(0.25, 0.75), ttable: NewTrustTable()}
}
