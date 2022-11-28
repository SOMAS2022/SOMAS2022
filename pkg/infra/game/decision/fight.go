package decision

type FightAction int64

const (
	Undecided FightAction = iota
	Defend
	Cower
	Attack
)
