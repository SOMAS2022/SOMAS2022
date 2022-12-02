package message

import (
	"infra/game/commons"

	"github.com/google/uuid"
)

type Message interface {
	sealedMessage()
}

type Inform interface {
	Message
	sealedInform()
}

type Request interface {
	Message
	sealedRequest()
}

type Proposal interface {
	Message
	sealedProposal()
}

type FightRequest interface {
	Request
	sealedFightRequest()
}

type FightInform interface {
	Inform
	sealedFightInform()
}

type TaggedMessage struct {
	sender  commons.ID
	message Message
	mID     uuid.UUID
}

func NewTaggedMessage(sender commons.ID, message Message, mID uuid.UUID) *TaggedMessage {
	return &TaggedMessage{sender: sender, message: message, mID: mID}
}

func (t TaggedMessage) Sender() commons.ID {
	return t.sender
}

func (t TaggedMessage) Message() Message {
	return t.message
}
