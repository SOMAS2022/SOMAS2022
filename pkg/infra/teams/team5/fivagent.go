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

	"github.com/benbjohnson/immutable"
)

type FivAgent struct {
	bravery     int
	preHealth   uint
	exploreRate float32
	qtable      *Qtable
}

func (fiv *FivAgent) FightResolution(agent agent.BaseAgent, prop commons.ImmutableList[proposal.Rule[decision.FightAction]]) immutable.Map[commons.ID, decision.FightAction] {
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

func (fiv *FivAgent) LootActionNoProposal(baseAgent agent.BaseAgent) immutable.SortedMap[commons.ItemID, struct{}] {
	return *immutable.NewSortedMap[commons.ItemID, struct{}](nil)
}

func (fiv *FivAgent) LootAction(
	baseAgent agent.BaseAgent,
	proposedLoot immutable.SortedMap[commons.ItemID, struct{}],
	acceptedProposal message.Proposal[decision.LootAction],
) immutable.SortedMap[commons.ItemID, struct{}] {
	return proposedLoot
}

func (fiv *FivAgent) CurrentQState(baseAgent agent.BaseAgent) string {
	initHealth := 1000.0
	initStamina := 2000.0

	mystate := baseAgent.AgentState()
	myview := baseAgent.View()

	myHealth := ""
	myStamina := ""
	relativeAT := ""
	relativeSH := ""

	switch {
	case mystate.Hp < uint(0.3*initHealth):
		myHealth = "Low"
	case uint(0.3*initHealth) <= mystate.Hp && mystate.Hp < uint(0.6*initHealth):
		myHealth = "Mid"
	case uint(0.6*initHealth) <= mystate.Hp:
		myHealth = "Hih"
	}

	switch {
	case mystate.Stamina < uint(0.3*initStamina):
		myStamina = "Low"
	case uint(0.3*initStamina) <= mystate.Stamina && mystate.Stamina < uint(0.6*initStamina):
		myStamina = "Mid"
	case uint(0.6*initStamina) <= mystate.Stamina:
		myStamina = "Hih"
	}

	numAlive := 0.0
	popATGreaterToCount := 0.0
	popSHGreaterToCount := 0.0
	othersStates := myview.AgentState()
	itr := othersStates.Iterator()
	for !itr.Done() {
		_, agState, ok := itr.Next()
		if ok && agState.Hp > 0 {
			numAlive += 1
			if agState.Attack+agState.BonusAttack < mystate.TotalAttack() {
				popATGreaterToCount += 1
			}
			if agState.Defense+agState.BonusDefense < mystate.TotalDefense() {
				popSHGreaterToCount += 1
			}
		}
	}

	switch {
	case popATGreaterToCount < 0.25*numAlive:
		relativeAT = "Weakee"
	case 0.25*numAlive <= popATGreaterToCount && popATGreaterToCount < 0.75*numAlive:
		relativeAT = "Ordin"
	case 0.75 <= popATGreaterToCount && popATGreaterToCount <= numAlive:
		relativeAT = "Master"
	}

	switch {
	case popSHGreaterToCount < 0.25*numAlive:
		relativeSH = "Weakee"
	case 0.25*numAlive <= popSHGreaterToCount && popSHGreaterToCount < 0.75*numAlive:
		relativeSH = "Ordin"
	case 0.75 <= popSHGreaterToCount && popSHGreaterToCount <= numAlive:
		relativeSH = "Master"
	}

	return myHealth + "-" + myStamina + "-" + relativeAT + "-" + relativeSH
}

func (fiv *FivAgent) Explore(qstate string) decision.FightAction {
	var sa SaPair
	var fightDecision decision.FightAction
	fight := rand.Intn(3)
	switch fight {
	case 0:
		sa = SaPair{state: qstate, action: "Cower"}
		fightDecision = decision.Cower
	case 1:
		sa = SaPair{state: qstate, action: "Attck"}
		fightDecision = decision.Attack
	default:
		sa = SaPair{state: qstate, action: "Defnd"}
		fightDecision = decision.Defend
	}
	_, exist := fiv.qtable.table[sa]
	if !exist {
		fiv.qtable.table[sa] = 0
	}
	fiv.qtable.saTaken = sa
	return fightDecision
}

func (fiv *FivAgent) Exploit(qstate string) decision.FightAction {
	maxQAction := fiv.qtable.GetMaxQAction(qstate)
	var sa SaPair
	var fightDecision decision.FightAction
	switch maxQAction {
	case "NoSaPairAvailable":
		return fiv.Explore(qstate)
	case "Cower":
		sa = SaPair{state: qstate, action: "Cower"}
		fightDecision = decision.Cower
	case "Attck":
		sa = SaPair{state: qstate, action: "Attck"}
		fightDecision = decision.Attack
	case "Defnd":
		sa = SaPair{state: qstate, action: "Defnd"}
		fightDecision = decision.Defend
	}
	fiv.qtable.saTaken = sa
	return fightDecision
}

func (fiv *FivAgent) FightActionNoProposal(baseAgent agent.BaseAgent) decision.FightAction {
	fiv.preHealth = baseAgent.AgentState().Hp
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

func (fiv *FivAgent) HandleLootInformation(m message.TaggedInformMessage[message.LootInform], agent agent.BaseAgent) {
}

func (fiv *FivAgent) HandleLootRequest(m message.TaggedRequestMessage[message.LootRequest]) message.LootInform {
	//TODO implement me
	panic("implement me")
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

func (fiv *FivAgent) LootAllocation(baseAgent agent.BaseAgent, proposal message.Proposal[decision.LootAction]) immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]] {
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

func (fiv *FivAgent) DonateToHpPool(baseAgent agent.BaseAgent) uint {
	return 0
}

func (fiv *FivAgent) UpdateInternalState(a agent.BaseAgent, _ *commons.ImmutableList[decision.ImmutableFightResult], _ *immutable.Map[decision.Intent, uint], log chan<- logging.AgentLog) {
	fiv.bravery += rand.Intn(10)
	log <- logging.AgentLog{
		Name: a.Name(),
		ID:   a.ID(),
		Properties: map[string]float32{
			"bravery": float32(fiv.bravery),
		},
	}
	percentHealthLoss := (float32(a.AgentState().Hp) - float32(fiv.preHealth)) / float32(fiv.preHealth) * 100
	cqState := fiv.CurrentQState(a)
	fSas := []SaPair{{state: cqState, action: "Cower"}, {state: cqState, action: "Attck"}, {state: cqState, action: "Defnd"}}
	fiv.qtable.Learn(percentHealthLoss, fiv.qtable.GetMaxFR(fSas))
	fiv.qtable.Print()
}

func (fiv *FivAgent) CreateManifesto(_ agent.BaseAgent) *decision.Manifesto {
	manifesto := decision.NewManifesto(false, false, 10, 5)
	return manifesto
}

func (fiv *FivAgent) HandleConfidencePoll(_ agent.BaseAgent) decision.Intent {
	switch rand.Intn(3) {
	case 0:
		return decision.Abstain
	case 1:
		return decision.Negative
	default:
		return decision.Positive
	}
}

func (fiv *FivAgent) HandleFightInformation(_ message.TaggedInformMessage[message.FightInform], baseAgent agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) {
	// baseAgent.Log(logging.Trace, logging.LogField{"bravery": r.bravery, "hp": baseAgent.AgentState().Hp}, "Cowering")
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

func (fiv *FivAgent) HandleFightRequest(_ message.TaggedRequestMessage[message.FightRequest], _ *immutable.Map[commons.ID, decision.FightAction]) message.FightInform {
	return nil
}

func (fiv *FivAgent) HandleElectionBallot(b agent.BaseAgent, _ *decision.ElectionParams) decision.Ballot {
	// Extract ID of alive agents
	view := b.View()
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

	// Randomly fill the ballot
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

func (fiv *FivAgent) HandleFightProposal(_ message.Proposal[decision.FightAction], _ agent.BaseAgent) decision.Intent {
	intent := rand.Intn(2)
	if intent == 0 {
		return decision.Positive
	} else {
		return decision.Negative
	}
}

func (fiv *FivAgent) HandleFightProposalRequest(
	_ message.Proposal[decision.FightAction],
	_ agent.BaseAgent,
	_ *immutable.Map[commons.ID, decision.FightAction],
) bool {
	switch rand.Intn(2) {
	case 0:
		return true
	default:
		return false
	}
}

func (fiv *FivAgent) HandleUpdateWeapon(_ agent.BaseAgent) decision.ItemIdx {
	// weapons := b.AgentState().weapons
	// return decision.ItemIdx(rand.Intn(weapons.Len() + 1))

	// 0th weapon has the greatest attack points
	return decision.ItemIdx(0)
}

func (fiv *FivAgent) HandleUpdateShield(_ agent.BaseAgent) decision.ItemIdx {
	// shields := b.AgentState().Shields
	// return decision.ItemIdx(rand.Intn(shields.Len() + 1))
	return decision.ItemIdx(0)
}

func (fiv *FivAgent) HandleTradeNegotiation(_ agent.BaseAgent, _ message.TradeInfo) message.TradeMessage {
	return message.TradeRequest{}
}

func NewFivAgent() agent.Strategy {
	return &FivAgent{bravery: rand.Intn(5), exploreRate: float32(0.25), qtable: NewQTable(0.2, 0.8)}
}
