package decision

type ProposalAction interface {
	FightAction | LootAction
}

type Rule[A ProposalAction] struct {
	action    A
	condition Condition
}

func (r Rule[ProposalAction]) Action() ProposalAction {
	return r.action
}

func (r Rule[ProposalAction]) Condition() Condition {
	return r.condition
}

func NewRule[A ProposalAction](action A, condition Condition) *Rule[A] {
	return &Rule[A]{action: action, condition: condition}
}
