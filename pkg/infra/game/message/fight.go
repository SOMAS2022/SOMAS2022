package message

import (
	"github.com/benbjohnson/immutable"
	"infra/game/commons"
	"infra/game/decision"
)

type FightProposalMessage struct {
	sender     commons.ID
	proposal   immutable.Map[commons.ID, decision.FightAction]
	proposalID commons.ProposalID
}

func (f FightProposalMessage) Proposal() immutable.Map[commons.ID, decision.FightAction] {
	return f.proposal
}

func (f FightProposalMessage) ProposalID() commons.ProposalID {
	return f.proposalID
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
		proposalID: taggedMessage.mID.String(),
	}
}
