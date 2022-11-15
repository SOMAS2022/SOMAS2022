package decision

type LootDecision struct{}

type HPPoolDecision struct{}

type FightAction int64

const (
	Attack FightAction = iota
	Defend
	Cower
	Undecided
)
