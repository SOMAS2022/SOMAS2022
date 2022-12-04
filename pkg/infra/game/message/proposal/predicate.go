package proposal

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
)

func ToPredicate[A decision.ProposalAction](rules commons.ImmutableList[Rule[A]]) func(state.State, state.AgentState) A {
	iterator := rules.Iterator()
	predicates := make([]func(s state.State, agentState state.AgentState) (A, bool), 0)
	for !iterator.Done() {
		rule, _ := iterator.Next()
		switch cond := rule.condition.(type) {
		case *ComparativeCondition:
			predicates = append(predicates, makePredicate(*cond, rule))
		}
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

func makePredicate[A decision.ProposalAction](cond ComparativeCondition, rule Rule[A]) func(s state.State, agentState state.AgentState) (A, bool) {
	return func(s state.State, agentState state.AgentState) (A, bool) {
		var attr uint
		switch cond.Attribute {
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
		switch cond.Comparator {
		case GreaterThan:
			if attr > cond.Value {
				return rule.action, true
			} else {
				return rule.action, false
			}
		default:
			if attr < cond.Value {
				return rule.action, true
			} else {
				return rule.action, false
			}
		}
	}
}
