package state

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state/proposal"
)

func ToSinglePredicate[A decision.ProposalAction](rules commons.ImmutableList[proposal.Rule[A]]) func(AgentState) A {
	iterator := rules.Iterator()
	predicates := make([]func(state AgentState) (A, bool), 0)
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

func ToMultiPredicate[A decision.ProposalAction](rules commons.ImmutableList[proposal.Rule[A]]) func(AgentState) map[A]struct{} {
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

func makePredicate(cond proposal.Condition) func(agentState AgentState) bool {
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
	case proposal.DefectorCondition:
		return defectorEval()
	default:
		return func(_ AgentState) bool {
			return true
		}
	}
}

func andEval(cond proposal.AndCondition) func(AgentState) bool {
	return func(agentState AgentState) bool {
		return makePredicate(cond.CondA())(agentState) && makePredicate(cond.CondB())(agentState)
	}
}

func orEval(cond proposal.OrCondition) func(AgentState) bool {
	return func(agentState AgentState) bool {
		return makePredicate(cond.CondA())(agentState) || makePredicate(cond.CondB())(agentState)
	}
}

func defectorEval() func(AgentState) bool {
	return func(agentState AgentState) bool {
		return agentState.Defector.IsDefector()
	}
}

func buildCompPredicate(condT proposal.ComparativeCondition) func(agentState AgentState) bool {
	return func(agentState AgentState) bool {
		var attr uint
		switch condT.Attribute {
		case proposal.Health:
			attr = agentState.Hp
		case proposal.Stamina:
			attr = agentState.Stamina
		case proposal.TotalAttack:
			attr = agentState.TotalAttack()
		case proposal.TotalDefence:
			attr = agentState.TotalDefense()
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
