package message

import (
	"infra/game/commons"
	"infra/game/decision"

	"github.com/benbjohnson/immutable"
)

type MapProposal[A decision.ProposalAction] struct {
	proposalID commons.ProposalID
	proposal   immutable.Map[commons.ID, A]
}

func (p MapProposal[A]) sealedProposal() {
}

func (p MapProposal[A]) sealedMessage() {
}

func (p MapProposal[A]) ProposalID() commons.ProposalID {
	return p.proposalID
}

func (p MapProposal[A]) Proposal() immutable.Map[commons.ID, A] {
	return p.proposal
}

func NewProposal[A decision.ProposalAction](proposalID commons.ProposalID, proposalMap immutable.Map[commons.ID, A]) *MapProposal[A] {
	return &MapProposal[A]{
		proposalID: proposalID,
		proposal:   proposalMap,
	}
}
