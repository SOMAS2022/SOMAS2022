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

	"github.com/benbjohnson/immutable"
)

func SocialRewardManifesto(s *SocialAgent, agent *agent.BaseAgent) float64 {

	//TODO
	// if we are doing badly then the social reward will be low,
	// current_state = agent.latestState.hp + agent.latestState.defense

	total_agents := len(s.agentsSurvivalLikelihood)
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
	deviation := s.standardDeviationAgentSurvival / 10
	if deviation > 1 {
		deviation = 1
	}
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
	// if a small number of poeple are doing badly then eta is high

	// if everyone is even then eta is in the middle
	// TODO add a more complex implementation
	return s.prevLeaderSurvivalEffect
}

func (s *SocialAgent) CreateManifesto(agent agent.BaseAgent) *decision.Manifesto {
	// TODO add count of agents making bad decision

	// TODO how should this be modified
	threshold_social_reward := 0.0 // between -1 and 1
	// TODO how should this be modified
	threshold_min_eta := 0.0
	// TODO how should this be modified
	support_threshold := 0.0

	social_reward := SocialRewardManifesto(s, &agent)
	if social_reward > threshold_social_reward {
		// Try to control
		return decision.NewManifesto(true, true, 5, 5)
	}

	most_likely_eta := MostLikelySurvivalLikelihoodOutcome(s)
	if most_likely_eta > threshold_min_eta {
		// Do nothing
		return decision.NewManifesto(false, false, 1, 1)
	}

	if SimilarAgentCount(s) < support_threshold {
		// Send out the least bad proposal that still gives us necessary power
		return decision.NewManifesto(true, true, 1, 1)
	}

	// Do nothing
	return decision.NewManifesto(false, false, 1, 1)
}

func (s *SocialAgent) HandleConfidencePoll(agent agent.BaseAgent) decision.Intent {
	if s.leaderRating < 0.4 {
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

		if s.survivalLikelihood < 0.1 {
			attitude -= s.agentHarshnessScore[id]
		}
		if attitude > 0.5 {
			ballot = append(ballot, id)
		}
	}

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

// TODO
func (s *SocialAgent) LootAllocation(ba agent.BaseAgent) immutable.Map[commons.ID, immutable.List[commons.ItemID]] {
	// TODO give loot that tries to even out the the overall scores of everyone

	// For everyone if their survival threshold is high enough give them some loot

	view := ba.View()
	weapons := ba.Loot().Weapons().Iterator()
	ids := commons.ImmutableMapKeys(view.AgentState())
	agents := view.AgentState()
	var max_attack uint = 0
	var max_health state.HealthRange = 0
	var max_defense uint = 0
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

		social_score := s.sumSocialCapital[ids[id_index]]
		survival_likelihood := s.agentsSurvivalLikelihood[ids[id_index]]
		agent_state, _ := agents.Get(ids[id_index])
		health := agent_state.Hp
		attack := agent_state.Attack

		if survival_likelihood > 0.1 {
			weapon_prob += 1.0
		}
		weapon_prob += float64(attack) / float64(max_attack)
		weapon_prob += float64(social_score) / float64(max_social)
		weapon_prob += float64(health) / float64(max_health) * 0.5

		if health < 5 {
			weapon_prob = 0
		}

		weapon_cumulative_prop = append(weapon_cumulative_prop, weapon_prob+last_weapon_prop)
		defense_cumulative_prob = append(defense_cumulative_prob, defense_prob+last_defense_prop)
		hp_cumulative_prob = append(hp_cumulative_prob, hp_prob+last_hp_prop)
		stamina_cumulative_prob = append(stamina_cumulative_prob, last_stamina_prob+stamina_prob)

		last_weapon_prop = weapon_prob
		last_defense_prop = defense_prob
		last_hp_prop = defense_prob
		last_stamina_prob = stamina_prob
	}

	// distribute weapons
	for !weapons.Done() {
		weapon, _ := weapons.Next()
		random := rand.Float64() * weapon_cumulative_prop[len(weapon_cumulative_prop)-1]

		change := len(weapon_cumulative_prop) / 4
		index := len(weapon_cumulative_prop) / 2
		if index == len(weapon_cumulative_prop)-1 {
			// pass loot to final person
		}
		if index > 2 {
			if weapon_cumulative_prop[index] > random && weapon_cumulative_prop[index-1] < random {
				// pass loot to that person
			} else if weapon_cumulative_prop[index] > random {
				index += change
			} else {
				index -= change
			}
		} else {
			// allocate to first person
		}

		change /= 2
	}

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
	mMapped := make(map[commons.ID]immutable.List[commons.ItemID])
	for id, itemIDS := range lootAllocation {
		mMapped[id] = commons.ListToImmutable(itemIDS)
	}
	return commons.MapToImmutable(mMapped)
}

// TODO
func (s *SocialAgent) FightResolution(_ agent.BaseAgent) commons.ImmutableList[proposal.Rule[decision.FightAction]] {
	// Punish agents with low social capital
	// If our agents score is really low may need to start acting selfishly
	// Otherwise act in the general social good

	// Also want to increase our social score at the same time~

	rules := make([]proposal.Rule[decision.FightAction], 0)

	rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Attack,
		proposal.NewComparativeCondition(proposal.Health, proposal.GreaterThan, 1000),
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

	return *commons.NewImmutableList(rules)
}
