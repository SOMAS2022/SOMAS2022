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
	proposalId uuid.UUID
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
		proposalId: taggedMessage.mId,
	}
}
