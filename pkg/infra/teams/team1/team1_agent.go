package team1

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
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

	geometricMeanSocialCapital          map[string]float64 // Mean of social capital
	agentsSurvivalLikelihood            map[string]float64 // A score indicating how good an  agents position is, 0 to 1
	survivalLikelihood                  float64            // The score for this agent
	standardDeviationSurvivalLikelihood float64            // Indication of the spread of the agent scores
	meanSurvivalLikelihood              float64            // Mean of the scores should be shifted so that it is 0.5
	prevLeaderSurvivalEffect            float64            // The effect of the previous leader on our standing -1 to 1
	leaderRating                        float64            // How happy we are with the current leader 0 to 1
	fightBiasMovingAverage              float64            // An indicator of the bias against us -1 to 1, -1 indicates a positive bias
	averageFightMovingAverage           float64            // Helper used to calculate the fight chance multiplier
	fightChanceBiasMultiplier           float64            // Estimation of the multiplier on the chance we will fight due to social standing
	socialWelfareScore                  float64            // A score of how happy we are with the overall social score of all agents, negative is bad
	shieldNeeded                        float64            // The total minimum amount of shield that needs to be used
	totalStaminaExcessRatio             float64            // The total stamina excess from attacking expected to be left in the game
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
