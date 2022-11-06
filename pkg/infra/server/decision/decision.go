package decision

type Decision interface {
	sealed()
}

type FightDecision struct{}

func (FightDecision) sealed() {}

type LootDecision struct{}

func (LootDecision) sealed() {}

type HPPoolDecision struct{}

func (HPPoolDecision) sealed() {}
