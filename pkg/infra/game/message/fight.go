package message

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message/proposal"
)

type FightProposalMessage struct {
	sender     commons.ID
	proposal   commons.ImmutableList[proposal.Rule[decision.FightAction]]
	proposalID commons.ProposalID
}

func NewFightProposalMessage(sender commons.ID, proposal commons.ImmutableList[proposal.Rule[decision.FightAction]], proposalID commons.ProposalID) *FightProposalMessage {
	return &FightProposalMessage{sender: sender, proposal: proposal, proposalID: proposalID}
}

func (f FightProposalMessage) sealedMessage() {
	panic("implement me")
}

func (f FightProposalMessage) Proposal() commons.ImmutableList[proposal.Rule[decision.FightAction]] {
	return f.proposal
}

func (f FightProposalMessage) ProposalID() commons.ProposalID {
	return f.proposalID
}
