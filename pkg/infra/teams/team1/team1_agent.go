package team1

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/message/proposal"
	"infra/game/state"
	"infra/logging"
	"infra/teams/team1/internal"
	"math/rand"
	"os"
	"strconv"

	"github.com/benbjohnson/immutable"
)

type SocialAgent struct {
	socialCapital map[string][4]float64 // agentID -> [Institutions, Networks, Trustworthiness, Honour]
	selfishness   float64               // Weighting of how selfish an agent is (0 -> not selfish, 1 -> very selfish)
	// Will gosip to all agents who's network value is above this
	gossipThreshold float64
	// Proportion of agents to talk badly about
	propHate float64
	// Proportion of agents to talk well about
	propAdmire float64
	// Propportion of agents to trade with
	propTrade float64

	graphID int // for logging
}

func (s *SocialAgent) FightResolution(
	agent agent.BaseAgent,
	prop commons.ImmutableList[proposal.Rule[decision.FightAction]],
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

func (s *SocialAgent) LootActionNoProposal(baseAgent agent.BaseAgent) immutable.SortedMap[commons.ItemID, struct{}] {
	return *immutable.NewSortedMap[commons.ItemID, struct{}](nil)
}

func (s *SocialAgent) LootAction(baseAgent agent.BaseAgent, proposedLoot immutable.SortedMap[commons.ItemID, struct{}], acceptedProposal message.Proposal[decision.LootAction]) immutable.SortedMap[commons.ItemID, struct{}] {
	return proposedLoot
}

func (s *SocialAgent) FightActionNoProposal(baseAgent agent.BaseAgent) decision.FightAction {
	qState := internal.BaseAgentToQState(baseAgent)

	// If we are training a Q function, maybe do an action other than the best action
	exploration := os.Getenv("EXPLORATION")
	if exploration != "" {
		epsilon, _ := strconv.ParseFloat(exploration, 64)

		if epsilon < rand.Float64() {
			// Do random action
			return decision.FightAction(rand.Intn(3))
		}
	}

	// Calculate best action based on current state and selfishness
	coopTable := internal.CooperationQ(qState)
	selfTable := internal.SelfishQ(qState)

	multipliedCoop := internal.ConstMulSlice(1.0-s.selfishness, coopTable[:])
	multipliedSelf := internal.ConstMulSlice(s.selfishness, selfTable[:])

	totalQSlice := internal.AddSlices(multipliedCoop, multipliedSelf)

	// Return index of best action (assumes array ordering in same order as decision.FightAction
	return decision.FightAction(internal.Argmax(totalQSlice))
}

func (s *SocialAgent) FightAction(baseAgent agent.BaseAgent, proposedAction decision.FightAction, acceptedProposal message.Proposal[decision.FightAction]) decision.FightAction {
	return s.FightActionNoProposal(baseAgent)
}

func (s *SocialAgent) HandleLootInformation(m message.TaggedInformMessage[message.LootInform], agent agent.BaseAgent) {
	//agent.AgentState().Hp
}

func (s *SocialAgent) HandleLootRequest(m message.TaggedRequestMessage[message.LootRequest]) message.LootInform {
	//TODO implement me
	panic("implement me")
}

func (s *SocialAgent) HandleLootProposal(_ message.Proposal[decision.LootAction], _ agent.BaseAgent) decision.Intent {
	switch rand.Intn(3) {
	case 0:
		return decision.Positive
	case 1:
		return decision.Negative
	default:
		return decision.Abstain
	}
}

func (s *SocialAgent) HandleLootProposalRequest(_ message.Proposal[decision.LootAction], _ agent.BaseAgent) bool {
	switch rand.Intn(2) {
	case 0:
		return true
	default:
		return false
	}
}

func (s *SocialAgent) LootAllocation(
	ba agent.BaseAgent,
	proposal message.Proposal[decision.LootAction],
	proposedAllocations immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]],
) immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]] {
	lootAllocation := make(map[commons.ID][]commons.ItemID)
	view := ba.View()
	ids := commons.ImmutableMapKeys(view.AgentState())
	iterator := ba.Loot().Weapons().Iterator()
	allocateRandomly(iterator, ids, lootAllocation)
	iterator = ba.Loot().Shields().Iterator()
	allocateRandomly(iterator, ids, lootAllocation)
	iterator = ba.Loot().HpPotions().Iterator()
	allocateRandomly(iterator, ids, lootAllocation)
	iterator = ba.Loot().StaminaPotions().Iterator()
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

func (s *SocialAgent) DonateToHpPool(baseAgent agent.BaseAgent) uint {
	return 0
}

func (s *SocialAgent) UpdateSelfishness(agent agent.BaseAgent) {
	// Find utility of agents own state
	selfUtility := internal.UtilityOfState(agent.AgentState())

	// Extract view, agentStates
	view := agent.View()
	agentState := view.AgentState()

	// Find list of all agents with higher utility than oneself
	var betterAgents []string
	itr := agentState.Iterator()
	for !itr.Done() {
		agentID, hiddenState, _ := itr.Next()

		if internal.UtilityOfHiddenState(hiddenState) > selfUtility {
			betterAgents = append(betterAgents, agentID)
		}
	}

	// If no agents have a higher utility than you, do nothing
	if len(betterAgents) == 0 {
		return
	}

	// Calculate average trustworthiness of agents with better state
	totalTrustworthiness := 0.0

	for _, agentID := range betterAgents {
		totalTrustworthiness += s.socialCapital[agentID][3]
	}

	averageTrustworthiness := totalTrustworthiness / float64(len(betterAgents))

	// If agent with better state than oneself has higher trustworthiness
	if averageTrustworthiness > s.socialCapital[agent.ID()][2] {
		s.selfishness -= 0.01
	} else { // If agents with better state than oneself has lower trustworthiness
		s.selfishness += 0.01
	}
}

func (s *SocialAgent) UpdateInternalState(self agent.BaseAgent, fightResult *commons.ImmutableList[decision.ImmutableFightResult], _ *immutable.Map[decision.Intent, uint], _ chan<- logging.AgentLog) {
	// Update socialCapital at end of each round
	itr := fightResult.Iterator()
	for !itr.Done() { // For each fight round
		fightDecisions, _ := itr.Next()

		s.updateSocialCapital(self, fightDecisions)
	}

	s.UpdateSelfishness(self)
}

func (s *SocialAgent) CreateManifesto(_ agent.BaseAgent) *decision.Manifesto {
	manifesto := decision.NewManifesto(false, true, 10, 50)
	return manifesto
}

func (s *SocialAgent) HandleConfidencePoll(_ agent.BaseAgent) decision.Intent {
	switch rand.Intn(3) {
	case 0:
		return decision.Abstain
	case 1:
		return decision.Negative
	default:
		return decision.Positive
	}
}

func (s *SocialAgent) HandleFightInformation(m message.TaggedInformMessage[message.FightInform], baseAgent agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) {
	// baseAgent.Log(logging.Trace, logging.LogField{"bravery": r.bravery, "hp": baseAgent.AgentState().Hp}, "Cowering")
	switch m.Message().(type) {
	case *message.StartFight:
		s.sendGossip(baseAgent)
	case message.ArrayInfo:
		s.receiveGossip(m.Message().(message.ArrayInfo), m.Sender())
	}
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

func (s *SocialAgent) HandleFightRequest(_ message.TaggedRequestMessage[message.FightRequest], _ *immutable.Map[commons.ID, decision.FightAction]) message.FightInform {
	return nil
}

func (s *SocialAgent) HandleElectionBallot(b agent.BaseAgent, _ *decision.ElectionParams) decision.Ballot {
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

func (s *SocialAgent) HandleFightProposal(_ message.Proposal[decision.FightAction], _ agent.BaseAgent) decision.Intent {
	intent := rand.Intn(2)
	if intent == 0 {
		return decision.Positive
	} else {
		return decision.Negative
	}
}

func (s *SocialAgent) HandleFightProposalRequest(
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

func (s *SocialAgent) HandleUpdateWeapon(_ agent.BaseAgent) decision.ItemIdx {
	// weapons := b.AgentState().weapons
	// return decision.ItemIdx(rand.Intn(weapons.Len() + 1))

	// 0th weapon has the greatest attack points
	return decision.ItemIdx(0)
}

func (s *SocialAgent) HandleUpdateShield(_ agent.BaseAgent) decision.ItemIdx {
	// shields := b.AgentState().Shields
	// return decision.ItemIdx(rand.Intn(shields.Len() + 1))
	return decision.ItemIdx(0)
}

func (s *SocialAgent) HandleTradeNegotiation(BA agent.BaseAgent, m message.TradeInfo) message.TradeMessage {

	agentState := BA.AgentState()

	bestWeaponDonation := uint(0)
	bestShieldDonation := uint(0)
	var bestWeaponDonationID string
	var bestShieldDonationID string
	for negId, neg := range m.Negotiations {
		//fmt.Println(neg.Agent2)
		if neg.Agent2 == BA.ID() {
			//fmt.Println("offer made to me")
			if offer, ok := neg.GetOffer(neg.Agent1); ok {
				if offer.ItemType == commons.Weapon && offer.Item.Value() > bestWeaponDonation {
					bestWeaponDonation = offer.Item.Value()
					bestWeaponDonationID = negId
				} else if offer.ItemType == commons.Shield && offer.Item.Value() > bestShieldDonation {
					bestShieldDonation = offer.Item.Value()
					bestShieldDonationID = negId
				}
			}
		}
	}

	if bestWeaponDonation > agentState.BonusAttack() {
		//fmt.Println("Accepted weapons offer")
		return message.TradeAccept{TradeID: bestWeaponDonationID}
	} else if bestShieldDonation > agentState.BonusDefense() {
		//fmt.Println("Accepted shield offer")
		return message.TradeAccept{TradeID: bestShieldDonationID}
	}

	if agentState.Weapons.Len() < 2 && agentState.Shields.Len() < 2 { // Cant trade due to no next best weapon
		return message.TradeRequest{}
	}

	sortedSC := internal.GetSortedAgentSubset(BA.ID(), s.socialCapital)

	// check what the second best weapon held is
	bestFreeWStats := uint(0)
	bestFreeWIdx := int(-1)
	it := agentState.Weapons.Iterator()
	for !it.Done() {
		i, w := it.Next()
		if w.Id() != agentState.WeaponInUse {
			if bestFreeWStats < w.Value() {
				bestFreeWStats = w.Value()
				bestFreeWIdx = i
			}
		}
	}

	bestFreeSStats := uint(0)
	bestFreeSIdx := int(-1)
	it = agentState.Shields.Iterator()
	for !it.Done() {
		i, w := it.Next()
		if w.Id() != agentState.ShieldInUse {
			if bestFreeSStats < w.Value() {
				bestFreeSStats = w.Value()
				bestFreeSIdx = i
			}
		}
	}

	if bestFreeWIdx == -1 && bestFreeSIdx == -1 {
		//fmt.Println(agentState.Weapons.Len(), agentState.Shields.Len())
		return message.TradeRequest{}
	}

nextAgent:
	for _, sci := range sortedSC {
		if internal.OverallPerception(sci.Arr) < 0 {
			break
		}
		// check if a trade negotiation is in place with that agent
		for _, neg := range m.Negotiations {
			if neg.Agent1 == BA.ID() && neg.Agent2 == sci.ID {
				continue nextAgent
			}
		}

		//fmt.Println(sci.ID)
		if agentState.Weapons.Len() > agentState.Shields.Len() {
			//fmt.Println("Offered weapon")
			TO, _ := message.NewTradeOffer(commons.Weapon, uint(bestFreeWIdx), agentState.Weapons, agentState.Shields)
			TD := message.NewTradeDemand(commons.Shield, 0)
			return message.TradeRequest{CounterPartyID: sci.ID, Offer: TO, Demand: TD}
		} else {
			//fmt.Println("Offered shield")
			TO, _ := message.NewTradeOffer(commons.Shield, uint(bestFreeSIdx), agentState.Weapons, agentState.Shields)
			TD := message.NewTradeDemand(commons.Shield, 0)
			return message.TradeRequest{CounterPartyID: sci.ID, Offer: TO, Demand: TD}
		}
	}

	return message.TradeRequest{}
}

func NewSocialAgent() agent.Strategy {
	return &SocialAgent{
		selfishness:     rand.Float64(),
		gossipThreshold: 0.5,
		propAdmire:      0.1,
		propHate:        0.1,
	}
}
