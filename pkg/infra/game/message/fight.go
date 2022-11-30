package message

import (
	"infra/game/commons"
	"infra/game/decision"

	"github.com/benbjohnson/immutable"
)

type FightProposalMessage struct {
	sender     commons.ID
	proposal   immutable.Map[commons.ID, decision.FightAction]
	proposalId commons.ProposalID
}

func (f FightProposalMessage) Proposal() immutable.Map[commons.ID, decision.FightAction] {
	return f.proposal
}

func (f FightProposalMessage) ProposalId() commons.ProposalID {
	return f.proposalId
}

type ProposalPayload struct {
	internalMap immutable.Map[commons.ID, decision.FightAction]
}

func (p ProposalPayload) isPayload() {
}

func NewFightProposalMessage(taggedMessage TaggedMessage) *FightProposalMessage {
	return &FightProposalMessage{
		sender:     taggedMessage.Sender(),
		proposal:   taggedMessage.message.Payload().(ProposalPayload).internalMap,
		proposalId: taggedMessage.mId.String(),
	}
}
