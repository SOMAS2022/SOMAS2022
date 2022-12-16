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
	"math"
	"math/rand"
	"os"
	"sort"
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

	proposalAccuracyThreshold float64

	//
	socialCapitalMean map[string]float64

	// helper for agent accuracy
	currentProposalAccuracyThreshold float64
	hasVotedThisRound                bool
	votedOnFirstRound                bool
	isFirstRound                     bool
}

func (s *SocialAgent) FightResolution(
	agent agent.BaseAgent,
	prop commons.ImmutableList[proposal.Rule[decision.FightAction]],
	proposedActions immutable.Map[commons.ID, decision.FightAction],
) immutable.Map[commons.ID, decision.FightAction] {
	view := agent.View()
	agents := view.AgentState()
	builder := immutable.NewMapBuilder[commons.ID, decision.FightAction](nil)

	// find the percentile social wellbeing of an agent
	type AgentCapital struct {
		capital float64
		id      string
	}
	capitals := make([]AgentCapital, 0)
	for _, id := range commons.ImmutableMapKeys(view.AgentState()) {
		capitals = append(capitals, AgentCapital{
			capital: s.socialCapitalMean[id],
			id:      id,
		})
	}
	sort.SliceStable(capitals, func(i, j int) bool {
		return capitals[i].capital < capitals[j].capital
	})

	for index := 0; index < len(capitals); index++ {
		agent_state, _ := agents.Get(capitals[index].id)
		qState := internal.HiddenAgentToQState(agent_state, view)
		coop_rewards := internal.CooperationQ(qState)
		coop_q_action := decision.FightAction(internal.Argmax(coop_rewards[:]))
		selfish_rewards := internal.SelfishQ(qState)
		selfish_q_action := decision.FightAction(internal.Argmax(selfish_rewards[:]))

		use_selfish := rand.Float64() < math.Pow(float64(index)/float64(len(capitals)), 8)
		var action decision.FightAction
		if use_selfish {
			action = selfish_q_action
		} else {
			action = coop_q_action
		}

		// add a degree of complete randomness
		if rand.Float64() > 0.9 {
			switch rand.Intn(3) {
			case 0:
				action = decision.Attack
			case 1:
				action = decision.Defend
			default:
				action = decision.Cower
			}
		}

		builder.Set(capitals[index].id, action)
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
	qState := internal.BaseAgentToQState(baseAgent)
	rewards_coop := internal.CooperationQ(qState)
	rewards_self := internal.SelfishQ(qState)
	multipliedCoop := internal.ConstMulSlice(1.0-s.selfishness, rewards_coop[:])
	multipliedSelf := internal.ConstMulSlice(s.selfishness, rewards_self[:])
	totalQSlice := internal.AddSlices(multipliedCoop, multipliedSelf)
	desired_action := decision.FightAction(internal.Argmax(totalQSlice))
	if desired_action == proposedAction {
		return desired_action
	}
	max := math.Max(math.Max(totalQSlice[0], totalQSlice[1]), totalQSlice[2])
	min := math.Min(math.Min(totalQSlice[0], totalQSlice[1]), totalQSlice[2])
	diff := max - min
	avg := (max + min) / 2.0
	totalQSlice = []float64{
		(totalQSlice[0] - avg) / diff,
		(totalQSlice[1] - avg) / diff,
		(totalQSlice[2] - avg) / diff,
	}

	if totalQSlice[proposedAction]+3*(1-s.selfishness) > 1 {
		return proposedAction
	}
	return desired_action
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

func SampleDistribution(distribution []float64) int {
	random := rand.Float64() * distribution[len(distribution)-1]
	// Add max iteration
	change := len(distribution) / 4
	index := len(distribution) / 2
	for count := 0; count < 100; count++ {
		if index == len(distribution)-1 || index == 0 {
			return index
		}
		if index >= 1 {
			if distribution[index] > random && distribution[index-1] < random {
				return index
			} else if distribution[index] > random {
				index -= change
			} else {
				index += change
			}
		}
		change /= 2
		if change < 1 {
			change = 1
		}
	}
	return index
}

func AllocateWithProbabilityDistribution(distribution []float64, iterator commons.Iterator[state.Item], ids []commons.ID, lootAllocation map[commons.ID][]commons.ItemID) {
	for !iterator.Done() {
		item, _ := iterator.Next()
		toBeAllocated := ids[SampleDistribution(distribution)]
		if l, ok := lootAllocation[toBeAllocated]; ok {
			l = append(l, item.Id())
			lootAllocation[toBeAllocated] = l
		} else {
			l := make([]commons.ItemID, 0)
			l = append(l, item.Id())
			lootAllocation[toBeAllocated] = l
		}
	}
}

func (s *SocialAgent) FindMaxStats(baseAgent agent.BaseAgent) struct {
	MaxAttack  float64
	MaxHealth  float64
	MaxStamina float64
	MaxDefense float64
	MaxSocial  float64
} {
	view := baseAgent.View()
	ids := commons.ImmutableMapKeys(view.AgentState())
	agents := view.AgentState()
	var max_attack uint = 0
	var max_health state.HealthRange = 0
	var max_defense uint = 0
	var max_stamina state.StaminaRange = 0
	max_social := 0.0

	// find the maximum values
	for id_index := range ids {
		agent_state, _ := agents.Get(ids[id_index])
		if agent_state.Attack > max_attack {
			max_attack = agent_state.Attack
		}
		if agent_state.Hp > max_health {
			max_health = agent_state.Hp
		}
		if agent_state.Defense > max_defense {
			max_defense = agent_state.Defense
		}
		if agent_state.Stamina > max_stamina {
			max_stamina = agent_state.Stamina
		}
		if s.socialCapitalMean[ids[id_index]] > max_social {
			max_social = s.socialCapitalMean[ids[id_index]]
		}
	}

	return struct {
		MaxAttack  float64
		MaxHealth  float64
		MaxStamina float64
		MaxDefense float64
		MaxSocial  float64
	}{
		MaxAttack:  float64(max_attack),
		MaxDefense: float64(max_defense),
		MaxHealth:  float64(max_health),
		MaxStamina: float64(max_stamina),
		MaxSocial:  max_social,
	}
}

func (s *SocialAgent) LootAllocation(
	baseAgent agent.BaseAgent,
	proposal message.Proposal[decision.LootAction],
	proposedAllocations immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]],
) immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]] {
	lootAllocation := make(map[commons.ID][]commons.ItemID)
	view := baseAgent.View()
	ids := commons.ImmutableMapKeys(view.AgentState())
	agents := view.AgentState()
	max_stats := s.FindMaxStats(baseAgent)

	var weapon_cumulative_prop []float64
	last_weapon_prop := 0.0
	var defense_cumulative_prob []float64
	last_defense_prop := 0.0
	var hp_cumulative_prob []float64
	last_hp_prop := 0.0
	var stamina_cumulative_prob []float64
	last_stamina_prob := 0.0

	// find cumulative probabilities of receiving different loot types
	for id_index := range ids {
		weapon_prob := 0.0
		defense_prob := 0.0
		hp_prob := 0.0
		stamina_prob := 0.0

		agent_state, _ := agents.Get(ids[id_index])
		weapon_prob += float64(agent_state.Hp)/max_stats.MaxHealth*0.9 + 9*float64(agent_state.Stamina)/max_stats.MaxStamina
		weapon_prob -= float64(agent_state.Attack) / 50
		weapon_prob -= float64(agent_state.Defense) / 100
		if weapon_prob < 0.0 {
			weapon_prob = 0.0
		}
		defense_prob += max_stats.MaxStamina / (float64(agent_state.Stamina) + 0.1)
		// logging.Log(logging.Error, nil, "---------")
		// logging.Log(logging.Error, nil, ids[id_index])
		defense_prob -= float64(agent_state.Defense) / 50
		defense_prob -= float64(agent_state.Attack) / 100
		// logging.Log(logging.Error, nil, fmt.Sprintf("%f", defense_prob))
		if defense_prob < 0.0 {
			defense_prob = 0.0
		}
		hp_prob += max_stats.MaxHealth/(float64(agent_state.Hp)+0.1) + s.socialCapitalMean[ids[id_index]]
		stamina_prob += max_stats.MaxStamina/(float64(agent_state.Stamina)+0.1) + float64(agent_state.Attack)/100

		weapon_cumulative_prop = append(weapon_cumulative_prop, weapon_prob+last_weapon_prop)
		defense_cumulative_prob = append(defense_cumulative_prob, defense_prob+last_defense_prop)
		hp_cumulative_prob = append(hp_cumulative_prob, hp_prob+last_hp_prop)
		stamina_cumulative_prob = append(stamina_cumulative_prob, last_stamina_prob+stamina_prob)
		last_weapon_prop = weapon_prob + last_weapon_prop
		last_defense_prop = defense_prob + last_defense_prop
		last_hp_prop = hp_prob + last_hp_prop
		last_stamina_prob = stamina_prob + last_stamina_prob
	}

	// distribute according to the cumulative prob distributions
	iterator := baseAgent.Loot().Weapons().Iterator()
	AllocateWithProbabilityDistribution(weapon_cumulative_prop, iterator, ids, lootAllocation)
	iterator = baseAgent.Loot().Shields().Iterator()
	AllocateWithProbabilityDistribution(defense_cumulative_prob, iterator, ids, lootAllocation)
	iterator = baseAgent.Loot().HpPotions().Iterator()
	AllocateWithProbabilityDistribution(hp_cumulative_prob, iterator, ids, lootAllocation)
	iterator = baseAgent.Loot().StaminaPotions().Iterator()
	AllocateWithProbabilityDistribution(stamina_cumulative_prob, iterator, ids, lootAllocation)

	mMapped := make(map[commons.ID]immutable.SortedMap[commons.ItemID, struct{}])
	for id, itemIDS := range lootAllocation {
		mMapped[id] = commons.ListToImmutableSortedSet(itemIDS)
	}
	return commons.MapToImmutable(mMapped)
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
	itr := fightResult.Iterator()
	for !itr.Done() { // For each fight round
		fightDecisions, _ := itr.Next()

		s.updateSocialCapital(self, fightDecisions)
	}

	// proposal voting reset
	if !s.hasVotedThisRound {
		s.proposalAccuracyThreshold *= 0.9
	}
	if s.votedOnFirstRound {
		s.proposalAccuracyThreshold *= 1.1
	}
	s.currentProposalAccuracyThreshold = s.proposalAccuracyThreshold
	s.hasVotedThisRound = false
	s.votedOnFirstRound = false
	s.isFirstRound = true

	// social capital mean
	var meanCapital = make(map[string]float64)

	for id, element := range s.socialCapital {
		meanCapital[id] = (element[0] + element[1] + element[2] + element[3]) / 4.0
	}
	s.socialCapitalMean = meanCapital

	s.UpdateSelfishness(self)
}

func (s *SocialAgent) CreateManifesto(_ agent.BaseAgent) *decision.Manifesto {
	manifesto := decision.NewManifesto(true, true, uint(rand.Intn(20)+5), uint(rand.Intn(30)+20))
	return manifesto
}

func (s *SocialAgent) HandleConfidencePoll(baseAgent agent.BaseAgent) decision.Intent {
	view := baseAgent.View()
	id := view.CurrentLeader()
	if rand.Float64() < math.Pow((s.socialCapitalMean[id]+1)/2.0, 2) {
		return decision.Negative
	}
	return decision.Positive
}

func (s *SocialAgent) HandleFightInformation(m message.TaggedInformMessage[message.FightInform], baseAgent agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) {
	switch m.Message().(type) {
	case *message.StartFight:
		prop := *commons.NewImmutableList(s.CreateFightProposal(baseAgent))
		_ = baseAgent.SendFightProposalToLeader(prop)
		s.sendGossip(baseAgent)
	case message.ArrayInfo:
		s.receiveGossip(m.Message().(message.ArrayInfo), m.Sender())
	}
}

func (s *SocialAgent) HandleFightRequest(_ message.TaggedRequestMessage[message.FightRequest], _ *immutable.Map[commons.ID, decision.FightAction]) message.FightInform {
	return nil
}

func (s *SocialAgent) HandleElectionBallot(b agent.BaseAgent, electionParams *decision.ElectionParams) decision.Ballot {
	var ballot decision.Ballot
	candidates := electionParams.CandidateList().Iterator()
	for !candidates.Done() {
		id, _, _ := candidates.Next()
		if rand.Float64() < (s.socialCapitalMean[id]+1)/2.0 {
			ballot = append(ballot, id)
		}
	}
	ballot = append(ballot, b.ID())
	return ballot
}

func (s *SocialAgent) CreateFightProposal(baseAgent agent.BaseAgent) []proposal.Rule[decision.FightAction] {
	// find the action each agent will make
	// check what the average of each state type is
	view := baseAgent.View()
	ids := commons.ImmutableMapKeys(view.AgentState())
	agents := view.AgentState()
	average_attack := 0.0
	average_defense := 0.0
	for id_index := 0; id_index < len(ids); id_index++ {
		agent_state, _ := agents.Get(ids[id_index])
		average_attack += float64(agent_state.Attack)
		average_defense += float64(agent_state.Defense)
	}
	average_attack /= float64(len(ids))
	average_defense /= float64(len(ids))
	halved_attack_average := average_attack / 2.0
	halved_defend_average := average_defense / 2.0
	// construct rules based on these agent averages, 36 different rules each corresponding to a range of the possible state space
	rules := make([]proposal.Rule[decision.FightAction], 0)
	for health_range := 1; health_range <= 3; health_range++ {
		health_val_min := 250 * health_range
		for stamina_range := 1; stamina_range <= 3; stamina_range++ {
			stamina_val_min := 500 * stamina_range
			for attack_quartile := 1.0; attack_quartile < 4.0; attack_quartile += 2.0 {
				attack_mid := halved_attack_average * attack_quartile
				for defend_quartile := 1.0; defend_quartile < 4.0; defend_quartile += 2.0 {
					defend_mid := halved_defend_average * defend_quartile

					// find what an agent with the current stats would do using the q function
					qState := internal.HiddenAgentToQState(state.HiddenAgentState{
						Hp:      state.HealthRange(health_val_min),
						Stamina: state.StaminaRange(stamina_val_min),
						Attack:  uint(attack_mid),
						Defense: uint(defend_mid),
					}, view)
					rewards := internal.CooperationQ(qState)
					q_action := decision.FightAction(internal.Argmax(rewards[:]))

					// make a rule that implements this q_action, four sets of ranges
					rules = append(rules, *proposal.NewRule(q_action, proposal.NewAndCondition(
						proposal.NewAndCondition( //
							proposal.NewAndCondition(
								proposal.NewComparativeCondition(proposal.TotalAttack, proposal.GreaterThan, proposal.Value(attack_mid-halved_attack_average)),
								proposal.NewComparativeCondition(proposal.TotalAttack, proposal.LessThan, proposal.Value(attack_mid+halved_attack_average)),
							), // attack
							proposal.NewAndCondition(
								proposal.NewComparativeCondition(proposal.TotalDefence, proposal.GreaterThan, proposal.Value(defend_mid-halved_defend_average)),
								proposal.NewComparativeCondition(proposal.TotalDefence, proposal.LessThan, proposal.Value(defend_mid+halved_defend_average)),
							), // defense
						),
						proposal.NewAndCondition(
							proposal.NewAndCondition(
								proposal.NewComparativeCondition(proposal.Health, proposal.GreaterThan, proposal.Value(health_val_min)),
								proposal.NewComparativeCondition(proposal.Health, proposal.LessThan, proposal.Value(health_val_min+250)),
							), // health
							proposal.NewAndCondition(
								proposal.NewComparativeCondition(proposal.Health, proposal.GreaterThan, proposal.Value(stamina_val_min)),
								proposal.NewComparativeCondition(proposal.Health, proposal.LessThan, proposal.Value(stamina_val_min+500)),
							), // stamina
						),
					)))
				}
			}
		}
	}
	return rules
}

func (s *SocialAgent) HandleFightProposal(prop message.Proposal[decision.FightAction], baseAgent agent.BaseAgent) decision.Intent {
	var result decision.Intent
	view := baseAgent.View()
	ids := commons.ImmutableMapKeys(view.AgentState())
	agents := view.AgentState()
	rules := prop.Rules()
	action_checker := proposal.ToSinglePredicate(rules)
	accuracy := 0.0
	for id_index := 0; id_index < 20; id_index++ {
		id_index := rand.Intn(len(ids))
		agent_state, _ := agents.Get(ids[id_index])
		proposal_action := action_checker(state.AgentState{
			Hp:      uint(agent_state.Hp),
			Stamina: uint(agent_state.Stamina),
			Attack:  agent_state.Attack,
			Defense: agent_state.Defense,
		})
		qState := internal.HiddenAgentToQState(agent_state, view)
		rewards := internal.CooperationQ(qState)
		q_action := decision.FightAction(internal.Argmax(rewards[:]))
		decision_match := q_action == proposal_action
		if decision_match {
			accuracy += 1.0
		}
	}
	// accuracy /= float64(len(ids))
	accuracy /= 20.0

	if accuracy > s.currentProposalAccuracyThreshold {
		result = decision.Positive
		s.hasVotedThisRound = true
		if s.isFirstRound {
			s.votedOnFirstRound = true
		}
	} else {
		result = decision.Negative
		s.currentProposalAccuracyThreshold *= 0.9
	}

	s.isFirstRound = false
	return result
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

	bestWeaponDonation, bestShieldDonation, bestWeaponDonationID, bestShieldDonationID := GetBestTrades(BA, m)

	tradeMessage, done := AcceptOffers(bestWeaponDonation, agentState, bestWeaponDonationID, bestShieldDonation, bestShieldDonationID)
	if done {
		return tradeMessage
	}

	if agentState.Weapons.Len() < 2 && agentState.Shields.Len() < 2 { // Cant trade due to no next best weapon
		return message.TradeRequest{}
	}

	sortedSC := internal.GetSortedAgentSubset(BA.ID(), s.socialCapital)

	bestFreeWIdx, bestFreeSIdx := GetSecondBestEquipment(agentState)

	if bestFreeWIdx == -1 && bestFreeSIdx == -1 {
		//fmt.Println(agentState.Weapons.Len(), agentState.Shields.Len())
		return message.TradeRequest{}
	}

	m2, done2 := ProposeTrade(BA, m, sortedSC, agentState, bestFreeWIdx, bestFreeSIdx)
	if done2 {
		return m2
	}

	return message.TradeRequest{}
}

func ProposeTrade(BA agent.BaseAgent, m message.TradeInfo, sortedSC []internal.SocialCapInfo, agentState state.AgentState, bestFreeWIdx int, bestFreeSIdx int) (message.TradeMessage, bool) {
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
			return message.TradeRequest{CounterPartyID: sci.ID, Offer: TO, Demand: TD}, true
		} else {
			//fmt.Println("Offered shield")
			TO, _ := message.NewTradeOffer(commons.Shield, uint(bestFreeSIdx), agentState.Weapons, agentState.Shields)
			TD := message.NewTradeDemand(commons.Shield, 0)
			return message.TradeRequest{CounterPartyID: sci.ID, Offer: TO, Demand: TD}, true
		}
	}
	return nil, false
}

func AcceptOffers(bestWeaponDonation uint, agentState state.AgentState, bestWeaponDonationID string, bestShieldDonation uint, bestShieldDonationID string) (message.TradeMessage, bool) {
	if bestWeaponDonation > agentState.BonusAttack() {
		//fmt.Println("Accepted weapons offer")
		return message.TradeAccept{TradeID: bestWeaponDonationID}, true
	} else if bestShieldDonation > agentState.BonusDefense() {
		//fmt.Println("Accepted shield offer")
		return message.TradeAccept{TradeID: bestShieldDonationID}, true
	}
	return nil, false
}

func GetSecondBestEquipment(agentState state.AgentState) (int, int) {
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
	return bestFreeWIdx, bestFreeSIdx
}

func GetBestTrades(BA agent.BaseAgent, m message.TradeInfo) (uint, uint, string, string) {
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
	return bestWeaponDonation, bestShieldDonation, bestWeaponDonationID, bestShieldDonationID
}

func NewSocialAgent() agent.Strategy {
	return &SocialAgent{
		selfishness:               rand.Float64(),
		gossipThreshold:           0.5,
		propAdmire:                0.1,
		propHate:                  0.1,
		proposalAccuracyThreshold: 0.8,
	}
}
