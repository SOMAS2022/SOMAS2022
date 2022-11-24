package decision

import "infra/game/commons"

type LootDecision struct{}

type HPPoolDecision struct{}

type FightAction int64

// Ballot used for leader election
// It is defined as an array of string so that it can work with different voting methods.
// e.g. 1 candidate in choose-one voting and >1 candidates in ranked voting
type Ballot []commons.ID

const (
	Attack FightAction = iota
	Defend
	Cower
	Undecided
)
