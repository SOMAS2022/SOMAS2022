package team1

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/message/proposal"
	"infra/teams/team1/internal"
	"math/rand"

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

	// manifesto states
	sumSocialCapital map[string]float64
	// TODO
	// agentHarshnessScore            map[string]float64
	agentsSurvivalLikelihood       map[string]float64
	survivalLikelihood             float64
	standardDeviationAgentSurvival float64
	meanSurvivalLikelihood         float64
	prevLeaderSurvivalEffect       float64
	leaderRating                   float64
}

func (s *SocialAgent) LootActionNoProposal(baseAgent agent.BaseAgent) immutable.SortedMap[commons.ItemID, struct{}] {
	return *immutable.NewSortedMap[commons.ItemID, struct{}](nil)
}

func (s *SocialAgent) LootAction(baseAgent agent.BaseAgent, proposedLoot immutable.SortedMap[commons.ItemID, struct{}]) immutable.SortedMap[commons.ItemID, struct{}] {
	return proposedLoot
}

func (s *SocialAgent) FightActionNoProposal(baseAgent agent.BaseAgent) decision.FightAction {
	qState := internal.BaseAgentToQState(baseAgent)

	// Calculate best action based on current state and selfishness
	coopTable := internal.CooperationQ(qState)

	// TODO: Maybe make non-deterministic
	// Return index of best action (assumes array ordering in same order as decision.FightAction
	return decision.FightAction(internal.Argmax(coopTable[:]))
}

func (s *SocialAgent) FightAction(baseAgent agent.BaseAgent, proposedAction decision.FightAction) decision.FightAction {
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

func (s *SocialAgent) DonateToHpPool(baseAgent agent.BaseAgent) uint {
	//return uint(rand.Intn(int(baseAgent.AgentState().Hp)))
	return 0
}

func (s *SocialAgent) UpdateInternalState(self agent.BaseAgent, fightResult *commons.ImmutableList[decision.ImmutableFightResult], decisions *immutable.Map[decision.Intent, uint]) {
	// Update socialCapital at end of each round
	itr := fightResult.Iterator()
	for !itr.Done() { // For each fight round
		fightDecisions, _ := itr.Next()

		s.updateSocialCapital(self, fightDecisions)
	}
	s.UpdateLeadershipState(self, fightResult, decisions)
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

func (r *SocialAgent) HandleTradeNegotiation(_ agent.BaseAgent, _ message.TradeInfo) message.TradeMessage {
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
