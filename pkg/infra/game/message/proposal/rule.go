package proposal

import "infra/game/decision"

type Rule[A decision.ProposalAction] struct {
	action    A
	condition Condition
}

func (r Rule[A]) Action() A {
	return r.action
}

func (r Rule[A]) Condition() Condition {
	return r.condition
}

func NewRule[A decision.ProposalAction](action A, condition Condition) *Rule[A] {
	return &Rule[A]{action: action, condition: condition}
}
