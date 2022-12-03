package message

import (
	"infra/game/commons"
	"infra/game/decision"

	"github.com/benbjohnson/immutable"
)

type FightProposalMessage struct {
	sender     commons.ID
	proposal   immutable.Map[commons.ID, decision.FightAction]
	proposalID commons.ProposalID
}

func NewFightProposalMessage(sender commons.ID, proposal immutable.Map[commons.ID, decision.FightAction], proposalID commons.ProposalID) *FightProposalMessage {
	return &FightProposalMessage{sender: sender, proposal: proposal, proposalID: proposalID}
}

func (f FightProposalMessage) sealedMessage() {
	panic("implement me")
}

func (f FightProposalMessage) Proposal() immutable.Map[commons.ID, decision.FightAction] {
	return f.proposal
}

func (f FightProposalMessage) ProposalID() commons.ProposalID {
	return f.proposalID
}
