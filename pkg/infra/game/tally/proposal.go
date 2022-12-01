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
