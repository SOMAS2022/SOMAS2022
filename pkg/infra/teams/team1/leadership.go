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
	deviation := s.standardDeviationAgentSurvival
	if deviation > 1 {
		deviation = 1
	}
	deviation *= 0.3
	// if there is a large number of agents with a bad social score we will do badly
	capital_threshold := 0.3
	total_bad_social_capital := 0
	for _, element := range s.sumSocialCapital {
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

func SimilarAgentCount(s *SocialAgent) float64 {
	// distance between our score and the median score
	count := 0.0
	for _, element := range s.agentsSurvivalLikelihood {
		if math.Abs(element-s.meanSurvivalLikelihood) < 0.1 {
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

	if SimilarAgentCount(s) < support_threshold {
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

func (s *SocialAgent) HandleElectionBallot(b agent.BaseAgent, _ *decision.ElectionParams) decision.Ballot {
	// Iterate through each agent and check the harshness score of that agent and the social capital and the similarity in survival to calculate what to do
	var ballot decision.Ballot

	for id, survival := range s.agentsSurvivalLikelihood {
		attitude := 0.0
		attitude += (1 - math.Abs(s.survivalLikelihood-survival))

		attitude += s.sumSocialCapital[id]

		if attitude >= 0.8 {
			ballot = append(ballot, id)
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

func allocate_with_distribution(distribution []float64, iterator commons.Iterator[state.Item], ids []commons.ID, lootAllocation map[commons.ID][]commons.ItemID) {
	for !iterator.Done() {
		item, _ := iterator.Next()
		random := rand.Float64() * distribution[len(distribution)-1]
		var toBeAllocated string
		// Add max iteration
		change := int(len(distribution) / 4)
		index := int(len(distribution) / 2)
		for count := 0; count < 100; count++ {
			if index == len(distribution)-1 || index == 0 {
				toBeAllocated = ids[index]
				break
			}
			if index >= 1 {
				if distribution[index] > random && distribution[index-1] < random {
					toBeAllocated = ids[index]
					// logging.Log(logging.Error, nil, "yayyyyyy")
					break
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
		if s.sumSocialCapital[ids[id_index]] > max_social {
			max_social = s.sumSocialCapital[ids[id_index]]
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
	allocate_with_distribution(weapon_cumulative_prop, iterator, ids, lootAllocation)
	iterator = ba.Loot().Shields().Iterator()
	allocate_with_distribution(defense_cumulative_prob, iterator, ids, lootAllocation)
	iterator = ba.Loot().HpPotions().Iterator()
	allocate_with_distribution(hp_cumulative_prob, iterator, ids, lootAllocation)
	iterator = ba.Loot().StaminaPotions().Iterator()
	allocate_with_distribution(stamina_cumulative_prob, iterator, ids, lootAllocation)

	mMapped := make(map[commons.ID]immutable.SortedMap[commons.ItemID, struct{}])
	for id, itemIDS := range lootAllocation {
		mMapped[id] = commons.ListToImmutableSortedSet(itemIDS)
	}
	return commons.MapToImmutable(mMapped)
}

// TODO
func (s *SocialAgent) FightResolution(agent agent.BaseAgent, prop commons.ImmutableList[proposal.Rule[decision.FightAction]]) immutable.Map[commons.ID, decision.FightAction] {

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

	for index := 0; index < len(ids); index++ {
		percentile := float64(index) / float64(len(ids))
		health_id := ids[percentiles[0][index][1]]
		attack_id := ids[percentiles[0][index][1]]
		stamina_id := ids[percentiles[0][index][1]]
		defense_id := ids[percentiles[0][index][1]]

		// adjust attack likelihood
		decision_likelihoods[health_id] = [3]float64{
			float64(percentile) * 3.0,
			decision_likelihoods[health_id][1],
			decision_likelihoods[health_id][2]}
		decision_likelihoods[stamina_id] = [3]float64{
			decision_likelihoods[health_id][0] + float64(percentile)*3.0,
			decision_likelihoods[health_id][1],
			decision_likelihoods[health_id][2]}
		decision_likelihoods[attack_id] = [3]float64{
			decision_likelihoods[health_id][0] + float64(percentile)*3.0,
			decision_likelihoods[health_id][1],
			decision_likelihoods[health_id][2]}
		decision_likelihoods[defense_id] = [3]float64{
			decision_likelihoods[health_id][0] - float64(percentile)*0.5,
			decision_likelihoods[health_id][1],
			decision_likelihoods[health_id][2]}

		decision_likelihoods[health_id] = [3]float64{
			decision_likelihoods[health_id][0],
			(1.0/float64(percentile) + 0.1) * 1,
			decision_likelihoods[health_id][2]}
		decision_likelihoods[stamina_id] = [3]float64{
			decision_likelihoods[health_id][0],
			(1.0/float64(percentile) + 0.1) * 1,
			decision_likelihoods[health_id][2]}

		decision_likelihoods[health_id] = [3]float64{
			decision_likelihoods[health_id][0],
			decision_likelihoods[health_id][1],
			decision_likelihoods[health_id][2] + float64(percentile)}
		decision_likelihoods[stamina_id] = [3]float64{
			decision_likelihoods[health_id][0],
			decision_likelihoods[health_id][1],
			decision_likelihoods[health_id][2] + float64(percentile)*2.0}
		decision_likelihoods[attack_id] = [3]float64{
			decision_likelihoods[health_id][0],
			decision_likelihoods[health_id][1],
			decision_likelihoods[health_id][2] + float64(percentile)*0.5}
		decision_likelihoods[defense_id] = [3]float64{
			decision_likelihoods[health_id][0],
			decision_likelihoods[health_id][1],
			decision_likelihoods[health_id][2] + float64(percentile)*2.5}
	}

	// convert to cumulative prob
	for index := 0; index < len(ids); index++ {
		id := ids[index]
		decision_likelihoods[id] = [3]float64{decision_likelihoods[id][0],
			decision_likelihoods[id][0] + decision_likelihoods[id][1],
			decision_likelihoods[id][0] + decision_likelihoods[id][1] + decision_likelihoods[id][2]}
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

func (s *SocialAgent) UpdateLeadershipState(self agent.BaseAgent, fightResult *commons.ImmutableList[decision.ImmutableFightResult], _ *immutable.Map[decision.Intent, uint]) {
	view := self.View()
	agents := view.AgentState()
	ids := commons.ImmutableMapKeys(view.AgentState())
	// sumSocialCapital map[string]float64
	var sumSocialCapital = make(map[string]float64)

	for id, element := range s.socialCapital {
		sumSocialCapital[id] = element[0] + element[1] + element[2] + element[3]
	}
	s.sumSocialCapital = sumSocialCapital

	// TODO
	// agentHarshnessScore map[string]float64
	// // TODO

	// agentsSurvivalLikelihood map[string]float64
	var survival_likelihood = make(map[string]float64)
	for id := range ids {
		agent_state, _ := agents.Get(ids[id])
		survival_likelihood[ids[id]] = 0.0
		survival_likelihood[ids[id]] += float64(agent_state.Hp) + float64(agent_state.Defense)
		survival_likelihood[ids[id]] -= float64(agent_state.Attack) * 0.3

	}

	// logging.Log(logging.Error, nil, fmt.Sprintf("%f", v))

	max_survival := 0.0
	mean_survival := 0.0
	deviation := 0.0
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

	// survivalLikelihood float64
	s.agentsSurvivalLikelihood = survival_likelihood
	s.survivalLikelihood = survival_likelihood[self.ID()]
	// standardDeviationAgentSurvival float64
	s.standardDeviationAgentSurvival = deviation
	// meanSurvivalLikelihood float64
	s.meanSurvivalLikelihood = mean_survival

}
