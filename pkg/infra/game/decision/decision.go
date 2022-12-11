package decision

import (
	"infra/game/commons"
)

type ProposalAction interface {
	FightAction | LootAction
}

type HPPoolDecision struct{}

// Intent is used for polling.
// Positive can mean true/agree/have confidence
// Negative can mean false/disagree/don't have confidence
// Abstain means ambivalence.
type Intent uint

const (
	Positive Intent = iota
	Negative
	Abstain
)

type HpPoolDonation struct {
	AgentID  commons.ID
	Donation uint
}

type ItemIdx uint
