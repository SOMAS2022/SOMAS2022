package message

import (
	"infra/game/commons"
	"infra/game/decision"

	"github.com/google/uuid"
)

type Proposal[A decision.ProposalAction] struct {
	proposalID commons.ProposalID
	proposerID commons.ID
	rules      commons.ImmutableList[decision.Rule[A]]
}

func (p Proposal[A]) ProposerID() commons.ID {
	return p.proposerID
}

func (p Proposal[A]) ProposalID() commons.ProposalID {
	return p.proposalID
}

func (p Proposal[A]) Rules() commons.ImmutableList[decision.Rule[A]] {
	return p.rules
}

func (p Proposal[A]) sealedMessage() {
}

func NewProposal[A decision.ProposalAction](rules commons.ImmutableList[decision.Rule[A]], proposerID commons.ID) *Proposal[A] {
	return &Proposal[A]{proposalID: uuid.NewString(), rules: rules, proposerID: proposerID}
}

func NewProposalInternal[A decision.ProposalAction](proposalID commons.ProposalID, rules commons.ImmutableList[decision.Rule[A]]) *Proposal[A] {
	return &Proposal[A]{proposalID: proposalID, rules: rules, proposerID: uuid.Nil.String()}
}
