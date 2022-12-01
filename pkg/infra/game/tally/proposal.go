package tally

import (
	"github.com/benbjohnson/immutable"
	"infra/game/commons"
	"infra/game/decision"
)

type Proposal[A decision.ProposalAction] struct {
	proposalID commons.ProposalID
	proposal   immutable.Map[commons.ID, A]
}

func (p Proposal[A]) Proposal() immutable.Map[commons.ID, A] {
	return p.proposal
}
