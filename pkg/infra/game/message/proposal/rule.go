package proposal

import "infra/game/decision"

type Rule[A decision.ProposalAction] struct {
	action    A
	condition Condition
}

func NewRule[A decision.ProposalAction](action A, condition Condition) *Rule[A] {
	return &Rule[A]{action: action, condition: condition}
}
