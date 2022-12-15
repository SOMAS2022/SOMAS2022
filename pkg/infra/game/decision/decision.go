package decision

import (
	"infra/game/commons"

	"github.com/benbjohnson/immutable"
)

type ProposalAction interface {
	FightAction | LootAction
}

type HPPoolDecision struct{}

type Manifesto struct {
	fightDecisionPower bool
	lootDecisionPower  bool
	termLength         uint
	overthrowThreshold uint
}

func (m Manifesto) FightDecisionPower() bool {
	return m.fightDecisionPower
}

func (m Manifesto) LootDecisionPower() bool {
	return m.lootDecisionPower
}

func (m Manifesto) TermLength() uint {
	return m.termLength
}

func (m Manifesto) OverthrowThreshold() uint {
	return m.overthrowThreshold
}

func NewManifesto(fightDecisionPower bool, lootDecisionPower bool, termLength uint, overthrowThreshold uint) *Manifesto {
	return &Manifesto{
		fightDecisionPower: fightDecisionPower,
		lootDecisionPower:  lootDecisionPower,
		termLength:         termLength,
		overthrowThreshold: overthrowThreshold,
	}
}

type ElectionParams struct {
	candidateList       *immutable.Map[commons.ID, Manifesto]
	strategy            VotingStrategy
	numberOfPreferences uint
}

func (e ElectionParams) CandidateList() *immutable.Map[commons.ID, Manifesto] {
	return e.candidateList
}

func NewElectionParams(candidateList map[commons.ID]Manifesto, strategy VotingStrategy, numberOfPreferences uint) *ElectionParams {
	builder := immutable.NewMapBuilder[commons.ID, Manifesto](nil)
	for id, manifesto := range candidateList {
		builder.Set(id, manifesto)
	}
	return &ElectionParams{candidateList: builder.Map(), strategy: strategy, numberOfPreferences: numberOfPreferences}
}

func (e ElectionParams) Strategy() VotingStrategy {
	return e.strategy
}

func (e ElectionParams) NumberOfPreferences() uint {
	return e.numberOfPreferences
}

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

// Ballot used for leader election
// It is defined as an array of string so that it can work with different voting methods.
// e.g. 1 candidate in choose-one voting and >1 candidates in ranked voting.
type Ballot []commons.ID

type VotingStrategy uint

const (
	SingleChoicePlurality = iota
	BordaCount
)

type HpPoolDonation struct {
	AgentID  commons.ID
	Donation uint
}

type ItemIdx uint
