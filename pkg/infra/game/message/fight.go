package message

import (
	"infra/game/commons"
	"infra/game/decision"

	"github.com/benbjohnson/immutable"
	"github.com/google/uuid"
)

type FightProposalMessage struct {
	sender     commons.ID
	proposal   immutable.Map[commons.ID, decision.FightAction]
	proposalID commons.ProposalID
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

func NewFightProposalMessage(senderID commons.ID, proposal immutable.Map[commons.ID, decision.FightAction]) *FightProposalMessage {
	newUuid, _ := uuid.NewUUID()

	return &FightProposalMessage{
		sender:     senderID,
		proposal:   proposal,
		proposalID: newUuid.String(),
	}
}
