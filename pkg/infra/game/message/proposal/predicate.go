package proposal

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
)

func ToSinglePredicate[A decision.ProposalAction](rules commons.ImmutableList[Rule[A]]) func(state.State, state.AgentState) A {
	iterator := rules.Iterator()
	predicates := make([]func(s state.State, agentState state.AgentState) (A, bool), 0)
	for !iterator.Done() {
		rule, _ := iterator.Next()
		pred := makePredicate(rule.condition)
		wrappedPredicate := func(s state.State, agentState state.AgentState) (A, bool) {
			return rule.action, pred(s, agentState)
		}
		predicates = append(predicates, wrappedPredicate)
	}
	if len(predicates) > 0 {
		return func(s state.State, agentState state.AgentState) A {
			for _, predicate := range predicates {
				action, match := predicate(s, agentState)
				if match {
					return action
				}
			}
			// todo: what to do if unallocated with parameterized class
			a, _ := predicates[len(predicates)-1](s, agentState)
			return a
		}
	}
	return nil
}

func ToMultiPredicate[A decision.ProposalAction](rules commons.ImmutableList[Rule[A]]) func(state.State, state.AgentState) map[A]struct{} {
	iterator := rules.Iterator()
	predicates := make([]func(s state.State, agentState state.AgentState) (A, bool), 0)
	for !iterator.Done() {
		rule, _ := iterator.Next()
		pred := makePredicate(rule.condition)
		wrappedPredicate := func(s state.State, agentState state.AgentState) (A, bool) {
			return rule.action, pred(s, agentState)
		}
		predicates = append(predicates, wrappedPredicate)
	}
	if len(predicates) > 0 {
		return func(s state.State, agentState state.AgentState) map[A]struct{} {
			res := make(map[A]struct{})
			for _, predicate := range predicates {
				action, match := predicate(s, agentState)
				if match {
					res[action] = struct{}{}
				}
			}
			return res
		}
	}
	return nil
}

func makePredicate(cond Condition) func(s state.State, agentState state.AgentState) bool {
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
		return func(_ state.State, _ state.AgentState) bool {
			return true
		}
	}
}

func andEval(cond AndCondition) func(state.State, state.AgentState) bool {
	return func(s state.State, agentState state.AgentState) bool {
		return makePredicate(cond.CondA())(s, agentState) && makePredicate(cond.CondB())(s, agentState)
	}
}

func orEval(cond OrCondition) func(state.State, state.AgentState) bool {
	return func(s state.State, agentState state.AgentState) bool {
		return makePredicate(cond.CondA())(s, agentState) || makePredicate(cond.CondB())(s, agentState)
	}
}

func defectorEval() func(state.State, state.AgentState) bool {
	return func(_ state.State, agentState state.AgentState) bool {
		return agentState.Defector.IsDefector()
	}
}

func buildCompPredicate(condT ComparativeCondition) func(s state.State, agentState state.AgentState) bool {
	return func(s state.State, agentState state.AgentState) bool {
		var attr uint
		switch condT.Attribute {
		case Health:
			attr = agentState.Hp
		case Stamina:
			attr = agentState.Stamina
		case TotalAttack:
			attr = agentState.TotalAttack(s)
		case TotalDefence:
			attr = agentState.TotalDefense(s)
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
