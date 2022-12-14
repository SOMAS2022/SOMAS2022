package team6

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message/proposal"
)

type BasicAgentState struct {
	Hp      uint
	Stamina uint
	Attack  uint
	Defense uint
}

// Currently converts rules to a single predicate rather than multiple predicates, and doesn't consider defecting
func ProposalSimilarity[A decision.ProposalAction](rules1 commons.ImmutableList[proposal.Rule[A]], rules2 commons.ImmutableList[proposal.Rule[A]]) float32 {
	hp := [10]uint{1000, 1000, 750, 750, 500, 500, 250, 250, 100, 100}
	stamina := [10]uint{2000, 2000, 1500, 1500, 1000, 1000, 500, 500, 200, 200}
	attack := [10]uint{20, 60, 20, 60, 30, 20, 40, 30, 80, 20}
	defense := [10]uint{20, 60, 20, 70, 20, 20, 40, 40, 20, 20}

	total := 10
	same := 0

	// Convert rules to predicates
	rule1 := ToSinglePredicate(rules1)
	rule2 := ToSinglePredicate(rules2)

	// Apply rules to agent states and compare outcomes
	for i := 0; i < 10; i++ {
		agentState := BasicAgentState{Hp: hp[i], Stamina: stamina[i], Attack: attack[i], Defense: defense[i]}

		decision1 := rule1(agentState)
		decision2 := rule2(agentState)

		if decision1 == decision2 {
			same++
		}
	}

	return float32(same) / float32(total)
}

func ToSinglePredicate[A decision.ProposalAction](rules commons.ImmutableList[proposal.Rule[A]]) func(BasicAgentState) A {
	iterator := rules.Iterator()
	predicates := make([]func(agentState BasicAgentState) (A, bool), 0)
	for !iterator.Done() {
		rule, _ := iterator.Next()
		pred := makePredicate(rule.Condition())
		wrappedPredicate := func(agentState BasicAgentState) (A, bool) {
			return rule.Action(), pred(agentState)
		}
		predicates = append(predicates, wrappedPredicate)
	}
	if len(predicates) > 0 {
		return func(agentState BasicAgentState) A {
			for _, predicate := range predicates {
				action, match := predicate(agentState)
				if match {
					return action
				}
			}
			a, _ := predicates[len(predicates)-1](agentState)
			return a
		}
	}
	return nil
}

func makePredicate(cond proposal.Condition) func(agentState BasicAgentState) bool {
	switch condT := cond.(type) {
	case *proposal.ComparativeCondition:
		return buildCompPredicate(*condT)
	case proposal.ComparativeCondition:
		return buildCompPredicate(condT)
	case *proposal.AndCondition:
		return andEval(*condT)
	case proposal.AndCondition:
		return andEval(condT)
	case *proposal.OrCondition:
		return orEval(*condT)
	case proposal.OrCondition:
		return orEval(condT)
	case proposal.DefectorCondition: // Essentially considers all agents as not defecting
		return func(_ BasicAgentState) bool {
			return false
		}
	default:
		return func(_ BasicAgentState) bool {
			return true
		}
	}
}

func andEval(cond proposal.AndCondition) func(BasicAgentState) bool {
	return func(agentState BasicAgentState) bool {
		return makePredicate(cond.CondA())(agentState) && makePredicate(cond.CondB())(agentState)
	}
}

func orEval(cond proposal.OrCondition) func(BasicAgentState) bool {
	return func(agentState BasicAgentState) bool {
		return makePredicate(cond.CondA())(agentState) || makePredicate(cond.CondB())(agentState)
	}
}

func buildCompPredicate(condT proposal.ComparativeCondition) func(agentState BasicAgentState) bool {
	return func(agentState BasicAgentState) bool {
		var attr uint
		switch condT.Attribute {
		case proposal.Health:
			attr = agentState.Hp
		case proposal.Stamina:
			attr = agentState.Stamina
		case proposal.TotalAttack:
			attr = agentState.Attack
		case proposal.TotalDefence:
			attr = agentState.Defense
		default:
			attr = agentState.Hp
		}
		switch condT.Comparator {
		case proposal.GreaterThan:
			if attr > condT.Value {
				return true
			} else {
				return false
			}
		default:
			if attr < condT.Value {
				return true
			} else {
				return false
			}
		}
	}
}
