package tally

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/tally/internal"

	"github.com/benbjohnson/immutable"
)

type AcceptedProposal[A decision.ProposalAction] struct {
	proposalId commons.ProposalID
	proposal   immutable.Map[commons.ID, A]
}

type Tally[A decision.ProposalAction] struct {
	proposalTally immutable.Map[commons.ProposalID, uint]
	proposalMap   map[commons.ProposalID]immutable.Map[commons.ID, A]
	currMax       internal.VoteCount
	vote          <-chan commons.ProposalID
}

func NewTally[A decision.ProposalAction](vote <-chan commons.ProposalID) *Tally[A] {
	return &Tally[A]{
		proposalTally: *immutable.NewMapBuilder[commons.ProposalID, uint](nil).Map(),
		proposalMap:   make(map[commons.ProposalID]immutable.Map[commons.ID, A]),
		vote:          vote}
}

func (t *Tally[A]) getProposal(id commons.ProposalID) {

}

//func newTally() *tally {
//	return &tally{proposalTally: make(map[commons.ProposalID]uint)}
//}
//
//func (t *tally) addProposal() {
//}
//
//func (t *tally) handleVote() {
//}
