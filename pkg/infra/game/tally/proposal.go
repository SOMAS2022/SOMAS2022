package tally

import (
	"infra/game/commons"
	"infra/game/decision"

	"github.com/benbjohnson/immutable"
)

type Proposal[A decision.ProposalAction] struct {
	proposalID commons.ProposalID
	proposal   immutable.Map[commons.ID, A]
}

func (p Proposal[A]) Proposal() immutable.Map[commons.ID, A] {
	return p.proposal
}

func NewProposal[A decision.ProposalAction](proposalID commons.ProposalID, proposalMap immutable.Map[commons.ID, A]) *Proposal[A] {
	return &Proposal[A]{
		proposalID: proposalID,
		proposal:   proposalMap,
	}
}
