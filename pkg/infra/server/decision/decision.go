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

type FightAction int64

const (
	Attack FightAction = iota
	Defend
	Cower
)
