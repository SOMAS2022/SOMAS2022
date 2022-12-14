package team1

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/message/proposal"
	"infra/game/state"
	"math"
	"math/rand"
	"sort"

	"github.com/benbjohnson/immutable"
)

var socialWelfareMinimum = -0.3
var attitudeWithImpositionThreshold = 0.8
var attitudeWithoutImpositionThreshold = 0.3

func SocialRewardManifesto(s *SocialAgent, agent *agent.BaseAgent) float64 {

	//TODO
	// if we are doing badly then the social reward will be low,
	// current_state = agent.latestState.hp + agent.latestState.defense

	total_agents := len(s.agentsSurvivalLikelihood)
	if total_agents == 0 {
		total_agents += 1
	}
	// If a small number of agents is doing badly then the social reward will be low
	// If we are doing badly and a small number are doing well then the social reward will be positive
	bad_threshold := 0.2
	good_threshold := 0.8
	total_bad := 0
	total_good := 0
	for _, element := range s.agentsSurvivalLikelihood {
		if element < bad_threshold {
			total_bad += 1
		}
		if element > good_threshold {
			total_good += 1
		}
	}
	total_bad = total_bad / total_agents
	total_good = total_good / total_agents
	// if everyone is even the social reward will be good
	deviation := s.standardDeviationSurvivalLikelihood
	if deviation > 1 {
		deviation = 1
	}
	deviation *= 0.3
	// if there is a large number of agents with a bad social score we will do badly
	capital_threshold := 0.3
	total_bad_social_capital := 0
	for _, element := range s.geometricMeanSocialCapital {
		if element < capital_threshold {
			total_bad_social_capital += 1
		}
	}
	total_bad_social_capital = total_bad_social_capital / total_agents

	total_score := 0.0
	total_score += (1 - deviation)
	total_score -= (float64(total_bad_social_capital) * 0.1)
	total_score += (1 - s.survivalLikelihood) * float64(total_good)
	total_score -= (s.survivalLikelihood) * float64(total_bad)

	return total_score / 2
}

func ScoreIsSimilar(s *SocialAgent, survivalLikelihood float64) bool {
	return math.Abs(survivalLikelihood-s.survivalLikelihood) < 0.1
}

func SimilarAgentRatio(s *SocialAgent) float64 {
	// distance between our score and the median score
	count := 0.0
	for _, element := range s.agentsSurvivalLikelihood {
		if ScoreIsSimilar(s, element) {
			count += 1.0
		}
	}
	return count / float64(len(s.agentsSurvivalLikelihood))
}

func MostLikelySurvivalLikelihoodOutcome(s *SocialAgent) float64 {
	// if we are doing badly compared to everyone else then eta will be low
	// if a small number of people are doing badly then eta is high
	// if everyone is even then eta is in the middle
	// TODO add a more complex implementation
	return s.prevLeaderSurvivalEffect
}

func (s *SocialAgent) CreateManifesto(agent agent.BaseAgent) *decision.Manifesto {
	total_agents := len(s.agentsSurvivalLikelihood)
	if total_agents == 0 {
		total_agents += 1
	}
	// TODO add count of agents making bad decision

	// TODO how should this be modified
	threshold_social_reward := 0.2 // between -1 and 1
	// TODO how should this be modified
	threshold_min_eta := 0.6
	// TODO how should this be modified
	support_threshold := 0.3

	social_reward := SocialRewardManifesto(s, &agent)
	if social_reward > threshold_social_reward {
		// Try to control
		//TODO PERCENT NUM AGENT
		return decision.NewManifesto(true, true, 5, uint(total_agents/2))
	}

	most_likely_eta := MostLikelySurvivalLikelihoodOutcome(s)
	if most_likely_eta > threshold_min_eta {
		// Do nothing
		//TODO PERCENT NUM AGENT
		return decision.NewManifesto(false, false, 1, uint(total_agents/2))
	}

	if SimilarAgentRatio(s) < support_threshold {
		// Send out the least bad proposal that still gives us necessary power
		//TODO PERCENT NUM AGENT
		return decision.NewManifesto(true, true, 1, uint(total_agents/2))
	}

	// Do nothing
	//TODO PERCENT NUM AGENT
	return decision.NewManifesto(false, false, 1, uint(total_agents/2))
}

func (s *SocialAgent) HandleConfidencePoll(agent agent.BaseAgent) decision.Intent {
	if s.leaderRating < 0.2 {
		return decision.Negative
	}
	if s.leaderRating < 0.6 {
		return decision.Abstain
	}
	return decision.Positive
}

func (s *SocialAgent) HandleElectionBallot(b agent.BaseAgent, electionParams *decision.ElectionParams) decision.Ballot {
	// Iterate through each agent and check the harshness score of that agent and the social capital and the similarity in survival to calculate what to do
	var ballot decision.Ballot

	requireLeaderForPunishment := s.socialWelfareScore < socialWelfareMinimum
	candidates := electionParams.CandidateList().Iterator()
	// for id, survival := range candidates {
	for !candidates.Done() {
		id, manifesto, _ := candidates.Next()
		survival := s.agentsSurvivalLikelihood[id]
		attitude := 0.0
		// Increase likelihood of voting if the agent is similar to us
		attitude += (1 - math.Abs(s.survivalLikelihood-survival))
		// Increase attitude if agent has a good social capital
		attitude += s.geometricMeanSocialCapital[id]
		// Scale attitude to be at most 1
		attitude /= 2

		if attitude >= attitudeWithImpositionThreshold && requireLeaderForPunishment {
			// vote with a probability related to the attitude if the attitude is high enough
			random := rand.Float64()
			if random < attitude && manifesto.FightImposition() && manifesto.LootImposition() {
				ballot = append(ballot, id)
			}
		} else if attitude > attitudeWithoutImpositionThreshold {
			random := rand.Float64()
			if random < attitude {
				ballot = append(ballot, id)
			}
		}
	}
	ballot = append(ballot, b.ID())
	return ballot
}

// TODO
func (s *SocialAgent) HandleFightProposal(_ message.Proposal[decision.FightAction], _ agent.BaseAgent) decision.Intent {
	intent := rand.Intn(2)
	if intent == 0 {
		return decision.Positive
	} else {
		return decision.Negative
	}
}

func SampleDistribution(distribution []float64) int {
	random := rand.Float64() * distribution[len(distribution)-1]
	// Add max iteration
	change := int(len(distribution) / 4)
	index := int(len(distribution) / 2)
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
		// random := rand.Float64() * distribution[len(distribution)-1]
		var toBeAllocated string
		// Add max iteration
		// change := int(len(distribution) / 4)
		// index := int(len(distribution) / 2)
		toBeAllocated = ids[SampleDistribution(distribution)]
		// for count := 0; count < 100; count++ {
		// 	if index == len(distribution)-1 || index == 0 {
		// 		toBeAllocated = ids[index]
		// 		break
		// 	}
		// 	if index >= 1 {
		// 		if distribution[index] > random && distribution[index-1] < random {
		// 			toBeAllocated = ids[index]
		// 			// logging.Log(logging.Error, nil, "yayyyyyy")
		// 			break
		// 		} else if distribution[index] > random {
		// 			index -= change
		// 		} else {
		// 			index += change
		// 		}
		// 	}
		// 	change /= 2
		// 	if change < 1 {
		// 		change = 1
		// 	}
		// }
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

// TODO
func (s *SocialAgent) LootAllocation(ba agent.BaseAgent) immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]] {
	lootAllocation := make(map[commons.ID][]commons.ItemID)

	view := ba.View()
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
		if s.geometricMeanSocialCapital[ids[id_index]] > max_social {
			max_social = s.geometricMeanSocialCapital[ids[id_index]]
		}
	}

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

		// social_score := s.sumSocialCapital[ids[id_index]]
		survival_likelihood := s.agentsSurvivalLikelihood[ids[id_index]]
		agent_state, _ := agents.Get(ids[id_index])

		if survival_likelihood > 0.05 {
			weapon_prob += 1.0
			// stamina_prob += 1.0
		}

		// weapon_prob += float64(max_attack) / ((float64(agent_state.Attack) * (float64(agent_state.Attack))) + 0.1)
		// weapon_prob += float64(social_score) / float64(max_social) * 0.1
		weapon_prob += float64(agent_state.Hp)/float64(max_health)*0.9 + 9*float64(agent_state.Stamina)/float64(max_stamina)

		defense_prob += float64(max_defense) / (math.Pow(float64(agent_state.Defense), 4) + 0.1)
		// defense_prob += float64(social_score) / float64(max_social) * 0.1
		// defense_prob += float64(agent_state.Hp) / float64(max_health)

		// hp_prob += float64(social_score) / float64(max_social) * 0.1
		hp_prob += float64(max_health) / (math.Pow(float64(agent_state.Hp), 4) + 0.1)

		// stamina_prob += float64(social_score) / float64(max_social) * 0.3
		stamina_prob += float64(max_stamina) / (math.Pow(float64(agent_state.Stamina), 4) + 0.1)

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
	iterator := ba.Loot().Weapons().Iterator()
	AllocateWithProbabilityDistribution(weapon_cumulative_prop, iterator, ids, lootAllocation)
	iterator = ba.Loot().Shields().Iterator()
	AllocateWithProbabilityDistribution(defense_cumulative_prob, iterator, ids, lootAllocation)
	iterator = ba.Loot().HpPotions().Iterator()
	AllocateWithProbabilityDistribution(hp_cumulative_prob, iterator, ids, lootAllocation)
	iterator = ba.Loot().StaminaPotions().Iterator()
	AllocateWithProbabilityDistribution(stamina_cumulative_prob, iterator, ids, lootAllocation)

	mMapped := make(map[commons.ID]immutable.SortedMap[commons.ItemID, struct{}])
	for id, itemIDS := range lootAllocation {
		mMapped[id] = commons.ListToImmutableSortedSet(itemIDS)
	}
	return commons.MapToImmutable(mMapped)
}

// TODO
func (s *SocialAgent) FightResolution(agent agent.BaseAgent, prop commons.ImmutableList[proposal.Rule[decision.FightAction]]) immutable.Map[commons.ID, decision.FightAction] {

	// // find out what our expected action should be
	// preferedDecision := decision.Attack

	// totalMinFutureDamage := 0.0
	// totalHealth := 0.0
	// totalStamina := 0.0

	// shieldNeeded := totalHealth - totalMinFutureDamage
	// if shieldNeeded < 0 {
	// 	shieldNeeded = 0
	// }

	// excessStamina := totalStamina - totalMinFutureDamage*1.5

	if s.totalStaminaExcessRatio < 0 {
		// sacrifice an agent
		return s.SacrificeAgent(agent, prop)
	}

	// make people attack and defend such that the correct percentage of people defend

	view := agent.View()
	ids := commons.ImmutableMapKeys(view.AgentState())
	agents := view.AgentState()

	// monster_attack := view.MonsterAttack()
	// monster_health := view.MonsterHealth()
	// find percentiles of agents
	var percentiles [4][][2]int // health, attack, stamina, defense (val, index)
	total_attack := 0
	for id_index := 0; id_index < len(ids); id_index++ {

		agent_state, _ := agents.Get(ids[id_index])
		percentiles[0] = append(percentiles[0], [2]int{int(agent_state.Hp), id_index})
		percentiles[1] = append(percentiles[1], [2]int{int(agent_state.Attack), id_index})
		total_attack += int(agent_state.Attack)
		percentiles[2] = append(percentiles[1], [2]int{int(agent_state.Stamina), id_index})
		percentiles[3] = append(percentiles[1], [2]int{int(agent_state.Defense), id_index})
	}

	for index := 0; index < 4; index++ {
		sort.SliceStable(percentiles[index], func(i, j int) bool {
			return percentiles[index][i][0] < percentiles[index][j][0]
		})
	}

	var decision_likelihoods = make(map[string][3]float64) // attack, cower, defend

	// find required attack
	// find total attack
	// find probability attack
	//TODO this should be calculated through adjusting for over-killing and ensuring their are enough defenders
	p_attack_required := 0.9

	cower_multiplier := 2 * s.totalStaminaExcessRatio
	for index := 0; index < len(ids); index++ {
		percentile := float64(index) / float64(len(ids))
		health_id := ids[percentiles[0][index][1]]
		attack_id := ids[percentiles[1][index][1]]
		stamina_id := ids[percentiles[2][index][1]]
		defense_id := ids[percentiles[3][index][1]]

		// adjust indicator likelihood
		decision_likelihoods[health_id] = [3]float64{
			decision_likelihoods[health_id][0] + float64(percentile)*3.0,                            // attack
			decision_likelihoods[health_id][1] + (1.0/float64(percentile)+0.1)*1 + cower_multiplier, // cower
			decision_likelihoods[health_id][2] + float64(percentile)}                                // defend

		decision_likelihoods[stamina_id] = [3]float64{
			decision_likelihoods[stamina_id][0] + float64(percentile)*2.0,         // attack
			decision_likelihoods[stamina_id][1] + (3.0/float64(percentile)+0.1)*1, // cower
			decision_likelihoods[stamina_id][2] + float64(percentile)*2.0}         // defend

		decision_likelihoods[attack_id] = [3]float64{
			decision_likelihoods[attack_id][0] + float64(percentile)*1.0, // attack
			decision_likelihoods[attack_id][1],                           // cower
			decision_likelihoods[attack_id][2] - 0.2*float64(percentile)} // defend

		decision_likelihoods[defense_id] = [3]float64{
			decision_likelihoods[defense_id][0] + float64(percentile)*3.0,         // attack
			decision_likelihoods[defense_id][1] + (1.0/float64(percentile)+0.1)*1, // cower
			decision_likelihoods[defense_id][2] + float64(percentile)}             // defend

	}

	// shift the attack probabilities so that the correct number of people attack hopefully
	total_p_attack := 0.0
	for index := 0; index < len(ids); index++ {
		total_p_attack += decision_likelihoods[ids[index]][0] / (decision_likelihoods[ids[index]][1] + decision_likelihoods[ids[index]][2])
	}
	total_p_attack /= float64(len(ids))
	for index := 0; index < len(ids); index++ {
		decision_likelihoods[ids[index]] = [3]float64{
			decision_likelihoods[ids[index]][0] * p_attack_required / total_p_attack,
			decision_likelihoods[ids[index]][1],
			decision_likelihoods[ids[index]][2],
		}
	}

	builder := immutable.NewMapBuilder[commons.ID, decision.FightAction](nil)
	for _, id := range commons.ImmutableMapKeys(view.AgentState()) {
		random := rand.Float64() * decision_likelihoods[id][2]
		// logging.Log(logging.Error, nil, fmt.Sprintf("%f", decision_likelihoods[id]))
		var fightAction decision.FightAction

		if random < decision_likelihoods[id][0] {
			fightAction = decision.Attack
		} else if random < decision_likelihoods[id][1] {
			fightAction = decision.Cower
		} else {
			fightAction = decision.Defend
		}

		builder.Set(id, fightAction)
	}
	return *builder.Map()
}

func (s *SocialAgent) SacrificeAgent(agent agent.BaseAgent, prop commons.ImmutableList[proposal.Rule[decision.FightAction]]) immutable.Map[commons.ID, decision.FightAction] {

	// make and agent probability of sacrifice higher if it has bad social standing
	// make it higher if it has low health
	// make it higher if it has high attack
	view := agent.View()
	ids := commons.ImmutableMapKeys(view.AgentState())
	agents := view.AgentState()
	var decision_likelihoods = make([]float64, 0) // attack, cower, defend

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
		if s.geometricMeanSocialCapital[ids[id_index]] > max_social {
			max_social = s.geometricMeanSocialCapital[ids[id_index]]
		}
	}

	for index := 0; index < len(ids); index++ {
		agent_state, _ := agents.Get(ids[index])
		likelihood := float64(agent_state.Attack) / float64(max_attack)
		likelihood += float64(max_health-agent_state.Hp) / float64(max_health)
		likelihood += math.Pow(1-s.geometricMeanSocialCapital[ids[index]], 2) * 2
		decision_likelihoods = append(decision_likelihoods, likelihood)
	}

	sacrificedAgentID := ids[SampleDistribution(decision_likelihoods)]

	builder := immutable.NewMapBuilder[commons.ID, decision.FightAction](nil)
	for _, id := range commons.ImmutableMapKeys(view.AgentState()) {
		var fightAction decision.FightAction
		if id == sacrificedAgentID {
			fightAction = decision.Attack
		} else {
			fightAction = decision.Cower
		}
		builder.Set(id, fightAction)
	}
	return *builder.Map()
}

// TODO
func (s *SocialAgent) HandleFightInformation(m message.TaggedInformMessage[message.FightInform], baseAgent agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) {

	// calculate the maximum damage done to each person
	// remove anyone who would not survive that damage
	// check if anymore people would die
	// if so remove them from the

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

func (s *SocialAgent) UpdateLeadershipState(self agent.BaseAgent, fightResult *commons.ImmutableList[decision.ImmutableFightResult], _ *immutable.Map[decision.Intent, uint]) {
	CalculateSurvivalLikelihoods(s, self)
	CalculateSurvivalLikelihoodStats(s, self)
	EstimateSocialWelfare(s)
	EstimateSocialBias(s, fightResult, self)

	// TODO
	s.totalStaminaExcessRatio = rand.Float64()
	s.shieldNeeded = rand.Float64()
}

func CalculateSurvivalLikelihoods(s *SocialAgent, self agent.BaseAgent) {
	view := self.View()
	agents := view.AgentState()
	ids := commons.ImmutableMapKeys(view.AgentState())
	var geometricMeanSocialCapital = make(map[string]float64)

	for id, element := range s.socialCapital {
		geometricMeanSocialCapital[id] = math.Pow(element[0]*element[1]*element[2]*element[3], 1.0/4.0)
	}
	s.geometricMeanSocialCapital = geometricMeanSocialCapital
	s.leaderRating = geometricMeanSocialCapital[view.CurrentLeader()]

	var survival_likelihood = make(map[string]float64)
	for id := range ids {
		agent_state, _ := agents.Get(ids[id])
		survival_likelihood[ids[id]] = 0.0
		survival_likelihood[ids[id]] += float64(agent_state.Hp) + float64(agent_state.Defense)
		survival_likelihood[ids[id]] -= float64(agent_state.Attack) * 0.3

	}
	s.agentsSurvivalLikelihood = survival_likelihood
}

func CalculateSurvivalLikelihoodStats(s *SocialAgent, self agent.BaseAgent) {
	max_survival := 0.0
	mean_survival := 0.0
	deviation := 0.0
	survival_likelihood := s.agentsSurvivalLikelihood
	for _, survival := range survival_likelihood {
		if survival > float64(max_survival) {
			max_survival = survival
		}
	}
	for id, survival := range survival_likelihood {
		survival_likelihood[id] = survival / max_survival
		mean_survival += survival / max_survival
	}
	mean_survival /= float64(len(survival_likelihood))
	for id, survival := range survival_likelihood {
		survival_likelihood[id] = (survival / mean_survival) * 0.5
		mean_survival += survival / max_survival
	}
	mean_survival = 0.5
	for _, survival := range survival_likelihood {
		deviation += math.Pow(survival-mean_survival, 2)
	}
	deviation /= float64(len(survival_likelihood))
	deviation = math.Pow(deviation, 0.5)

	s.agentsSurvivalLikelihood = survival_likelihood
	s.survivalLikelihood = survival_likelihood[self.ID()]
	s.standardDeviationSurvivalLikelihood = deviation
	s.meanSurvivalLikelihood = mean_survival
}

func EstimateSocialWelfare(s *SocialAgent) {
	s.socialWelfareScore = 0
	exponent := 2.0
	for _, socialCapital := range s.geometricMeanSocialCapital {
		// Since the score is between -1 and 1 shift so as to be between 0 and 2
		s.socialWelfareScore += math.Pow(2-socialCapital, exponent)
	}
	s.socialWelfareScore /= float64(len(s.geometricMeanSocialCapital))
	// Shift the output score to be between -1 and 1 again
	s.socialWelfareScore = 2 - math.Pow(s.socialWelfareScore, 1/exponent)
}

func EstimateSocialBias(s *SocialAgent, fightResult *commons.ImmutableList[decision.ImmutableFightResult], self agent.BaseAgent) {
	fight_result := fightResult.Get(fightResult.Len() - 1)
	attacking_agents := fight_result.AttackingAgents()
	made_to_fight := false           // Check if we were made to fight
	num_similar_made_to_fight := 0.0 // Find the number of agents similar to us made to fight

	for index := 0; index < attacking_agents.Len(); index++ {
		made_to_fight = made_to_fight || (attacking_agents.Get(index) == self.ID())
		if ScoreIsSimilar(s, s.agentsSurvivalLikelihood[attacking_agents.Get(index)]) {
			num_similar_made_to_fight += 1.0
		}
	}
	// Find the ratio of agents similar to us mde to fight
	ratio_agents_made_to_fight := num_similar_made_to_fight / (SimilarAgentRatio(s)*float64(len(s.agentsSurvivalLikelihood)) + 0.1)

	// Adjust this number to change the many rounds to average over
	e := 0.1
	if made_to_fight {
		s.fightBiasMovingAverage = (1-e)*s.fightBiasMovingAverage + (1-ratio_agents_made_to_fight)*e
	} else {
		s.fightBiasMovingAverage = (1-e)*s.fightBiasMovingAverage - ratio_agents_made_to_fight*e
	}

	s.averageFightMovingAverage = (1-e)*s.averageFightMovingAverage + ratio_agents_made_to_fight
	// Find the multiplier on our chance of fighting due to the bias
	s.fightChanceBiasMultiplier = (s.fightBiasMovingAverage + s.averageFightMovingAverage) / s.averageFightMovingAverage
	// logging.Log(logging.Error, nil, fmt.Sprintf("%f", (s.fightChanceBiasMultiplier)))
}
