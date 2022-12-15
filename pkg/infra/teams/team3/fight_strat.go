package team3

import (
	"math/rand"

	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/message/proposal"
	"infra/logging"

	"github.com/benbjohnson/immutable"
)

// HP pool donation
func (a *AgentThree) DonateToHpPool(baseAgent agent.BaseAgent) uint {
	donation := rand.Intn(2)
	// If our health is > 50% and we feel generous then donate some (max 20%) HP
	if donation == 1 && a.HP > PERCENTAGE {
		return uint(rand.Intn((a.HP * 20) / 100))
	} else {
		return 0
	}
}

func (a *AgentThree) FightAction(
	baseAgent agent.BaseAgent,
	proposedAction decision.FightAction,
	acceptedProposal message.Proposal[decision.FightAction],
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

func (a *AgentThree) FightResolution(baseAgent agent.BaseAgent, prop commons.ImmutableList[proposal.Rule[decision.FightAction]]) immutable.Map[commons.ID, decision.FightAction] {
	view := baseAgent.View()
	builder := immutable.NewMapBuilder[commons.ID, decision.FightAction](nil)
	for _, id := range commons.ImmutableMapKeys(view.AgentState()) {
		var fightAction decision.FightAction

		// Check for our agent and assign what we want to do
		if id == baseAgent.ID() {
			action := a.CurrentAction(baseAgent)
			fightAction = action
			baseAgent.Log(logging.Trace, logging.LogField{"hp": a.HP, "choice": action, "util": a.utilityScore[view.CurrentLeader()]}, "Intent")
		} else {
			switch rand.Intn(3) {
			case 0:
				fightAction = decision.Attack
			case 1:
				fightAction = decision.Defend
			default:
				fightAction = decision.Cower
			}
		}
		builder.Set(id, fightAction)
	}
	return *builder.Map()
}

// Send proposal to leader
func (a *AgentThree) HandleFightInformation(_ message.TaggedInformMessage[message.FightInform], baseAgent agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) {
	// baseAgent.Log(logging.Trace, logging.LogField{"bravery": r.bravery, "hp": baseAgent.AgentState().Hp}, "Cowering")
	makesProposal := rand.Intn(100)

	baseAgent.Log(logging.Trace, logging.LogField{"hp": a.HP, "decision": a.CurrentAction(baseAgent)}, "HP")
	baseAgent.Log(logging.Trace, logging.LogField{"history": a.fightDecisionsHistory}, "Fight")

	// Well, not everytime. Just sometimes
	if makesProposal > 80 {
		rules := make([]proposal.Rule[decision.FightAction], 0)

		rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Attack,
			proposal.NewAndCondition(*proposal.NewComparativeCondition(proposal.Health, proposal.GreaterThan, 500),
				*proposal.NewComparativeCondition(proposal.Stamina, proposal.GreaterThan, 1000)),
		))

		rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Defend,
			proposal.NewComparativeCondition(proposal.TotalDefence, proposal.GreaterThan, 50),
		))

		rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Cower,
			proposal.NewComparativeCondition(proposal.Health, proposal.LessThan, 100),
		))

		rules = append(rules, *proposal.NewRule[decision.FightAction](decision.Attack,
			proposal.NewComparativeCondition(proposal.Stamina, proposal.GreaterThan, baseAgent.AgentState().Stamina),
		))

		prop := *commons.NewImmutableList(rules)
		_ = baseAgent.SendFightProposalToLeader(prop)
	}
}

// Calculate our agents action
func (a *AgentThree) CurrentAction(baseAgent agent.BaseAgent) decision.FightAction {
	// // !!!!!!!!!!!!!!!!!!!!! need to implement
	// StartingMonsterHP := 1000
	// view := baseAgent.View()
	// agentState := baseAgent.AgentState()

	// currentLevel := int(view.CurrentLevel())
	// // edge case - alg 9
	// if float64(agentState.Hp) < 0.6*AverageArray(GetHealthAllAgents(baseAgent)) || float64(agentState.Stamina) < 0.6*AverageArray(GetStaminaAllAgents(baseAgent)) {
	// 	return decision.Cower
	// }
	// // change decision, already not edge case - alg 10
	// // every 5 levels, alpha +1
	// alpha := int(currentLevel / 5)
	// damageTaken := GetStartingHP() - int(agentState.Hp)
	// attackDealt := (StartingMonsterHP - int(view.MonsterHealth())) / StartingMonsterHP

	// if attackDealt <= damageTaken && currentLevel > alpha+5 {
	// 	return decision.Attack
	// } else if attackDealt > damageTaken && currentLevel > alpha+5 {
	// 	return decision.Defend
	// }

	// // catchall, execution will never get here
	// return decision.Cower
	return decision.Attack
}

// Vote on proposal
func (a *AgentThree) HandleFightProposal(m message.Proposal[decision.FightAction], baseAgent agent.BaseAgent) decision.Intent {
	agree := true

	// rules := m.Rules()
	// itr := rules.Iterator()
	// for !itr.Done() {
	// 	rule, _ := itr.Next()
	// 	// baseAgent.Log(logging.Trace, logging.LogField{"rule": rule}, "Rule Proposal")
	// }

	// Selfish, only agree if our decision is ok
	if agree {
		return decision.Positive
	} else {
		return decision.Negative
	}
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
