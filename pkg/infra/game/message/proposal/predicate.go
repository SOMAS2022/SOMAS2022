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
		makePredicate(rule)
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
		makePredicate(rule)
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

func makePredicate[A decision.ProposalAction](rule Rule[A]) func(s state.State, agentState state.AgentState) (A, bool) {
	switch condT := rule.condition.(type) {
	case ComparativeCondition:
		return buildCompPredicate(rule, condT)
	default:
		return func(s state.State, agentState state.AgentState) (A, bool) {
			//todo: use state/agentState to check if in recent defector set
			return rule.action, true
		}
	}
}

func buildCompPredicate[A decision.ProposalAction](rule Rule[A], condT ComparativeCondition) func(s state.State, agentState state.AgentState) (A, bool) {
	return func(s state.State, agentState state.AgentState) (A, bool) {
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
				return rule.action, true
			} else {
				return rule.action, false
			}
		default:
			if attr < condT.Value {
				return rule.action, true
			} else {
				return rule.action, false
			}
		}
	}
}
