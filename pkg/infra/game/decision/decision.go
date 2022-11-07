package decision

type Decision interface {
	decisionSealed()
}

type FightDecision struct {
	Action FightAction
}

func (FightDecision) decisionSealed() {}

type LootDecision struct{}

func (LootDecision) decisionSealed() {}

type HPPoolDecision struct{}

func (HPPoolDecision) decisionSealed() {}

type FightAction interface {
	actionSealed()
}

type Cower struct{}

func (Cower) actionSealed() {}

type Fight struct {
	Attack uint
	Defend uint
}

func (Fight) actionSealed() {}
