package tally

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/tally/internal"

	"github.com/benbjohnson/immutable"
)

type Proposal[A decision.ProposalAction] struct {
	proposalId commons.ProposalID
	proposal   immutable.Map[commons.ID, A]
}

func (p Proposal[A]) Proposal() immutable.Map[commons.ID, A] {
	return p.proposal
}

type Tally[A decision.ProposalAction] struct {
	proposalTally map[commons.ProposalID]uint
	proposalMap   map[commons.ProposalID]immutable.Map[commons.ID, A]
	currMax       internal.VoteCount
	votes         <-chan commons.ProposalID
	proposals     <-chan Proposal[A]
	closure       <-chan struct{}
}

func (t *Tally[A]) ProposalTally() map[commons.ProposalID]uint {
	return t.proposalTally
}

func (t *Tally[A]) ProposalMap() map[commons.ProposalID]immutable.Map[commons.ID, A] {
	return t.proposalMap
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

func NewProposal[A decision.ProposalAction](proposalId commons.ProposalID, proposalMap immutable.Map[commons.ID, A]) *Proposal[A] {
	return &Proposal[A]{
		proposalId: proposalId,
		proposal:   proposalMap,
	}
}

// HandleMessages call from goroutine
func (t *Tally[A]) HandleMessages() {
	for {
		select {
		case proposal := <-t.proposals:
			t.proposalMap[proposal.proposalId] = proposal.proposal
			t.proposalTally[proposal.proposalId] = 0
		case vote := <-t.votes:
			t.proposalTally[vote]++
			if t.currMax.Count < t.proposalTally[vote] {
				t.currMax.Id = vote
				t.currMax.Count = t.proposalTally[vote]
			}
		case <-t.closure:
			return
		}
	}
}

// GetMax call from thread after goroutine closes
func (t *Tally[A]) GetMax() Proposal[A] {
	return *NewProposal[A](t.currMax.Id, t.proposalMap[t.currMax.Id])
}
