package tally

import (
	"github.com/google/uuid"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/tally/internal"

	"github.com/benbjohnson/immutable"
)

type Proposal[A decision.ProposalAction] struct {
	proposalId commons.ProposalID
	proposal   immutable.Map[commons.ID, A]
}

type Tally[A decision.ProposalAction] struct {
	proposalTally map[commons.ProposalID]uint
	proposalMap   map[commons.ProposalID]immutable.Map[commons.ID, A]
	currMax       internal.VoteCount
	votes         <-chan commons.ProposalID
	proposals     <-chan Proposal[A]
	closure       <-chan struct{}
}

func NewTally[A decision.ProposalAction](votes <-chan commons.ProposalID,
	proposals <-chan Proposal[A],
	closure <-chan struct{}) *Tally[A] {
	return &Tally[A]{
		proposalTally: make(map[commons.ProposalID]uint),
		proposalMap:   make(map[commons.ProposalID]immutable.Map[commons.ID, A]),
		votes:         votes,
		proposals:     proposals,
		closure:       closure,
	}
}

func NewProposal[A decision.ProposalAction](proposalMap map[commons.ID]A) Proposal[A] {
	builder := immutable.NewMapBuilder[commons.ID, A](nil)
	for id, a := range proposalMap {
		builder.Set(id, a)
	}
	return Proposal[A]{
		proposalId: uuid.NewString(),
		proposal:   *builder.Map(),
	}
}

// call from goroutine
func (t *Tally[A]) handleMessages() {
	for {
		select {
		case proposal := <-t.proposals:
			t.proposalMap[proposal.proposalId] = proposal.proposal
			t.proposalTally[proposal.proposalId] = 0
		case vote := <-t.votes:
			t.proposalTally[vote]++
		case <-t.closure:
			return
		}
	}
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
