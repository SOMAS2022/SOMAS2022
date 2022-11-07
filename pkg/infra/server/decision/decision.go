package decision

type Decision interface {
	decisionSealed()
}

type FightDecision struct {
	Choice FightChoice
}

func (FightDecision) decisionSealed() {}

type LootDecision struct{}

func (LootDecision) decisionSealed() {}

type HPPoolDecision struct{}

func (HPPoolDecision) decisionSealed() {}

type FightChoice interface {
	fightChoiceSealed()
}

type Attack struct{}

func (Attack) fightChoiceSealed() {}

type Defend struct{}

func (Defend) fightChoiceSealed() {}

type Cower struct{}

func (Cower) fightChoiceSealed() {}
