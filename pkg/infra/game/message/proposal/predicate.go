package proposal

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
)

func ToSinglePredicate[A decision.ProposalAction](rules commons.ImmutableList[Rule[A]]) func(state.AgentState) A {
	iterator := rules.Iterator()
	predicates := make([]func(agentState state.AgentState) (A, bool), 0)
	for !iterator.Done() {
		rule, _ := iterator.Next()
		pred := makePredicate(rule.condition)
		wrappedPredicate := func(agentState state.AgentState) (A, bool) {
			return rule.action, pred(agentState)
		}
		predicates = append(predicates, wrappedPredicate)
	}
	if len(predicates) > 0 {
		return func(agentState state.AgentState) A {
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

func ToMultiPredicate[A decision.ProposalAction](rules commons.ImmutableList[Rule[A]]) func(state.AgentState) map[A]struct{} {
	iterator := rules.Iterator()
	predicates := make([]func(agentState state.AgentState) (A, bool), 0)
	for !iterator.Done() {
		rule, _ := iterator.Next()
		pred := makePredicate(rule.condition)
		wrappedPredicate := func(agentState state.AgentState) (A, bool) {
			return rule.action, pred(agentState)
		}
		predicates = append(predicates, wrappedPredicate)
	}
	if len(predicates) > 0 {
		return func(agentState state.AgentState) map[A]struct{} {
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

func makePredicate(cond Condition) func(agentState state.AgentState) bool {
	switch condT := cond.(type) {
	case *ComparativeCondition:
		return buildCompPredicate(*condT)
	case ComparativeCondition:
		return buildCompPredicate(condT)
	case *AndCondition:
		return andEval(*condT)
	case AndCondition:
		return andEval(condT)
	case *OrCondition:
		return orEval(*condT)
	case OrCondition:
		return orEval(condT)
	case DefectorCondition:
		return defectorEval()
	default:
		return func(_ state.AgentState) bool {
			return true
		}
	}
}

func andEval(cond AndCondition) func(state.AgentState) bool {
	return func(agentState state.AgentState) bool {
		return makePredicate(cond.CondA())(agentState) && makePredicate(cond.CondB())(agentState)
	}
}

func orEval(cond OrCondition) func(state.AgentState) bool {
	return func(agentState state.AgentState) bool {
		return makePredicate(cond.CondA())(agentState) || makePredicate(cond.CondB())(agentState)
	}
}

func defectorEval() func(state.AgentState) bool {
	return func(agentState state.AgentState) bool {
		return agentState.Defector.IsDefector()
	}
}

func buildCompPredicate(condT ComparativeCondition) func(agentState state.AgentState) bool {
	return func(agentState state.AgentState) bool {
		var attr uint
		switch condT.Attribute {
		case Health:
			attr = agentState.Hp
		case Stamina:
			attr = agentState.Stamina
		case TotalAttack:
			attr = agentState.TotalAttack()
		case TotalDefence:
			attr = agentState.TotalDefense()
		default:
			attr = agentState.Hp
		}
		switch condT.Comparator {
		case GreaterThan:
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
