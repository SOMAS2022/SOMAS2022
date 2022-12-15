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

	graphID int // for logging

	proposalAccuracyThreshold float64

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

// TODO
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

// Update social capital at end of each round
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

func (s *SocialAgent) CreateFightProposal(baseAgent agent.BaseAgent) []proposal.Rule[decision.FightAction] {
	// find the action each agent will make
	// check what the average of each state type is
	type AgentAction struct {
		action      decision.FightAction
		agent_state state.HiddenAgentState
	}
	agent_actions := make([]AgentAction, 0)
	view := baseAgent.View()
	ids := commons.ImmutableMapKeys(view.AgentState())
	agents := view.AgentState()
	average_attack := 0.0
	average_defense := 0.0
	for id_index := 0; id_index < len(ids); id_index++ {
		agent_state, _ := agents.Get(ids[id_index])
		average_attack += float64(agent_state.Attack)
		average_defense += float64(agent_state.Defense)
		// qState := internal.HiddenAgentToQState(agent_state, view)
		// rewards := internal.CooperationQ(qState)
		// q_action := decision.FightAction(internal.Argmax(rewards[:]))
		// agent_actions = append(agent_actions, AgentAction{action: q_action, agent_state: agent_state})
	}

	average_attack /= float64(len(ids))
	average_defense /= float64(len(ids))

	// find the average of each stat
	// agent_averages := [3][4]float64{{0.0, 0.0, 0.0, 0.0}} // defend cower attack, health attack stamina defense
	// for _, agent_action := range agent_actions {
	// 	agent_averages[agent_action.action] = [4]float64{
	// 		agent_averages[agent_action.action][0] + float64(agent_action.agent_state.Hp),
	// 		agent_averages[agent_action.action][1] + float64(agent_action.agent_state.Attack),
	// 		agent_averages[agent_action.action][2] + float64(agent_action.agent_state.Stamina),
	// 		agent_averages[agent_action.action][3] + float64(agent_action.agent_state.Defense),
	// 	}
	// }
	// for index := 0; index < 3; index++ {
	// 	agent_averages[index] = [4]float64{
	// 		agent_averages[index][0] / float64(len(agent_actions)),
	// 		agent_averages[index][1] / float64(len(agent_actions)),
	// 		agent_averages[index][2] / float64(len(agent_actions)),
	// 		agent_averages[index][3] / float64(len(agent_actions)),
	// 	}
	// }
	halved_attack_average := 0.0
	halved_defend_average := 0.0
	// construct rules based on these agent averages, 36 different rules each corresponding to a range of the possible state space
	rules := make([]proposal.Rule[decision.FightAction], 0)
	for health_range := 1; health_range <= 3; health_range++ {
		for stamina_range := 1; stamina_range <= 3; stamina_range++ {
			health_val_min := 250 * health_range
			stamina_val_min := 500 * stamina_range
			for attack_quartile := 1.0; attack_quartile < 4.0; attack_quartile += 2.0 {
				for defend_quartile := 1.0; defend_quartile < 4.0; defend_quartile += 2.0 {
					// create a rule for the current quartile health and attack rule
					attack_mid := halved_attack_average * attack_quartile
					defend_mid := halved_defend_average * defend_quartile
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

}

// TODO
func (s *SocialAgent) HandleFightProposal(prop message.Proposal[decision.FightAction], baseAgent agent.BaseAgent) decision.Intent {
	var result decision.Intent
	view := baseAgent.View()
	ids := commons.ImmutableMapKeys(view.AgentState())
	agents := view.AgentState()
	rules := prop.Rules()
	action_checker := proposal.ToSinglePredicate(rules)
	accuracy := 0.0
	for id_index := 0; id_index < len(ids); id_index++ {
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
	accuracy /= float64(len(ids))

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

func (s *SocialAgent) HandleTradeNegotiation(_ agent.BaseAgent, _ message.TradeInfo) message.TradeMessage {
	return message.TradeRequest{}
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
