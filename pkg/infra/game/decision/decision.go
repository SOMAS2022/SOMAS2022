package decision

import "infra/game/commons"

type LootDecision struct{}

type HPPoolDecision struct{}

type FightAction int64

const (
	Attack FightAction = iota
	Defend
	Cower
	Undecided
)

type ElectionParams struct {
	candidateList       []commons.ID
	strategy            VotingStrategy
	numberOfPreferences uint
}

// Intent is used for polling.
// Positive can mean true/agree/have confidence
// Negative can mean false/disagree/don't have confidence
// Abstain means ambivalenc
type Intent uint

const (
	Positive Intent = iota
	Negative
	Abstain
)

// Ballot used for leader election
// It is defined as an array of string so that it can work with different voting methods.
// e.g. 1 candidate in choose-one voting and >1 candidates in ranked voting
type Ballot []commons.ID

type VotingManifesto struct {
	Running         bool
	FightImposition bool
	LootImposition  bool
	ResignThreshold uint
	TermLength      uint
}

type VotingStrategy uint

const (
	SingleChoicePlurality = iota
)
