package state

import (
	"infra/game/commons"
	"infra/game/decision"
)

func ToSinglePredicate[A decision.ProposalAction](rules commons.ImmutableList[decision.Rule[A]]) func(AgentState) A {
	iterator := rules.Iterator()
	predicates := make([]func(agentState AgentState) (A, bool), 0)
	for !iterator.Done() {
		rule, _ := iterator.Next()
		pred := makePredicate(rule.Condition())
		wrappedPredicate := func(agentState AgentState) (A, bool) {
			return rule.Action(), pred(agentState)
		}
		predicates = append(predicates, wrappedPredicate)
	}
	if len(predicates) > 0 {
		return func(agentState AgentState) A {
			for _, predicate := range predicates {
				action, match := predicate(agentState)
				if match {
					return action
				}
			}
			// todo: what to do if unallocated with parameterized class
			a, _ := predicates[len(predicates)-1](agentState)
			return a
		}
	}
	return nil
}

func ToMultiPredicate[A decision.ProposalAction](rules commons.ImmutableList[decision.Rule[A]]) func(AgentState) map[A]struct{} {
	iterator := rules.Iterator()
	predicates := make([]func(agentState AgentState) (A, bool), 0)
	for !iterator.Done() {
		rule, _ := iterator.Next()
		pred := makePredicate(rule.Condition())
		wrappedPredicate := func(agentState AgentState) (A, bool) {
			return rule.Action(), pred(agentState)
		}
		predicates = append(predicates, wrappedPredicate)
	}
	if len(predicates) > 0 {
		return func(agentState AgentState) map[A]struct{} {
			res := make(map[A]struct{})
			for _, predicate := range predicates {
				action, match := predicate(agentState)
				if match {
					res[action] = struct{}{}
				}
			}
			return res
		}
	}
	return nil
}

func makePredicate(cond decision.Condition) func(agentState AgentState) bool {
	switch condT := cond.(type) {
	case *decision.ComparativeCondition:
		return buildCompPredicate(*condT)
	case decision.ComparativeCondition:
		return buildCompPredicate(condT)
	case *decision.AndCondition:
		return andEval(*condT)
	case decision.AndCondition:
		return andEval(condT)
	case *decision.OrCondition:
		return orEval(*condT)
	case decision.OrCondition:
		return orEval(condT)
	case decision.DefectorCondition:
		return defectorEval()
	default:
		return func(_ AgentState) bool {
			return true
		}
	}
}

func andEval(cond decision.AndCondition) func(AgentState) bool {
	return func(agentState AgentState) bool {
		return makePredicate(cond.CondA())(agentState) && makePredicate(cond.CondB())(agentState)
	}
}

func orEval(cond decision.OrCondition) func(AgentState) bool {
	return func(agentState AgentState) bool {
		return makePredicate(cond.CondA())(agentState) || makePredicate(cond.CondB())(agentState)
	}
}

func defectorEval() func(AgentState) bool {
	return func(agentState AgentState) bool {
		return agentState.Defector.IsDefector()
	}
}

func buildCompPredicate(condT decision.ComparativeCondition) func(agentState AgentState) bool {
	return func(agentState AgentState) bool {
		var attr uint
		switch condT.Attribute {
		case decision.Health:
			attr = agentState.Hp
		case decision.Stamina:
			attr = agentState.Stamina
		case decision.TotalAttack:
			attr = agentState.TotalAttack()
		case decision.TotalDefence:
			attr = agentState.TotalDefense()
		default:
			attr = agentState.Hp
		}
		switch condT.Comparator {
		case decision.GreaterThan:
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
