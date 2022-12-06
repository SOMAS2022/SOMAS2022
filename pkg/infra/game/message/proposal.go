package message

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message/proposal"

	"github.com/google/uuid"
)

type Proposal[A decision.ProposalAction] struct {
	proposalID commons.ProposalID
	rules      commons.ImmutableList[proposal.Rule[A]]
}

func (p Proposal[A]) ProposalID() commons.ProposalID {
	return p.proposalID
}

func (p Proposal[A]) Rules() commons.ImmutableList[proposal.Rule[A]] {
	return p.rules
}

func (p Proposal[A]) sealedMessage() {
}

func NewProposal[A decision.ProposalAction](rules commons.ImmutableList[proposal.Rule[A]]) *Proposal[A] {
	return &Proposal[A]{proposalID: uuid.NewString(), rules: rules}
}

func NewProposalInternal[A decision.ProposalAction](proposalID commons.ProposalID, rules commons.ImmutableList[proposal.Rule[A]]) *Proposal[A] {
	return &Proposal[A]{proposalID: proposalID, rules: rules}
}
