package team3

import (
	"math/rand"

	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/message/proposal"

	"github.com/benbjohnson/immutable"
)

var (
	initHP        int
	initMonsterHP int
)

// HP pool donation
func (a *AgentThree) DonateToHpPool(baseAgent agent.BaseAgent) uint {
	AS := baseAgent.AgentState()
	donation := rand.Intn(2)
	// If our health is > 50% and we feel generous then donate some (max 20%) HP
	if donation == 1 {
		if int(AS.Hp) > int(0.8*float64(GetStartingHP())) {
			return uint(rand.Intn((int(AS.Hp) * 30)) / 100)
		} else if int(AS.Hp) > int(0.5*float64(GetStartingHP())) {
			return uint(rand.Intn((int(AS.Hp) * 10)) / 100)
		} else {
			return 0
		}
	} else {
		return 0
	}
}

func (a *AgentThree) FightAction(
	baseAgent agent.BaseAgent,
	_ decision.FightAction,
	_ message.Proposal[decision.FightAction],
) decision.FightAction {
	return a.FightActionNoProposal(baseAgent)
}

func (a *AgentThree) FightActionNoProposal(baseAgent agent.BaseAgent) decision.FightAction {
	agentState := baseAgent.AgentState()
	// alg 8
	if float64(agentState.Hp) < 1.05*AverageArray(GetHealthAllAgents(baseAgent)) || float64(agentState.Stamina) < 1.05*AverageArray(GetStaminaAllAgents(baseAgent)) {
		return decision.Cower
	} else if agentState.BonusDefense() <= agentState.BonusAttack() {
		return decision.Attack
	} else {
		return decision.Defend
	}
}

// Send proposal to leader
func (a *AgentThree) HandleFightInformation(_ message.TaggedInformMessage[message.FightInform], baseAgent agent.BaseAgent, fightactionMap *immutable.Map[commons.ID, decision.FightAction]) {
	// baseAgent.Log(logging.Trace, logging.LogField{"bravery": r.bravery, "hp": baseAgent.AgentState().Hp}, "Cowering")
	// AS := baseAgent.AgentState()

	// baseAgent.Log(logging.Trace, logging.LogField{"hp": AS.Hp, "decision": a.CurrentAction(baseAgent)}, "HP")
	// baseAgent.Log(logging.Trace, logging.LogField{"history": a.fightDecisionsHistory}, "Fight")

	// id := baseAgent.ID()
	// choice, _ := fightactionMap.Get(id)
	// HPThreshold1, StaminaThreshold1, AttackThreshold1, DefenseThreshold1 := a.thresholdDecision(baseAgent, choice)

	makesProposal := rand.Intn(100)

	if makesProposal > 90 {
		rules := make([]proposal.Rule[decision.FightAction], 0)

		rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Attack,
			proposal.NewAndCondition(*proposal.NewComparativeCondition(proposal.Health, proposal.GreaterThan, 500),
				*proposal.NewComparativeCondition(proposal.Stamina, proposal.GreaterThan, 500)),
		))

		rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Defend,
			proposal.NewComparativeCondition(proposal.TotalDefence, proposal.GreaterThan, 200),
		))

		rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Cower,
			proposal.NewComparativeCondition(proposal.Health, proposal.GreaterThan, 1),
		))

		rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Attack,
			proposal.NewComparativeCondition(proposal.Stamina, proposal.GreaterThan, 10),
		))

		prop := *commons.NewImmutableList(rules)
		_ = baseAgent.SendFightProposalToLeader(prop)
	}
}

// Calculate our agents action
func (a *AgentThree) CurrentAction(baseAgent agent.BaseAgent) decision.FightAction {
	view := baseAgent.View()
	agentState := baseAgent.AgentState()

	currentLevel := int(view.CurrentLevel())
	var attackDealt int
	// only sample at start
	if currentLevel == 1 {
		initHP = int(agentState.Hp)
		initMonsterHP = int(view.MonsterHealth())
	}
	// edge case - alg 9
	if float64(agentState.Hp) < 0.6*AverageArray(GetHealthAllAgents(baseAgent)) || float64(agentState.Stamina) < 0.6*AverageArray(GetStaminaAllAgents(baseAgent)) {
		return decision.Cower
	}
	// change decision, already not edge case - alg 10
	// every 3 levels, alpha +1, alpha init at 3
	alpha := (currentLevel / 3) + 3

	if currentLevel > alpha+3 {
		damageTaken := initHP - int(agentState.Hp)
		if initMonsterHP == 0 {
			attackDealt = 0
		} else {
			attackDealt = (initMonsterHP - int(view.MonsterHealth())) / initMonsterHP
		}

		// re-init vars
		initHP = int(agentState.Hp)
		initMonsterHP = int(view.MonsterHealth())

		if attackDealt <= damageTaken {
			return decision.Attack
		} else if attackDealt > damageTaken {
			return decision.Defend
		}
	}
	// catchall
	return decision.Attack
}

// Vote on proposal
func (a *AgentThree) HandleFightProposal(m message.Proposal[decision.FightAction], baseAgent agent.BaseAgent) decision.Intent {
	intent := rand.Intn(2)
	// rules := m.Rules()
	// itr := rules.Iterator()
	// for !itr.Done() {
	// 	rule, _ := itr.Next()
	// 	// baseAgent.Log(logging.Trace, logging.LogField{"rule": rule}, "Rule Proposal")
	// }
	if intent == 0 {
		return decision.Positive
	} else {
		return decision.Negative
	}
}

func (a *AgentThree) thresholdDecision(baseAgent agent.BaseAgent, choice decision.FightAction) (float64, float64, float64, float64) {
	// view := baseAgent.View()
	agentState := baseAgent.AgentState()
	HPThreshold1, StaminaThreshold1, AttackThreshold1, DefenseThreshold1 := 0.0, 0.0, 0.0, 0.0

	var agentFought bool = false

	// iterate until we get most recent history
	i := 0
	itr := a.fightDecisionsHistory.Iterator()
	for !itr.Done() {
		res, _ := itr.Next()
		i += 1

		if i == a.fightDecisionsHistory.Len()-1 {
			agents := res.AttackingAgents()
			itr2 := agents.Iterator()
			// search for our agent in fight list
			for !itr.Done() {
				_, attackingAgentID := itr2.Next()
				if attackingAgentID == baseAgent.ID() {
					agentFought = true
				}
			}
		}
	}

	if choice == decision.Cower {
		if agentState.Hp >= uint(AverageArray(GetHealthAllAgents(baseAgent))) {
			HPThreshold1 = 1.7 * AverageArray(GetHealthAllAgents(baseAgent))
		}
		if agentState.Stamina >= uint(AverageArray(GetStaminaAllAgents(baseAgent))) {
			StaminaThreshold1 = 1.7 * AverageArray(GetHealthAllAgents(baseAgent))
		}

		if agentState.Hp < uint(AverageArray(GetHealthAllAgents(baseAgent))) {
			HPThreshold1 = 0.7 * AverageArray(GetHealthAllAgents(baseAgent))
		}
		if agentState.Stamina < uint(AverageArray(GetStaminaAllAgents(baseAgent))) {
			StaminaThreshold1 = 0.7 * AverageArray(GetHealthAllAgents(baseAgent))
		}

		AttackThreshold1 = 1.1 * AverageArray(GetAttackAllAgents(baseAgent))
		DefenseThreshold1 = 1.1 * AverageArray(GetDefenceAllAgents(baseAgent))
	}

	if agentFought {
		HPThreshold1 = AverageArray(GetHealthAllAgents(baseAgent))
		StaminaThreshold1 = AverageArray(GetHealthAllAgents(baseAgent))
		AttackThreshold1 = 0.4 * AverageArray(GetAttackAllAgents(baseAgent))
		DefenseThreshold1 = 0.4 * AverageArray(GetDefenceAllAgents(baseAgent))
	}
	return HPThreshold1, StaminaThreshold1, AttackThreshold1, DefenseThreshold1
}

func (a *AgentThree) HandleUpdateWeapon(baseAgent agent.BaseAgent) decision.ItemIdx {
	// weapons := b.AgentState().Weapons
	// return decision.ItemIdx(rand.Intn(weapons.Len() + 1))

	// 0th weapon has greatest attack points
	return decision.ItemIdx(0)
}

func (a *AgentThree) HandleUpdateShield(baseAgent agent.BaseAgent) decision.ItemIdx {
	// shields := b.AgentState().Shields
	// return decision.ItemIdx(rand.Intn(shields.Len() + 1))
	return decision.ItemIdx(0)
}
